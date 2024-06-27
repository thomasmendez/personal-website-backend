package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/thomasmendez/personal-website-backend/api/database"
)

type Service struct {
	DB        *database.Database
	TableName string
	Routes    *[]RouteHandler
}

type RouteHandler struct {
	Route   string
	Method  string
	Handler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

func NewService() *Service {
	var awsSession *session.Session
	var tableName string

	env := os.Getenv("ENV")

	tableName = os.Getenv("TABLE_NAME")
	if tableName == "" {
		log.Fatal("error in configuration: TABLE_NAME env not provided")
	}

	if env != "Local" {
		if env != "Dev" && env != "Stg" && env != "Prd" {
			log.Fatalf("error in configuration: ENV must be 'Dev', 'Stg', or 'Prd' \n Currently: %v", env)
		}
		if env == "Stg" || env == "Prd" {
			if os.Getenv("ORIGIN") == "" {
				log.Fatalf("error in configuration: ENV 'Stg' requires 'ORIGIN' variable")
			}
			if os.Getenv("HEADERS") == "" {
				log.Fatalf("error in configuration: ENV 'Stg' requires 'HEADERS' variable")
			}
			if os.Getenv("METHODS") == "" {
				log.Fatalf("error in configuration: ENV 'Stg' requires 'METHODS' variable")
			}
		}
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	} else {
		region := os.Getenv("REGION")
		if region == "" {
			log.Fatal("error in configuration: REGION env not provided")
		}
		awsSession = session.Must(session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Endpoint: aws.String("http://dynamodb:8000"),
				Region:   aws.String(region),
			},
			SharedConfigState: session.SharedConfigEnable,
		}))
	}

	s := &Service{
		DB:        database.NewDatabase(awsSession),
		TableName: tableName,
	}

	s.Routes = addRoutes(s)

	return s
}

func (s *Service) HandleRoute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, route := range *s.Routes {
		if request.Path == route.Route && request.HTTPMethod == route.Method {
			proxyResponse, err := route.Handler(ctx, request)
			proxyResponse.Headers = s.addProxyHeaders(os.Getenv("ENV"))
			return proxyResponse, err
		}
	}

	errRes := ErrorResponse{
		Message: "Route not found",
	}
	res, _ := json.Marshal(errRes)

	return events.APIGatewayProxyResponse{
		Headers:    s.addProxyHeaders(os.Getenv("ENV")),
		StatusCode: http.StatusNotFound,
		Body:       string(res),
	}, nil
}

func (s *Service) addProxyHeaders(env string) map[string]string {
	switch env {
	case "Dev":
		return map[string]string{
			"Access-Control-Allow-Origin":  "http://localhost:5173",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	case "Stg":
		return map[string]string{
			"Access-Control-Allow-Origin":  os.Getenv("ORIGIN"),
			"Access-Control-Allow-Headers": os.Getenv("HEADERS"),
			"Access-Control-Allow-Methods": os.Getenv("METHODS"),
		}
	case "Prd":
		return map[string]string{
			"Access-Control-Allow-Origin":  os.Getenv("ORIGIN"),
			"Access-Control-Allow-Headers": os.Getenv("HEADERS"),
			"Access-Control-Allow-Methods": os.Getenv("METHODS"),
		}
	default:
		return map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	}
}

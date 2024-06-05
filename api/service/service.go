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
			// proxyResponse.Headers = s.addProxyHeaders("dev")
			return proxyResponse, err
		}
	}

	errRes := ErrorResponse{
		Message: "Route not found",
	}
	res, _ := json.Marshal(errRes)

	return events.APIGatewayProxyResponse{
		// Headers:    s.addProxyHeaders("dev"),
		StatusCode: http.StatusNotFound,
		Body:       string(res),
	}, nil
}

// func (s *Service) addProxyHeaders(env string) map[string]string {
// 	switch env {
// 	case "dev":
// 		return map[string]string{
// 			"Access-Control-Allow-Origin":  "*",
// 			"Access-Control-Allow-Headers": "*",
// 			"Access-Control-Allow-Methods": "*",
// 		}
// 	default:
// 		return map[string]string{
// 			"Access-Control-Allow-Origin":  "*",
// 			"Access-Control-Allow-Headers": "*",
// 			"Access-Control-Allow-Methods": "*",
// 		}
// 	}
// }

package services

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/thomasmendez/personal-website-backend/api/database"
)

type Service struct {
	DB     *database.Database
	Routes *[]RouteHandler
}

type RouteHandler struct {
	Route   string
	Method  string
	Handler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewService() *Service {
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
			Region:   aws.String("us-west-2"),
		},
	}))

	s := &Service{
		DB: database.NewDatabase(awsSession),
	}

	s.Routes = &[]RouteHandler{
		{
			Route:   "/api/v1/work",
			Method:  http.MethodGet,
			Handler: s.getWorkHandler,
		},
		{
			Route:   "/api/v1/work",
			Method:  http.MethodPost,
			Handler: s.postWorkHandler,
		},
		{
			Route:   "/api/v1/work",
			Method:  http.MethodPut,
			Handler: s.updateWorkHandler,
		},
	}

	return s
}

func (s *Service) HandleRoute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, route := range *s.Routes {
		if request.Path == route.Route && request.HTTPMethod == route.Method {
			proxyResponse, err := route.Handler(ctx, request)
			proxyResponse.Headers = s.addProxyHeaders("dev")
			return proxyResponse, err
		}
	}

	errRes := ErrorResponse{
		Message: "Route not found",
	}
	res, _ := json.Marshal(errRes)

	return events.APIGatewayProxyResponse{
		Headers:    s.addProxyHeaders("dev"),
		StatusCode: http.StatusNotFound,
		Body:       string(res),
	}, nil
}

func (s *Service) addProxyHeaders(env string) map[string]string {
	switch env {
	case "dev":
		return map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	default:
		return map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "*",
			"Access-Control-Allow-Methods": "*",
		}
	}
}

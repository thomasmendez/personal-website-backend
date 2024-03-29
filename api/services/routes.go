package services

import (
	"context"
	"encoding/json"
	"log"
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
	Handler func() (events.APIGatewayProxyResponse, error)
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
	}

	return s
}

func (s *Service) HandleRoute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, route := range *s.Routes {
		if request.Path == route.Route && request.HTTPMethod == route.Method {
			return route.Handler()
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       http.StatusText(http.StatusNotFound),
	}, nil
}

func (s *Service) getWorkHandler() (events.APIGatewayProxyResponse, error) {

	work, err := s.DB.GetWork()

	if err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string("done"),
		}, err
	}

	workJson, err := json.Marshal(work)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(workJson),
	}, err
}

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
	"github.com/thomasmendez/personal-website-backend/api/models"
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
	}

	return s
}

func (s *Service) HandleRoute(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	for _, route := range *s.Routes {
		if request.Path == route.Route && request.HTTPMethod == route.Method {
			return route.Handler(ctx, request)
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       http.StatusText(http.StatusNotFound),
	}, nil
}

func (s *Service) getWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

func (s *Service) postWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var newWork models.Work
	err := json.Unmarshal([]byte(request.Body), &newWork)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	work, err := s.DB.PostWork(newWork)

	if err != nil {
		log.Print(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string("done"),
		}, err
	}

	workJson, err := json.Marshal(work)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(workJson),
	}, err
}

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (s *Service) getWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	work, err := s.DB.GetWork()

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: "There was an error in getting work",
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	workJson, err := json.Marshal(work)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
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
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting work with startDate of: %s", newWork.StartDate),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	workJson, err := json.Marshal(work)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(workJson),
	}, err
}

func (s *Service) updateWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var updateWork models.Work
	err := json.Unmarshal([]byte(request.Body), &updateWork)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	work, err := s.DB.PostWork(updateWork)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in updating work with startDate of: %s", updateWork.StartDate),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	workJson, err := json.Marshal(work)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(workJson),
	}, err
}
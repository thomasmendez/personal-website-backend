package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasmendez/personal-website-backend/api/database"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (s *Service) getWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	work, err := database.GetWork(s.DB)

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
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	work, err := database.PostWork(s.DB, newWork)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting work with sortValue of: %s", newWork.SortValue),
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
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	work, err := database.PostWork(s.DB, updateWork)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in updating work with sortValue of: %s", updateWork.SortValue),
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

func (s *Service) deleteWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var deleteWork models.Work
	err := json.Unmarshal([]byte(request.Body), &deleteWork)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	var existingWork models.Work
	err = database.GetItem(s.DB, deleteWork.PersonalWebsiteType, deleteWork.SortValue, &existingWork)

	if !reflect.DeepEqual(deleteWork, existingWork) {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       resError(http.StatusNotFound),
		}, err
	}

	err = database.DeleteItem(s.DB, deleteWork.PersonalWebsiteType, deleteWork.SortValue)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in deleting work with sortValue of: %s", deleteWork.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Resource was successfully deleted",
	}, err
}

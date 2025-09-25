package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasmendez/personal-website-backend/api/database"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (s *Service) getWorkHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	work, err := database.GetWork(ctx, s.DB.Client, s.TableName)

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

	err = s.validateWork(newWork)
	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting work: %s", err),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(res),
		}, nil
	}

	work, err := database.PostWork(ctx, s.DB.Client, s.TableName, newWork)

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

func (s *Service) validateWork(work models.Work) error {
	if work.PersonalWebsiteType == "" {
		return errors.New("personalWebsiteType cannot be empty")
	}
	if work.SortValue == "" {
		return errors.New("sortValue cannot be empty")
	}
	if work.JobDescription == nil {
		return errors.New("jobDescription cannot be empty")
	}
	return nil
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

	err = s.validateWork(updateWork)
	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting work: %s", err),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(res),
		}, nil
	}

	work, err := database.UpdateWork(ctx, s.DB.Client, s.TableName, updateWork)

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
	err = database.GetItem(ctx, s.DB.Client, s.TableName, deleteWork.PersonalWebsiteType, deleteWork.SortValue, &existingWork)

	if !reflect.DeepEqual(deleteWork, existingWork) {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       resError(http.StatusNotFound),
		}, err
	}

	err = database.DeleteItem(ctx, s.DB.Client, s.TableName, deleteWork.PersonalWebsiteType, deleteWork.SortValue)

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

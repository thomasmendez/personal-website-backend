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

func (s *Service) getSkillsToolsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	skillsTools, err := database.GetSkillsTools(s.DB, s.TableName)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: "There was an error in getting skillsTools",
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	skillsToolsJson, err := json.Marshal(skillsTools)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(skillsToolsJson),
	}, err
}

func (s *Service) postSkillsToolsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newSkillsTools models.SkillsTools
	err := json.Unmarshal([]byte(request.Body), &newSkillsTools)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	err = s.validateSkillTools(newSkillsTools)
	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting skillsTools: %s", err),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(res),
		}, nil
	}

	skillsTools, err := database.PostSkillsTools(s.DB, s.TableName, newSkillsTools)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting skillsTools with sortValue of: %s", newSkillsTools.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	skillsToolsJson, err := json.Marshal(skillsTools)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(skillsToolsJson),
	}, err
}

func (s *Service) updateSkillsToolsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var updateSkillsTools models.SkillsTools
	err := json.Unmarshal([]byte(request.Body), &updateSkillsTools)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	err = s.validateSkillTools(updateSkillsTools)
	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting skillsTools: %s", err),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(res),
		}, nil
	}

	skillsTools, err := database.UpdateSkillsTools(s.DB, s.TableName, updateSkillsTools)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in updating skillsTools with sortValue of: %s", updateSkillsTools.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	skillsToolsJson, err := json.Marshal(skillsTools)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(skillsToolsJson),
	}, err
}

func (s *Service) validateSkillTools(skillsTools models.SkillsTools) error {
	if skillsTools.PersonalWebsiteType == "" {
		return errors.New("personalWebsiteType cannot be empty")
	}
	if skillsTools.SortValue == "" {
		return errors.New("sortValue cannot be empty")
	}
	if len(skillsTools.Categories) == 0 {
		return errors.New("categories cannot be empty")
	}
	for _, category := range skillsTools.Categories {
		if category.Category == "" {
			return errors.New("category cannot be empty")
		}
		if len(category.List) == 0 {
			return errors.New("list cannot be empty")
		}
	}
	return nil
}

func (s *Service) deleteSkillsToolsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var deleteSkillsTools models.SkillsTools
	err := json.Unmarshal([]byte(request.Body), &deleteSkillsTools)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	var existingSkillsTools models.SkillsTools
	err = database.GetItem(s.DB, s.TableName, deleteSkillsTools.PersonalWebsiteType, deleteSkillsTools.SortValue, &existingSkillsTools)

	if !reflect.DeepEqual(deleteSkillsTools, existingSkillsTools) {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       resError(http.StatusNotFound),
		}, err
	}

	err = database.DeleteItem(s.DB, s.TableName, deleteSkillsTools.PersonalWebsiteType, deleteSkillsTools.SortValue)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in deleting skillsTools with sortValue of: %s", deleteSkillsTools.SortValue),
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

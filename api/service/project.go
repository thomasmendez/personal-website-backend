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

func (s *Service) getProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projects, err := database.GetProjects(s.DB)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: "There was an error in getting projects",
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	projectsJson, err := json.Marshal(projects)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(projectsJson),
	}, err
}

func (s *Service) postProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var newProject models.Project
	err := json.Unmarshal([]byte(request.Body), &newProject)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	project, err := database.PostProject(s.DB, newProject)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting project with sortValue: %s", newProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	projectJson, err := json.Marshal(project)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(projectJson),
	}, err
}

func (s *Service) updateProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updateProject models.Project
	err := json.Unmarshal([]byte(request.Body), &updateProject)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	project, err := database.UpdateProject(s.DB, updateProject)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in updating project with sortValue of: %s", updateProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	projectJson, err := json.Marshal(project)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(projectJson),
	}, err
}

func (s *Service) deleteProjectHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var deleteProject models.Project
	err := json.Unmarshal([]byte(request.Body), &deleteProject)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	var existingProject models.Project
	err = database.GetItem(s.DB, deleteProject.PersonalWebsiteType, deleteProject.SortValue, &existingProject)

	if !reflect.DeepEqual(deleteProject, existingProject) {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Resource not found",
		}, err
	}

	err = database.DeleteItem(s.DB, deleteProject.PersonalWebsiteType, deleteProject.SortValue)

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in deleting project with sortValue of: %s", deleteProject.SortValue),
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

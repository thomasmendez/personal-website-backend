package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	work, err := database.PostWork(s.DB, newWork)

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

	work, err := database.PostWork(s.DB, updateWork)

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
		StatusCode: http.StatusOK,
		Body:       string(workJson),
	}, err
}

func (s *Service) getSkillsToolsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	skillsTools, err := database.GetSkillsTools(s.DB)

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
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	skillsTools, err := database.PostSkillsTools(s.DB, newSkillsTools)

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
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	skillsTools, err := database.UpdateSkillsTools(s.DB, updateSkillsTools)

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
			Message: fmt.Sprintf("There was an error in updating project with startDate of: %s", updateProject.StartDate),
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

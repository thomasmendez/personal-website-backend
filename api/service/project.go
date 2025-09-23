package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/thomasmendez/personal-website-backend/api/database"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (s *Service) getProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projects, err := database.GetProjects(s.DB, s.TableName)

	if err != nil {
		log.Printf("error in getting projects: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       resError(http.StatusInternalServerError),
		}, err
	}

	// Generate presigned URL for mediaLink
	for i, project := range projects {
		if project.MediaLinkIsS3Bucket() {
			if fileName, err := project.GetFileNameFromMediaLink(); fileName != "" {
				presignedReq, err := s.S3.GeneratePresignedURL(ctx, fileName)
				if err != nil {
					log.Printf("error in generating presigned URL: %v", err)
					log.Printf("skipping generation of presigned URL for project %s", project.SortValue)
				} else {
					projects[i].MediaLink = &presignedReq.URL
				}
			} else {
				log.Printf("error in getting filename from mediaLink: %v", err)
				log.Printf("skipping generation of presigned URL for project %s", project.SortValue)
			}
		}
	}

	projectsJson, err := json.Marshal(projects)

	if err != nil {
		log.Printf("error in serializing projects: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       resError(http.StatusInternalServerError),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(projectsJson),
	}, err
}

func (s *Service) postProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newProject *models.Project
	var imageFile models.FileData
	var err error

	log.Printf("POST project request: %v", request)
	if request.IsBase64Encoded || strings.Contains(getContentType(request.Headers), "'multipart/form-data") {
		fmt.Println("parsing form data")
		newProject, imageFile, err = parseFormData[models.Project](request)
		if err != nil {
			log.Printf("error parsing form data: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
	} else if !request.IsBase64Encoded && strings.Contains(getContentType(request.Headers), "application/json") {
		fmt.Println("parsing json")
		err = json.Unmarshal([]byte(request.Body), &newProject)
		if err != nil {
			log.Printf("error in deserializing json: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
		if newProject.MediaLink != nil {
			log.Printf("error: mediaLink has invalid content, please use multipart/form-data")
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
	} else {
		fmt.Println("POST project request format is invalid")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	// Upload image to S3 if it exists
	var presignedReq *v4.PresignedHTTPRequest
	if imageFile.Filename != "" && imageFile.Content != nil && imageFile.ContentType != "" {
		log.Printf("uploading image file: %s to S3", imageFile.Filename)
		mediaLink, err := s.S3.SendFileToS3(ctx, imageFile)
		if err != nil {
			fmt.Println("failed to upload to S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}
		newProject.MediaLink = &mediaLink
		// get presign url for response
		presignedReq, err = s.S3.GeneratePresignedURL(ctx, imageFile.Filename)
		if err != nil {
			log.Printf("error in generating presigned URL: %v", err)
			log.Printf("skipping generation of presigned URL for project %s", newProject.SortValue)
		}
	} else {
		log.Printf("no valid image file provided in request")
	}

	log.Printf("adding new project: %v to database", newProject)
	project, err := database.PostProject(s.DB, s.TableName, *newProject)

	if err != nil {
		log.Printf("error in inserting project: %v", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("error in inserting project: %s", newProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	// add presigned url to project response if created
	if presignedReq != nil {
		project.MediaLink = &presignedReq.URL
	}

	projectJson, err := json.Marshal(project)

	if err != nil {
		log.Printf("error in serializing project: %v", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("error in project response for: %s", newProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       string(projectJson),
	}, err
}

func (s *Service) updateProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var updateProject *models.Project
	var imageFile models.FileData
	var err error

	log.Printf("UPDATE project request: %v", request)
	if request.IsBase64Encoded || strings.Contains(getContentType(request.Headers), "'multipart/form-data") {
		log.Printf("parsing form data")
		updateProject, imageFile, err = parseFormData[models.Project](request)
		if err != nil {
			log.Printf("error parsing form data: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
	} else if !request.IsBase64Encoded && strings.Contains(getContentType(request.Headers), "application/json") {
		log.Printf("parsing json")
		err = json.Unmarshal([]byte(request.Body), &updateProject)
		if err != nil {
			log.Printf("error in deserializing json: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
		if updateProject.MediaLink != nil {
			if !strings.HasPrefix(*updateProject.MediaLink, "http") {
				log.Printf("error: mediaLink has invalid content, please use multipart/form-data")
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusBadRequest,
					Body:       resError(http.StatusBadRequest),
				}, nil
			}
		}
	} else {
		fmt.Println("PUT project request format is invalid")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	if imageFile.Filename != "" && imageFile.Content != nil && imageFile.ContentType != "" {
		log.Printf("uploading image file: %s to S3", imageFile.Filename)
		mediaLink, err := s.S3.SendFileToS3(ctx, imageFile)
		if err != nil {
			log.Printf("failed to upload to S3: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}
		updateProject.MediaLink = &mediaLink
	}

	log.Printf("updating project: %v", updateProject)
	project, err := database.UpdateProject(s.DB, s.TableName, *updateProject)

	if err != nil {
		log.Printf("error in updating project: %v", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("error in updating project with sortValue of: %s", updateProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	projectJson, err := json.Marshal(project)

	if err != nil {
		log.Printf("error in serializing project: %v", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("error in updating project response with sortValue of: %s", updateProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(projectJson),
	}, err
}

func (s *Service) deleteProjectHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var deleteProject *models.Project

	log.Printf("DELETE project request: %v", request)
	err := json.Unmarshal([]byte(request.Body), &deleteProject)
	if err != nil {
		log.Printf("error in deserializing json: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	var existingProject models.Project
	err = database.GetItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue, &existingProject)
	if err != nil {
		log.Printf("error in getting project: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       resError(http.StatusInternalServerError),
		}, err
	}

	log.Printf("existing project: %v", existingProject)
	if existingProject.MediaLinkIsS3Bucket() {
		if fileName, err := existingProject.GetFileNameFromMediaLink(); fileName != "" {
			if exists, err := s.S3.FileExistsInS3(ctx, fileName); exists {
				err = s.S3.DeleteFileFromS3(ctx, fileName)
				if err != nil {
					log.Printf("error in deleting file from S3: %v", err)
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       resError(http.StatusInternalServerError),
					}, err
				}
			} else {
				if err != nil {
					var nf *types.NoSuchKey
					if errors.As(err, &nf) {
						log.Printf("file %s does not exist in S3", fileName)
					}
					log.Printf("error in getting file from S3: %v", err)
				}
			}
		} else {
			log.Printf("error in getting filename from mediaLink %s: %v", *existingProject.MediaLink, err)
			log.Printf("skipping deletion of file from S3")
		}
	} else {
		log.Printf("mediaLink for project %s is not a valid S3 bucket link", *existingProject.MediaLink)
	}

	log.Printf("deleting project: %v", deleteProject)
	err = database.DeleteItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue)

	if err != nil {
		log.Printf("error in deleting project: %v", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("error in deleting project: %s", deleteProject.SortValue),
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

func getContentType(headers map[string]string) string {
	// Try different case variations
	if ct := headers["Content-Type"]; ct != "" {
		return ct
	}
	if ct := headers["content-type"]; ct != "" { // this is dumb
		return ct
	}
	if ct := headers["Content-type"]; ct != "" { // also dumb
		return ct
	}
	return ""
}

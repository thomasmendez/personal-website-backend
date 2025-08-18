package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/thomasmendez/personal-website-backend/api/database"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

type FileData struct {
	Filename    string
	Content     []byte
	ContentType string
}

func (s *Service) getProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	projects, err := database.GetProjects(s.DB, s.TableName)

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

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: "There was an error deserializing projects",
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(projectsJson),
	}, err
}

func (s *Service) postProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	formData, imageFile, err := parseFormData(request)
	if err != nil {
		fmt.Println("Error parsing form data: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	newProject, err := createProjectFromFormData(formData)
	if err != nil {
		fmt.Println("Error creating project from form data: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	if imageFile.Filename != "" {
		mediaLink, err := sendFileToS3(ctx, imageFile)
		if err != nil {
			fmt.Println("failed to upload to S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}

		newProject.MediaLink = &mediaLink
	}

	project, err := database.PostProject(s.DB, s.TableName, newProject)

	if err != nil {
		fmt.Println("There was an error in inserting project with sortValue: %w", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting project: %s", newProject.SortValue),
		}
		res, _ := json.Marshal(errRes)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       string(res),
		}, err
	}

	projectJson, err := json.Marshal(project)

	if err != nil {
		fmt.Println("There was an error deserializing project with sortValue: %w", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in inserting project: %s", newProject.SortValue),
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
	formData, imageFile, err := parseFormData(request)
	if err != nil {
		fmt.Println("Error parsing form data: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	updateProject, err := createProjectFromFormData(formData)
	if err != nil {
		fmt.Println("Error creating project from form data: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	if imageFile.Filename != "" {
		mediaLink, err := sendFileToS3(ctx, imageFile)
		if err != nil {
			fmt.Println("failed to upload to S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}

		updateProject.MediaLink = &mediaLink
	}

	project, err := database.UpdateProject(s.DB, s.TableName, updateProject)

	if err != nil {
		fmt.Println("There was an error in updating project with sortValue: %w", err)
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

	if err != nil {
		fmt.Println("There was an error in updating project with sortValue: %w", err)
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error in updating project with sortValue of: %s", updateProject.SortValue),
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
	var deleteProject models.Project
	err := json.Unmarshal([]byte(request.Body), &deleteProject)
	if err != nil {
		fmt.Println("There was an error in deleting project with sortValue: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	var existingProject models.Project
	err = database.GetItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue, &existingProject)

	if !reflect.DeepEqual(deleteProject, existingProject) {
		fmt.Println("There was an error in deleting project with sortValue: %w", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       resError(http.StatusNotFound),
		}, err
	}

	if existingProject.MediaLink != nil && *existingProject.MediaLink != "" {
		err = deleteFileFromS3(ctx, *existingProject.MediaLink)
		if err != nil {
			fmt.Println("There was an error in deleting file from S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}
	}

	err = database.DeleteItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue)

	if err != nil {
		fmt.Println("There was an error in deleting project with sortValue: %w", err)
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

// delete file from s3
func deleteFileFromS3(ctx context.Context, mediaLink string) error {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	fileName := strings.Split(mediaLink, "/")[len(strings.Split(mediaLink, "/"))-1]

	_, err = s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// send file to s3
func sendFileToS3(ctx context.Context, file FileData) (string, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(file.Filename),
		Body:        strings.NewReader(string(file.Content)),
		ContentType: aws.String(file.ContentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("BUCKET_NAME"), file.Filename), nil
}

// create project from form data
func createProjectFromFormData(formData map[string][]byte) (models.Project, error) {
	var tasks []string
	var teamRoles *[]string
	var cloudServices *[]string
	var tools []string

	if formData["tasks"] != nil {
		parsedTasks, err := parseJSONArray(formData["tasks"])
		if err != nil {
			return models.Project{}, fmt.Errorf("error parsing tasks: %w", err)
		} else {
			tasks = parsedTasks
		}
	} else {
		return models.Project{}, fmt.Errorf("missing tasks")
	}

	if formData["teamRoles"] != nil {
		parsedRoles, err := parseJSONArray(formData["teamRoles"])
		if err != nil {
			return models.Project{}, fmt.Errorf("error parsing teamRoles: %w", err)
		} else {
			teamRoles = &parsedRoles
		}
	}

	if formData["cloudServices"] != nil {
		parsedServices, err := parseJSONArray(formData["cloudServices"])
		if err != nil {
			return models.Project{}, fmt.Errorf("error parsing cloudServices: %w", err)
		} else {
			cloudServices = &parsedServices
		}
	}

	if formData["tools"] != nil {
		parsedTools, err := parseJSONArray(formData["tools"])
		if err != nil {
			return models.Project{}, fmt.Errorf("error parsing tools: %w", err)
		} else {
			tools = parsedTools
		}
	} else {
		return models.Project{}, fmt.Errorf("missing tools")
	}

	var teamSize *string
	if formData["teamSize"] != nil {
		sizePtr := string(formData["teamSize"])
		teamSize = &sizePtr
	}

	var notes *string
	if formData["notes"] != nil {
		notesPtr := string(formData["notes"])
		notes = &notesPtr
	}

	var link *string
	if formData["link"] != nil {
		linkPtr := string(formData["link"])
		link = &linkPtr
	}

	var linkType *string
	if formData["linkType"] != nil {
		linkTypePtr := string(formData["linkType"])
		linkType = &linkTypePtr
	}

	var mediaLink *string
	if formData["mediaLink"] != nil {
		mediaLinkPtr := string(formData["mediaLink"])
		mediaLink = &mediaLinkPtr
	}

	if formData["personalWebsiteType"] == nil || formData["sortValue"] == nil || formData["category"] == nil || formData["name"] == nil || formData["description"] == nil || formData["featuresDescription"] == nil || formData["role"] == nil || formData["duration"] == nil || formData["startDate"] == nil || formData["endDate"] == nil {
		return models.Project{}, fmt.Errorf("missing required fields")
	}

	newProject := models.Project{
		PersonalWebsiteType: string(formData["personalWebsiteType"]),
		SortValue:           string(formData["sortValue"]),
		Category:            string(formData["category"]),
		Name:                string(formData["name"]),
		Description:         string(formData["description"]),
		FeaturesDescription: string(formData["featuresDescription"]),
		Role:                string(formData["role"]),
		Tasks:               tasks,
		TeamSize:            teamSize,
		TeamRoles:           teamRoles,
		CloudServices:       cloudServices,
		Tools:               tools,
		Duration:            string(formData["duration"]),
		StartDate:           string(formData["startDate"]),
		EndDate:             string(formData["endDate"]),
		Notes:               notes,
		Link:                link,
		LinkType:            linkType,
		MediaLink:           mediaLink,
	}
	return newProject, nil
}

// parse JSON array from multipart form data
func parseJSONArray(data []byte) ([]string, error) {
	if data == nil {
		return []string{}, nil
	}

	var result []string
	if err := json.Unmarshal(data, &result); err != nil {
		// If JSON parsing fails, fall back to comma-separated parsing
		// This handles cases where data might be sent as comma-separated values
		str := strings.TrimSpace(string(data))
		if str == "" {
			return []string{}, nil
		}
		return strings.Split(str, ","), nil
	}
	return result, nil
}

// parse form data from request body
func parseFormData(request events.APIGatewayProxyRequest) (map[string][]byte, FileData, error) {
	if !request.IsBase64Encoded && !strings.Contains(request.Headers["Content-Type"], "'multipart/form-data") {
		return nil, FileData{}, fmt.Errorf("request is not base64 encoded and does not contain 'multipart/form-data' in 'Content-Type'")
	}

	var bodyBytes []byte
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return nil, FileData{}, fmt.Errorf("failed to decode base64 body: %w", err)
	}

	_, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		return nil, FileData{}, fmt.Errorf("failed to parse content type: %w", err)
	}

	boundary := params["boundary"]
	if boundary == "" {
		return nil, FileData{}, fmt.Errorf("no boundary found in content type")
	}

	reader := multipart.NewReader(strings.NewReader(string(bodyBytes)), boundary)

	formData := make(map[string][]byte)
	files := make(map[string]FileData)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, FileData{}, fmt.Errorf("failed to get next part: %w", err)
		}

		content, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			return nil, FileData{}, fmt.Errorf("failed to read part content: %w", err)
		}

		fieldName := part.FormName()
		filename := part.FileName()

		if filename != "" {
			if content == nil {
				return nil, FileData{}, fmt.Errorf("file content is nil")
			}
			contentType, err := detectContentType(content)
			if err != nil {
				return nil, FileData{}, fmt.Errorf("failed to detect content type: %w", err)
			}
			files[fieldName] = FileData{
				Filename:    filename,
				Content:     content,
				ContentType: contentType,
			}
		} else {
			formData[fieldName] = content
		}

		part.Close()
	}

	imageFile, exist := files["mediaLink"]
	if exist {
		if imageFile.Content == nil {
			return nil, FileData{}, fmt.Errorf("'mediaLink' file content is nil")
		}
	}

	return formData, imageFile, nil
}

func detectContentType(data []byte) (string, error) {
	if len(data) < 8 {
		return "application/octet-stream", nil
	}

	// Check magic bytes for common image formats
	switch {
	case data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF:
		return "image/jpeg", nil
	case data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47:
		return "image/png", nil
	case data[0] == 0x47 && data[1] == 0x49 && data[2] == 0x46:
		return "image/gif", nil
	case data[0] == 0x52 && data[1] == 0x49 && data[2] == 0x46 && data[3] == 0x46 &&
		data[8] == 0x57 && data[9] == 0x45 && data[10] == 0x42 && data[11] == 0x50:
		return "image/webp", nil
	default:
		return "", fmt.Errorf("unknown content type")
	}
}

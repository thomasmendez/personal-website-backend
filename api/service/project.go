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

func (s *Service) postProjectsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !request.IsBase64Encoded && !strings.Contains(request.Headers["Content-Type"], "'multipart/form-data") {
		log.Printf("request.IsBase64Encoded: %v", request.IsBase64Encoded)
		log.Printf("request.Headers[\"Content-Type\"]: %v", request.Headers["Content-Type"])
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	var bodyBytes []byte
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		log.Printf("failed to decode base64 body: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	_, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		log.Printf("failed to parse content type: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	boundary := params["boundary"]
	if boundary == "" {
		log.Printf("no boundary found in content type")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	// Create a multipart reader
	reader := multipart.NewReader(strings.NewReader(string(bodyBytes)), boundary)

	formData := make(map[string][]byte)
	files := make(map[string]FileData)

	// Parse each part
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("failed to get next part: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}

		// Read the content of this part
		content, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			log.Printf("failed to read part content: %v", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}

		fieldName := part.FormName()
		filename := part.FileName()

		log.Printf("fieldName: %v", fieldName)
		log.Printf("filename: %v", filename)
		log.Printf("contentType: %v", part.Header.Get("Content-Type"))

		if filename != "" {
			// This is a file
			files[fieldName] = FileData{
				Filename:    filename,
				Content:     content,
				ContentType: part.Header.Get("Content-Type"),
			}
		} else {
			// This is a regular form field
			log.Printf("content: %v", string(content))
			formData[fieldName] = content
		}

		part.Close()
	}

	log.Printf("formData: %v", formData)
	log.Printf("files: %v", files)

	imageFile, exist := files["mediaLink"]
	if !exist {
		log.Printf("'mediaLink' file not found")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	if imageFile.Content == nil {
		// should still make the update for the metadata
		log.Printf("'mediaLink' file content is nil")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	var tasks []string
	var teamRoles []string
	var cloudServices []string
	var tools []string

	if formData["tasks"] != nil {
		parsedTasks, err := parseJSONArray(formData["tasks"])
		if err != nil {
			log.Printf("Error parsing tasks: %v", err)
		} else {
			tasks = parsedTasks
		}
	}

	if formData["teamRoles"] != nil {
		parsedRoles, err := parseJSONArray(formData["teamRoles"])
		if err != nil {
			log.Printf("Error parsing teamRoles: %v", err)
		} else {
			teamRoles = parsedRoles
		}
	}

	if formData["cloudServices"] != nil {
		parsedServices, err := parseJSONArray(formData["cloudServices"])
		if err != nil {
			log.Printf("Error parsing cloudServices: %v", err)
		} else {
			cloudServices = parsedServices
		}
	}

	if formData["tools"] != nil {
		parsedTools, err := parseJSONArray(formData["tools"])
		if err != nil {
			log.Printf("Error parsing tools: %v", err)
		} else {
			tools = parsedTools
		}
	}

	var teamSize string
	if formData["teamSize"] != nil {
		teamSize = string(formData["teamSize"])
	}

	var notes string
	if formData["notes"] != nil {
		notes = string(formData["notes"])
	}

	var link string
	if formData["link"] != nil {
		link = string(formData["link"])
	}

	var linkType string
	if formData["linkType"] != nil {
		linkType = string(formData["linkType"])
	}

	var mediaLink string
	if formData["mediaLink"] != nil {
		mediaLink = string(formData["mediaLink"])
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
		TeamSize:            &teamSize,
		TeamRoles:           &teamRoles,
		CloudServices:       &cloudServices,
		Tools:               tools,
		Duration:            string(formData["duration"]),
		StartDate:           string(formData["startDate"]),
		EndDate:             string(formData["endDate"]),
		Notes:               &notes,
		Link:                &link,
		LinkType:            &linkType,
		MediaLink:           &mediaLink,
	}

	// send request to s3 bucket
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       resError(http.StatusInternalServerError),
		}, err
	}

	s3Client := s3.NewFromConfig(cfg)

	fmt.Println("Uploading file to S3 bucket...")
	fmt.Println(os.Getenv("BUCKET_NAME"))
	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET_NAME")),
		Key:         aws.String(imageFile.Filename),
		Body:        strings.NewReader(string(imageFile.Content)),
		ContentType: aws.String(imageFile.ContentType),
	})
	if err != nil {
		log.Printf("failed to upload to S3: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       resError(http.StatusInternalServerError),
		}, err
	}
	fmt.Println("Uploaded file to S3 bucket!")
	// Generate S3 URL
	mediaLink = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", os.Getenv("BUCKET_NAME"), imageFile.Filename)
	newProject.MediaLink = &mediaLink

	project, err := database.PostProject(s.DB, s.TableName, newProject)

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

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error deserializing project with sortValue: %s", newProject.SortValue),
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
	var updateProject models.Project
	err := json.Unmarshal([]byte(request.Body), &updateProject)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	project, err := database.UpdateProject(s.DB, s.TableName, updateProject)

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

	if err != nil {
		log.Print(err.Error())
		errRes := ErrorResponse{
			Message: fmt.Sprintf("There was an error deserializing project with sortValue: %s", project.SortValue),
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
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, err
	}

	var existingProject models.Project
	err = database.GetItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue, &existingProject)

	if !reflect.DeepEqual(deleteProject, existingProject) {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       resError(http.StatusNotFound),
		}, err
	}

	err = database.DeleteItem(s.DB, s.TableName, deleteProject.PersonalWebsiteType, deleteProject.SortValue)

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

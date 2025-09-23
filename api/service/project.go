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
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasmendez/personal-website-backend/api/bucket"
	"github.com/thomasmendez/personal-website-backend/api/database"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

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

	// Generate presigned URL for mediaLink
	for i, project := range projects {
		if project.MediaLink != nil {
			if strings.Contains(*project.MediaLink, ".s3.amazonaws.com") {
				presignedReq, err := bucket.GeneratePresignedURL(ctx, s.S3.Client, s.BucketName, *project.MediaLink)
				if err != nil {
					log.Printf("error in generating presigned URL: %v", err)
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       resError(http.StatusInternalServerError),
					}, err
				}
				projects[i].MediaLink = &presignedReq.URL
			} else {
				log.Printf("mediaLink for project %s is not a valid S3 link, return as is", project.SortValue)
			}
		}
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
	var newProject *models.Project
	var imageFile models.FileData
	var err error

	fmt.Printf("request: %v", request)
	if request.IsBase64Encoded || strings.Contains(getContentType(request.Headers), "'multipart/form-data") {
		newProject, imageFile, err = parseFormData[models.Project](request)
		if err != nil {
			fmt.Println("Error parsing form data: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
	} else if !request.IsBase64Encoded && strings.Contains(getContentType(request.Headers), "application/json") {
		err = json.Unmarshal([]byte(request.Body), &newProject)
		if err != nil {
			log.Printf("error in serializing json: %v", err)
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
		fmt.Println("Invalid request format")
		log.Printf("request: %v", request)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	log.Printf("newProject: %v after parse", newProject)

	if imageFile.Filename != "" && imageFile.Content != nil && imageFile.ContentType != "" {
		mediaLink, err := bucket.SendFileToS3(ctx, s.S3.Client, s.BucketName, imageFile)
		if err != nil {
			fmt.Println("failed to upload to S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}
		newProject.MediaLink = &mediaLink
	}

	project, err := database.PostProject(s.DB, s.TableName, *newProject)

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
	var updateProject *models.Project
	var imageFile models.FileData
	var err error
	if request.IsBase64Encoded || strings.Contains(getContentType(request.Headers), "'multipart/form-data") {
		updateProject, imageFile, err = parseFormData[models.Project](request)
		if err != nil {
			fmt.Println("Error parsing form data: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusBadRequest,
				Body:       resError(http.StatusBadRequest),
			}, err
		}
	} else if !request.IsBase64Encoded && strings.Contains(getContentType(request.Headers), "application/json") {
		err = json.Unmarshal([]byte(request.Body), &updateProject)
		if err != nil {
			log.Printf("error in serializing json: %v", err)
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
		fmt.Println("Invalid request format")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       resError(http.StatusBadRequest),
		}, nil
	}

	if imageFile.Filename != "" && imageFile.Content != nil && imageFile.ContentType != "" {
		mediaLink, err := bucket.SendFileToS3(ctx, s.S3.Client, s.BucketName, imageFile)
		if err != nil {
			fmt.Println("failed to upload to S3: %w", err)
			return events.APIGatewayProxyResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       resError(http.StatusInternalServerError),
			}, err
		}
		updateProject.MediaLink = &mediaLink
	}

	project, err := database.UpdateProject(s.DB, s.TableName, *updateProject)

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
		err = bucket.DeleteFileFromS3(ctx, s.S3.Client, s.BucketName, *existingProject.MediaLink)
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

// parse form data from request body
func parseFormData[T any](request events.APIGatewayProxyRequest) (*T, models.FileData, error) {
	var bodyBytes []byte
	bodyBytes, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		return nil, models.FileData{}, fmt.Errorf("failed to decode base64 body: %w", err)
	}

	_, params, err := mime.ParseMediaType(getContentType(request.Headers))
	if err != nil {
		return nil, models.FileData{}, fmt.Errorf("failed to parse content type: %w", err)
	}

	boundary := params["boundary"]
	if boundary == "" {
		return nil, models.FileData{}, fmt.Errorf("no boundary found in content type")
	}

	reader := multipart.NewReader(strings.NewReader(string(bodyBytes)), boundary)

	var result T
	file := models.FileData{}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, models.FileData{}, fmt.Errorf("failed to get next part: %w", err)
		}

		content, err := io.ReadAll(part)
		if err != nil {
			part.Close()
			return nil, models.FileData{}, fmt.Errorf("failed to read part content: %w", err)
		}

		fieldName := part.FormName()
		filename := part.FileName()

		if filename != "" {
			if content == nil {
				return nil, models.FileData{}, fmt.Errorf("file content is nil")
			}
			contentType, err := detectContentType(content)
			if err != nil {
				return nil, models.FileData{}, fmt.Errorf("failed to detect content type: %w", err)
			}
			file = models.FileData{
				Filename:    filename,
				Content:     content,
				ContentType: contentType,
			}
			fmt.Printf("file: %s, filename: %s, contentType: %s\n", fieldName, filename, contentType)
		} else {
			log.Printf("fieldName: %s", fieldName)
			log.Printf("content: %v", content)

			// if filename is empty set the map of the result to nil
			v := reflect.ValueOf(&result).Elem()
			structName := strings.ToUpper(fieldName[:1]) + fieldName[1:]
			field := v.FieldByName(structName)
			if !field.CanSet() {
				log.Printf("invalid or non-settable field: %s", fieldName)
				continue // Skip invalid or non-settable fields
			}
			if err := setFieldValue(field, content); err != nil {
				return nil, models.FileData{}, fmt.Errorf("failed to set field value: %w", err)
			}
		}

		part.Close()
	}
	return &result, file, nil
}

func setFieldValue(field reflect.Value, value []byte) error {
	log.Printf("field kind: %v", field.Kind())
	log.Printf("field type: %v", field.Type())
	log.Printf("value: %v", value)

	switch field.Kind() {
	case reflect.String:
		log.Printf("setting string field: %s", string(value))
		field.SetString(string(value))
	case reflect.Slice:
		log.Printf("setting slice field: %v", value)
		return setSliceFromBytes(field, value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if intVal, err := strconv.ParseInt(string(value), 10, 64); err == nil {
			field.SetInt(intVal)
		} else {
			return err
		}
	case reflect.Float32, reflect.Float64:
		if floatVal, err := strconv.ParseFloat(string(value), 64); err == nil {
			field.SetFloat(floatVal)
		} else {
			return err
		}
	case reflect.Bool:
		if boolVal, err := strconv.ParseBool(string(value)); err == nil {
			field.SetBool(boolVal)
		} else {
			return err
		}
	case reflect.Ptr:
		if field.Type() == reflect.TypeOf((*string)(nil)) {
			if value == nil {
				field.Set(reflect.New(field.Type().Elem()))
				return nil
			}
			valuePtr := string(value)
			field.Set(reflect.ValueOf(&valuePtr))
		}
		if field.Type() == reflect.TypeOf((*[]string)(nil)) {
			if value == nil {
				field.Set(reflect.New(field.Type().Elem()))
				return nil
			}
			log.Printf("setting *[]string slice field: %v", value)
			return setPointerSliceFromBytes(field, value)
		}
	default:
		return fmt.Errorf("unsupported type: %v", field.Kind())
	}
	return nil
}

// Handle setting pointer slices from byte data
func setPointerSliceFromBytes(field reflect.Value, value []byte) error {
	// Check if field is a pointer
	if field.Type().Kind() != reflect.Ptr {
		return fmt.Errorf("field is not a pointer: %v", field.Type().Kind())
	}

	// Get the element type the pointer points to (should be a slice)
	sliceType := field.Type().Elem()
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("pointer does not point to a slice: %v", sliceType.Kind())
	}

	// Get the slice element type
	elemType := sliceType.Elem()

	// Handle *[]byte directly
	if elemType.Kind() == reflect.Uint8 {
		slice := reflect.New(sliceType).Elem()
		slice.SetBytes(value)
		field.Set(slice.Addr())
		return nil
	}

	str := string(value)
	log.Printf("str: %v", str)

	// Handle empty/null cases
	if str == "" || str == "null" {
		field.Set(reflect.Zero(field.Type())) // Set to nil
		return nil
	}

	var parts []string

	// Try to parse as JSON array first
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		var jsonArray []string
		if err := json.Unmarshal(value, &jsonArray); err == nil {
			parts = jsonArray
		} else {
			// If JSON parsing fails, fall back to comma-separated parsing
			parts = strings.Split(str, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
		}
	} else {
		// Parse as comma-separated values
		parts = strings.Split(str, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
	}

	// Create a new addressable slice by allocating memory for it
	slicePtr := reflect.New(sliceType)
	slice := slicePtr.Elem()
	slice.Set(reflect.MakeSlice(sliceType, len(parts), len(parts)))

	log.Printf("parts: %v", parts)
	log.Printf("slice: %v", slice)

	for i, part := range parts {
		part = strings.TrimSpace(part)
		elem := slice.Index(i)

		switch elemType.Kind() {
		case reflect.String:
			log.Printf("setting string slice field: %s", part)
			elem.SetString(part)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(part, 10, 64); err == nil {
				elem.SetInt(intVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to int: %v", part, err)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(part, 10, 64); err == nil {
				elem.SetUint(uintVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to uint: %v", part, err)
			}
		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(part, 64); err == nil {
				elem.SetFloat(floatVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to float: %v", part, err)
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(part); err == nil {
				elem.SetBool(boolVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to bool: %v", part, err)
			}
		default:
			return fmt.Errorf("unsupported slice element type: %v", elemType.Kind())
		}
	}

	// Set the field to point to the slice
	field.Set(slicePtr)
	return nil
}

// Handle setting slices from byte data
func setSliceFromBytes(field reflect.Value, value []byte) error {
	elemType := field.Type().Elem()

	// Handle []byte directly
	if elemType.Kind() == reflect.Uint8 {
		field.SetBytes(value)
		return nil
	}

	// Parse as comma-separated values for other slice types
	str := string(value)
	log.Printf("str: %v", str)
	if str == "" {
		field.Set(reflect.MakeSlice(field.Type(), 0, 0))
		return nil
	}

	var parts []string

	// Try to parse as JSON array first
	str = strings.TrimSpace(str)
	if strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
		var jsonArray []string
		if err := json.Unmarshal(value, &jsonArray); err == nil {
			parts = jsonArray
		} else {
			// If JSON parsing fails, fall back to comma-separated parsing
			parts = strings.Split(str, ",")
			for i := range parts {
				parts[i] = strings.TrimSpace(parts[i])
			}
		}
	} else {
		// Parse as comma-separated values
		parts = strings.Split(str, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
	}

	slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))

	log.Printf("parts: %v", parts)
	log.Printf("slice: %v", slice)

	for i, part := range parts {
		part = strings.TrimSpace(part)
		elem := slice.Index(i)

		switch elemType.Kind() {
		case reflect.String:
			log.Printf("setting string slice field: %s", part)
			elem.SetString(part)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if intVal, err := strconv.ParseInt(part, 10, 64); err == nil {
				elem.SetInt(intVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to int: %v", part, err)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if uintVal, err := strconv.ParseUint(part, 10, 64); err == nil {
				elem.SetUint(uintVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to uint: %v", part, err)
			}
		case reflect.Float32, reflect.Float64:
			if floatVal, err := strconv.ParseFloat(part, 64); err == nil {
				elem.SetFloat(floatVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to float: %v", part, err)
			}
		case reflect.Bool:
			if boolVal, err := strconv.ParseBool(part); err == nil {
				elem.SetBool(boolVal)
			} else {
				return fmt.Errorf("cannot convert '%s' to bool: %v", part, err)
			}
		default:
			return fmt.Errorf("unsupported slice element type: %v", elemType.Kind())
		}
	}

	field.Set(slice)
	return nil
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

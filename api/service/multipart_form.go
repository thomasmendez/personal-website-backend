package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

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
	log.Printf("Parsing as comma-separated values for other slice types")
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

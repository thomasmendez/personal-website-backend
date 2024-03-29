package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/thomasmendez/personal-website-backend/api/services"
)

// var dynamoDbClient *dynamodb.DynamoDB

const tableName = "PersonalWebsiteTable"

// var database *database.Database

type ErrorResponse struct {
	Message string `json:"message"`
}
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// for prd
	// session := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	srv := services.NewService()

	lambda.Start(srv.HandleRoute)
}

// func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	// switch request.HTTPMethod {
// 	// case http.MethodGet:
// 	// 	return handleGetRequest(ctx, request)
// 	// case http.MethodPost:
// 	// 	return handlePostRequest(ctx, request)
// 	// case http.MethodPut:
// 	// 	return handlePutRequest(ctx, request)
// 	// case http.MethodDelete:
// 	// 	return handleDeleteRequest(ctx, request)
// 	// default:
// 	// 	return events.APIGatewayProxyResponse{
// 	// 		StatusCode: http.StatusMethodNotAllowed,
// 	// 		Body:       http.StatusText(http.StatusMethodNotAllowed),
// 	// 	}, nil
// 	// }
// }

// func handleGetRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

// 	if !strings.HasPrefix(request.Path, "/api/work") {
// 		log.Printf("Invalid path: %s", request.Path)
// 		return events.APIGatewayProxyResponse{StatusCode: 404}, fmt.Errorf("not found")
// 	}

// 	pathParts := strings.Split(request.Path, "/")
// 	log.Printf("path parts: %v", pathParts)
// 	// if len(pathParts) < 4 {
// 	// 	log.Printf("Invalid path: %s", request.Path)
// 	// 	return events.APIGatewayProxyResponse{StatusCode: 404}, fmt.Errorf("not found")
// 	// }
// 	jobIDStr := pathParts[len(pathParts)-1]

// 	if jobIDStr == "work" {

// 		input := &dynamodb.ScanInput{
// 			TableName: aws.String("PersonalWebsiteTable"),
// 		}

// 		result, err := dynamoDbClient.Scan(input)
// 		if err != nil {
// 			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error scanning table: %s", err.Error())}
// 			responseBody, _ := json.Marshal(errorResponse)
// 			return events.APIGatewayProxyResponse{
// 				StatusCode: 500,
// 				Body:       string(responseBody),
// 			}, nil
// 		}

// 		// Iterate through the result and construct a list of JobDescription objects
// 		jobDescriptions := make([]Item, 0)
// 		for _, item := range result.Items {
// 			job := Item{
// 				ID:   aws.StringValue(item["id"].S),
// 				Name: aws.StringValue(item["name"].S),
// 			}
// 			jobDescriptions = append(jobDescriptions, job)
// 		}

// 		// Marshal the list of job descriptions into JSON format for the response body
// 		responseBody, err := json.Marshal(jobDescriptions)
// 		if err != nil {
// 			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error marshalling response: %s", err.Error())}
// 			responseBody, _ := json.Marshal(errorResponse)
// 			return events.APIGatewayProxyResponse{
// 				StatusCode: 500,
// 				Body:       string(responseBody),
// 			}, nil
// 		}

// 		// Return the list of job descriptions in the response
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 200,
// 			Body:       string(responseBody),
// 		}, nil
// 	}

// 	// Validate that the ID is an integer
// 	jobID, err := strconv.Atoi(jobIDStr)
// 	if err != nil {
// 		log.Printf("Error parsing job ID: %v", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("invalid job ID")
// 	}

// 	log.Printf("jobId: %v", jobID)

// 	// Build DynamoDB query input
// 	input := &dynamodb.GetItemInput{
// 		TableName: aws.String(tableName),
// 		Key: map[string]*dynamodb.AttributeValue{
// 			// "id": &types.AttributeValueMemberS{Value: *aws.String(strconv.Itoa(jobID))},
// 			"id": {
// 				S: aws.String("1"),
// 			},
// 			"name": {
// 				S: aws.String("John"),
// 			},
// 		},
// 	}

// 	result, err := dynamoDbClient.GetItem(input)
// 	if err != nil {
// 		log.Printf("Error calling DynamoDB GetItem API: %v", err)

// 		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
// 			errorResponse := ErrorResponse{Message: "Item not found"}
// 			responseBody, _ := json.Marshal(errorResponse)
// 			return events.APIGatewayProxyResponse{
// 				StatusCode: 404,
// 				Body:       string(responseBody),
// 			}, nil
// 		}

// 		errorResponse := ErrorResponse{Message: fmt.Sprintf("Error retrieving item: %s", err.Error())}
// 		responseBody, _ := json.Marshal(errorResponse)
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: 500,
// 			Body:       string(responseBody),
// 		}, nil
// 	}

// 	// Check if the item was found
// 	if result.Item == nil {
// 		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Item not found"}, nil
// 	}

// 	// item := Item{}

// 	// err = dynamodbattribute.UnmarshalMap(result.Item, &item)
// 	// if err != nil {
// 	// 	panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
// 	// }

// 	// Convert DynamoDB result to JSON
// 	jsonData, err := json.Marshal(result.Item)
// 	if err != nil {
// 		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error converting DynamoDB result to JSON"}, err
// 	}

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: 200,
// 		Body:       string(jsonData),
// 	}, nil
// }

// func handlePostRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	var newItem Item
// 	err := json.Unmarshal([]byte(request.Body), &newItem)
// 	if err != nil {
// 		log.Printf("err: %v", err)
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusBadRequest,
// 			Body:       "Bad Request: Invalid JSON",
// 		}, nil
// 	}

// 	item := map[string]*dynamodb.AttributeValue{
// 		"id":   {S: aws.String(newItem.ID)},
// 		"name": {S: aws.String(newItem.Name)},
// 	}

// 	input := &dynamodb.PutItemInput{
// 		Item:      item,
// 		TableName: aws.String(tableName),
// 	}

// 	// Insert item into DynamoDB table
// 	_, err = dynamoDbClient.PutItem(input)
// 	if err != nil {
// 		fmt.Println("Error inserting item:", err)
// 		return events.APIGatewayProxyResponse{
// 			StatusCode: http.StatusInternalServerError,
// 			Body:       "Error inserting data into DynamoDB",
// 		}, err
// 	}

// 	return events.APIGatewayProxyResponse{
// 		StatusCode: http.StatusCreated,
// 		Body:       "Item Created",
// 	}, nil
// }

// func handlePutRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	var updateItem Item
// 	if err := json.Unmarshal([]byte(request.Body), &updateItem); err != nil {
// 		log.Printf("Error parsing request body: %v", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 400}, err
// 	}

// 	updateExpression := "SET #name = :name"
// 	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
// 		":name": {S: aws.String(updateItem.Name)},
// 	}
// 	expressionAttributeNames := map[string]*string{
// 		"#name": aws.String("name"),
// 	}

// 	// Construct UpdateItemInput
// 	input := &dynamodb.UpdateItemInput{
// 		TableName:                 aws.String(tableName),
// 		Key:                       map[string]*dynamodb.AttributeValue{"id": {S: aws.String(updateItem.ID)}},
// 		UpdateExpression:          aws.String(updateExpression),
// 		ExpressionAttributeValues: expressionAttributeValues,
// 		ExpressionAttributeNames:  expressionAttributeNames,
// 	}

// 	// Update item in DynamoDB table
// 	updateItemOutput, err := dynamoDbClient.UpdateItem(input)
// 	if err != nil {
// 		log.Printf("Error updating item: %v", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 500}, err
// 	}
// 	log.Printf("updateItemOutput: %v", updateItemOutput)
// 	// responseBody := fmt.Sprintf("Item with ID %s updated successfully", idToUpdate)
// 	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "PUT Request"}, nil
// }

// func handleDeleteRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

// 	var deleteItem Item
// 	if err := json.Unmarshal([]byte(request.Body), &deleteItem); err != nil {
// 		log.Printf("Error parsing request body: %v", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 400}, err
// 	}

// 	// Prepare the key of the item to delete
// 	key := map[string]*dynamodb.AttributeValue{
// 		"id": {S: aws.String(deleteItem.ID)}, // Assuming 'id' is the primary key
// 	}

// 	// Construct DeleteItemInput
// 	input := &dynamodb.DeleteItemInput{
// 		Key:       key,
// 		TableName: aws.String(tableName),
// 	}

// 	// Delete item from DynamoDB table
// 	deleteItemOutput, err := dynamoDbClient.DeleteItem(input)
// 	if err != nil {
// 		log.Printf("Error deleting item: %v", err)
// 		return events.APIGatewayProxyResponse{StatusCode: 500}, err
// 	}

// 	log.Printf("delete item output: %v", deleteItemOutput)

// 	// responseBody := fmt.Sprintf("Item with ID %s deleted successfully", idToDelete)
// 	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "DELETE Request"}, nil
// }

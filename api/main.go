package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var dynamoDbClient *dynamodb.DynamoDB

const tableName = "PersonalWebsiteTable"

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
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Endpoint: aws.String("http://dynamodb:8000"),
			Region:   aws.String("us-west-2"),
		},
	}))
	dynamoDbClient = dynamodb.New(session)
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// input := &dynamodb.ListTablesInput{}

	// result, err := dynamoDbClient.ListTables(input)
	// if err != nil {
	// 	if aerr, ok := err.(awserr.Error); ok {
	// 		switch aerr.Code() {
	// 		case dynamodb.ErrCodeInternalServerError:
	// 			fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
	// 		default:
	// 			fmt.Println(aerr.Error())
	// 		}
	// 	} else {
	// 		// Print the error, cast err to awserr.Error to get the Code and
	// 		// Message from an error.
	// 		fmt.Println(err.Error())
	// 	}
	// 	return events.APIGatewayProxyResponse{StatusCode: 500}, fmt.Errorf("internal server error")
	// }

	// for _, n := range result.TableNames {
	// 	fmt.Println(*n)
	// }

	// assign the last read tablename as the start for our next call to the ListTables function
	// the maximum number of table names returned in a call is 100 (default), which requires us to make
	// multiple calls to the ListTables function to retrieve all table names
	// input.ExclusiveStartTableName = result.LastEvaluatedTableName

	// if result.LastEvaluatedTableName == nil {
	// 	break
	// }

	switch request.HTTPMethod {
	case http.MethodGet:
		return handleGetRequest(ctx, request)
	case http.MethodPost:
		return handlePostRequest(ctx, request)
	case http.MethodPut:
		return handlePutRequest(ctx, request)
	case http.MethodDelete:
		return handleDeleteRequest(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       http.StatusText(http.StatusMethodNotAllowed),
		}, nil
	}
}

func handleGetRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	if !strings.HasPrefix(request.Path, "/api/work") {
		log.Printf("Invalid path: %s", request.Path)
		return events.APIGatewayProxyResponse{StatusCode: 404}, fmt.Errorf("not found")
	}

	pathParts := strings.Split(request.Path, "/")
	log.Printf("path parts: %v", pathParts)
	// if len(pathParts) < 4 {
	// 	log.Printf("Invalid path: %s", request.Path)
	// 	return events.APIGatewayProxyResponse{StatusCode: 404}, fmt.Errorf("not found")
	// }
	jobIDStr := pathParts[len(pathParts)-1]

	if jobIDStr == "work" {

		input := &dynamodb.ScanInput{
			TableName: aws.String("PersonalWebsiteTable"),
		}

		result, err := dynamoDbClient.Scan(input)
		if err != nil {
			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error scanning table: %s", err.Error())}
			responseBody, _ := json.Marshal(errorResponse)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       string(responseBody),
			}, nil
		}

		// Iterate through the result and construct a list of JobDescription objects
		jobDescriptions := make([]Item, 0)
		for _, item := range result.Items {
			job := Item{
				ID:   aws.StringValue(item["id"].S),
				Name: aws.StringValue(item["name"].S),
			}
			jobDescriptions = append(jobDescriptions, job)
		}

		// Marshal the list of job descriptions into JSON format for the response body
		responseBody, err := json.Marshal(jobDescriptions)
		if err != nil {
			errorResponse := ErrorResponse{Message: fmt.Sprintf("Error marshalling response: %s", err.Error())}
			responseBody, _ := json.Marshal(errorResponse)
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       string(responseBody),
			}, nil
		}

		// Return the list of job descriptions in the response
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       string(responseBody),
		}, nil
	}

	// Validate that the ID is an integer
	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		log.Printf("Error parsing job ID: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("invalid job ID")
	}

	log.Printf("jobId: %v", jobID)

	// Build DynamoDB query input
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			// "id": &types.AttributeValueMemberS{Value: *aws.String(strconv.Itoa(jobID))},
			"id": {
				S: aws.String("1"),
			},
			"name": {
				S: aws.String("John"),
			},
		},
	}

	result, err := dynamoDbClient.GetItem(input)
	if err != nil {
		log.Printf("Error calling DynamoDB GetItem API: %v", err)

		if _, ok := err.(*dynamodb.ResourceNotFoundException); ok {
			errorResponse := ErrorResponse{Message: "Item not found"}
			responseBody, _ := json.Marshal(errorResponse)
			return events.APIGatewayProxyResponse{
				StatusCode: 404,
				Body:       string(responseBody),
			}, nil
		}

		errorResponse := ErrorResponse{Message: fmt.Sprintf("Error retrieving item: %s", err.Error())}
		responseBody, _ := json.Marshal(errorResponse)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       string(responseBody),
		}, nil
	}

	// Check if the item was found
	if result.Item == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "Item not found"}, nil
	}

	// item := Item{}

	// err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	// }

	// Convert DynamoDB result to JSON
	jsonData, err := json.Marshal(result.Item)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Error converting DynamoDB result to JSON"}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonData),
	}, nil
}

func handlePostRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var newItem Item
	err := json.Unmarshal([]byte(request.Body), &newItem)
	if err != nil {
		log.Printf("err: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Bad Request: Invalid JSON",
		}, nil
	}

	item := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String(newItem.ID)},
		"name": {S: aws.String(newItem.Name)},
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	// Insert item into DynamoDB table
	_, err = dynamoDbClient.PutItem(input)
	if err != nil {
		fmt.Println("Error inserting item:", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error inserting data into DynamoDB",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Body:       "Item Created",
	}, nil
}

func handlePutRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "PUT Request"}, nil
}

func handleDeleteRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "DELETE Request"}, nil
}

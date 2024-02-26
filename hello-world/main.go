package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodGet:
		return handleGET(request)
	case http.MethodPost:
		return handlePOST(request)
	case http.MethodPut:
		return handlePUT(request)
	case http.MethodDelete:
		return handleDELETE(request)
	default:
		return events.APIGatewayProxyResponse{StatusCode: http.StatusMethodNotAllowed}, nil
	}
}

func handleGET(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Implement your GET logic
	// Use DynamoDB SDK to fetch data from the table

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "GET Request"}, nil
}

func handlePOST(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	var data map[string]interface{}
	err := json.Unmarshal([]byte(reqBody), &data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid JSON format",
		}, nil
	}

	value := data["key"]

	// Do something with the value...

	// Return a response
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("Received value: %v", value),
	}

	return response, nil
}

func handlePUT(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Implement your PUT logic
	// Use DynamoDB SDK to update data in the table

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "PUT Request"}, nil
}

func handleDELETE(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Implement your DELETE logic
	// Use DynamoDB SDK to delete data from the table

	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: "DELETE Request"}, nil
}

func main() {
	lambda.Start(handler)
}

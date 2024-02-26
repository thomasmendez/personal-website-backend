package main

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	tests := []struct {
		name           string
		request        events.APIGatewayProxyRequest
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "GET Request",
			request:        events.APIGatewayProxyRequest{HTTPMethod: http.MethodGet},
			expectedStatus: http.StatusOK,
			expectedBody:   "GET Request",
		},
		{
			name: "POST Request with valid JSON",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       `{"key": "value"}`,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Received value: value",
		},
		{
			name: "POST Request with invalid JSON",
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "invalid-json",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON format",
		},
		{
			name:           "PUT Request",
			request:        events.APIGatewayProxyRequest{HTTPMethod: http.MethodPut},
			expectedStatus: http.StatusOK,
			expectedBody:   "PUT Request",
		},
		{
			name:           "DELETE Request",
			request:        events.APIGatewayProxyRequest{HTTPMethod: http.MethodDelete},
			expectedStatus: http.StatusOK,
			expectedBody:   "DELETE Request",
		},
		{
			name:           "Unsupported Method",
			request:        events.APIGatewayProxyRequest{HTTPMethod: "INVALID"},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := handler(test.request)
			assert.Nil(t, err)
			assert.Equal(t, test.expectedStatus, response.StatusCode)
			assert.Equal(t, test.expectedBody, response.Body)
		})
	}
}

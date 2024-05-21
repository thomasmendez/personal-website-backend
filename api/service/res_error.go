package service

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func resError(errorStatusCode int) (errRes string) {
	var message string
	switch errorStatusCode {
	case http.StatusBadRequest:
		message = "Bad Request: Invalid JSON"
	case http.StatusNotFound:
		message = "Resource not found"
	default:
		message = "Unknown error"
	}
	res, _ := json.Marshal(ErrorResponse{
		Message: message,
	})
	return string(res)
}

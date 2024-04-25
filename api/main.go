package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	service "github.com/thomasmendez/personal-website-backend/api/service"
)

func main() {
	srv := service.NewService()
	lambda.Start(srv.HandleRoute)
}

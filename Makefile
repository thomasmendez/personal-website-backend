.PHONY: build

db:
	docker compose up -d
db-create-table:
	aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
start:
	sam.cmd local start-api --docker-network dynamodb-backend
build:
	sam.cmd build
test:
	cd api && go test ./...
test-coverage:
	cd api && go test -coverprofile=coverage.out ./...
	cd api && go tool cover -html=coverage.out -o coverage.html
test-integration:
	cd api && INTEGRATION=1 go test ./...
build-go:
	cd api && GOARCH=arm64 GOOS=linux go build -o bootstrap main.go
build-lambda-windows:
	/c/Users/owner/go/bin/build-lambda-zip.exe -o ./api/lambda-handler.zip ./api/bootstrap
# personal-website-backend

Backend for updating details for personal website. Backend is a deployed Go Lambda function that updates a DynamoDB table.

## Requirements

* AWS CLI already configured with Administrator permission
* [Docker installed](https://www.docker.com/community-edition)
* SAM CLI - [Install the SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)
* [Golang](https://golang.org) - Currently using v.1.22 (can use [Chocolatey](https://community.chocolatey.org/packages/golang) or [Homebrew](https://formulae.brew.sh/formula/go))
* [Bruno](https://www.usebruno.com/) (Optional) - Open source API client. Similar to Postman, but is offline-only and will never require an account. 
* [Make](https://www.gnu.org/software/make/) (Optional) - To run makefile commands 

## Local development

1. **Start DynamoDb Docker Container**
    ```shell
    docker compose up -d
    ```

2. **Create DynamoDb Tables**
    ```shell
    aws dynamodb create-table --cli-input-json file://json/create-table.json --endpoint-url http://localhost:8000
    ```

3. **Build the go executable for the [lambda linux environment](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)**
    ```shell
    GOARCH=arm64 GOOS=linux go build -o bootstrap main.go
    ```

4. **Zip project (windows)**

    `C:\Users\owner\go\bin\build-lambda-zip.exe -o lambda-handler.zip bootstrap` using the provided `build-lambda-zip` package. If needed, you can install with `go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest`. See [To create a .zip deployment package (Windows)](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)

5. **Invoke function locally through local API Gateway**
    ```shell
    sam.cmd local start-api --docker-network dynamodb-backend --template-file=template.yaml
    ```

    *Note: Use `sam.cmd` when running AWS SAM on windows*

    To rebuild and apply local changes, use `sam.cmd build --use-container` 

6. **Run CRUD Integration Tests**
    ```shell
    cd api && INTEGRATION=1 go test ./...
    ```

    *Note: Integration test will take about a minute to complete*

7. **(Optional) Run CRUD Requests Individually**
    
    Can also run CRUD API requests by importing the `/bruno` collection

## Deployment

Ensure you follow the same steps you did to build the executable and zipping the project

1. **Build the go executable for the [lambda linux environment](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)**
    ```shell
    GOARCH=arm64 GOOS=linux go build -o bootstrap main.go
    ```

2. **Zip project (windows)**

    `C:\Users\owner\go\bin\build-lambda-zip.exe -o lambda-handler.zip bootstrap` using the provided `build-lambda-zip` package. If needed, you can install with `go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest`. See [To create a .zip deployment package (Windows)](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)

3. **Deploy CloudFromation stack template**
    ```shell
    sam.cmd deploy --guided --template-file=template.yaml
    ```

## Helpful Commands

### DynamoDB Commands

**View DynamoDb Tables**
```shell
aws dynamodb list-tables --endpoint-url http://localhost:8000
```

**Add DynamoDb Data**

To write one item
```shell
aws dynamodb put-item --cli-input-json file://json/work/add-item.json --endpoint-url http://localhost:8000
```
To write multiple
```shell
aws dynamodb batch-write-item --cli-input-json file://json/work/add-items.json --endpoint-url http://localhost:8000
```

**View DynamoDb Data**
```shell
aws dynamodb scan --table-name PersonalWebsiteTable --endpoint-url http://localhost:8000
```

**Query DynamoDb Data (Recent to Oldest)**
```shell
aws dynamodb query --cli-input-json file://json/work/query-recent-items.json --endpoint-url http://localhost:8000
```

**Delete DynamoDb Data**
```shell
aws dynamodb delete-item --cli-input-json file://json/work/delete-item.json --endpoint-url http://localhost:8000
```

### Testing Commands

Go to `api` directory to run tests

**Unit Tests**
```shell
go test ./...
```

**Coverage Report**

```shell
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**Integration Tests**
```shell
INTEGRATION=1 go test ./...
```

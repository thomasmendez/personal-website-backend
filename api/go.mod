require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.39.1
	github.com/aws/aws-sdk-go-v2/config v1.31.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.20.12
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.50.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.87.0
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.0 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.18.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodbstreams v1.30.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.11.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.28.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.37.0 // indirect
	github.com/aws/smithy-go v1.23.0 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
)

replace gopkg.in/yaml.v2 => gopkg.in/yaml.v2 v2.2.8

module github.com/thomasmendez/personal-website-backend/api

go 1.22

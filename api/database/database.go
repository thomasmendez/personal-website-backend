package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "PersonalWebsiteTable"

type Database struct {
	DB *dynamodb.DynamoDB
}

func NewDatabase(awsSession *session.Session) (database *Database) {
	dynamoDbClient := dynamodb.New(awsSession)
	return &Database{dynamoDbClient}
}

package database

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const tableName = "PersonalWebsiteTable"

type Database struct {
	*dynamodb.DynamoDB
}

func NewDatabase(awsSession *session.Session) (database *Database) {
	return &Database{dynamodb.New(awsSession)}
}

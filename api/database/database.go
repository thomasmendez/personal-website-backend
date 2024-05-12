package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

const tableName = "PersonalWebsiteTable"

type Database struct {
	*dynamodb.DynamoDB
}

func NewDatabase(awsSession *session.Session) (database *Database) {
	return &Database{dynamodb.New(awsSession)}
}

func GetItem(svc dynamodbiface.DynamoDBAPI, personalWebsiteType string, sortValue string, item interface{}) (err error) {
	inputGet := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String(personalWebsiteType)},
			"sortValue":           {S: aws.String(sortValue)},
		},
		TableName: aws.String(tableName),
	}
	result, err := svc.GetItem(inputGet)
	if err != nil {
		log.Print(err)
		return err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Print(err)
		return err
	}
	return err
}

func DeleteItem(svc dynamodbiface.DynamoDBAPI, personalWebsiteType string, sortValue string) (err error) {
	key := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String(personalWebsiteType)},
		"sortValue":           {S: aws.String(sortValue)},
	}
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key:       key,
	}
	_, err = svc.DeleteItem(input)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

package database

import (
	"log"
	"reflect"

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

func unmarshalDynamodbMapSlice(queryOutput dynamodb.QueryOutput, slicePtr interface{}) error {
	items := queryOutput.Items
	sliceValue := reflect.ValueOf(slicePtr).Elem()
	elementType := sliceValue.Type().Elem()

	for _, item := range items {
		newItem := reflect.New(elementType).Interface()
		if err := dynamodbattribute.UnmarshalMap(item, newItem); err != nil {
			return err
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(newItem).Elem()))
	}

	return nil
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
		log.Printf("error in DynamoDB GetItem func: %v", err)
		return err
	}
	err = dynamodbattribute.UnmarshalMap(result.Item, item)
	if err != nil {
		log.Printf("error in DynamoDB UnmarshalMap func: %v", err)
		return err
	}
	return nil
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
		log.Printf("error in DynamoDB UnmarshalMap func: %v", err)
		return err
	}
	return nil
}

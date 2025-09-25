package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeyWork = "Work"

func GetWork(ctx context.Context, svc *dynamodb.Client, tableName string) (work []models.Work, err error) {
	work = make([]models.Work, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey and sortValue > :startDateValue"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":partitionKey": &types.AttributeValueMemberS{
				Value: partitionKeyWork,
			},
			":startDateValue": &types.AttributeValueMemberS{
				Value: "1970-01-01",
			},
		},
		ScanIndexForward: aws.Bool(false),
	}
	queryOutput, err := svc.Query(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB Query func: %v", err)
		return work, err
	}
	err = unmarshalDynamodbMapSlice(queryOutput, &work)
	return work, err
}

func PostWork(ctx context.Context, svc *dynamodb.Client, tableName string, newWork models.Work) (work models.Work, err error) {
	item, err := attributevalue.MarshalMap(newWork)
	if err != nil {
		log.Printf("error marshalling newWork: %v", err)
		return work, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB PutItem func: %v", err)
		return work, err
	}
	err = GetItem(ctx, svc, tableName, newWork.PersonalWebsiteType, newWork.SortValue, &work)
	return work, err
}

func UpdateWork(ctx context.Context, svc *dynamodb.Client, tableName string, updateWork models.Work) (work models.Work, err error) {
	err = UpdateItem(ctx, svc, tableName, updateWork, updateWork.PersonalWebsiteType, updateWork.SortValue)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return work, err
	}

	err = GetItem(ctx, svc, tableName, updateWork.PersonalWebsiteType, updateWork.SortValue, &work)
	return work, err
}

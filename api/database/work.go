package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	item := map[string]types.AttributeValue{
		"personalWebsiteType": &types.AttributeValueMemberS{Value: partitionKeyWork},
		"sortValue":           &types.AttributeValueMemberS{Value: newWork.SortValue},
		"jobTitle":            &types.AttributeValueMemberS{Value: newWork.JobTitle},
		"company":             &types.AttributeValueMemberS{Value: newWork.Company},
		"location": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"city":  &types.AttributeValueMemberS{Value: newWork.Location.City},
				"state": &types.AttributeValueMemberS{Value: newWork.Location.State},
			},
		},
		"startDate":      &types.AttributeValueMemberS{Value: newWork.StartDate},
		"endDate":        &types.AttributeValueMemberS{Value: newWork.EndDate},
		"jobRole":        &types.AttributeValueMemberS{Value: newWork.JobRole},
		"jobDescription": &types.AttributeValueMemberSS{Value: newWork.JobDescription},
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
	updateExpression := "SET #jobTitle = :jobTitleVal, #company = :companyVal, #location = :locationVal, #startDate = :startDateVal, #endDate = :endDateVal, #jobRole = :jobRoleVal, #jobDescription = :jobDescriptionVal"
	expressionAttributeNames := map[string]string{
		"#jobTitle":       "jobTitle",
		"#company":        "company",
		"#location":       "location",
		"#startDate":      "startDate",
		"#endDate":        "endDate",
		"#jobRole":        "jobRole",
		"#jobDescription": "jobDescription",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
		":jobTitleVal": &types.AttributeValueMemberS{Value: updateWork.JobTitle},
		":companyVal":  &types.AttributeValueMemberS{Value: updateWork.Company},
		":locationVal": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"city":  &types.AttributeValueMemberS{Value: updateWork.Location.City},
				"state": &types.AttributeValueMemberS{Value: updateWork.Location.State},
			},
		},
		":startDateVal":      &types.AttributeValueMemberS{Value: updateWork.StartDate},
		":endDateVal":        &types.AttributeValueMemberS{Value: updateWork.EndDate},
		":jobRoleVal":        &types.AttributeValueMemberS{Value: updateWork.JobRole},
		":jobDescriptionVal": &types.AttributeValueMemberSS{Value: updateWork.JobDescription},
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"personalWebsiteType": &types.AttributeValueMemberS{Value: partitionKeyWork},
			"sortValue":           &types.AttributeValueMemberS{Value: updateWork.SortValue},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	_, err = svc.UpdateItem(ctx, updateInput)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return work, err
	}
	err = GetItem(ctx, svc, tableName, updateWork.PersonalWebsiteType, updateWork.SortValue, &work)
	return work, err
}

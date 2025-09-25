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

const partitionKeySkillsTools = "SkillsTools"

func GetSkillsTools(ctx context.Context, svc *dynamodb.Client, tableName string) (skillsTools []models.SkillsTools, err error) {
	skillsTools = make([]models.SkillsTools, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":partitionKey": &types.AttributeValueMemberS{
				Value: partitionKeySkillsTools,
			},
		},
	}
	queryOutput, err := svc.Query(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB Query func: %v", err)
		return skillsTools, err
	}
	err = unmarshalDynamodbMapSlice(queryOutput, &skillsTools)
	return skillsTools, err
}

func PostSkillsTools(ctx context.Context, svc *dynamodb.Client, tableName string, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	item, err := attributevalue.MarshalMap(newSkillsTools)
	if err != nil {
		log.Printf("error marshalling newSkillsTools: %v", err)
		return skillsTools, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB PutItem func: %v", err)
		return skillsTools, err
	}
	err = GetItem(ctx, svc, tableName, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

func UpdateSkillsTools(ctx context.Context, svc *dynamodb.Client, tableName string, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	err = UpdateItem(ctx, svc, tableName, newSkillsTools, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return skillsTools, err
	}
	err = GetItem(ctx, svc, tableName, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

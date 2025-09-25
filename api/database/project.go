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

const partitionKeyProjects = "Projects"

func GetProjects(ctx context.Context, svc *dynamodb.Client, tableName string) (projects []models.Project, err error) {
	projects = make([]models.Project, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":partitionKey": &types.AttributeValueMemberS{
				Value: partitionKeyProjects,
			},
		},
	}
	queryOutput, err := svc.Query(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB Query func: %v", err)
		return projects, err
	}
	err = unmarshalDynamodbMapSlice(queryOutput, &projects)
	return projects, err
}

func PostProject(ctx context.Context, svc *dynamodb.Client, tableName string, newProject models.Project) (project models.Project, err error) {
	item, err := attributevalue.MarshalMap(newProject)
	if err != nil {
		log.Printf("error marshalling newProject: %v", err)
		return project, err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(ctx, input)
	if err != nil {
		log.Printf("error in DynamoDB PutItem func: %v", err)
		return project, err
	}
	err = GetItem(ctx, svc, tableName, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

func UpdateProject(ctx context.Context, svc *dynamodb.Client, tableName string, newProject models.Project) (project models.Project, err error) {
	err = UpdateItem(ctx, svc, tableName, newProject, newProject.PersonalWebsiteType, newProject.SortValue)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return project, err
	}
	err = GetItem(ctx, svc, tableName, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeyProjects = "Projects"

func GetProjects(svc dynamodbiface.DynamoDBAPI) (projects []models.Project, err error) {
	projects = make([]models.Project, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionKey": {
				S: aws.String(partitionKeyProjects),
			},
		},
	}
	queryOutput, err := svc.Query(input)
	if err != nil {
		return projects, err
	}
	for _, item := range queryOutput.Items {
		var projectsItem models.Project
		err := dynamodbattribute.UnmarshalMap(item, &projectsItem)
		if err != nil {
			return projects, err
		}
		projects = append(projects, projectsItem)
	}
	return projects, nil
}

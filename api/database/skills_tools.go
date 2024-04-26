package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func GetSkillsTools(svc dynamodbiface.DynamoDBAPI) (skillsTools []models.SkillsTools, err error) {
	skillsTools = make([]models.SkillsTools, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionKey": {
				S: aws.String("SkillsTools"),
			},
		},
	}

	queryOutput, err := svc.Query(input)
	if err != nil {
		return skillsTools, err
	}

	for _, item := range queryOutput.Items {
		var skillsToolsItem models.SkillsTools
		err := dynamodbattribute.UnmarshalMap(item, &skillsToolsItem)
		if err != nil {
			return skillsTools, err
		}
		skillsTools = append(skillsTools, skillsToolsItem)
	}

	return skillsTools, nil
}

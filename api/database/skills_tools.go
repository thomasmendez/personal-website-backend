package database

import (
	"log"

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

func PostSkillsTools(svc dynamodbiface.DynamoDBAPI, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String("SkillsTools")},
		"sortValue":           {S: aws.String(newSkillsTools.SortValue)},
		"skillsToolsCategory": {S: aws.String(newSkillsTools.SkillsToolsCategory)},
		"skillsToolsType":     {S: aws.String(newSkillsTools.SkillsToolsType)},
	}
	skillsToolsList := make([]*string, len(newSkillsTools.SkillsToolsList))
	for i, skillTool := range newSkillsTools.SkillsToolsList {
		skillsToolsList[i] = aws.String(skillTool)
	}
	item["skillsToolsList"] = &dynamodb.AttributeValue{SS: skillsToolsList}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Print(err)
		return skillsTools, err
	}
	inputGet := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String("SkillsTools")},
			"sortValue":           {S: aws.String(skillsTools.SortValue)},
		},
		TableName: aws.String(tableName),
	}

	result, err := svc.GetItem(inputGet)
	if err != nil {
		log.Print(err)
		return skillsTools, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &skillsTools)
	if err != nil {
		return skillsTools, err
	}

	return skillsTools, nil
}

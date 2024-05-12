package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeySkillsTools = "SkillsTools"

func GetSkillsTools(svc dynamodbiface.DynamoDBAPI) (skillsTools []models.SkillsTools, err error) {
	skillsTools = make([]models.SkillsTools, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionKey": {
				S: aws.String(partitionKeySkillsTools),
			},
		},
	}
	queryOutput, err := svc.Query(input)
	if err != nil {
		log.Printf("error in DynamoDB Query func: %v", err)
		return skillsTools, err
	}
	err = unmarshalDynamodbMapSlice(*queryOutput, &skillsTools)
	return skillsTools, err
}

func PostSkillsTools(svc dynamodbiface.DynamoDBAPI, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String(partitionKeySkillsTools)},
		"sortValue":           {S: aws.String(newSkillsTools.SortValue)},
		"category":            {S: aws.String(newSkillsTools.Category)},
		"type":                {S: aws.String(newSkillsTools.Type)},
		"list":                {SS: aws.StringSlice(newSkillsTools.List)},
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("error in DynamoDB PutItem func: %v", err)
		return skillsTools, err
	}
	err = GetItem(svc, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

func UpdateSkillsTools(svc dynamodbiface.DynamoDBAPI, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	updateExpression := "SET #category = :categoryVal, #type = :typeVal, #list = :listVal"
	expressionAttributeNames := map[string]*string{
		"#category": aws.String("category"),
		"#type":     aws.String("type"),
		"#list":     aws.String("list"),
	}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":categoryVal": {S: aws.String(newSkillsTools.Category)},
		":typeVal":     {S: aws.String(newSkillsTools.Type)},
		":listVal":     {SS: aws.StringSlice(newSkillsTools.List)},
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String(partitionKeySkillsTools)},
			"sortValue":           {S: aws.String(newSkillsTools.SortValue)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return skillsTools, err
	}
	err = GetItem(svc, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

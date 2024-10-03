package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeySkillsTools = "SkillsTools"

func GetSkillsTools(svc dynamodbiface.DynamoDBAPI, tableName string) (skillsTools []models.SkillsTools, err error) {
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

func PostSkillsTools(svc dynamodbiface.DynamoDBAPI, tableName string, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	item, err := dynamodbattribute.MarshalMap(newSkillsTools)
	if err != nil {
		log.Printf("error marshalling newSkillsTools: %v", err)
		return skillsTools, err
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
	err = GetItem(svc, tableName, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

func UpdateSkillsTools(svc dynamodbiface.DynamoDBAPI, tableName string, newSkillsTools models.SkillsTools) (skillsTools models.SkillsTools, err error) {
	// Marshal the Categories field into DynamoDB-compatible values
	categoriesAttrVal, err := dynamodbattribute.Marshal(newSkillsTools.Categories)
	if err != nil {
		log.Printf("error marshalling Categories: %v", err)
		return skillsTools, err
	}
	// Update expression for the fields you want to update
	updateExpression := "SET #categories = :categoriesVal"
	// Expression attribute names (used to avoid reserved keywords)
	expressionAttributeNames := map[string]*string{
		"#categories": aws.String("categories"),
	}
	// Expression attribute values (setting the values to be updated)
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":categoriesVal": categoriesAttrVal,
	}
	// Define the primary key (personalWebsiteType and sortValue)
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
	err = GetItem(svc, tableName, newSkillsTools.PersonalWebsiteType, newSkillsTools.SortValue, &skillsTools)
	return skillsTools, err
}

package database

import (
	"log"
	"strconv"

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

func PostProject(svc dynamodbiface.DynamoDBAPI, newProject models.Project) (project models.Project, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String(newProject.PersonalWebsiteType)},
		"sortValue":           {S: aws.String(newProject.SortValue)},
		"category":            {S: aws.String(newProject.Category)},
		"name":                {S: aws.String(newProject.Name)},
		"description":         {S: aws.String(newProject.Description)},
		"featuresDescription": {S: aws.String(newProject.FeaturesDescription)},
		"role":                {S: aws.String(newProject.Role)},
		"tasks":               {SS: aws.StringSlice(newProject.Tasks)},
		"tools":               {SS: aws.StringSlice(newProject.Tools)},
		"duration":            {S: aws.String(newProject.Duration)},
		"startDate":           {S: aws.String(newProject.StartDate)},
		"endDate":             {S: aws.String(newProject.EndDate)},
	}
	if newProject.TeamSize != nil {
		item["teamSize"] = &dynamodb.AttributeValue{N: aws.String(strconv.Itoa(*newProject.TeamSize))}
	} else {
		item["teamSize"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.TeamRoles != nil {
		item["teamRoles"] = &dynamodb.AttributeValue{SS: aws.StringSlice(*newProject.TeamRoles)}
	} else {
		item["teamRoles"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.CloudServices != nil {
		item["cloudServices"] = &dynamodb.AttributeValue{SS: aws.StringSlice(*newProject.CloudServices)}
	} else {
		item["cloudServices"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.CloudServices != nil {
		item["notes"] = &dynamodb.AttributeValue{S: aws.String(*newProject.Notes)}
	} else {
		item["notes"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.CloudServices != nil {
		item["link"] = &dynamodb.AttributeValue{S: aws.String(*newProject.Link)}
	} else {
		item["link"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.CloudServices != nil {
		item["linkType"] = &dynamodb.AttributeValue{S: aws.String(*newProject.LinkType)}
	} else {
		item["linkType"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if newProject.CloudServices != nil {
		item["mediaLink"] = &dynamodb.AttributeValue{S: aws.String(*newProject.MediaLink)}
	} else {
		item["mediaLink"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Print(err)
		return project, err
	}
	err = GetItem(svc, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

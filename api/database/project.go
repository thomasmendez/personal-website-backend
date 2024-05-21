package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeyProjects = "Projects"

func GetProjects(svc dynamodbiface.DynamoDBAPI, tableName string) (projects []models.Project, err error) {
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
		log.Printf("error in DynamoDB Query func: %v", err)
		return projects, err
	}
	err = unmarshalDynamodbMapSlice(*queryOutput, &projects)
	return projects, err
}

func PostProject(svc dynamodbiface.DynamoDBAPI, tableName string, newProject models.Project) (project models.Project, err error) {
	item := projectItem(newProject)
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Printf("error in DynamoDB PutItem func: %v", err)
		return project, err
	}
	err = GetItem(svc, tableName, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

func projectItem(project models.Project) (item map[string]*dynamodb.AttributeValue) {
	item = map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String(project.PersonalWebsiteType)},
		"sortValue":           {S: aws.String(project.SortValue)},
		"category":            {S: aws.String(project.Category)},
		"name":                {S: aws.String(project.Name)},
		"description":         {S: aws.String(project.Description)},
		"featuresDescription": {S: aws.String(project.FeaturesDescription)},
		"role":                {S: aws.String(project.Role)},
		"tasks":               {SS: aws.StringSlice(project.Tasks)},
		"tools":               {SS: aws.StringSlice(project.Tools)},
		"duration":            {S: aws.String(project.Duration)},
		"startDate":           {S: aws.String(project.StartDate)},
		"endDate":             {S: aws.String(project.EndDate)},
	}
	if project.TeamSize != nil {
		item["teamSize"] = &dynamodb.AttributeValue{N: aws.String(*project.TeamSize)}
	} else {
		item["teamSize"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.TeamRoles != nil {
		item["teamRoles"] = &dynamodb.AttributeValue{SS: aws.StringSlice(*project.TeamRoles)}
	} else {
		item["teamRoles"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.CloudServices != nil {
		item["cloudServices"] = &dynamodb.AttributeValue{SS: aws.StringSlice(*project.CloudServices)}
	} else {
		item["cloudServices"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.Notes != nil {
		item["notes"] = &dynamodb.AttributeValue{S: aws.String(*project.Notes)}
	} else {
		item["notes"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.Link != nil {
		item["link"] = &dynamodb.AttributeValue{S: aws.String(*project.Link)}
	} else {
		item["link"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.LinkType != nil {
		item["linkType"] = &dynamodb.AttributeValue{S: aws.String(*project.LinkType)}
	} else {
		item["linkType"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	if project.MediaLink != nil {
		item["mediaLink"] = &dynamodb.AttributeValue{S: aws.String(*project.MediaLink)}
	} else {
		item["mediaLink"] = &dynamodb.AttributeValue{NULL: aws.Bool(true)}
	}
	return item
}

func UpdateProject(svc dynamodbiface.DynamoDBAPI, tableName string, newProject models.Project) (project models.Project, err error) {
	item := projectItem(newProject)

	updateExpression := "SET #category = :categoryVal, #name = :nameVal, #description = :descriptionVal, #featuresDescription = :featuresDescriptionVal, #role = :roleVal, #tasks = :tasksVal, #teamSize = :teamSizeVal, #teamRoles = :teamRolesVal, #cloudServices = :cloudServicesVal, #tools = :toolsVal, #duration = :durationVal, #startDate = :startDateVal, #endDate = :endDateVal, #notes = :notesVal, #link = :linkVal, #linkType = :linkTypeVal, #mediaLink = :mediaLinkVal"
	expressionAttributeNames := map[string]*string{
		"#category":            aws.String("category"),
		"#name":                aws.String("name"),
		"#description":         aws.String("description"),
		"#featuresDescription": aws.String("featuresDescription"),
		"#role":                aws.String("role"),
		"#tasks":               aws.String("tasks"),
		"#teamSize":            aws.String("teamSize"),
		"#teamRoles":           aws.String("teamRoles"),
		"#cloudServices":       aws.String("cloudServices"),
		"#tools":               aws.String("tools"),
		"#duration":            aws.String("duration"),
		"#startDate":           aws.String("startDate"),
		"#endDate":             aws.String("endDate"),
		"#notes":               aws.String("notes"),
		"#link":                aws.String("link"),
		"#linkType":            aws.String("linkType"),
		"#mediaLink":           aws.String("mediaLink"),
	}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":categoryVal":            item["category"],
		":nameVal":                item["name"],
		":descriptionVal":         item["description"],
		":featuresDescriptionVal": item["featuresDescription"],
		":roleVal":                item["role"],
		":tasksVal":               item["tasks"],
		":teamSizeVal":            item["teamSize"],
		":teamRolesVal":           item["teamRoles"],
		":cloudServicesVal":       item["cloudServices"],
		":toolsVal":               item["tools"],
		":durationVal":            item["duration"],
		":startDateVal":           item["startDate"],
		":endDateVal":             item["endDate"],
		":notesVal":               item["notes"],
		":linkVal":                item["link"],
		":linkTypeVal":            item["linkType"],
		":mediaLinkVal":           item["mediaLink"],
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String(partitionKeyProjects)},
			"sortValue":           {S: aws.String(newProject.SortValue)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return project, err
	}
	err = GetItem(svc, tableName, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

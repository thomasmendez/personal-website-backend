package database

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	item := projectItem(newProject)
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

func projectItem(project models.Project) (item map[string]types.AttributeValue) {
	item = map[string]types.AttributeValue{
		"personalWebsiteType": &types.AttributeValueMemberS{Value: project.PersonalWebsiteType},
		"sortValue":           &types.AttributeValueMemberS{Value: project.SortValue},
		"category":            &types.AttributeValueMemberS{Value: project.Category},
		"name":                &types.AttributeValueMemberS{Value: project.Name},
		"description":         &types.AttributeValueMemberS{Value: project.Description},
		"featuresDescription": &types.AttributeValueMemberS{Value: project.FeaturesDescription},
		"role":                &types.AttributeValueMemberS{Value: project.Role},
		"tasks":               &types.AttributeValueMemberSS{Value: project.Tasks},
		"tools":               &types.AttributeValueMemberSS{Value: project.Tools},
		"duration":            &types.AttributeValueMemberS{Value: project.Duration},
		"startDate":           &types.AttributeValueMemberS{Value: project.StartDate},
		"endDate":             &types.AttributeValueMemberS{Value: project.EndDate},
	}
	if project.TeamSize != nil {
		item["teamSize"] = &types.AttributeValueMemberN{Value: *project.TeamSize}
	} else {
		item["teamSize"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.TeamRoles != nil {
		item["teamRoles"] = &types.AttributeValueMemberSS{Value: *project.TeamRoles}
	} else {
		item["teamRoles"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.CloudServices != nil {
		item["cloudServices"] = &types.AttributeValueMemberSS{Value: *project.CloudServices}
	} else {
		item["cloudServices"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.Notes != nil {
		item["notes"] = &types.AttributeValueMemberS{Value: *project.Notes}
	} else {
		item["notes"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.Link != nil {
		item["link"] = &types.AttributeValueMemberS{Value: *project.Link}
	} else {
		item["link"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.LinkType != nil {
		item["linkType"] = &types.AttributeValueMemberS{Value: *project.LinkType}
	} else {
		item["linkType"] = &types.AttributeValueMemberNULL{Value: true}
	}
	if project.MediaLink != nil {
		item["mediaLink"] = &types.AttributeValueMemberS{Value: *project.MediaLink}
	} else {
		item["mediaLink"] = &types.AttributeValueMemberNULL{Value: true}
	}
	return item
}

func UpdateProject(ctx context.Context, svc *dynamodb.Client, tableName string, newProject models.Project) (project models.Project, err error) {
	item := projectItem(newProject)

	updateExpression := "SET #category = :categoryVal, #name = :nameVal, #description = :descriptionVal, #featuresDescription = :featuresDescriptionVal, #role = :roleVal, #tasks = :tasksVal, #teamSize = :teamSizeVal, #teamRoles = :teamRolesVal, #cloudServices = :cloudServicesVal, #tools = :toolsVal, #duration = :durationVal, #startDate = :startDateVal, #endDate = :endDateVal, #notes = :notesVal, #link = :linkVal, #linkType = :linkTypeVal, #mediaLink = :mediaLinkVal"
	expressionAttributeNames := map[string]string{
		"#category":            "category",
		"#name":                "name",
		"#description":         "description",
		"#featuresDescription": "featuresDescription",
		"#role":                "role",
		"#tasks":               "tasks",
		"#teamSize":            "teamSize",
		"#teamRoles":           "teamRoles",
		"#cloudServices":       "cloudServices",
		"#tools":               "tools",
		"#duration":            "duration",
		"#startDate":           "startDate",
		"#endDate":             "endDate",
		"#notes":               "notes",
		"#link":                "link",
		"#linkType":            "linkType",
		"#mediaLink":           "mediaLink",
	}
	expressionAttributeValues := map[string]types.AttributeValue{
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
		Key: map[string]types.AttributeValue{
			"personalWebsiteType": &types.AttributeValueMemberS{Value: partitionKeyProjects},
			"sortValue":           &types.AttributeValueMemberS{Value: newProject.SortValue},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	_, err = svc.UpdateItem(ctx, updateInput)
	if err != nil {
		log.Printf("error in DynamoDB UpdateItem func: %v", err)
		return project, err
	}
	err = GetItem(ctx, svc, tableName, newProject.PersonalWebsiteType, newProject.SortValue, &project)
	return project, err
}

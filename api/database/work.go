package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

const partitionKeyWork = "Work"

func GetWork(svc dynamodbiface.DynamoDBAPI) (work []models.Work, err error) {
	work = make([]models.Work, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey and sortValue > :startDateValue"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionKey": {
				S: aws.String(partitionKeyWork),
			},
			":startDateValue": {
				S: aws.String("1970-01-01"),
			},
		},
		ScanIndexForward: aws.Bool(false),
	}
	queryOutput, err := svc.Query(input)
	if err != nil {
		return work, err
	}
	for _, item := range queryOutput.Items {
		var workItem models.Work
		err := dynamodbattribute.UnmarshalMap(item, &workItem)
		if err != nil {
			return work, err
		}
		work = append(work, workItem)
	}
	return work, nil
}

func PostWork(svc dynamodbiface.DynamoDBAPI, newWork models.Work) (work models.Work, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String(partitionKeyWork)},
		"sortValue":           {S: aws.String(newWork.SortValue)},
		"jobTitle":            {S: aws.String(newWork.JobTitle)},
		"company":             {S: aws.String(newWork.Company)},
		"location": {
			M: map[string]*dynamodb.AttributeValue{
				"city":  {S: aws.String(newWork.Location.City)},
				"state": {S: aws.String(newWork.Location.State)},
			},
		},
		"startDate":      {S: aws.String(newWork.StartDate)},
		"endDate":        {S: aws.String(newWork.EndDate)},
		"jobRole":        {S: aws.String(newWork.JobRole)},
		"jobDescription": {SS: aws.StringSlice(newWork.JobDescription)},
	}
	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}
	_, err = svc.PutItem(input)
	if err != nil {
		log.Print(err)
		return work, err
	}
	err = GetItem(svc, newWork.PersonalWebsiteType, newWork.SortValue, &work)
	return work, err
}

func UpdateWork(svc dynamodbiface.DynamoDBAPI, newWork models.Work) (work models.Work, err error) {
	updateExpression := "SET #jobTitle = :jobTitleVal, #company = :companyVal, #location = :locationVal, #startDate = :startDateVal, #endDate = :endDateVal, #jobRole = :jobRoleVal, #jobDescription = :jobDescriptionVal"
	expressionAttributeNames := map[string]*string{
		"#jobTitle":       aws.String("jobTitle"),
		"#company":        aws.String("company"),
		"#location":       aws.String("location"),
		"#startDate":      aws.String("startDate"),
		"#endDate":        aws.String("endDate"),
		"#jobRole":        aws.String("jobRole"),
		"#jobDescription": aws.String("jobDescription"),
	}
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":jobTitleVal": {S: aws.String(newWork.JobTitle)},
		":companyVal":  {S: aws.String(newWork.Company)},
		":locationVal": {M: map[string]*dynamodb.AttributeValue{
			"city":  {S: aws.String(newWork.Location.City)},
			"state": {S: aws.String(newWork.Location.State)},
		}},
		":startDateVal":      {S: aws.String(newWork.StartDate)},
		":endDateVal":        {S: aws.String(newWork.EndDate)},
		":jobRoleVal":        {S: aws.String(newWork.JobRole)},
		":jobDescriptionVal": {SS: aws.StringSlice(newWork.JobDescription)},
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String(partitionKeyWork)},
			"sortValue":           {S: aws.String(newWork.SortValue)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		log.Print(err)
		return work, err
	}
	err = GetItem(svc, newWork.PersonalWebsiteType, newWork.SortValue, &work)
	return work, err
}

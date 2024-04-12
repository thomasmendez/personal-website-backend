package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (db *Database) GetWork() (work []models.Work, err error) {
	work = make([]models.Work, 0)
	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("personalWebsiteType = :partitionKey and sortValue > :startDateValue"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":partitionKey": {
				S: aws.String("Job"),
			},
			":startDateValue": {
				S: aws.String("1970-01-01"),
			},
		},
		ScanIndexForward: aws.Bool(false),
	}

	queryOutput, err := db.DB.Query(input)
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

func (db *Database) PostWork(newWork models.Work) (work models.Work, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String("Job")},
		"sortValue":           {S: aws.String(newWork.SortValue)},
		"jobTitle":            {S: aws.String(newWork.JobTitle)},
		"company":             {S: aws.String(newWork.Company)},
		"location": {
			M: map[string]*dynamodb.AttributeValue{
				"city":  {S: aws.String(newWork.Location.City)},
				"state": {S: aws.String(newWork.Location.State)},
			},
		},
		"startDate": {S: aws.String(newWork.StartDate)},
		"endDate":   {S: aws.String(newWork.EndDate)},
		"jobRole":   {S: aws.String(newWork.JobRole)},
	}

	jobDescription := make([]*string, len(newWork.JobDescription))
	for i, desc := range newWork.JobDescription {
		jobDescription[i] = aws.String(desc)
	}
	item["jobDescription"] = &dynamodb.AttributeValue{SS: jobDescription}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(tableName),
	}

	_, err = db.DB.PutItem(input)
	if err != nil {
		log.Print(err)
		return work, err
	}

	inputGet := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String("Job")},
			"sortValue":           {S: aws.String(newWork.SortValue)},
		},
		TableName: aws.String(tableName),
	}

	result, err := db.DB.GetItem(inputGet)
	if err != nil {
		log.Print(err)
		return work, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &work)
	if err != nil {
		return work, err
	}

	return work, nil
}

func (db *Database) UpdateWork(newWork models.Work) (work models.Work, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String("Job")},
		"sortValue":           {S: aws.String(newWork.SortValue)},
		"jobTitle":            {S: aws.String(newWork.JobTitle)},
		"company":             {S: aws.String(newWork.Company)},
		"location": {
			M: map[string]*dynamodb.AttributeValue{
				"city":  {S: aws.String(newWork.Location.City)},
				"state": {S: aws.String(newWork.Location.State)},
			},
		},
		"startDate": {S: aws.String(newWork.StartDate)},
		"endDate":   {S: aws.String(newWork.EndDate)},
		"jobRole":   {S: aws.String(newWork.JobRole)},
	}

	jobDescription := make([]*string, len(newWork.JobDescription))
	for i, desc := range newWork.JobDescription {
		jobDescription[i] = aws.String(desc)
	}
	item["jobDescription"] = &dynamodb.AttributeValue{SS: jobDescription}

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
		":jobTitleVal":       {S: aws.String(newWork.JobTitle)},
		":companyVal":        {S: aws.String(newWork.Company)},
		":locationVal":       {M: item["location"].M},
		":startDateVal":      {S: aws.String(newWork.StartDate)},
		":endDateVal":        {S: aws.String(newWork.EndDate)},
		":jobRoleVal":        {S: aws.String(newWork.JobRole)},
		":jobDescriptionVal": item["jobDescription"],
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String("Job")},
			"sortValue":           {S: aws.String(newWork.StartDate)},
		},
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}

	_, err = db.DB.UpdateItem(updateInput)
	if err != nil {
		log.Print(err)
		return work, err
	}

	inputGet := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"personalWebsiteType": {S: aws.String("Job")},
			"sortValue":           {S: aws.String(newWork.StartDate)},
		},
		TableName: aws.String(tableName),
	}

	result, err := db.DB.GetItem(inputGet)
	if err != nil {
		log.Print(err)
		return work, err
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &work)
	if err != nil {
		return work, err
	}

	return work, nil
}

// func (db *Database) GetWork() (work []models.Work, err error) {
// 	input := &dynamodb.ScanInput{
// 		TableName: aws.String("PersonalWebsiteTable"),
// 	}

// 	scanOutput, err := db.DB.Scan(input)
// 	if err != nil {
// 		return work, err
// 	}

// 	// Iterate through the result and construct a list of JobDescription objects
// 	for _, item := range scanOutput.Items {

// 		locationMap := item["location"].M
// 		dateMap := item["date"].M

// 		location := models.Location{
// 			City:  aws.StringValue(locationMap["city"].S),
// 			State: aws.StringValue(locationMap["state"].S),
// 		}

// 		date := models.Date{
// 			StartDate: aws.StringValue(dateMap["startDate"].S),
// 			EndDate:   aws.StringValue(dateMap["endDate"].S),
// 		}

// 		workItem := models.Work{
// 			JobTitle:       aws.StringValue(item["jobTitle"].S),
// 			Company:        aws.StringValue(item["company"].S),
// 			Location:       location,
// 			Date:           date,
// 			JobRole:        aws.StringValue(item["jobRole"].S),
// 			JobDescription: aws.StringValueSlice(item["jobDescription"].SS),
// 		}
// 		work = append(work, workItem)
// 	}

// 	return work, err
// }

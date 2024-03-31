package database

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (db *Database) GetWork() (work []models.Work, err error) {
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

	// Unmarshal the results
	for _, item := range queryOutput.Items {
		var workItem models.Work
		err := dynamodbattribute.UnmarshalMap(item, &workItem)
		if err != nil {
			return nil, err
		}
		work = append(work, workItem)
	}

	return work, nil
}

func (db *Database) PostWork(newWork models.Work) (work models.Work, err error) {
	item := map[string]*dynamodb.AttributeValue{
		"personalWebsiteType": {S: aws.String("Job")},
		"sortValue":           {S: aws.String(newWork.StartDate)},
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
		log.Print("here3")
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
		log.Print("here2")
		log.Print(err)
		return work, err
	}

	log.Print(result)

	work, err = ParseDynamoDBItemToWork(result.Item)
	if err != nil {
		log.Print("here1")
		log.Print(err)
		return work, err
	}

	return work, nil
}

func ParseDynamoDBItemToWork(item map[string]*dynamodb.AttributeValue) (work models.Work, err error) {
	if jobTitleAttr, ok := item["jobTitle"]; ok {
		work.JobTitle = aws.StringValue(jobTitleAttr.S)
	}
	if companyAttr, ok := item["company"]; ok {
		work.Company = aws.StringValue(companyAttr.S)
	}
	if locationAttr, ok := item["location"]; ok {
		if cityAttr, ok := locationAttr.M["city"]; ok {
			work.Location.City = aws.StringValue(cityAttr.S)
		}
		if stateAttr, ok := locationAttr.M["state"]; ok {
			work.Location.State = aws.StringValue(stateAttr.S)
		}
	}
	if startDateAttr, ok := item["startDate"]; ok {
		work.StartDate = aws.StringValue(startDateAttr.S)
	}
	if endDateAttr, ok := item["endDate"]; ok {
		work.EndDate = aws.StringValue(endDateAttr.S)
	}
	if jobRoleAttr, ok := item["jobRole"]; ok {
		work.JobRole = aws.StringValue(jobRoleAttr.S)
	}
	if jobDescriptionAttr, ok := item["jobDescription"]; ok {
		if jobDescriptionAttr.SS != nil {
			work.JobDescription = make([]string, len(jobDescriptionAttr.SS))
			for i, desc := range jobDescriptionAttr.SS {
				work.JobDescription[i] = aws.StringValue(desc)
			}
		}
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

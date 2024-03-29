package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (db *Database) GetWork() (work []models.Work, err error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String("PersonalWebsiteTable"),
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

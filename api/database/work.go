package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func (db *Database) GetWork() (work []models.Work, err error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("PersonalWebsiteTable"),
	}

	scanOutput, err := db.DB.Scan(input)
	if err != nil {
		return work, err
	}

	// Iterate through the result and construct a list of JobDescription objects
	for _, item := range scanOutput.Items {

		locationMap := item["location"].M
		dateMap := item["date"].M

		location := models.Location{
			City:  aws.StringValue(locationMap["city"].S),
			State: aws.StringValue(locationMap["state"].S),
		}

		date := models.Date{
			StartDate: aws.StringValue(dateMap["startDate"].S),
			EndDate:   aws.StringValue(dateMap["endDate"].S),
		}

		workItem := models.Work{
			JobTitle:       aws.StringValue(item["jobTitle"].S),
			Company:        aws.StringValue(item["company"].S),
			Location:       location,
			Date:           date,
			JobRole:        aws.StringValue(item["jobRole"].S),
			JobDescription: aws.StringValueSlice(item["jobDescription"].SS),
		}
		work = append(work, workItem)
	}

	return work, err
}

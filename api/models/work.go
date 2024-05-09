package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Work struct {
	PersonalWebsiteType string   `json:"personalWebsiteType"`
	SortValue           string   `json:"sortValue"`
	JobTitle            string   `json:"jobTitle"`
	Company             string   `json:"company"`
	Location            Location `json:"location"`
	StartDate           string   `json:"startDate"`
	EndDate             string   `json:"endDate"`
	JobRole             string   `json:"jobRole"`
	JobDescription      []string `json:"jobDescription"`
}

type Location struct {
	City  string `json:"city"`
	State string `json:"state"`
}

// Work model used for test cases
var TestWork = Work{
	PersonalWebsiteType: "Work",
	SortValue:           "2020-01-01",
	JobTitle:            "Software Engineer",
	Company:             "ABC Inc",
	Location: Location{
		City:  "New York",
		State: "NY",
	},
	StartDate:      "2020-01-01",
	EndDate:        "2020-12-31",
	JobRole:        "Backend Developer",
	JobDescription: []string{"Developed backend systems", "Optimized database queries"},
}

// Work Item model used for dynamodb
var TestWorkItem = map[string]*dynamodb.AttributeValue{
	"personalWebsiteType": {S: aws.String("Work")},
	"sortValue":           {S: aws.String("2020-01-01")},
	"jobTitle":            {S: aws.String("Software Engineer")},
	"company":             {S: aws.String("ABC Inc")},
	"location": {
		M: map[string]*dynamodb.AttributeValue{
			"city":  {S: aws.String("New York")},
			"state": {S: aws.String("NY")},
		},
	},
	"startDate":      {S: aws.String("2020-01-01")},
	"endDate":        {S: aws.String("2020-12-31")},
	"jobRole":        {S: aws.String("Backend Developer")},
	"jobDescription": {SS: []*string{aws.String("Developed backend systems"), aws.String("Optimized database queries")}},
}

package tests

import (
	"github.com/thomasmendez/personal-website-backend/api/models"
)

// Work model used for test cases
var TestWork = models.Work{
	PersonalWebsiteType: "Work",
	SortValue:           "2020-01-01",
	JobTitle:            "Software Engineer",
	Company:             "ABC Inc",
	Location: models.Location{
		City:  "New York",
		State: "NY",
	},
	StartDate:      "2020-01-01",
	EndDate:        "2020-12-31",
	JobRole:        "Backend Developer",
	JobDescription: []string{"Developed backend systems", "Optimized database queries"},
}

// Work Item model used for dynamodb
// var TestWorkItem = map[string]*dynamodb.AttributeValue{
// 	"personalWebsiteType": {S: aws.String("Work")},
// 	"sortValue":           {S: aws.String("2020-01-01")},
// 	"jobTitle":            {S: aws.String("Software Engineer")},
// 	"company":             {S: aws.String("ABC Inc")},
// 	"location": {
// 		M: map[string]*dynamodb.AttributeValue{
// 			"city":  {S: aws.String("New York")},
// 			"state": {S: aws.String("NY")},
// 		},
// 	},
// 	"startDate":      {S: aws.String("2020-01-01")},
// 	"endDate":        {S: aws.String("2020-12-31")},
// 	"jobRole":        {S: aws.String("Backend Developer")},
// 	"jobDescription": {SS: []*string{aws.String("Developed backend systems"), aws.String("Optimized database queries")}},
// }

// func AssertWork(t *testing.T, expectedWork models.Work, actualWork models.Work) {
// 	assert.Equal(t, expectedWork.PersonalWebsiteType, actualWork.PersonalWebsiteType)
// 	assert.Equal(t, expectedWork.SortValue, actualWork.SortValue)
// 	assert.Equal(t, expectedWork.JobTitle, actualWork.JobTitle)
// 	assert.Equal(t, expectedWork.Company, actualWork.Company)
// 	assert.Equal(t, expectedWork.Location.City, actualWork.Location.City)
// 	assert.Equal(t, expectedWork.Location.State, actualWork.Location.State)
// 	assert.Equal(t, expectedWork.StartDate, actualWork.StartDate)
// 	assert.Equal(t, expectedWork.EndDate, actualWork.EndDate)
// 	assert.Equal(t, expectedWork.JobRole, actualWork.JobRole)
// 	assert.Equal(t, expectedWork.JobDescription, actualWork.JobDescription)
// }

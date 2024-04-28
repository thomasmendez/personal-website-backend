package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func TestWorkGet(t *testing.T) {

	mockDB := &mockDynamoDB{}

	for _, test := range []struct {
		label         string
		mockQueryFunc func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
		expectedWork  []models.Work
		expectedError error
	}{
		{
			label: "valid query output",
			mockQueryFunc: func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
				mockOutput := &dynamodb.QueryOutput{
					Items: []map[string]*dynamodb.AttributeValue{
						{
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
						},
					},
				}
				return mockOutput, nil
			},
			expectedWork: []models.Work{
				{
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
				},
			},
			expectedError: nil,
		},
		{
			label: "query error",
			mockQueryFunc: func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
				return nil, errors.New("error querying database")
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			mockDB.QueryFunc = test.mockQueryFunc

			result, err := GetWork(mockDB)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error querying database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, 1)
			for i, work := range result {
				assert.Equal(t, test.expectedWork[i].PersonalWebsiteType, work.PersonalWebsiteType)
				assert.Equal(t, test.expectedWork[i].SortValue, work.SortValue)
				assert.Equal(t, test.expectedWork[i].JobTitle, work.JobTitle)
				assert.Equal(t, test.expectedWork[i].Company, work.Company)
				assert.Equal(t, test.expectedWork[i].Location.City, work.Location.City)
				assert.Equal(t, test.expectedWork[i].Location.State, work.Location.State)
				assert.Equal(t, test.expectedWork[i].StartDate, work.StartDate)
				assert.Equal(t, test.expectedWork[i].EndDate, work.EndDate)
				assert.Equal(t, test.expectedWork[i].JobRole, work.JobRole)
				assert.Equal(t, test.expectedWork[i].JobDescription, work.JobDescription)
			}
		})
	}
}

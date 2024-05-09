package database

import (
	"errors"
	"testing"

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
						models.TestWorkItem,
					},
				}
				return mockOutput, nil
			},
			expectedWork:  []models.Work{models.TestWork},
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
				assertWork(t, test.expectedWork[i], work)
			}
		})
	}
}

func assertWork(t *testing.T, expectedWork models.Work, actualWork models.Work) {
	assert.Equal(t, expectedWork.PersonalWebsiteType, actualWork.PersonalWebsiteType)
	assert.Equal(t, expectedWork.SortValue, actualWork.SortValue)
	assert.Equal(t, expectedWork.JobTitle, actualWork.JobTitle)
	assert.Equal(t, expectedWork.Company, actualWork.Company)
	assert.Equal(t, expectedWork.Location.City, actualWork.Location.City)
	assert.Equal(t, expectedWork.Location.State, actualWork.Location.State)
	assert.Equal(t, expectedWork.StartDate, actualWork.StartDate)
	assert.Equal(t, expectedWork.EndDate, actualWork.EndDate)
	assert.Equal(t, expectedWork.JobRole, actualWork.JobRole)
	assert.Equal(t, expectedWork.JobDescription, actualWork.JobDescription)
}

func TestPostWork(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label         string
		newWork       models.Work
		mockPutFunc   func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		mockGetFunc   func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedWork  models.Work
		expectedError error
	}{
		{
			label:   "valid query output",
			newWork: models.TestWork,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: models.TestWorkItem,
				}
				return mockOutput, nil
			},
			expectedWork:  models.TestWork,
			expectedError: nil,
		},
		{
			label:   "query error",
			newWork: models.TestWork,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("error getting item from database")
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			mockDB.PutFunc = test.mockPutFunc
			mockDB.GetFunc = test.mockGetFunc

			result, err := PostWork(mockDB, test.newWork)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error getting item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			assertWork(t, test.expectedWork, result)
		})
	}
}

func TestUpdateWork(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label          string
		updateWork     models.Work
		mockUpdateFunc func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
		mockGetFunc    func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedWork   models.Work
		expectedError  error
	}{
		{
			label:      "valid query output",
			updateWork: models.TestWork,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: models.TestWorkItem,
				}
				return mockOutput, nil
			},
			expectedWork:  models.TestWork,
			expectedError: nil,
		},
		{
			label:      "query error",
			updateWork: models.TestWork,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("error updating item from database")
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			mockDB.UpdateFunc = test.mockUpdateFunc
			mockDB.GetFunc = test.mockGetFunc

			result, err := UpdateWork(mockDB, test.updateWork)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error updating item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			assertWork(t, test.expectedWork, result)
		})
	}
}

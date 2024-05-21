package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
	"github.com/thomasmendez/personal-website-backend/api/tests"
)

func TestGetSkillsTools(t *testing.T) {

	mockDB := &mockDynamoDB{}

	for _, test := range []struct {
		label          string
		mockQueryFunc  func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
		expectedResult []models.SkillsTools
		expectedError  error
	}{
		{
			label: "valid query output",
			mockQueryFunc: func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
				mockOutput := &dynamodb.QueryOutput{
					Items: []map[string]*dynamodb.AttributeValue{
						tests.TestSkillsToolsItem,
					},
				}
				return mockOutput, nil
			},
			expectedResult: []models.SkillsTools{
				{
					PersonalWebsiteType: "SkillsTools",
					SortValue:           "Programming Languages",
					Category:            "Tools",
					Type:                "Programming Languages",
					List:                []string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"},
				},
			},
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

			result, err := GetSkillsTools(mockDB, "personalWebsiteTableDev")

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error querying database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, 1)

			for i, skillsTools := range result {
				tests.AssertSkillsTools(t, test.expectedResult[i], skillsTools)
			}
		})
	}
}

func TestPostSkillsTools(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label               string
		newSkillsTools      models.SkillsTools
		mockPutFunc         func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		mockGetFunc         func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedSkillsTools models.SkillsTools
		expectedError       error
	}{
		{
			label:          "valid query output",
			newSkillsTools: tests.TestSkillsTools,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: tests.TestSkillsToolsItem,
				}
				return mockOutput, nil
			},
			expectedSkillsTools: tests.TestSkillsTools,
			expectedError:       nil,
		},
		{
			label:          "query error",
			newSkillsTools: tests.TestSkillsTools,
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

			result, err := PostSkillsTools(mockDB, "personalWebsiteTableDev", test.newSkillsTools)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error getting item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			tests.AssertSkillsTools(t, test.expectedSkillsTools, result)
		})
	}
}

func TestUpdateSkillsTools(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label               string
		newSkillsTools      models.SkillsTools
		mockUpdateFunc      func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
		mockGetFunc         func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedSkillsTools models.SkillsTools
		expectedError       error
	}{
		{
			label:          "valid query output",
			newSkillsTools: tests.TestSkillsTools,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"personalWebsiteType": {S: aws.String("SkillsTools")},
						"sortValue":           {S: aws.String("Programming Languages")},
						"Category":            {S: aws.String("Tools")},
						"Type":                {S: aws.String("Programming Languages")},
						"List":                {SS: aws.StringSlice([]string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"})},
					},
				}
				return mockOutput, nil
			},
			expectedSkillsTools: tests.TestSkillsTools,
			expectedError:       nil,
		},
		{
			label:          "query error",
			newSkillsTools: tests.TestSkillsTools,
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

			result, err := UpdateSkillsTools(mockDB, "personalWebsiteTableDev", test.newSkillsTools)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error updating item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			tests.AssertSkillsTools(t, test.expectedSkillsTools, result)
		})
	}
}

func TestDeleteSkillsTools(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label               string
		deleteSkillsTools   models.SkillsTools
		mockDeleteFunc      func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
		mockGetFunc         func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedSkillsTools models.SkillsTools
		expectedError       error
	}{
		{
			label:             "valid query output",
			deleteSkillsTools: tests.TestSkillsTools,
			mockDeleteFunc: func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: tests.TestSkillsToolsItem,
				}
				return mockOutput, nil
			},
			expectedSkillsTools: tests.TestSkillsTools,
			expectedError:       nil,
		},
		{
			label:             "query error",
			deleteSkillsTools: tests.TestSkillsTools,
			mockDeleteFunc: func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("error deleting item from database")
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			mockDB.DeleteFunc = test.mockDeleteFunc
			mockDB.GetFunc = test.mockGetFunc

			err := DeleteItem(mockDB, "personalWebsiteTableDev", test.deleteSkillsTools.PersonalWebsiteType, test.deleteSkillsTools.SortValue)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, "error deleting item from database", err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

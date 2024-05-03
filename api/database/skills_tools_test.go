package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
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
						{
							"personalWebsiteType": {S: aws.String("SkillsTools")},
							"sortValue":           {S: aws.String("Programming Languages")},
							"skillsToolsCategory": {S: aws.String("Tools")},
							"skillsToolsType":     {S: aws.String("Programming Languages")},
							"skillsToolsList":     {SS: aws.StringSlice([]string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"})},
						},
					},
				}
				return mockOutput, nil
			},
			expectedResult: []models.SkillsTools{
				{
					PersonalWebsiteType: "SkillsTools",
					SortValue:           "Programming Languages",
					SkillsToolsCategory: "Tools",
					SkillsToolsType:     "Programming Languages",
					SkillsToolsList:     []string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"},
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

			result, err := GetSkillsTools(mockDB)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error querying database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, 1)

			for i, skillsToolsResult := range result {
				assert.Equal(t, test.expectedResult[i].PersonalWebsiteType, skillsToolsResult.PersonalWebsiteType)
				assert.Equal(t, test.expectedResult[i].SortValue, skillsToolsResult.SortValue)
				assert.Equal(t, test.expectedResult[i].SkillsToolsCategory, skillsToolsResult.SkillsToolsCategory)
				assert.Equal(t, test.expectedResult[i].SkillsToolsType, skillsToolsResult.SkillsToolsType)
				assert.Equal(t, test.expectedResult[i].SkillsToolsList, skillsToolsResult.SkillsToolsList)
			}
		})
	}
}

var expectedSkillsTools = models.SkillsTools{
	PersonalWebsiteType: "SkillsTools",
	SortValue:           "Programming Languages",
	SkillsToolsCategory: "Tools",
	SkillsToolsType:     "Programming Languages",
	SkillsToolsList:     []string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"},
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
			newSkillsTools: expectedSkillsTools,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"personalWebsiteType": {S: aws.String("SkillsTools")},
						"sortValue":           {S: aws.String("Programming Languages")},
						"skillsToolsCategory": {S: aws.String("Tools")},
						"skillsToolsType":     {S: aws.String("Programming Languages")},
						"skillsToolsList":     {SS: aws.StringSlice([]string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"})},
					},
				}
				return mockOutput, nil
			},
			expectedSkillsTools: expectedSkillsTools,
			expectedError:       nil,
		},
		{
			label:          "query error",
			newSkillsTools: expectedSkillsTools,
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

			result, err := PostSkillsTools(mockDB, test.newSkillsTools)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error getting item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedSkillsTools.PersonalWebsiteType, result.PersonalWebsiteType)
			assert.Equal(t, test.expectedSkillsTools.SortValue, result.SortValue)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsCategory, result.SkillsToolsCategory)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsType, result.SkillsToolsType)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsList, result.SkillsToolsList)
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
			newSkillsTools: expectedSkillsTools,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: map[string]*dynamodb.AttributeValue{
						"personalWebsiteType": {S: aws.String("SkillsTools")},
						"sortValue":           {S: aws.String("Programming Languages")},
						"skillsToolsCategory": {S: aws.String("Tools")},
						"skillsToolsType":     {S: aws.String("Programming Languages")},
						"skillsToolsList":     {SS: aws.StringSlice([]string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"})},
					},
				}
				return mockOutput, nil
			},
			expectedSkillsTools: expectedSkillsTools,
			expectedError:       nil,
		},
		{
			label:          "query error",
			newSkillsTools: expectedSkillsTools,
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

			result, err := UpdateSkillsTools(mockDB, test.newSkillsTools)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error updating item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedSkillsTools.PersonalWebsiteType, result.PersonalWebsiteType)
			assert.Equal(t, test.expectedSkillsTools.SortValue, result.SortValue)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsCategory, result.SkillsToolsCategory)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsType, result.SkillsToolsType)
			assert.Equal(t, test.expectedSkillsTools.SkillsToolsList, result.SkillsToolsList)
		})
	}
}

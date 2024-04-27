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
					PersonalWebsite:     "SkillsTools",
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
				assert.Equal(t, test.expectedResult[i].PersonalWebsite, skillsToolsResult.PersonalWebsite)
				assert.Equal(t, test.expectedResult[i].SortValue, skillsToolsResult.SortValue)
				assert.Equal(t, test.expectedResult[i].SkillsToolsCategory, skillsToolsResult.SkillsToolsCategory)
				assert.Equal(t, test.expectedResult[i].SkillsToolsType, skillsToolsResult.SkillsToolsType)
				assert.Equal(t, test.expectedResult[i].SkillsToolsList, skillsToolsResult.SkillsToolsList)
			}
		})
	}
}

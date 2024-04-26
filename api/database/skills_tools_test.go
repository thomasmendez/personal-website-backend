package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestGetSkillsTools(t *testing.T) {

	mockDB := &mockDynamoDB{}

	t.Run("Success case: valid query output", func(t *testing.T) {
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

		mockDB.QueryFunc = func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
			return mockOutput, nil
		}

		result, err := GetSkillsTools(mockDB)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "SkillsTools", result[0].PersonalWebsite)
		assert.Equal(t, "Programming Languages", result[0].SortValue)
	})

	t.Run("Error case: query error", func(t *testing.T) {
		mockDB.QueryFunc = func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("error querying database")
		}

		result, err := GetSkillsTools(mockDB)

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

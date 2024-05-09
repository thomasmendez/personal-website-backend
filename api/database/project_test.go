package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func TestGetProjects(t *testing.T) {

	mockDB := &mockDynamoDB{}

	for _, test := range []struct {
		label          string
		mockQueryFunc  func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error)
		expectedResult []models.Project
		expectedError  error
	}{
		{
			label: "valid query output",
			mockQueryFunc: func(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
				mockOutput := &dynamodb.QueryOutput{
					Items: []map[string]*dynamodb.AttributeValue{
						models.TestProjectItem,
					},
				}
				return mockOutput, nil
			},
			expectedResult: []models.Project{
				models.TestProject,
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

			result, err := GetProjects(mockDB)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error querying database", err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, 1)

			for i, project := range result {
				assertProject(t, test.expectedResult[i], project)
			}
		})
	}
}

func assertProject(t *testing.T, expectedProject models.Project, actualProject models.Project) {
	assert.Equal(t, expectedProject.PersonalWebsiteType, actualProject.PersonalWebsiteType)
	assert.Equal(t, expectedProject.SortValue, actualProject.SortValue)
	assert.Equal(t, expectedProject.Category, actualProject.Category)
	assert.Equal(t, expectedProject.Name, actualProject.Name)
	assert.Equal(t, expectedProject.Description, actualProject.Description)
	assert.Equal(t, expectedProject.FeaturesDescription, actualProject.FeaturesDescription)
	assert.Equal(t, expectedProject.Role, actualProject.Role)
	assert.Equal(t, expectedProject.Tasks, actualProject.Tasks)
	assert.Equal(t, expectedProject.TeamSize, actualProject.TeamSize)
	assert.Equal(t, expectedProject.TeamRoles, actualProject.TeamRoles)
	assert.Equal(t, expectedProject.CloudServices, actualProject.CloudServices)
	assert.Equal(t, expectedProject.Tools, actualProject.Tools)
	assert.Equal(t, expectedProject.Duration, actualProject.Duration)
	assert.Equal(t, expectedProject.StartDate, actualProject.StartDate)
	assert.Equal(t, expectedProject.EndDate, actualProject.EndDate)
	assert.Equal(t, expectedProject.Notes, actualProject.Notes)
	assert.Equal(t, expectedProject.Link, actualProject.Link)
	assert.Equal(t, expectedProject.LinkType, actualProject.LinkType)
	assert.Equal(t, expectedProject.MediaLink, actualProject.MediaLink)
}

func TestPostProject(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label           string
		newProject      models.Project
		mockPutFunc     func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		mockGetFunc     func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedProject models.Project
		expectedError   error
	}{
		{
			label:      "valid query output",
			newProject: models.TestProject,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: models.TestProjectItemNil,
				}
				return mockOutput, nil
			},
			expectedProject: models.TestProjectNil,
			expectedError:   nil,
		},
		{
			label:      "query error",
			newProject: models.TestProjectNil,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				return nil, errors.New("error inserting item")
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			mockDB.PutFunc = test.mockPutFunc
			mockDB.GetFunc = test.mockGetFunc

			result, err := PostProject(mockDB, test.newProject)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error inserting item", err.Error())
				return
			}

			assert.NoError(t, err)
			assertProject(t, test.expectedProject, result)
		})
	}
}

func TestUpdateProject(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label           string
		updateProject   models.Project
		mockUpdateFunc  func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
		mockGetFunc     func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedProject models.Project
		expectedError   error
	}{
		{
			label:         "valid query output",
			updateProject: models.TestProjectNil,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: models.TestProjectItemNil,
				}
				return mockOutput, nil
			},
			expectedProject: models.TestProjectNil,
			expectedError:   nil,
		},
		{
			label:         "query error",
			updateProject: models.TestProjectNil,
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

			result, err := UpdateProject(mockDB, test.updateProject)

			if err != nil {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, "error updating item from database", err.Error())
				return
			}

			assert.NoError(t, err)
			assertProject(t, test.expectedProject, result)
		})
	}
}

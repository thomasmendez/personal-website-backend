package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
	"github.com/thomasmendez/personal-website-backend/api/tests"
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
						tests.TestProjectItem,
					},
				}
				return mockOutput, nil
			},
			expectedResult: []models.Project{
				tests.TestProject,
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
				tests.AssertProject(t, test.expectedResult[i], project)
			}
		})
	}
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
			newProject: tests.TestProject,
			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: tests.TestProjectItemNil,
				}
				return mockOutput, nil
			},
			expectedProject: tests.TestProjectNil,
			expectedError:   nil,
		},
		{
			label:      "query error",
			newProject: tests.TestProjectNil,
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
			tests.AssertProject(t, test.expectedProject, result)
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
			updateProject: tests.TestProjectNil,
			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: tests.TestProjectItemNil,
				}
				return mockOutput, nil
			},
			expectedProject: tests.TestProjectNil,
			expectedError:   nil,
		},
		{
			label:         "query error",
			updateProject: tests.TestProjectNil,
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
			tests.AssertProject(t, test.expectedProject, result)
		})
	}
}

func TestDeleteProject(t *testing.T) {
	mockDB := &mockDynamoDB{}
	for _, test := range []struct {
		label           string
		deleteProject   models.Project
		mockDeleteFunc  func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
		mockGetFunc     func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
		expectedProject models.Project
		expectedError   error
	}{
		{
			label:         "valid query output",
			deleteProject: tests.TestProject,
			mockDeleteFunc: func(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
				return nil, nil
			},
			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
				mockOutput := &dynamodb.GetItemOutput{
					Item: tests.TestProjectItem,
				}
				return mockOutput, nil
			},
			expectedProject: tests.TestProject,
			expectedError:   nil,
		},
		{
			label:         "query error",
			deleteProject: tests.TestProject,
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

			err := DeleteItem(mockDB, test.deleteProject.PersonalWebsiteType, test.deleteProject.SortValue)

			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, "error deleting item from database", err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

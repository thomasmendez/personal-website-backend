package database

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

func TestGetProjects(t *testing.T) {

	mockDB := &mockDynamoDB{}
	teamSize := 1
	teamRoles := []string{"Frontend Developer", "Backend Developer"}
	cloudServices := []string{"AWS"}
	notes := "Site is still in development stages"
	link := "http://my-url"
	linkType := "YouTube"
	mediaLink := "http://link-to-media-file"

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
						{
							"personalWebsiteType": {S: aws.String("Projects")},
							"sortValue":           {S: aws.String("Project Title")},
							"category":            {S: aws.String("Software Engineering")},
							"name":                {S: aws.String("Personal Website")},
							"description":         {S: aws.String("My personal website created in the cloud")},
							"role":                {S: aws.String("Project Lead")},
							"tasks":               {SS: aws.StringSlice([]string{"Develop backend microservices"})},
							"teamSize":            {N: aws.String("1")},
							"teamRoles":           {SS: aws.StringSlice([]string{"Frontend Developer", "Backend Developer"})},
							"cloudServices":       {SS: aws.StringSlice([]string{"AWS"})},
							"tools":               {SS: aws.StringSlice([]string{"React", "Go"})},
							"duration":            {S: aws.String("6 Months")},
							"startDate":           {S: aws.String("Jan 2024")},
							"endDate":             {S: aws.String("Dec 2024")},
							"notes":               {S: aws.String("Site is still in development stages")},
							"link":                {S: aws.String("http://my-url")},
							"linkType":            {S: aws.String("YouTube")},
							"mediaLink":           {S: aws.String("http://link-to-media-file")},
						},
					},
				}
				return mockOutput, nil
			},
			expectedResult: []models.Project{
				{
					PersonalWebsiteType: "Projects",
					SortValue:           "Project Title",
					Category:            "Software Engineering",
					Name:                "Personal Website",
					Description:         "My personal website created in the cloud",
					Role:                "Project Lead",
					Tasks:               []string{"Develop backend microservices"},
					TeamSize:            &teamSize,
					TeamRoles:           &teamRoles,
					CloudServices:       &cloudServices,
					Tools:               []string{"React", "Go"},
					Duration:            "6 Months",
					StartDate:           "Jan 2024",
					EndDate:             "Dec 2024",
					Notes:               &notes,
					Link:                &link,
					LinkType:            &linkType,
					MediaLink:           &mediaLink,
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
				assert.Equal(t, test.expectedResult[i].PersonalWebsiteType, project.PersonalWebsiteType)
				assert.Equal(t, test.expectedResult[i].SortValue, project.SortValue)
				assert.Equal(t, test.expectedResult[i].Category, project.Category)
				assert.Equal(t, test.expectedResult[i].Name, project.Name)
				assert.Equal(t, test.expectedResult[i].Description, project.Description)
				assert.Equal(t, test.expectedResult[i].FeaturesDescription, project.FeaturesDescription)
				assert.Equal(t, test.expectedResult[i].Role, project.Role)
				assert.Equal(t, test.expectedResult[i].Tasks, project.Tasks)
				assert.Equal(t, test.expectedResult[i].TeamSize, project.TeamSize)
				assert.Equal(t, test.expectedResult[i].TeamRoles, project.TeamRoles)
				assert.Equal(t, test.expectedResult[i].CloudServices, project.CloudServices)
				assert.Equal(t, test.expectedResult[i].Tools, project.Tools)
				assert.Equal(t, test.expectedResult[i].Duration, project.Duration)
				assert.Equal(t, test.expectedResult[i].StartDate, project.StartDate)
				assert.Equal(t, test.expectedResult[i].EndDate, project.EndDate)
				assert.Equal(t, test.expectedResult[i].Notes, project.Notes)
				assert.Equal(t, test.expectedResult[i].Link, project.Link)
				assert.Equal(t, test.expectedResult[i].LinkType, project.LinkType)
				assert.Equal(t, test.expectedResult[i].MediaLink, project.MediaLink)
			}
		})
	}
}

// var expectedProject = models.Project{
// 	PersonalWebsiteType: "Projects",
// 	SortValue:           "Project Title",
// 	Category:            "Software Engineering",
// 	Name:                "Personal Website",
// 	Description:         "My personal website created in the cloud",
// 	Role:                "Project Lead",
// 	Tasks:               []string{"Develop backend microservices"},
// 	TeamSize:            nil,
// 	TeamRoles:           nil,
// 	CloudServices:       nil,
// 	Tools:               []string{"React", "Go"},
// 	Duration:            "6 Months",
// 	StartDate:           "Jan 2024",
// 	EndDate:             "Dec 2024",
// 	Notes:               nil,
// 	Link:                nil,
// 	MediaLink:           nil,
// }

// func TestPostProject(t *testing.T) {
// 	mockDB := &mockDynamoDB{}
// 	for _, test := range []struct {
// 		label           string
// 		newProject      models.Project
// 		mockPutFunc     func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
// 		mockGetFunc     func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
// 		expectedProject models.Project
// 		expectedError   error
// 	}{
// 		{
// 			label:      "valid query output",
// 			newProject: expectedProject,
// 			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
// 				return nil, nil
// 			},
// 			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
// 				mockOutput := &dynamodb.GetItemOutput{
// 					Item: map[string]*dynamodb.AttributeValue{
// 						"personalWebsiteType": {S: aws.String("Projects")},
// 						"sortValue":           {S: aws.String("Project Title")},
// 						"category":            {S: aws.String("Software Engineering")},
// 						"name":                {S: aws.String("Project Title")},
// 						"description":         {SS: aws.StringSlice([]string{"Develop backend microservices"})},
// 					},
// 				}
// 				return mockOutput, nil
// 			},
// 			expectedProject: expectedProject,
// 			expectedError:   nil,
// 		},
// 		{
// 			label:      "query error",
// 			newProject: expectedProject,
// 			mockPutFunc: func(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
// 				return nil, nil
// 			},
// 			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
// 				return nil, errors.New("error getting item from database")
// 			},
// 		},
// 	} {
// 		t.Run(test.label, func(t *testing.T) {
// 			mockDB.PutFunc = test.mockPutFunc
// 			mockDB.GetFunc = test.mockGetFunc

// 			result, err := PostProject(mockDB, test.newProject)

// 			if err != nil {
// 				assert.Error(t, err)
// 				assert.Empty(t, result)
// 				assert.Equal(t, "error getting item from database", err.Error())
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Equal(t, test.expectedProject.PersonalWebsiteType, result.PersonalWebsiteType)
// 			assert.Equal(t, test.expectedProject.SortValue, result.SortValue)
// 			assert.Equal(t, test.expectedProject.Category, result.Category)
// 			assert.Equal(t, test.expectedProject.Name, result.Name)
// 			assert.Equal(t, test.expectedProject.Description, result.Description)
// 			assert.Equal(t, test.expectedProject.FeaturesDescription, result.FeaturesDescription)
// 			assert.Equal(t, test.expectedProject.Role, result.Role)
// 			assert.Equal(t, test.expectedProject.Tasks, result.Tasks)
// 			assert.Equal(t, test.expectedProject.TeamSize, result.TeamSize)
// 			assert.Equal(t, test.expectedProject.TeamRoles, result.TeamRoles)
// 			assert.Equal(t, test.expectedProject.CloudServices, result.CloudServices)
// 			assert.Equal(t, test.expectedProject.Tools, result.Tools)
// 			assert.Equal(t, test.expectedProject.Duration, result.Duration)
// 			assert.Equal(t, test.expectedProject.StartDate, result.StartDate)
// 			assert.Equal(t, test.expectedProject.EndDate, result.EndDate)
// 			assert.Equal(t, test.expectedProject.Notes, result.Notes)
// 			assert.Equal(t, test.expectedProject.Link, result.Link)
// 			assert.Equal(t, test.expectedProject.LinkType, result.LinkType)
// 			assert.Equal(t, test.expectedProject.MediaLink, result.MediaLink)
// 		})
// 	}
// }

// func TestUpdateProject(t *testing.T) {
// 	mockDB := &mockDynamoDB{}
// 	for _, test := range []struct {
// 		label           string
// 		newProject      models.Project
// 		mockUpdateFunc  func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error)
// 		mockGetFunc     func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
// 		expectedProject models.Project
// 		expectedError   error
// 	}{
// 		{
// 			label:      "valid query output",
// 			newProject: expectedProject,
// 			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
// 				return nil, nil
// 			},
// 			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
// 				mockOutput := &dynamodb.GetItemOutput{
// 					Item: map[string]*dynamodb.AttributeValue{
// 						"personalWebsiteType": {S: aws.String("Projects")},
// 						"sortValue":           {S: aws.String("Project Title")},
// 						"category":            {S: aws.String("Software Engineering")},
// 						"name":                {S: aws.String("Project Title")},
// 						"description":         {SS: aws.StringSlice([]string{"Develop backend microservices"})},
// 					},
// 				}
// 				return mockOutput, nil
// 			},
// 			expectedProject: expectedProject,
// 			expectedError:   nil,
// 		},
// 		{
// 			label:      "query error",
// 			newProject: expectedProject,
// 			mockUpdateFunc: func(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
// 				return nil, nil
// 			},
// 			mockGetFunc: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
// 				return nil, errors.New("error updating item from database")
// 			},
// 		},
// 	} {
// 		t.Run(test.label, func(t *testing.T) {
// 			mockDB.UpdateFunc = test.mockUpdateFunc
// 			mockDB.GetFunc = test.mockGetFunc

// 			result, err := UpdateProject(mockDB, test.newProject)

// 			if err != nil {
// 				assert.Error(t, err)
// 				assert.Empty(t, result)
// 				assert.Equal(t, "error updating item from database", err.Error())
// 				return
// 			}

// 			assert.NoError(t, err)
// 			assert.Equal(t, test.expectedProject.PersonalWebsiteType, result.PersonalWebsiteType)
// 			assert.Equal(t, test.expectedProject.SortValue, result.SortValue)
// 			assert.Equal(t, test.expectedProject.Category, result.Category)
// 			assert.Equal(t, test.expectedProject.Name, result.Name)
// 			assert.Equal(t, test.expectedProject.Description, result.Description)
// 			assert.Equal(t, test.expectedProject.FeaturesDescription, result.FeaturesDescription)
// 			assert.Equal(t, test.expectedProject.Role, result.Role)
// 			assert.Equal(t, test.expectedProject.Tasks, result.Tasks)
// 			assert.Equal(t, test.expectedProject.TeamSize, result.TeamSize)
// 			assert.Equal(t, test.expectedProject.TeamRoles, result.TeamRoles)
// 			assert.Equal(t, test.expectedProject.CloudServices, result.CloudServices)
// 			assert.Equal(t, test.expectedProject.Tools, result.Tools)
// 			assert.Equal(t, test.expectedProject.Duration, result.Duration)
// 			assert.Equal(t, test.expectedProject.StartDate, result.StartDate)
// 			assert.Equal(t, test.expectedProject.EndDate, result.EndDate)
// 			assert.Equal(t, test.expectedProject.Notes, result.Notes)
// 			assert.Equal(t, test.expectedProject.Link, result.Link)
// 			assert.Equal(t, test.expectedProject.LinkType, result.LinkType)
// 			assert.Equal(t, test.expectedProject.MediaLink, result.MediaLink)
// 		})
// 	}
// }

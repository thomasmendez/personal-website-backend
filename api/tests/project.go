package tests

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

// Project model used for test cases
var TestProject = models.Project{
	PersonalWebsiteType: "Projects",
	SortValue:           "Project Title",
	Category:            "Software Engineering",
	Name:                "Personal Website",
	Description:         "My personal website created in the cloud",
	FeaturesDescription: "User is able to view my work",
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
}

var teamSize = "1"
var teamRoles = []string{"Frontend Developer", "Backend Developer"}
var cloudServices = []string{"AWS"}
var notes = "Site is still in development stages"
var link = "http://my-url"
var linkType = "YouTube"
var mediaLink = "http://link-to-media-file"

// Project Item model used for dynamodb
var TestProjectItem = map[string]*dynamodb.AttributeValue{
	"personalWebsiteType": {S: aws.String("Projects")},
	"sortValue":           {S: aws.String("Project Title")},
	"category":            {S: aws.String("Software Engineering")},
	"name":                {S: aws.String("Personal Website")},
	"description":         {S: aws.String("My personal website created in the cloud")},
	"featuresDescription": {S: aws.String("User is able to view my work")},
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
}

// Project model used for test cases with empty values
var TestProjectNil = models.Project{
	PersonalWebsiteType: "Projects",
	SortValue:           "Project Title",
	Category:            "Software Engineering",
	Name:                "Personal Website",
	Description:         "My personal website created in the cloud",
	FeaturesDescription: "User is able to view my work",
	Role:                "Project Lead",
	Tasks:               []string{"Develop backend microservices"},
	TeamSize:            nil,
	TeamRoles:           nil,
	CloudServices:       nil,
	Tools:               []string{"React", "Go"},
	Duration:            "6 Months",
	StartDate:           "Jan 2024",
	EndDate:             "Dec 2024",
	Notes:               nil,
	Link:                nil,
	LinkType:            nil,
	MediaLink:           nil,
}

// Project Item model used for dynamodb with empty values
var TestProjectItemNil = map[string]*dynamodb.AttributeValue{
	"personalWebsiteType": {S: aws.String("Projects")},
	"sortValue":           {S: aws.String("Project Title")},
	"category":            {S: aws.String("Software Engineering")},
	"name":                {S: aws.String("Personal Website")},
	"description":         {S: aws.String("My personal website created in the cloud")},
	"featuresDescription": {S: aws.String("User is able to view my work")},
	"role":                {S: aws.String("Project Lead")},
	"tasks":               {SS: aws.StringSlice([]string{"Develop backend microservices"})},
	"teamSize":            {NULL: aws.Bool(true)},
	"teamRoles":           {NULL: aws.Bool(true)},
	"cloudServices":       {NULL: aws.Bool(true)},
	"tools":               {SS: aws.StringSlice([]string{"React", "Go"})},
	"duration":            {S: aws.String("6 Months")},
	"startDate":           {S: aws.String("Jan 2024")},
	"endDate":             {S: aws.String("Dec 2024")},
	"notes":               {NULL: aws.Bool(true)},
	"link":                {NULL: aws.Bool(true)},
	"linkType":            {NULL: aws.Bool(true)},
	"mediaLink":           {NULL: aws.Bool(true)},
}

func AssertProject(t *testing.T, expectedProject models.Project, actualProject models.Project) {
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

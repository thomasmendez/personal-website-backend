package tests

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/thomasmendez/personal-website-backend/api/models"
)

// SkillsTools model used for test cases
var TestSkillsTools = models.SkillsTools{
	PersonalWebsiteType: "SkillsTools",
	SortValue:           "Programming Languages",
	Category:            "Tools",
	Type:                "Programming Languages",
	List:                []string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"},
}

// SkillsTools Item model used for dynamodb
var TestSkillsToolsItem = map[string]*dynamodb.AttributeValue{
	"personalWebsiteType": {S: aws.String("SkillsTools")},
	"sortValue":           {S: aws.String("Programming Languages")},
	"category":            {S: aws.String("Tools")},
	"type":                {S: aws.String("Programming Languages")},
	"list":                {SS: aws.StringSlice([]string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"})},
}

func AssertSkillsTools(t *testing.T, expectedSkillsTools models.SkillsTools, actualSkillsTools models.SkillsTools) {
	assert.Equal(t, expectedSkillsTools.PersonalWebsiteType, actualSkillsTools.PersonalWebsiteType)
	assert.Equal(t, expectedSkillsTools.SortValue, actualSkillsTools.SortValue)
	assert.Equal(t, expectedSkillsTools.Category, actualSkillsTools.Category)
	assert.Equal(t, expectedSkillsTools.Type, actualSkillsTools.Type)
	assert.Equal(t, expectedSkillsTools.List, actualSkillsTools.List)
}

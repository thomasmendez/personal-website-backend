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
	SkillsToolsCategory: "Tools",
	SkillsToolsType:     "Programming Languages",
	SkillsToolsList:     []string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"},
}

// SkillsTools Item model used for dynamodb
var TestSkillsToolsItem = map[string]*dynamodb.AttributeValue{
	"personalWebsiteType": {S: aws.String("SkillsTools")},
	"sortValue":           {S: aws.String("Programming Languages")},
	"skillsToolsCategory": {S: aws.String("Tools")},
	"skillsToolsType":     {S: aws.String("Programming Languages")},
	"skillsToolsList":     {SS: aws.StringSlice([]string{"Go", "Python", "JavaScript", "Java", "Swift", "C#"})},
}

func AssertSkillsTools(t *testing.T, expectedSkillsTools models.SkillsTools, actualSkillsTools models.SkillsTools) {
	assert.Equal(t, expectedSkillsTools.PersonalWebsiteType, actualSkillsTools.PersonalWebsiteType)
	assert.Equal(t, expectedSkillsTools.SortValue, actualSkillsTools.SortValue)
	assert.Equal(t, expectedSkillsTools.SkillsToolsCategory, actualSkillsTools.SkillsToolsCategory)
	assert.Equal(t, expectedSkillsTools.SkillsToolsType, actualSkillsTools.SkillsToolsType)
	assert.Equal(t, expectedSkillsTools.SkillsToolsList, actualSkillsTools.SkillsToolsList)
}

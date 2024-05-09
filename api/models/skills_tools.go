package models

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type SkillsTools struct {
	PersonalWebsiteType string   `json:"personalWebsiteType"`
	SortValue           string   `json:"sortValue"`
	SkillsToolsCategory string   `json:"skillsToolsCategory"`
	SkillsToolsType     string   `json:"skillsToolsType"`
	SkillsToolsList     []string `json:"skillsToolsList"`
}

// SkillsTools model used for test cases
var TestSkillsTools = SkillsTools{
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

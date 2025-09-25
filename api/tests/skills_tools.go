package tests

import (
	"reflect"
	"testing"

	"github.com/thomasmendez/personal-website-backend/api/models"
)

// SkillsTools model used for test cases
var TestSkillsTools = models.SkillsTools{
	PersonalWebsiteType: "SkillsTools",
	SortValue:           "Tools",
	Categories: []models.Category{
		{
			Category: "Programming Languages",
			List:     []string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"},
		},
		{
			Category: "Cloud Services",
			List:     []string{"AWS", "Azure", "Google Cloud Platform", "Digital Ocean"},
		},
	},
}

// SkillsTools Item model used for dynamodb
// var categories = []*dynamodb.AttributeValue{
// 	{
// 		M: map[string]*dynamodb.AttributeValue{
// 			"category": {S: aws.String("Programming Languages")},
// 			"list":     {SS: aws.StringSlice([]string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"})},
// 		},
// 	},
// 	{
// 		M: map[string]*dynamodb.AttributeValue{
// 			"category": {S: aws.String("Cloud Services")},
// 			"list":     {SS: aws.StringSlice([]string{"AWS", "Azure", "Google Cloud Platform", "Digital Ocean"})},
// 		},
// 	},
// }

// var TestSkillsToolsItem = map[string]*dynamodb.AttributeValue{
// 	"personalWebsiteType": {S: aws.String("SkillsTools")},
// 	"sortValue":           {S: aws.String("Tools")},
// 	"categories": {
// 		L: categories,
// 	},
// }

func AssertSkillsTools(t *testing.T, expectedSkillsTools models.SkillsTools, actualSkillsTools models.SkillsTools) {
	if !reflect.DeepEqual(expectedSkillsTools, actualSkillsTools) {
		t.Errorf("expected %v, got %v", expectedSkillsTools, actualSkillsTools)
	}
}

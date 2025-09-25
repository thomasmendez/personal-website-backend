package models

type SkillsTools struct {
	PersonalWebsiteType string     `json:"personalWebsiteType" dynamodbav:"personalWebsiteType"`
	SortValue           string     `json:"sortValue" dynamodbav:"sortValue"`
	Categories          []Category `json:"categories" dynamodbav:"categories"`
}

type Category struct {
	Category string   `json:"category" dynamodbav:"category"`
	List     []string `json:"list" dynamodbav:"list"`
}

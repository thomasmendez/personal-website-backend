package models

type Work struct {
	PersonalWebsiteType string   `json:"personalWebsiteType" dynamodbav:"personalWebsiteType"`
	SortValue           string   `json:"sortValue" dynamodbav:"sortValue"`
	JobTitle            string   `json:"jobTitle" dynamodbav:"jobTitle"`
	Company             string   `json:"company" dynamodbav:"company"`
	Location            Location `json:"location" dynamodbav:"location"`
	StartDate           string   `json:"startDate" dynamodbav:"startDate"`
	EndDate             string   `json:"endDate" dynamodbav:"endDate"`
	JobRole             string   `json:"jobRole" dynamodbav:"jobRole"`
	JobDescription      []string `json:"jobDescription" dynamodbav:"jobDescription"`
}

type Location struct {
	City  string `json:"city" dynamodbav:"city"`
	State string `json:"state" dynamodbav:"state"`
}

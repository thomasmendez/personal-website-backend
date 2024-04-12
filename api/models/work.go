package models

type Work struct {
	PersonalWebsiteType string   `json:"personalWebsiteType"`
	SortValue           string   `json:"sortValue"`
	JobTitle            string   `json:"jobTitle"`
	Company             string   `json:"company"`
	Location            Location `json:"location"`
	StartDate           string   `json:"startDate"`
	EndDate             string   `json:"endDate"`
	JobRole             string   `json:"jobRole"`
	JobDescription      []string `json:"JobDescription"`
}

type Location struct {
	City  string `json:"city"`
	State string `json:"state"`
}

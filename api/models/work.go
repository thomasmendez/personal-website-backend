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
	JobDescription      []string `json:"jobDescription"`
}

type Location struct {
	City  string `json:"city"`
	State string `json:"state"`
}

// Work model used for test cases
var TestWork = Work{
	PersonalWebsiteType: "Work",
	SortValue:           "2020-01-01",
	JobTitle:            "Software Engineer",
	Company:             "ABC Inc",
	Location: Location{
		City:  "New York",
		State: "NY",
	},
	StartDate:      "2020-01-01",
	EndDate:        "2020-12-31",
	JobRole:        "Backend Developer",
	JobDescription: []string{"Developed backend systems", "Optimized database queries"},
}

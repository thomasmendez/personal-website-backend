package models

type Project struct {
	PersonalWebsiteType string    `json:"personalWebsiteType"`
	SortValue           string    `json:"sortValue"`
	Category            string    `json:"category"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	FeaturesDescription string    `json:"featuresDescription"`
	Role                string    `json:"role"`
	Tasks               []string  `json:"tasks"`
	TeamSize            *string   `json:"teamSize"`
	TeamRoles           *[]string `json:"teamRoles"`
	CloudServices       *[]string `json:"cloudServices"`
	Tools               []string  `json:"tools"`
	Duration            string    `json:"duration"`
	StartDate           string    `json:"startDate"`
	EndDate             string    `json:"endDate"`
	Notes               *string   `json:"notes"`
	Link                *string   `json:"link"`
	LinkType            *string   `json:"linkType"`
	MediaLink           *string   `json:"mediaLink"`
}

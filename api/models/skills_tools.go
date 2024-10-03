package models

type SkillsTools struct {
	PersonalWebsiteType string     `json:"personalWebsiteType"`
	SortValue           string     `json:"sortValue"` // value should be "Skills" or "Tools"
	Categories          []Category `json:"categories"`
}

type Category struct {
	Category string   `json:"category"`
	List     []string `json:"list"`
}

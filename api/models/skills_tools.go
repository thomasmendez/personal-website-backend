package models

type SkillsTools struct {
	PersonalWebsiteType string   `json:"personalWebsiteType"`
	SortValue           string   `json:"sortValue"`
	SkillsToolsCategory string   `json:"skillsToolsCategory"`
	SkillsToolsType     string   `json:"skillsToolsType"`
	SkillsToolsList     []string `json:"skillsToolsList"`
}

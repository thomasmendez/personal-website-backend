package models

type SkillsTools struct {
	PersonalWebsite     string   `json:"PersonalWebsiteType"`
	SortValue           string   `json:"SortValue"`
	SkillsToolsCategory string   `json:"SkillsToolsCategory"`
	SkillsToolsType     string   `json:"SkillsToolsType"`
	SkillsToolsList     []string `json:"SkillsToolsList"`
}

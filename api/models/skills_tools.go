package models

type SkillsTools struct {
	PersonalWebsiteType string   `json:"personalWebsiteType"`
	SortValue           string   `json:"sortValue"`
	Category            string   `json:"category"` // value should be "Skills" or "Tools"
	Type                string   `json:"type"`
	List                []string `json:"list"`
}

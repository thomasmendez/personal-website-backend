package models

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

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

func (p *Project) MediaLinkIsS3Bucket() bool {
	if p.MediaLink == nil {
		return false
	}
	return strings.Contains(*p.MediaLink, "s3.amazonaws.com")
}

func (p *Project) GetFileNameFromMediaLink() (string, error) {
	if p.MediaLink == nil {
		return "", fmt.Errorf("no mediaLink found in Project")
	}
	parsedURL, err := url.Parse(*p.MediaLink)
	if err != nil {
		return "", fmt.Errorf("failed to parse url: %w", err)
	}

	filename := path.Base(parsedURL.Path)

	// Check if we actually got a filename (not empty or just "/")
	if filename == "." || filename == "/" || filename == "" {
		return "", fmt.Errorf("no filename found in MediaLink")
	}

	return filename, nil
}

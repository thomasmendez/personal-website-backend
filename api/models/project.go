package models

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type Project struct {
	PersonalWebsiteType string    `json:"personalWebsiteType" dynamodbav:"personalWebsiteType"`
	SortValue           string    `json:"sortValue" dynamodbav:"sortValue"`
	Category            string    `json:"category" dynamodbav:"category"`
	Name                string    `json:"name" dynamodbav:"name"`
	Description         string    `json:"description" dynamodbav:"description"`
	FeaturesDescription string    `json:"featuresDescription" dynamodbav:"featuresDescription"`
	Role                string    `json:"role" dynamodbav:"role"`
	Tasks               []string  `json:"tasks" dynamodbav:"tasks"`
	TeamSize            *string   `json:"teamSize" dynamodbav:"teamSize"`
	TeamRoles           *[]string `json:"teamRoles" dynamodbav:"teamRoles"`
	CloudServices       *[]string `json:"cloudServices" dynamodbav:"cloudServices"`
	Tools               []string  `json:"tools" dynamodbav:"tools"`
	Duration            string    `json:"duration" dynamodbav:"duration"`
	StartDate           string    `json:"startDate" dynamodbav:"startDate"`
	EndDate             string    `json:"endDate" dynamodbav:"endDate"`
	Notes               *string   `json:"notes" dynamodbav:"notes"`
	Link                *string   `json:"link" dynamodbav:"link"`
	LinkType            *string   `json:"linkType" dynamodbav:"linkType"`
	MediaLink           *string   `json:"mediaLink" dynamodbav:"mediaLink"`
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

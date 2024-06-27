package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/thomasmendez/personal-website-backend/api/models"
	"github.com/thomasmendez/personal-website-backend/api/tests"
)

func TestProjectApi(t *testing.T) {
	integrationTest(t)

	var latestProjectsResponse models.Project
	for _, test := range []struct {
		label           string
		route           string
		method          string
		reqBodyProjects func() *models.Project
		assertFunc      func(expectedStruct interface{}, resBody []byte)
	}{
		{
			label:  "Post Projects",
			route:  "/api/v1/projects",
			method: http.MethodPost,
			reqBodyProjects: func() *models.Project {
				return &tests.TestProject
			},
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var Projects models.Project
				err := json.Unmarshal(resBody, &Projects)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				tests.AssertProject(t, expectedStruct.(models.Project), Projects)
				latestProjectsResponse = Projects
			},
		},
		{
			label:  "Update Projects",
			route:  "/api/v1/projects",
			method: http.MethodPut,
			reqBodyProjects: func() *models.Project {
				teamSize := "1"
				teamRoles := []string{"Backend Developer", "Frontend Developer"}
				cloudServices := []string{"AWS", "Azure", "GCP"}
				notes := "Site is still in development stages"
				link := "http://my-url"
				linkType := "Youtube"
				mediaLink := "http://link-to-media-file"
				// modify previous response for update
				latestProjectsResponse.Category = "Software Engineering"
				latestProjectsResponse.Name = "Social Media Site"
				latestProjectsResponse.Description = "A social media site"
				latestProjectsResponse.FeaturesDescription = "A user is able to communicate with other users"
				latestProjectsResponse.Role = "Project Lead"
				latestProjectsResponse.Tasks = []string{"Develop backend", "Develop frontend"}
				latestProjectsResponse.TeamSize = &teamSize
				latestProjectsResponse.TeamRoles = &teamRoles
				latestProjectsResponse.CloudServices = &cloudServices
				latestProjectsResponse.Tools = []string{"Go", "React"}
				latestProjectsResponse.Duration = "6 Months"
				latestProjectsResponse.StartDate = "Jan 2024"
				latestProjectsResponse.EndDate = "Dec 2024"
				latestProjectsResponse.Notes = &notes
				latestProjectsResponse.Link = &link
				latestProjectsResponse.LinkType = &linkType
				latestProjectsResponse.MediaLink = &mediaLink
				return &latestProjectsResponse
			},
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var Projects models.Project
				err := json.Unmarshal(resBody, &Projects)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				tests.AssertProject(t, expectedStruct.(models.Project), Projects)
				latestProjectsResponse = Projects
			},
		},
		{
			label:           "Get Projects list",
			route:           "/api/v1/projects",
			method:          http.MethodGet,
			reqBodyProjects: nil,
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var Projects []models.Project
				err := json.Unmarshal(resBody, &Projects)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				for i, result := range Projects {
					tests.AssertProject(t, expectedStruct.([]models.Project)[i], result)
				}
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			// arrange
			httpClient := &http.Client{}
			url := "http://127.0.0.1:3000" + test.route

			var reqBodyProjects *models.Project
			if test.reqBodyProjects != nil {
				reqBodyProjects = test.reqBodyProjects()
			}
			reqBodyJson, err := json.Marshal(&reqBodyProjects)
			log.Print(string(reqBodyJson))
			if err != nil {
				t.Fatalf("failed to marshal Projects request: %v", err)
			}
			req, err := http.NewRequest(test.method, url, nil)
			if test.reqBodyProjects != nil {
				req, err = http.NewRequest(test.method, url, bytes.NewBuffer(reqBodyJson))
			}
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			// act
			res, err := httpClient.Do(req)
			if err != nil {
				t.Fatalf("failed to send request: %v", err)
			}
			defer res.Body.Close()
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("error in reading body: %v", err)
			}
			if res.StatusCode != 200 && res.StatusCode != 201 {
				t.Logf("Test request %v: %v", test.label, string(body))
				t.Fatalf("error status code: %v", res.StatusCode)
			}

			// check if it is a slice or array and then assert
			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				t.Fatalf("Error: %v", err)
				return
			}

			value := reflect.ValueOf(data)
			if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
				test.assertFunc([]models.Project{latestProjectsResponse}, body)
			} else {
				test.assertFunc(*reqBodyProjects, body)
			}
		})
	}
}

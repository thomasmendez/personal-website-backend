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

func TestSkillsToolsApi(t *testing.T) {
	// integrationTest(t)

	var latestSkillsToolsResponse models.SkillsTools
	for _, test := range []struct {
		label              string
		route              string
		method             string
		reqBodySkillsTools func() *models.SkillsTools
		assertFunc         func(expectedStruct interface{}, resBody []byte)
	}{
		{
			label:  "Post SkillsTools",
			route:  "/api/v1/skillsTools",
			method: http.MethodPost,
			reqBodySkillsTools: func() *models.SkillsTools {
				return &tests.TestSkillsTools
			},
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var skillsTools models.SkillsTools
				err := json.Unmarshal(resBody, &skillsTools)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				tests.AssertSkillsTools(t, expectedStruct.(models.SkillsTools), skillsTools)
				latestSkillsToolsResponse = skillsTools
			},
		},
		// {
		// 	label:  "Update SkillsTools",
		// 	route:  "/api/v1/skillsTools",
		// 	method: http.MethodPut,
		// 	reqBodySkillsTools: func() *models.SkillsTools {
		// 		// modify previous response for update
		// 		latestSkillsToolsResponse.Category = "Tools"
		// 		latestSkillsToolsResponse.Type = "Programming Languages"
		// 		latestSkillsToolsResponse.List = []string{"C#", "Go", "Java", "JavaScript", "Python", "Swift"}
		// 		return &latestSkillsToolsResponse
		// 	},
		// 	assertFunc: func(expectedStruct interface{}, resBody []byte) {
		// 		var skillsTools models.SkillsTools
		// 		err := json.Unmarshal(resBody, &skillsTools)
		// 		if err != nil {
		// 			t.Fatalf("error in unmarshal: %v", err)
		// 		}
		// 		tests.AssertSkillsTools(t, expectedStruct.(models.SkillsTools), skillsTools)
		// 		latestSkillsToolsResponse = skillsTools
		// 	},
		// },
		// {
		// 	label:              "Get SkillsTools list",
		// 	route:              "/api/v1/skillsTools",
		// 	method:             http.MethodGet,
		// 	reqBodySkillsTools: nil,
		// 	assertFunc: func(expectedStruct interface{}, resBody []byte) {
		// 		var skillsTools []models.SkillsTools
		// 		err := json.Unmarshal(resBody, &skillsTools)
		// 		if err != nil {
		// 			t.Fatalf("error in unmarshal: %v", err)
		// 		}
		// 		for i, result := range skillsTools {
		// 			tests.AssertSkillsTools(t, expectedStruct.([]models.SkillsTools)[i], result)
		// 		}
		// 	},
		// },
	} {
		t.Run(test.label, func(t *testing.T) {
			// arrange
			httpClient := &http.Client{}
			url := "http://127.0.0.1:3000" + test.route

			var reqBodySkillsTools *models.SkillsTools
			if test.reqBodySkillsTools != nil {
				reqBodySkillsTools = test.reqBodySkillsTools()
			}
			reqBodyJson, err := json.Marshal(&reqBodySkillsTools)
			log.Print(string(reqBodyJson))
			if err != nil {
				t.Fatalf("failed to marshal skillsTools request: %v", err)
			}
			req, err := http.NewRequest(test.method, url, nil)
			if test.reqBodySkillsTools != nil {
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
				t.Log(string(body))
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
				test.assertFunc([]models.SkillsTools{latestSkillsToolsResponse}, body)
			} else {
				test.assertFunc(*reqBodySkillsTools, body)
			}
		})
	}
}

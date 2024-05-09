package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/thomasmendez/personal-website-backend/api/models"
	"github.com/thomasmendez/personal-website-backend/api/tests"
)

func TestWorkApi(t *testing.T) {
	integrationTest(t)

	var latestWorkResponse models.Work
	for _, test := range []struct {
		label       string
		route       string
		method      string
		reqBodyWork func() *models.Work
		assertFunc  func(expectedStruct interface{}, resBody []byte)
	}{
		{
			label:  "Post Work",
			route:  "/api/v1/work",
			method: http.MethodPost,
			reqBodyWork: func() *models.Work {
				return &tests.TestWork
			},
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var work models.Work
				err := json.Unmarshal(resBody, &work)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				tests.AssertWork(t, expectedStruct.(models.Work), work)
				latestWorkResponse = work
			},
		},
		{
			label:  "Update Work",
			route:  "/api/v1/work",
			method: http.MethodPut,
			reqBodyWork: func() *models.Work {
				// modify previous response for update
				latestWorkResponse.JobTitle = "Senior Software Engineer"
				latestWorkResponse.Company = "New ABC Inc"
				latestWorkResponse.Location.City = "San Francisco"
				latestWorkResponse.Location.State = "CA"
				latestWorkResponse.StartDate = "2020-01-01"
				latestWorkResponse.EndDate = "2021-12-31"
				latestWorkResponse.JobRole = "Frontend Developer"
				latestWorkResponse.JobDescription = []string{"Created UI Themes", "Developed SPA Applications"}
				return &latestWorkResponse
			},
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var work models.Work
				err := json.Unmarshal(resBody, &work)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				tests.AssertWork(t, expectedStruct.(models.Work), work)
				latestWorkResponse = work
			},
		},
		{
			label:       "Get Work list",
			route:       "/api/v1/work",
			method:      http.MethodGet,
			reqBodyWork: nil,
			assertFunc: func(expectedStruct interface{}, resBody []byte) {
				var work []models.Work
				err := json.Unmarshal(resBody, &work)
				if err != nil {
					t.Fatalf("error in unmarshal: %v", err)
				}
				for i, result := range work {
					tests.AssertWork(t, expectedStruct.([]models.Work)[i], result)
				}
			},
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			// arrange
			httpClient := &http.Client{}
			url := "http://127.0.0.1:3000" + test.route
			var reqBodyWork *models.Work
			if test.reqBodyWork != nil {
				reqBodyWork = test.reqBodyWork()
			}
			reqBodyJson, err := json.Marshal(&reqBodyWork)
			if err != nil {
				t.Fatalf("failed to marshal work request: %v", err)
			}
			req, err := http.NewRequest(test.method, url, nil)
			if test.reqBodyWork != nil {
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

			// check if it is an array
			var data interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				t.Fatalf("Error: %v", err)
				return
			}

			// Get the reflect.Value of the unmarshalled data
			value := reflect.ValueOf(data)

			// Check if it's a slice or an array
			switch value.Kind() {
			case reflect.Slice:
				t.Log("It's a slice")
			case reflect.Array:
				t.Log("It's an array")
			default:
				t.Log("It's neither a slice nor an array")
			}

			// var work interface{}
			if value.Kind() == reflect.Slice || value.Kind() == reflect.Array {
				test.assertFunc([]models.Work{latestWorkResponse}, body)
			} else {
				test.assertFunc(*reqBodyWork, body)
			}
		})
	}
}

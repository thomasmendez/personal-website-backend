package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/thomasmendez/personal-website-backend/api/models"
	"github.com/thomasmendez/personal-website-backend/api/tests"
)

func TestWorkApi(t *testing.T) {
	integrationTest(t)

	var postWorkResponse models.Work
	for _, test := range []struct {
		label       string
		route       string
		method      string
		reqBodyWork func() *models.Work
	}{
		{
			label:  "post Work",
			route:  "/api/v1/work",
			method: http.MethodPost,
			reqBodyWork: func() *models.Work {
				return &tests.TestWork
			},
		},
		{
			label:  "update Work",
			route:  "/api/v1/work",
			method: http.MethodPut,
			reqBodyWork: func() *models.Work {
				postWorkResponse.JobTitle = "Senior Software Engineer"
				postWorkResponse.Company = "New ABC Inc"
				postWorkResponse.Location.City = "San Francisco"
				postWorkResponse.Location.State = "CA"
				postWorkResponse.StartDate = "2020-01-01"
				postWorkResponse.EndDate = "2021-12-31"
				postWorkResponse.JobRole = "Frontend Developer"
				postWorkResponse.JobDescription = []string{"Created UI Themes", "Developed SPA Applications"}
				return &postWorkResponse
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

			// assert
			var work models.Work
			err = json.Unmarshal(body, &work)
			if err != nil {
				t.Fatalf("error in unmarshal: %v", err)
			}
			if test.method == http.MethodPost {
				tests.AssertWork(t, *reqBodyWork, work)
			}
			if test.method == http.MethodPut {
				tests.AssertWork(t, *reqBodyWork, work)
			}
			if test.method == http.MethodPost {
				postWorkResponse = work
			}
		})
	}
}

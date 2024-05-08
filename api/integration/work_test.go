package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/thomasmendez/personal-website-backend/api/models"
)

func TestWorkApi(t *testing.T) {
	integrationTest(t)

	for _, test := range []struct {
		label       string
		route       string
		method      string
		reqBodyWork *models.Work
	}{
		{
			label:       "post Work",
			route:       "/api/v1/work",
			method:      http.MethodPost,
			reqBodyWork: &models.ExpectedWork,
		},
	} {
		t.Run(test.label, func(t *testing.T) {
			// arrange
			httpClient := &http.Client{}
			url := "http://127.0.0.1:3000" + test.route
			reqBodyJson, err := json.Marshal(&test.reqBodyWork)
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
			body, _ := io.ReadAll(res.Body)
			actualResponse := string(body)
			t.Log(actualResponse)

			// assert
		})
	}
}

func integrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, set environment variable INTEGRATION=1")
	}
}

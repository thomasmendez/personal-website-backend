package integration

import (
	"os"
	"testing"
)

func integrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, set environment variable INTEGRATION=1")
	}
}

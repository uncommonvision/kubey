package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealth(t *testing.T) {
	// Set Gin to test mode so it doesn't output to stdout.
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call the handler.
	Health(c)

	if w.Code != 200 {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp["status"] != "ok" {
		t.Fatalf("expected status 'ok', got %v", resp["status"])
	}

	// "time" should be a non-zero number (Unix timestamp).
	if tval, ok := resp["time"].(float64); !ok || tval <= 0 {
		t.Fatalf("expected a positive timestamp, got %v", resp["time"])
	}
}

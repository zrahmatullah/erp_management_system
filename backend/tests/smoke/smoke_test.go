package smoke_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func getBaseURL() string {
	baseURL := os.Getenv("API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return strings.TrimRight(baseURL, "/")
}

// TestSmoke_ServerHealth verifies that the server is alive and reporting healthy
func TestSmoke_ServerHealth(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(getBaseURL() + "/health")
	if err != nil {
		t.Fatalf("Smoke Test FAILED: Server is unreachable at %s: %v", getBaseURL(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", resp.StatusCode)
	}

	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if status, ok := body["status"].(string); !ok || status != "healthy" {
		t.Errorf("Expected health status 'healthy', got '%v'", body["status"])
	}
}

// TestSmoke_SecurityHeaders verifies critical OWASP security headers
func TestSmoke_SecurityHeaders(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(getBaseURL() + "/health")
	if err != nil {
		t.Fatalf("Smoke Test FAILED: %v", err)
	}
	defer resp.Body.Close()

	headers := []struct {
		name     string
		expected string
	}{
		{"X-Content-Type-Options", "nosniff"},
		{"X-Frame-Options", "DENY"},
		{"Referrer-Policy", "strict-origin-when-cross-origin"},
	}

	for _, h := range headers {
		val := resp.Header.Get(h.name)
		if val != h.expected {
			t.Errorf("Header %s = '%s', expected '%s'", h.name, val, h.expected)
		}
	}

	if csp := resp.Header.Get("Content-Security-Policy"); csp == "" {
		t.Errorf("Content-Security-Policy header is missing")
	}
}

// TestSmoke_AccessControlGuard verifies protected routes reject unauthenticated requests
func TestSmoke_AccessControlGuard(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(getBaseURL() + "/api/v1/master/units")
	if err != nil {
		t.Fatalf("Smoke Test FAILED: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized for unauthenticated request, got %d", resp.StatusCode)
	}
}

// TestSmoke_AuthSanity verifies authentication endpoint responds properly
func TestSmoke_AuthSanity(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}

	// 1. Malformed payload should yield 400 Bad Request
	badReq, _ := http.NewRequest(http.MethodPost, getBaseURL()+"/api/v1/auth/login", bytes.NewBufferString("{invalid_json"))
	badReq.Header.Set("Content-Type", "application/json")
	respBad, err := client.Do(badReq)
	if err != nil {
		t.Fatalf("Smoke Test FAILED: %v", err)
	}
	defer respBad.Body.Close()

	if respBad.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for malformed payload, got %d", respBad.StatusCode)
	}

	// 2. Valid credentials should authenticate successfully
	loginPayload := map[string]string{
		"email":    "admin@cafe-erp.com",
		"password": "Admin@123",
	}
	bodyBytes, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest(http.MethodPost, getBaseURL()+"/api/v1/auth/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Smoke Test FAILED: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK for valid admin login, got %d", resp.StatusCode)
	}

	var authData struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&authData); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if authData.Data.Token == "" {
		t.Errorf("Expected non-empty JWT token in login response")
	}
}


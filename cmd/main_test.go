package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	code := m.Run()
	os.Exit(code)
}

func TestPing(t *testing.T) {
	router := setupRouter()
	w := performRequest(router, "GET", "/ping", nil)

	assertStatusCode(t, w, http.StatusOK)
	assertJSONResponse(t, w, map[string]string{"message": "pong"})
}

func TestEnvironmentVariablesLoading(t *testing.T) {
	// Skip loading .env file in CI environment
	if os.Getenv("CI") != "true" {
		err := godotenv.Load()
		if err != nil {
			t.Logf("Error loading .env file: %v", err)
			t.Log("Setting mock values for environment variables")
			setMockEnvVars()
		} else {
			t.Log("Successfully loaded .env file")
		}
	} else {
		t.Log("Running in CI environment, setting mock values")
		setMockEnvVars()
	}

	requiredEnvVars := []string{
		"APP_ID",
		"APP_CERTIFICATE",
		"CUSTOMER_ID",
		"CUSTOMER_SECRET",
		"AGORA_BASE_URL",
		"AGORA_CLOUD_RECORDING_URL",
		"AGORA_RTT_URL",
		"STORAGE_VENDOR",
		"STORAGE_REGION",
		"STORAGE_BUCKET",
		"STORAGE_BUCKET_ACCESS_KEY",
		"STORAGE_BUCKET_SECRET_KEY",
	}

	for _, envVar := range requiredEnvVars {
		value, exists := os.LookupEnv(envVar)
		if !exists {
			t.Errorf("Required environment variable %s is not set", envVar)
		} else if value == "" {
			t.Errorf("Required environment variable %s is empty", envVar)
		}
	}

	// Test specific format for certain variables
	if appID := os.Getenv("APP_ID"); len(appID) != 32 {
		t.Errorf("APP_ID should be 32 characters long, got %d", len(appID))
	}

	if appCert := os.Getenv("APP_CERTIFICATE"); len(appCert) != 32 {
		t.Errorf("APP_CERTIFICATE should be 32 characters long, got %d", len(appCert))
	}

	// Test numeric values
	numericVars := map[string]string{
		"STORAGE_VENDOR": os.Getenv("STORAGE_VENDOR"),
		"STORAGE_REGION": os.Getenv("STORAGE_REGION"),
	}
	for varName, value := range numericVars {
		if _, err := strconv.Atoi(value); err != nil {
			t.Errorf("%s should be a numeric value, got %s", varName, value)
		}
	}

	// Test URL format
	if baseURL := os.Getenv("AGORA_BASE_URL"); !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		t.Errorf("AGORA_BASE_URL should start with http:// or https://, got %s", baseURL)
	}
}

func setMockEnvVars() {
	os.Setenv("APP_ID", "12345678901234567890123456789012")
	os.Setenv("APP_CERTIFICATE", "12345678901234567890123456789012")
	os.Setenv("CUSTOMER_ID", "mock_customer_id")
	os.Setenv("CUSTOMER_SECRET", "mock_customer_secret")
	os.Setenv("AGORA_BASE_URL", "https://api.agora.io/")
	os.Setenv("AGORA_CLOUD_RECORDING_URL", "mock_recording_url")
	os.Setenv("AGORA_RTT_URL", "mock_rtt_url")
	os.Setenv("STORAGE_VENDOR", "1")
	os.Setenv("STORAGE_REGION", "1")
	os.Setenv("STORAGE_BUCKET", "mock_bucket")
	os.Setenv("STORAGE_BUCKET_ACCESS_KEY", "mock_access_key")
	os.Setenv("STORAGE_BUCKET_SECRET_KEY", "mock_secret_key")
}

func TestGetBasicAuth(t *testing.T) {
	testCases := []struct {
		name           string
		customerID     string
		customerSecret string
		expected       string
	}{
		{"Valid Credentials", "testID", "testSecret", "Basic dGVzdElEOnRlc3RTZWNyZXQ="},
		{"Empty Credentials", "", "", "Basic Og=="},
		{"Special Characters", "test:ID", "test@Secret", "Basic dGVzdDpJRDp0ZXN0QFNlY3JldA=="},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getBasicAuth(tc.customerID, tc.customerSecret)
			if result != tc.expected {
				t.Errorf("getBasicAuth(%q, %q) = %q, want %q", tc.customerID, tc.customerSecret, result, tc.expected)
			}
		})
	}
}

func TestServerSetup(t *testing.T) {
	router := setupRouter()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("ListenAndServe() error = %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("Server shutdown error: %v", err)
	}
}

// Helper functions

func setupRouter() *gin.Engine {
	router := gin.New()
	router.GET("/ping", Ping)
	return router
}

func performRequest(r http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(body))
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func assertStatusCode(t *testing.T, res *httptest.ResponseRecorder, expected int) {
	t.Helper()
	if res.Code != expected {
		t.Errorf("Expected status %d; got %d", expected, res.Code)
	}
}

func assertJSONResponse(t *testing.T, res *httptest.ResponseRecorder, expected map[string]string) {
	t.Helper()
	var response map[string]string
	err := json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Unable to parse response body: %v", err)
	}
	for k, v := range expected {
		if response[k] != v {
			t.Errorf("Expected %s to be '%s', but got '%s'", k, v, response[k])
		}
	}
}

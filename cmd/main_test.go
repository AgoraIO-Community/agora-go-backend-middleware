package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/cloud_recording_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/http_headers"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/real_time_transcription_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/rtmp_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func TestMain(t *testing.T) {
	os.Clearenv()
	setMockEnvVars()

	// Use the mock setup
	router, err := setupRouter()
	if err != nil {
		t.Fatalf("Failed to setup router: %v", err)
	}

	// Create a test server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Send a request to the test server
	resp, err := http.Get(ts.URL + "/ping")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if message, exists := response["message"]; !exists || message != "pong" {
		t.Errorf("Expected message to be 'pong', got '%s'", message)
	}
}

func TestPing(t *testing.T) {
	router := setupPingTestRouter()
	w := performRequest(router, "GET", "/ping", nil)

	assertStatusCode(t, w, http.StatusOK)
	assertJSONResponse(t, w, map[string]string{"message": "pong"})
}

func TestEnvironmentVariablesLoading(t *testing.T) {
	os.Clearenv()
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
		"AGORA_RTMP_URL",
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
	os.Setenv("APP_ID", "a1b2c3d4e5f60718293a4b5c6d7e8f90")
	os.Setenv("APP_CERTIFICATE", "f9e8d7c6b5a40918273e6d5c4b3a2f1c")
	os.Setenv("CUSTOMER_ID", "1234567890abcdef1234567890abcdef")
	os.Setenv("CUSTOMER_SECRET", "abcdef1234567890abcdef1234567890")
	os.Setenv("AGORA_BASE_URL", "https://api.agora.io/")
	os.Setenv("AGORA_CLOUD_RECORDING_URL", "v1/apps/{appId}/cloud_recording")
	os.Setenv("AGORA_RTT_URL", "v1/projects/{appId}/rtsc/speech-to-text")
	os.Setenv("AGORA_RTMP_URL", "v1/projects/{{appId}}/rtmp-converters")
	os.Setenv("STORAGE_VENDOR", "1")
	os.Setenv("STORAGE_REGION", "1")
	os.Setenv("STORAGE_BUCKET", "agora_middleware_mock_bucket")
	os.Setenv("STORAGE_BUCKET_ACCESS_KEY", "AGORAIO1234567890AXYZ")
	os.Setenv("STORAGE_BUCKET_SECRET_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
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
	testCases := []struct {
		name       string
		envVars    map[string]string
		expected   int
		setupError bool
	}{
		{
			name:       "No Env Vars Loaded",
			envVars:    map[string]string{},
			expected:   http.StatusInternalServerError,
			setupError: true,
		},
		{
			name: "Only APP_ID and APP_CERT",
			envVars: map[string]string{
				"APP_ID":          "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"APP_CERTIFICATE": "f9e8d7c6b5a40918273e6d5c4b3a2f1c",
			},
			expected:   http.StatusOK,
			setupError: false,
		},
		{
			name: "AGORA_BASE_URL without other vars",
			envVars: map[string]string{
				"APP_ID":          "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"APP_CERTIFICATE": "f9e8d7c6b5a40918273e6d5c4b3a2f1c",
				"AGORA_BASE_URL":  "https://api.agora.io/",
			},
			expected:   http.StatusInternalServerError,
			setupError: true,
		},
		{
			name: "RTMP only, skipping Cloud Recording and RTT",
			envVars: map[string]string{
				"APP_ID":          "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"APP_CERTIFICATE": "f9e8d7c6b5a40918273e6d5c4b3a2f1c",
				"CUSTOMER_ID":     "1234567890abcdef1234567890abcdef",
				"CUSTOMER_SECRET": "abcdef1234567890abcdef1234567890",
				"AGORA_BASE_URL":  "https://api.agora.io/",
				"AGORA_RTMP_URL":  "v1/projects/{{appId}}/rtmp-converters",
			},
			expected:   http.StatusOK,
			setupError: false,
		},
		{
			name: "AGORA_BASE_URL with all required vars",
			envVars: map[string]string{
				"APP_ID":                    "a1b2c3d4e5f60718293a4b5c6d7e8f90",
				"APP_CERTIFICATE":           "f9e8d7c6b5a40918273e6d5c4b3a2f1c",
				"CUSTOMER_ID":               "1234567890abcdef1234567890abcdef",
				"CUSTOMER_SECRET":           "abcdef1234567890abcdef1234567890",
				"AGORA_BASE_URL":            "https://api.agora.io/",
				"AGORA_CLOUD_RECORDING_URL": "v1/apps/{appId}/cloud_recording",
				"AGORA_RTT_URL":             "v1/projects/{appId}/rtsc/speech-to-text",
				"AGORA_RTMP_URL":            "v1/projects/{{appId}}/rtmp-converters",
				"STORAGE_VENDOR":            "1",
				"STORAGE_REGION":            "1",
				"STORAGE_BUCKET":            "agora_middleware_mock_bucket",
				"STORAGE_BUCKET_ACCESS_KEY": "AGORAIO1234567890XXXX",
				"STORAGE_BUCKET_SECRET_KEY": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
			},
			expected:   http.StatusOK,
			setupError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			// Set test-specific environment variables
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			// Attempt to setup the router
			router, err := setupRouter()

			if tc.setupError {
				if err == nil {
					t.Errorf("Expected setup error, but got none")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected setup error: %v", err)
				}

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

				if resp.StatusCode != tc.expected {
					t.Errorf("Expected status %v; got %v", tc.expected, resp.Status)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					t.Errorf("Server shutdown error: %v", err)
				}
			}
		})
	}
}

func TestSetupRouterErrors(t *testing.T) {
	testCases := []struct {
		name          string
		envVars       map[string]string
		expectedError string
	}{
		{
			name:          "Missing APP_ID",
			envVars:       map[string]string{},
			expectedError: "FATAL ERROR: ENV not properly configured, APP ID and APP CERTIFICATE are required",
		},
		{
			name: "Missing APP_CERTIFICATE",
			envVars: map[string]string{
				"APP_ID": "testAppID",
			},
			expectedError: "FATAL ERROR: ENV not properly configured, APP ID and APP CERTIFICATE are required",
		},
		{
			name: "Missing Basic Auth with AGORA_BASE_URL",
			envVars: map[string]string{
				"APP_ID":          "testAppID",
				"APP_CERTIFICATE": "testAppCertificate",
				"AGORA_BASE_URL":  "https://api.agora.io/",
			},
			expectedError: "FATAL ERROR: ENV not properly configured for Basic Auth",
		},
		// TODO: Add more error cases
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear existing environment variables
			os.Clearenv()

			// Set test environment variables
			for k, v := range tc.envVars {
				os.Setenv(k, v)
			}

			// Setup router
			_, err := setupRouter()

			// Check for expected error
			if err == nil || !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectedError, err)
			}
		})
	}
}

func TestGracefulShutdown(t *testing.T) {
	os.Clearenv()
	setMockEnvVars()

	server := setupServer()

	// Create a channel to signal when the server has started
	started := make(chan bool)

	go func() {
		// Signal that the server is about to start
		started <- true
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server error: %v", err)
		}
	}()

	// Wait for the server to start
	<-started

	// Give the server a moment to fully initialize
	time.Sleep(100 * time.Millisecond)

	// Verify the server is running
	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Fatalf("Server is not running: %v", err)
	}
	resp.Body.Close()

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("Server Shutdown failed: %v", err)
	}

	// Wait a moment for the shutdown to complete
	time.Sleep(100 * time.Millisecond)

	// Try to make a request, it should fail
	client := &http.Client{
		Timeout: time.Second, // Set a short timeout
	}
	_, err = client.Get("http://localhost:8080/ping")
	if err == nil {
		t.Error("Expected error after server shutdown, got none")
	}
}

// Helper functions
func setupPingTestRouter() *gin.Engine {
	router := gin.New()
	router.GET("/ping", Ping)

	return router
}

func setupRouter() (*gin.Engine, error) {
	router := gin.New()
	router.GET("/ping", Ping)

	// Retrieve all configuration values from environment variables.
	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	customerIDEnv, customerIDExists := os.LookupEnv("CUSTOMER_ID")
	customerSecretEnv, customerSecretExists := os.LookupEnv("CUSTOMER_SECRET")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")
	baseURLEnv, baseURLExists := os.LookupEnv("AGORA_BASE_URL")
	cloudRecordingURLEnv, cloudRecordingURLExists := os.LookupEnv("AGORA_CLOUD_RECORDING_URL")
	realTimeTranscriptionURLEnv, realTimeTranscriptionURLExists := os.LookupEnv("AGORA_RTT_URL")
	rtmpURLEnv, rtmpURLExists := os.LookupEnv("AGORA_RTMP_URL")
	cloudPlayerURLEnv, cloudPlayerURLExists := os.LookupEnv("AGORA_CLOUD_PLAYER_URL")
	storageVendorEnv, vendorExists := os.LookupEnv("STORAGE_VENDOR")
	storageRegionEnv, regionExists := os.LookupEnv("STORAGE_REGION")
	storageBucketEnv, bucketExists := os.LookupEnv("STORAGE_BUCKET")
	storageAccessKeyEnv, accessKeyExists := os.LookupEnv("STORAGE_BUCKET_ACCESS_KEY")
	storageSecretKeyEnv, secretKeyExists := os.LookupEnv("STORAGE_BUCKET_SECRET_KEY")

	// log.Printf("appIDExists: %v, appIDEnv: %v", appIDExists, appIDEnv)
	// log.Printf("appCertExists: %v, appCertEnv: %v", appCertExists, appCertEnv)
	// log.Printf("customerIDExists: %v, customerIDEnv: %v", customerIDExists, customerIDEnv)
	// log.Printf("customerSecretExists: %v, customerSecretEnv: %v", customerSecretExists, customerSecretEnv)
	// log.Printf("baseURLExists: %v, baseURLEnv: %v", baseURLExists, baseURLEnv)
	// log.Printf("cloudRecordingURLExists: %v, cloudRecordingURLEnv: %v", cloudRecordingURLExists, cloudRecordingURLEnv)
	// log.Printf("realTimeTranscriptionURLExists: %v, realTimeTranscriptionURLEnv: %v", realTimeTranscriptionURLExists, realTimeTranscriptionURLEnv)
	// log.Printf("rtmpURLExists: %v, rtmpURLEnv: %v", rtmpURLExists, rtmpURLEnv)
	// log.Printf("vendorExists: %v, storageVendorEnv: %v", vendorExists, storageVendorEnv)
	// log.Printf("regionExists: %v, storageRegionEnv: %v", regionExists, storageRegionEnv)
	// log.Printf("bucketExists: %v, storageBucketEnv: %v", bucketExists, storageBucketEnv)
	// log.Printf("accessKeyExists: %v, storageAccessKeyEnv: %v", accessKeyExists, storageAccessKeyEnv)
	// log.Printf("secretKeyExists: %v, storageSecretKeyEnv: %v", secretKeyExists, storageSecretKeyEnv)

	// Check for the presence of core environment variables
	if !appIDExists || !appCertExists {
		return nil, fmt.Errorf("FATAL ERROR: ENV not properly configured, APP ID and APP CERTIFICATE are required.")
	}

	// Set up the Gin HTTP router with headers for CORS, caching, and timestamp.
	var httpHeaders = http_headers.NewHttpHeaders(corsAllowOrigin)
	router.Use(httpHeaders.NoCache())
	router.Use(httpHeaders.CORShttpHeaders())
	router.Use(httpHeaders.Timestamp())

	// Initialize services & register routes.
	tokenService := token_service.NewTokenService(appIDEnv, appCertEnv)
	tokenService.RegisterRoutes(router)

	if baseURLExists && baseURLEnv != "" {
		// Check for Basic Auth settings if baseURL is provided
		if !customerIDExists || !customerSecretExists {
			return nil, fmt.Errorf("FATAL ERROR: ENV not properly configured for Basic Auth, check .env file for all required variables")
		}
		// get basicAuth key
		basicAuthKey := getBasicAuth(customerIDEnv, customerSecretEnv)

		if cloudRecordingURLExists || realTimeTranscriptionURLExists {
			if !vendorExists || !regionExists || !bucketExists || !accessKeyExists || !secretKeyExists {
				return nil, fmt.Errorf("FATAL ERROR: ENV not properly configured for cloud storage, check .env file for all required variables")
			}

			storageVendorInt, storageVendorErr := strconv.Atoi(storageVendorEnv)
			if storageVendorErr != nil {
				return nil, fmt.Errorf("FATAL ERROR: Invalid STORAGE_VENDOR not properly configured")
			}

			storageRegionInt, storageRegionErr := strconv.Atoi(storageRegionEnv)
			if storageRegionErr != nil {
				return nil, fmt.Errorf("FATAL ERROR: Invalid STORAGE_REGION not properly configured")
			}
			// Configure storage settings based on environment variables.
			storageConfig := cloud_recording_service.StorageConfig{
				Vendor:    storageVendorInt,
				Region:    storageRegionInt,
				Bucket:    storageBucketEnv,
				AccessKey: storageAccessKeyEnv,
				SecretKey: storageSecretKeyEnv,
			}

			if cloudRecordingURLExists {
				// Init Cloud Recording Service
				cloudRecordingUrl := baseURLEnv + strings.Replace(cloudRecordingURLEnv, "{appId}", appIDEnv, 1) // replace the place-holder value with appID from Env
				cloudRecordingService := cloud_recording_service.NewCloudRecordingService(appIDEnv, cloudRecordingUrl, basicAuthKey, tokenService, storageConfig)
				cloudRecordingService.RegisterRoutes(router)
			}

			if realTimeTranscriptionURLExists {
				// Init Real Time Transcription Service
				realTimeTranscriptionUrl := baseURLEnv + strings.Replace(realTimeTranscriptionURLEnv, "{appId}", appIDEnv, 1) //replace the place-holder value with appID
				realTimeTranscriptionService := real_time_transcription_service.NewRTTService(appIDEnv, realTimeTranscriptionUrl, basicAuthKey, tokenService, storageConfig)
				realTimeTranscriptionService.RegisterRoutes(router)
			}
		}

		if rtmpURLExists || cloudPlayerURLExists {
			// support just rtmp or cloudplayer
			rtmpURL, cloudPlayerURL := "", ""
			// Check if rtmp and cloud player urls exist and replace the place-holder value with appID from Env
			if rtmpURLExists {
				rtmpURL = strings.Replace(rtmpURLEnv, "{appId}", appIDEnv, 1)
			}
			if cloudPlayerURLExists {
				cloudPlayerURL = strings.Replace(cloudPlayerURLEnv, "{appId}", appIDEnv, 1)
			}
			rtmpService := rtmp_service.NewRtmpService(appIDEnv, baseURLEnv, rtmpURL, cloudPlayerURL, basicAuthKey, tokenService)
			rtmpService.RegisterRoutes(router)
		}
	} else {
		log.Println("WARNING: baseURLEnv Not Found - SKIPPING Cloud Recording, RTT and RTMP services ")
	}

	return router, nil
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

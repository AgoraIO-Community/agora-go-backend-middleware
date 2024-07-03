package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// Setup code
	gin.SetMode(gin.TestMode)
	code := m.Run()
	// Teardown code
	os.Exit(code)
}

func TestPing(t *testing.T) {
	// Set Gin to Test Mode
	gin.SetMode(gin.TestMode)

	// Setup the router
	router := gin.Default()
	router.GET("/ping", Ping)

	// Create a request to send to the above route
	req, _ := http.NewRequest("GET", "/ping", nil)

	// Create a response recorder
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Check to see if the response was what you expected
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	expected := "pong"
	if response["message"] != expected {
		t.Fatalf("Expected message to be '%s' but it was '%s'", expected, response["message"])
	}
}

func TestGetBasicAuth(t *testing.T) {
	customerID := "testID"
	customerSecret := "testSecret"
	expected := "Basic dGVzdElEOnRlc3RTZWNyZXQ="

	result := getBasicAuth(customerID, customerSecret)
	if result != expected {
		t.Errorf("getBasicAuth(%q, %q) = %q, want %q", customerID, customerSecret, result, expected)
	}
}

func TestServerSetup(t *testing.T) {
	// Create a Gin router
	router := gin.New()
	router.GET("/ping", Ping) // Add the ping route

	// Test server configuration
	server := &http.Server{
		Addr:    ":8080",
		Handler: router, // Use the Gin router as the handler
	}

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("ListenAndServe() error = %v", err)
		}
	}()

	// Give the server a moment to start
	time.Sleep(100 * time.Millisecond)

	// Test if the server is running by making a request
	resp, err := http.Get("http://localhost:8080/ping")
	if err != nil {
		t.Fatalf("Could not send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Server shutdown error: %v", err)
	}
}

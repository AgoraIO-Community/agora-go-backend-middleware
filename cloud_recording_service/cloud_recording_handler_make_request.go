package cloud_recording_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// makeRequest is a helper function to make HTTP requests with basic authentication.
//
// Parameters:
//   - method: string - The HTTP method to use for the request (e.g., "GET", "POST").
//   - url: string - The URL to send the request to.
//   - body: optional interface{} - The request body to send (should be provided for methods like POST, PUT).
//
// Returns:
//   - []byte: The response body from the server.
//   - error: An error if there are any issues during the request.
//
// Behavior:
//   - Creates a new HTTP request with the specified method, URL, and body (as needed).
//   - Sets the Authorization and Content-Type headers.
//   - Sends the request using an HTTP client.
//   - Reads and returns the response body, or an error if the request fails.
func (s *CloudRecordingService) makeRequest(method, url string, body interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	if method == "GET" {
		// Create request without body for GET
		req, err = http.NewRequest(method, url, nil)
	} else if body != nil {
		// Marshal the request body into JSON
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %v", err)
		}

		// Create the HTTP request with a body for other methods
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set the Content-Type header for requests with a body
		req.Header.Set("Content-Type", "application/json")
	} else {
		return nil, fmt.Errorf("error creating request for method: %s - request body missing", method)
	}

	// Set the Authorization header
	req.Header.Set("Authorization", s.basicAuth)

	// Send the request with a 10-second timeout
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Check for a successful status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

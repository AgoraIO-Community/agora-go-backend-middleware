package real_time_transcription_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// makeRequest is a utility function that creates and sends HTTP requests with basic authentication.
// It is capable of handling different HTTP methods and supports optional request bodies.
//
// Parameters:
//   - method: string - The HTTP method (e.g., "GET", "POST") to be used for making the request.
//   - url: string - The endpoint URL to which the request is sent.
//   - body: interface{} (optional) - The payload for the request, required for methods like "POST" and "PUT".
//
// Returns:
//   - []byte: The raw response body received from the server.
//   - error: Non-nil error if there are issues during the request creation, sending, or processing.
//
// Behavior:
//   - Initializes an HTTP request with the provided method, URL, and body (if applicable).
//   - Sets necessary HTTP headers including Authorization for authentication and Content-Type for JSON content.
//   - Executes the request with a 10-second timeout using the http.Client.
//   - Validates the HTTP response status and reads the response body.
//   - Returns the response body or an error if the request was not successful.
func (s *RTTService) makeRequest(method, url string, body interface{}) ([]byte, error) {
	var req *http.Request
	var err error

	if method == "GET" || (method == "DELETE" && body == nil) {
		// Create a GET / DELETE request without a body.
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, fmt.Errorf("error making %s request: %v", method, err)
		}
	} else if body != nil {
		// Marshal the provided body into JSON for the request.
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body into JSON: %v", err)
		}

		// Create a request with a JSON body for non-GET methods.
		req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

	} else {
		// Return an error if a body is expected for non-GET requests but not provided.
		return nil, fmt.Errorf("request body missing for method %s", method)
	}

	// Set the 'Authorization' header for all requests.
	req.Header.Set("Authorization", s.basicAuth)

	// Set the 'Content-Type' header as it's required by all endpoints.
	req.Header.Set("Content-Type", "application/json")

	// Create and configure an HTTP client with a timeout.
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Check the HTTP response status code.
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

package real_time_transcription_service

import (
	"encoding/json"
	"fmt"
)

// HandleAcquireResourceReq constructs a URL, marshals the request payload, sends it to the Agora cloud recording API,
// and processes the response to acquire a resource for cloud recording.
//
// Parameters:
//   - acquireReq: AcquireResourceRequest - The structured data containing the details necessary for acquiring a resource.
//
// Returns:
//   - string: A unique identifier (resource ID) for the acquired resource from the Agora cloud recording API.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Converts the acquireReq object into JSON format for the API request.
//   - Constructs the URL for sending the acquisition request to the Agora cloud recording API.
//   - Utilizes makeRequest to perform the POST operation with the constructed URL and marshaled data.
//   - Interprets the API's JSON response to extract the resource ID if the operation succeeds.
//
// Notes:
//   - Assumes the availability of s.baseURL for constructing the request URL.
func (s *RTTService) HandleAcquireResourceReq(acquireReq AcquireBuilderTokenRequest) (json.RawMessage, error) {

	// Construct the URL for the POST request to acquire a cloud recording resource.
	url := fmt.Sprintf("%s/builderTokens", s.baseURL)

	// Send the POST request to the Agora cloud recording API.
	body, err := s.makeRequest("POST", url, acquireReq)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a struct to validate the response
	var response AcquireBuilderTokenResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	// Append a timestamp to the Agora response for auditing and record-keeping purposes.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

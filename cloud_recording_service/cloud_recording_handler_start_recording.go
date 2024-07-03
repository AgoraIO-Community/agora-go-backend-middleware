package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleStartRecordingReq initiates a cloud recording session using Agora's cloud recording service.
// It validates the request parameters, constructs the request URL, and sends the start recording request
// to the Agora API using the makeRequest helper function.
//
// Parameters:
//   - startReq: StartRecordingRequest - Contains the configuration settings for the recording session.
//   - resourceId: string - The resource ID previously acquired to identify the resource for the recording.
//   - modeType: string - Specifies the recording mode (e.g., individual, mix) to be used.
//
// Returns:
//   - json.RawMessage: The raw JSON response containing the recording ID (sid) from Agora.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for starting a new recording session based on the provided parameters.
//   - Sends a POST request with the start recording configuration to the Agora endpoint.
//   - Parses the JSON response to extract and validate the recording ID.
//   - Appends a timestamp to the response for record-keeping before returning the modified response.
//
// Notes:
//   - Assumes the presence of s.baseURL and s.appID for constructing the request URL.
//   - Utilizes s.makeRequest for sending the HTTP request and handling the response.
func (s *CloudRecordingService) HandleStartRecordingReq(startReq StartRecordingRequest, resourceId string, modeType string) (json.RawMessage, error) {
	// Construct the URL for the start recording endpoint.
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/mode/%s/start", s.baseURL, s.appID, resourceId, modeType)

	// Send a POST request to the start recording endpoint.
	body, err := s.makeRequest("POST", url, startReq)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a struct to validate the response
	var response StartRecordingResponse
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

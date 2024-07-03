package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleUpdateLayout processes the request to update the video layout during an ongoing cloud recording session.
// It constructs the request URL, validates the request data, and sends the update request to the Agora cloud recording API.
//
// Parameters:
//   - updateReq: UpdateLayoutRequest - The request payload containing the new layout settings.
//   - resourceId: string - The unique identifier for the resource (channel) that is being recorded.
//   - recordingId: string - The unique identifier for the ongoing recording session.
//   - modeType: string - Specifies the type of recording session, such as "individual" or "mixed".
//
// Returns:
//   - json.RawMessage: The raw JSON response from the cloud recording service, which includes details of the update operation.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for the update layout endpoint using the provided identifiers.
//   - Sends a POST request with the updated layout settings to the Agora cloud recording API.
//   - Parses the JSON response to validate its structure and confirm the successful application of the layout update.
//   - Appends a timestamp to the response for auditing purposes before returning the final response.
//
// Notes:
//   - Assumes the presence of s.baseURL to construct the request URL.
//   - The function uses s.makeRequest to handle the HTTP request and response handling efficiently.
func (s *CloudRecordingService) HandleUpdateLayout(updateReq UpdateLayoutRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error) {
	// Build the URL for the update layout endpoint.
	url := fmt.Sprintf("%s/resourceid/%s/sid/%s/mode/%s/updateLayout", s.baseURL, resourceId, recordingId, modeType)

	fmt.Println("HandleAcquireResourceReq with url: ", url)

	// Send a POST request to the update layout endpoint with the new settings.
	body, err := s.makeRequest("POST", url, updateReq)
	if err != nil {
		return nil, err
	}

	// Parse the response body to ensure it conforms to the expected structure.
	var response UpdateRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into UpdateRecordingResponse: %v", err)
	}

	// Append a timestamp to the Agora response for auditing and record-keeping.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error appending timestamp to response: %v", err)
	}

	return timestampBody, nil
}

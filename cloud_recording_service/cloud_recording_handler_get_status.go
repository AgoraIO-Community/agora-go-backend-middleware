package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// GetStatus handles the get status request.
// It constructs the URL and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Retrieves the resource ID, SID, and mode from the URL parameters.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) HandleGetStatus(resourceId string, recordingId string, modeType string) (json.RawMessage, error) {
	// build update recording endpoint
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/query", s.baseURL, s.appID, resourceId, recordingId, modeType)

	// send request to update recording endpoint
	body, err := s.makeRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	// Parse the response body to ensure it conforms to the expected structure
	var response ActiveRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing response body into StopRecordingResponse: %v", err)
	}

	// Validate the FileList conforms to the expected structure based on FileListMode.
	_, err = response.ServerResponse.UnmarshalFileList()
	if err != nil {
		return nil, fmt.Errorf("error parsing ServerResponse: %v", err)
	}

	// Add timestamp to Agora response
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

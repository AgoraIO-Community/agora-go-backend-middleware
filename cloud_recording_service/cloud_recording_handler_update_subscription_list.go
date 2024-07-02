package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// UpdateSubscriptionList handles the update subscription list request.
// It validates the request, constructs the URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a StartRecordingRequest struct.
//   - Validates the request fields.
//   - Constructs the URL and authentication header for the API request.
//   - Sends the request to the Agora cloud recording API and returns the response.
//
// Notes:
//   - This function assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing the API request.
func (s *CloudRecordingService) HandleUpdateSubscriptionList(updateReq UpdateRecordingRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error) {

	// build update recording endpoint
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/update", s.baseURL, s.appID, resourceId, recordingId, modeType)

	// send request to update recording endpoint
	body, err := s.makeRequest("POST", url, updateReq)
	if err != nil {
		return []byte{}, err
	}

	// Parse the response body to ensure it conforms to the expected structure
	var response UpdateRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing response body into StopRecordingResponse: %v", err)
	}

	// Add timestamp to Agora response
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

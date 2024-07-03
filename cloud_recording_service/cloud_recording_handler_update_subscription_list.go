package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleUpdateSubscriptionList processes the request to update the subscription list for a cloud recording session.
// It validates the provided request parameters, constructs the request URL, and sends the request to the Agora cloud recording API.
//
// Parameters:
//   - updateReq: UpdateSubscriptionRequest - The request payload containing the new subscription details.
//   - resourceId: string - The unique identifier for the resource (channel) that is being recorded.
//   - recordingId: string - The unique identifier for the ongoing recording session.
//   - modeType: string - Specifies the type of recording session, such as "individual" or "mixed".
//
// Returns:
//   - json.RawMessage: The raw JSON response from the cloud recording service, which includes details of the update operation.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for the update subscription endpoint using the provided identifiers.
//   - Sends a POST request with the updated subscription details to the Agora cloud recording API.
//   - Parses the JSON response to validate its structure and confirm the successful application of the subscription update.
//   - Appends a timestamp to the response for auditing purposes before returning the final response.
//
// Notes:
//   - Assumes the presence of s.baseURL, s.appID, s.customerID, and s.customerCertificate to construct the API request.
//   - Utilizes s.makeRequest to handle the HTTP request and response efficiently.
func (s *CloudRecordingService) HandleUpdateSubscriptionList(updateReq UpdateSubscriptionRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error) {
	// Construct the URL for the update subscription endpoint.
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/update", s.baseURL, s.appID, resourceId, recordingId, modeType)

	// Send a POST request to the update subscription endpoint with the new details.
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

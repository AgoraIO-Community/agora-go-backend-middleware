package cloud_recording_service

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
//   - Assumes the availability of s.baseURL, s.appID, s.customerID, and s.customerCertificate for constructing
//     the API request.
func (s *CloudRecordingService) HandleAcquireResourceReq(acquireReq AcquireResourceRequest) (string, error) {
	// Construct the URL for the POST request to acquire a cloud recording resource.
	url := fmt.Sprintf("%s/%s/cloud_recording/acquire", s.baseURL, s.appID)

	// Send the POST request to the Agora cloud recording API.
	body, err := s.makeRequest("POST", url, acquireReq)
	if err != nil {
		return "", err
	}

	// Parse the response body to extract the resource ID.
	var response struct {
		ResourceId string `json:"resourceId"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing response: %v", err)
	}

	return response.ResourceId, nil
}

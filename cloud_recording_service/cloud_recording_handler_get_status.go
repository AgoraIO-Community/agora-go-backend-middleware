package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleGetStatus constructs the URL and sends a GET request to the Agora cloud recording API
// to retrieve the status of a specific cloud recording session.
//
// Parameters:
//   - resourceId: string - Unique identifier for the resource in Agora Cloud Recording.
//   - recordingId: string - Session ID associated with the recording.
//   - modeType: string - Recording mode (e.g., individual, mix).
//
// Returns:
//   - json.RawMessage: JSON formatted response from the Agora cloud recording API.
//   - error: Error object if an issue occurs during the API call.
//
// Notes:
//   - Assumes availability of s.baseURL for constructing the request URL.
//   - Uses s.makeRequest to send the HTTP request and handles the response.
func (s *CloudRecordingService) HandleGetStatus(resourceId string, recordingId string, modeType string) (json.RawMessage, error) {

	// Construct the URL for the GET request to the cloud recording status endpoint.
	url := fmt.Sprintf("%s/resourceid/%s/sid/%s/mode/%s/query", s.baseURL, resourceId, recordingId, modeType)

	// Send the GET request to the Agora cloud recording API.
	body, err := s.makeRequest("GET", url, nil)
	if err != nil {
		return []byte{}, err
	}

	// Parse the response body to verify it conforms to the expected structure.
	var response ActiveRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []byte{}, fmt.Errorf("error parsing response body into ActiveRecordingResponse: %v", err)
	}

	// Validate the structure of the FileList based on the specified FileListMode.
	_, err = response.ServerResponse.UnmarshalFileList()
	if err != nil {
		return nil, fmt.Errorf("error parsing ServerResponse: %v", err)
	}

	// Append a timestamp to the response for auditing purposes.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error appending timestamp to response: %v", err)
	}

	return timestampBody, nil
}

package cloud_recording_service

import (
	"encoding/json"
	"fmt"
)

// HandleStopRecording processes the request to stop an ongoing cloud recording session in Agora's cloud recording service.
// It constructs the appropriate URL, validates the request parameters, and utilizes makeRequest to communicate with the Agora API.
//
// Parameters:
//   - stopReq: StopRecordingRequest - Object containing the necessary details to stop the recording.
//   - resourceId: string - The unique identifier for the resource (channel) that is being recorded.
//   - recordingId: string - The unique identifier for the ongoing recording session.
//   - modeType: string - Specifies the type of recording session, such as "individual" or "mixed".
//
// Returns:
//   - json.RawMessage: The raw JSON response from the cloud recording service, which includes the status of the stop request.
//   - error: Detailed error if the operation fails during any stage including validation, API request, or response parsing.
//
// Behavior:
//   - Constructs the URL for the stop recording endpoint using the provided identifiers.
//   - Sends a POST request to the Agora cloud recording API to stop the recording.
//   - Parses the received JSON response to validate its structure and confirm the successful termination of the recording.
//   - Adds a timestamp to the response for auditing purposes before returning the final response.
//
// Notes:
//   - The function assumes the availability of s.baseURL, s.appID, s.customerID, and s.customerCertificate to construct the API request.
//   - This function throws errors if any identifiers or request parameters are invalid or nil, ensuring robust error handling.
func (s *CloudRecordingService) HandleStopRecording(stopReq StopRecordingRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error) {
	// Construct the URL for the stop recording endpoint.
	url := fmt.Sprintf("%s/%s/cloud_recording/resourceid/%s/sid/%s/mode/%s/stop", s.baseURL, s.appID, resourceId, recordingId, modeType)

	// Send a POST request to the stop recording endpoint.
	body, err := s.makeRequest("POST", url, stopReq)
	if err != nil {
		return nil, err
	}

	// Parse the response body to validate its structure.
	var response ActiveRecordingResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body into StopRecordingResponse: %v", err)
	}

	// Validate the structure of the FileList from the response based on FileListMode.
	_, err = response.ServerResponse.UnmarshalFileList()
	if err != nil {
		return nil, fmt.Errorf("error validating ServerResponse: %v", err)
	}

	// Append a timestamp to the response for auditing and record-keeping.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error appending timestamp to response: %v", err)
	}

	return timestampBody, nil
}

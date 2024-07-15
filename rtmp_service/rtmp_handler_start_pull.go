package rtmp_service

import (
	"encoding/json"
	"fmt"
)

// HandleStartPullReq initiates an RTMP push request using Agora's Media Push service.
// It validates the request parameters, constructs the request URL, and sends the start recording request
// to the Agora API using the makeRequest helper function.
//
// Parameters:
//   - startReq: RtmpPushRequest - Contains the configuration settings for the RTMP push request.
//   - region: string - The region ID previously acquired to identify the resource for the recording.
//   - regionHintIp: *string - Optional parameter to provide a specific IP hint for the region.
//   - requestID: string - The unique request ID for tracing the request.
//
// Returns:
//   - json.RawMessage: The raw JSON response containing the recording ID (sid) from Agora.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for starting a new recording session based on the provided parameters.
//   - Appends the regionHintIp to the URL if provided and valid.
//   - Sends a POST request with the start recording configuration to the Agora endpoint.
//   - Parses the JSON response to extract and validate the recording ID.
//   - Appends a timestamp to the response for record-keeping before returning the modified response.
//
// Notes:
//   - Assumes the presence of s.baseURL & s.cloudPlayerURL for constructing the request URL.
//   - Utilizes s.makeRequest for sending the HTTP request and handling the response.
//   - Utilizes s.isValidIPv4 for validating the regionHintIp.
//   - Utilizes s.AddTimestamp to append a timestamp to the response.
func (s *RtmpService) HandleStartPullReq(startReq CloudPlayerStartRequest, region string, streamOriginIp *string, requestID string) (json.RawMessage, error) {
	// Construct the URL for the start recording endpoint.
	url := fmt.Sprintf("%s/%s/%s/players", s.baseURL, region, s.cloudPlayerURL)

	// Append regionHintIp if available and valid IPv4 address.
	if streamOriginIp != nil && s.isValidIPv4(*streamOriginIp) {
		url = fmt.Sprintf("%s?streamIp=%s", url, *streamOriginIp)
	}

	fmt.Println("HandleStartPullReq with url: ", url)

	// Send a POST request to the start recording endpoint.
	body, err := s.makeRequest("POST", url, startReq, requestID)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a struct to validate the response.
	var response StartCloudPlayerResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing cloud player response: %v", err)
	}

	// Append a timestamp to the Agora response for auditing and record-keeping purposes.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped cloud player response: %v", err)
	}

	return timestampBody, nil
}

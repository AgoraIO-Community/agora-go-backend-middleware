package rtmp_service

import (
	"encoding/json"
	"fmt"
)

// HandleStartPushReq initiates an RTMP push request using Agora's Media Push service.
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
//   - Assumes the presence of s.baseURL for constructing the request URL.
//   - Utilizes s.makeRequest for sending the HTTP request and handling the response.
//   - Utilizes s.isValidIPv4 for validating the regionHintIp.
func (s *RtmpService) HandleStartPushReq(startReq RtmpPushRequest, region string, regionHintIp *string, requestID string) (json.RawMessage, error) {
	// Construct the URL for the start recording endpoint.
	url := fmt.Sprintf("%s/%s/%s/rtmp-converters", s.baseURL, region, s.rtmpURL)

	// Append regionHintIp if available and valid IPv4 address.
	if regionHintIp != nil && s.isValidIPv4(*regionHintIp) {
		url = fmt.Sprintf("%s?regionHintIp=%s", url, *regionHintIp)
	}

	fmt.Println("HandleStartPushReq with url: ", url)

	// Send a POST request to the start recording endpoint.
	body, err := s.makeRequest("POST", url, startReq, requestID)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a struct to validate the response.
	var response StartRtmpResponse
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

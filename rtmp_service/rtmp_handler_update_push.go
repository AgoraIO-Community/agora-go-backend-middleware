package rtmp_service

import (
	"encoding/json"
	"fmt"
)

// HandleUpdatePushReq updates an existing RTMP push request using Agora's Media Push service.
// It constructs the request URL and sends the update request to the Agora API using the makeRequest helper function.
//
// Parameters:
//   - updateReq: RtmpPushRequest - Contains the configuration settings for the update RTMP request.
//   - converterId: string - The ID of the Converter returned in the start push request.
//   - region: string - The region ID for the rtmp resource.
//   - requestID: string - The unique request ID for tracing the request.
//
// Returns:
//   - json.RawMessage: The raw JSON response indicating the status of the update request.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for updating the rtmp session based on the provided parameters.
//   - Sends a PATCH request to the Agora endpoint to update the rtmp resource.
//   - Creates a success response for the client as the successful response won't have a body.
//   - Appends a timestamp to the response for record-keeping before returning the modified response.
//
// Notes:
//   - Assumes the presence of s.baseURL and s.rtmpURL for constructing the request URL.
//   - Utilizes s.makeRequest for sending the HTTP request and handling the response.
//   - Utilizes s.AddTimestamp to append a timestamp to the response.
func (s *RtmpService) HandleUpdatePushReq(updateReq RtmpPushRequest, converterId string, region string, requestID string) (json.RawMessage, error) {
	// Construct the URL for the update rtmp endpoint.
	url := fmt.Sprintf("%s/%s/%s/rtmp-converters/%s", s.baseURL, region, s.rtmpURL, converterId)

	fmt.Println("HandleUpdatePushReq with url: ", url)

	// Send a PATCH request to the update rtmp endpoint.
	body, err := s.makeRequest("PATCH", url, updateReq, requestID)
	if err != nil {
		return nil, err
	}

	// Parse the response body into a struct to validate the response.
	var response StartRtmpResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("error parsing cloud player response: %v", err)
	}

	// Append a timestamp to the response for auditing and record-keeping purposes.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

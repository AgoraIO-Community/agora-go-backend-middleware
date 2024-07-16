package rtmp_service

import (
	"encoding/json"
	"fmt"
)

// HandleStopPullReq stops an RTMP push request using Agora's Media Push service.
// It constructs the request URL and sends the stop request to the Agora API using the makeRequest helper function.
//
// Parameters:
//   - playerId: string - The ID of the Cloud Player returned in the start pull request.
//   - region: string - The region ID previously acquired to identify the resource for the recording.
//   - requestID: string - The unique request ID for tracing the request.
//
// Returns:
//   - json.RawMessage: The raw JSON response indicating the status of the stop request.
//   - error: Error object detailing any issues encountered during the API call.
//
// Behavior:
//   - Constructs the URL for stopping the cloud player session based on the provided parameters.
//   - Sends a DELETE request to the Agora endpoint to stop the cloud player service.
//   - Creates a success response for the client as the successful response won't have a body.
//   - Appends a timestamp to the response for record-keeping before returning the modified response.
//
// Notes:
//   - Assumes the presence of s.baseURL & s.cloudPlayerURL for constructing the request URL.
//   - Utilizes s.makeRequest for sending the HTTP request and handling the response.
//   - Utilizes s.AddTimestamp to append a timestamp to the response.
func (s *RtmpService) HandleStopPullReq(playerId string, region string, requestID string) (json.RawMessage, error) {
	// Construct the URL for the stop recording endpoint.
	url := fmt.Sprintf("%s%s/%s/players/%s", s.baseURL, region, s.cloudPlayerURL, playerId)

	fmt.Println("HandleStopPullReq with url: ", url)

	// Send a DELETE request to the stop recording endpoint.
	_, err := s.makeRequest("DELETE", url, nil, requestID)
	if err != nil {
		return nil, err
	}

	// Successful response won't have body so create a success response for client.
	response := CloudPlayerUpdateResponse{
		Status: "Success",
	}

	// Append a timestamp to the response for auditing and record-keeping purposes.
	timestampBody, err := s.AddTimestamp(&response)
	if err != nil {
		return nil, fmt.Errorf("error encoding timestamped response: %v", err)
	}

	return timestampBody, nil
}

package real_time_transcription_service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// AddTimestamp adds a current timestamp to any response object that supports the Timestampable interface.
// It then marshals the updated object back into JSON format for further use or storage.
func (s *RTTService) AddTimestamp(response Timestampable) (json.RawMessage, error) {
	// Set the current timestamp in UTC and RFC3339 format.
	now := time.Now().UTC().Format(time.RFC3339)
	response.SetTimestamp(now)

	// Marshal the response with the added timestamp back to JSON.
	timestampedBody, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("error marshaling final response with timestamp: %v", err)
	}
	return timestampedBody, nil
}

// generateUID generates a unique user identifier for use within cloud recording sessions.
// This function ensures the UID is never zero, which is reserved, by generating a random
// number between 1 and the maximum possible 32-bit integer value.
func (s *RTTService) generateUID() string {
	// Generate a random number starting from 1 to avoid 0, which is reserved.
	uid := rand.Intn(4294967294) + 1

	// Convert the integer UID to a string format and return it.
	return strconv.Itoa(uid)
}

package cloud_recording_service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// generateUID is a helper function to generate unique user id.
func generateUID() string {
	// Generate a random number between 1 and 2^32 - 1 (4294967295)
	// - starts from 1 to avoid 0, which is reserved.
	uid := rand.Intn(4294967294) + 1

	// Convert the integer to a string
	return strconv.Itoa(uid)
}

// Contains checks if a string is present in a slice.
func Contains(list []string, item string) bool {
	for _, a := range list {
		if a == item {
			return true
		}
	}
	return false
}

func (s *CloudRecordingService) AddTimestamp(response Timestampable) (json.RawMessage, error) {
	// Set the current timestamp
	now := time.Now().UTC().Format(time.RFC3339)
	response.SetTimestamp(now)

	// Marshal the response back to JSON
	timestampedBody, err := json.Marshal(response)
	if err != nil {
		return []byte{}, fmt.Errorf("error marshaling final response: %v", err)
	}
	return timestampedBody, nil
}

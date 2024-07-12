package rtmp_service

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// ValidateRegion checks if a specific string is present within a slice of strings.
// This is useful for determining if a particular item exists within a list.
func (s *RtmpService) ValidateRegion(regionToCheck string) bool {
	validRegions := []string{"na", "eu", "ap", "cn"}
	for _, region := range validRegions {
		if region == regionToCheck {
			return true
		}
	}
	return false
}

// AddTimestamp adds a current timestamp to any response object that supports the Timestampable interface.
// It then marshals the updated object back into JSON format for further use or storage.
func (s *RtmpService) AddTimestamp(response Timestampable) (json.RawMessage, error) {
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

// isValidIPv4 checks if a given string is a valid IPv4 address.
func (s *RtmpService) isValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

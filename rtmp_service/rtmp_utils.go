package rtmp_service

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
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

// generateUID generates a unique user identifier for use within cloud recording sessions.
// This function ensures the UID is never zero, which is reserved, by generating a random
// number between 1 and the maximum possible 32-bit integer value.
func (s *RtmpService) GenerateUID() string {
	// Generate a random number starting from 1 to avoid 0, which is reserved.
	uid := rand.Intn(4294967294) + 1

	// Convert the integer UID to a string format and return it.
	return strconv.Itoa(uid)
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

// checks if a given string is a valid IdleTimeout time.
func (s *RtmpService) ValidateIdleTimeOut(idleTimeOut *int) *int {
	isInRange, err := s.checkIntInRange(*idleTimeOut, 5, 600)
	if isInRange {
		return idleTimeOut
	}
	// default to 300
	if err != nil {
		log.Printf("warning: Using defautl IdleTimeOut,  error validating given value: %v", err)
	} else {
		log.Printf("warning: IdleTimeOut out of range using default: %v", err)
	}
	newIdleTimeOut := 300
	return &newIdleTimeOut
}

// checkIntInRange parses the input string to an int and checks if it's between min and max.
func (s *RtmpService) checkIntInRange(input int, min int, max int) (bool, error) {
	// Check if the int is between min and max
	if input >= min && input <= max {
		return true, nil
	}

	return false, nil
}

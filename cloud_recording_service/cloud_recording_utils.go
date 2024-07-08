package cloud_recording_service

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

// generateUID generates a unique user identifier for use within cloud recording sessions.
// This function ensures the UID is never zero, which is reserved, by generating a random
// number between 1 and the maximum possible 32-bit integer value.
func (s *CloudRecordingService) GenerateUID() string {
	// Generate a random number starting from 1 to avoid 0, which is reserved.
	uid := rand.Intn(4294967294) + 1

	// Convert the integer UID to a string format and return it.
	return strconv.Itoa(uid)
}

// ValidateRecordingMode checks if a specific string is present within a slice of strings.
// This is useful for determining if a particular item exists within a list.
func (s *CloudRecordingService) ValidateRecordingMode(valideModes []string, modeToCheck string) bool {
	for _, mode := range valideModes {
		if mode == modeToCheck {
			return true
		}
	}
	return false
}

// AddTimestamp adds a current timestamp to any response object that supports the Timestampable interface.
// It then marshals the updated object back into JSON format for further use or storage.
func (s *CloudRecordingService) AddTimestamp(response Timestampable) (json.RawMessage, error) {
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

// UnmarshalFileList interprets the file list from the server response, handling different formats based on the FileListMode.
// It supports 'string' and 'json' modes, returning the file list as either a slice of FileDetail or FileListEntry respectively.
func (sr *ServerResponse) UnmarshalFileList() (interface{}, error) {
	if sr.FileListMode == nil || sr.FileList == nil {
		// Ensure FileListMode and FileList are not nil before proceeding.
		return nil, fmt.Errorf("FileListMode or FileList are empty, cannot proceed with unmarshaling")
	}
	switch *sr.FileListMode {
	case "string":
		// Parse the file list as a slice of FileDetail structures.
		var fileList []FileDetail
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileDetail: %v", err)
		}
		return fileList, nil
	case "json":
		// Parse the file list as a slice of FileListEntry structures.
		var fileList []FileListEntry
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileListEntry: %v", err)
		}
		return fileList, nil
	default:
		// Handle unknown FileListMode by returning an error.
		return nil, fmt.Errorf("unknown FileListMode: %s", *sr.FileListMode)
	}
}

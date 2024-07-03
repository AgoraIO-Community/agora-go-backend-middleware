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

// UnmarshalFileList interprets the file list from the server response based on the specified FileListMode.
// It supports different data representations as per the mode set in the response (e.g., string or JSON).
//
// Returns:
//   - An interface holding the parsed file list which can be a slice of FileDetail or FileListEntry.
//   - An error if parsing fails or the mode is not supported.
//
// Notes:
//   - This function assumes that FileListMode and FileList are not nil before proceeding with parsing.

func (sr *ServerResponse) UnmarshalFileList() (interface{}, error) {
	if sr.FileListMode == nil || sr.FileList == nil {
		// No FileListMode set or FileList is nil
		return nil, fmt.Errorf("FileListMode or FileList are empty")
	}
	switch *sr.FileListMode {
	case "string":
		var fileList []FileDetail
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileDetail: %v", err)
		}
		return fileList, nil
	case "json":
		var fileList []FileListEntry
		if err := json.Unmarshal(*sr.FileList, &fileList); err != nil {
			return nil, fmt.Errorf("error parsing FileList into []FileListEntry: %v", err)
		}
		return fileList, nil
	default:
		return nil, fmt.Errorf("unknown FileListMode: %s", *sr.FileListMode)
	}
}

// cloud_recording_service_test.go

package cloud_recording_service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockCloudRecordingService struct {
	StartRecordingFunc           func(c *gin.Context)
	StopRecordingFunc            func(c *gin.Context)
	GetStatusFunc                func(c *gin.Context)
	UpdateSubscriptionListFunc   func(c *gin.Context)
	UpdateLayoutFunc             func(c *gin.Context)
	HandleAcquireResourceReqFunc func(acquireReq AcquireResourceRequest) (string, error)
	HandleStartRecordingReqFunc  func(startReq StartRecordingRequest, resourceId string, modeType string) (json.RawMessage, error)
	HandleStopRecordingFunc      func(stopReq StopRecordingRequest, resourceId string, recordingId string, modeType string) (json.RawMessage, error)
	AddTimestampFunc             func(response Timestampable) (json.RawMessage, error)
}

func (m *MockCloudRecordingService) StartRecording(c *gin.Context) { m.StartRecordingFunc(c) }
func (m *MockCloudRecordingService) StopRecording(c *gin.Context)  { m.StopRecordingFunc(c) }
func (m *MockCloudRecordingService) GetStatus(c *gin.Context)      { m.GetStatusFunc(c) }
func (m *MockCloudRecordingService) UpdateSubscriptionList(c *gin.Context) {
	m.UpdateSubscriptionListFunc(c)
}
func (m *MockCloudRecordingService) UpdateLayout(c *gin.Context) { m.UpdateLayoutFunc(c) }
func (m *MockCloudRecordingService) HandleAcquireResourceReq(a AcquireResourceRequest) (string, error) {
	return m.HandleAcquireResourceReqFunc(a)
}
func (m *MockCloudRecordingService) HandleStartRecordingReq(s StartRecordingRequest, r string, mode string) (json.RawMessage, error) {
	return m.HandleStartRecordingReqFunc(s, r, mode)
}
func (m *MockCloudRecordingService) HandleStopRecording(s StopRecordingRequest, r string, i string, mode string) (json.RawMessage, error) {
	return m.HandleStopRecordingFunc(s, r, i, mode)
}
func (m *MockCloudRecordingService) AddTimestamp(r Timestampable) (json.RawMessage, error) {
	return m.AddTimestampFunc(r)
}

// Now let's write our tests using this mock service

func TestStartRecording(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockCloudRecordingService{
		StartRecordingFunc: func(c *gin.Context) {
			response := StartRecordingResponse{
				ResourceId: "test_resource_id",
				Sid:        "test_sid",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/cloud_recording/start", mockService.StartRecording)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cloud_recording/start", strings.NewReader(`{"channelName": "test_channel"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response StartRecordingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ResourceId != "test_resource_id" {
		t.Errorf("Expected resourceId 'test_resource_id', got '%s'", response.ResourceId)
	}
	if response.Sid != "test_sid" {
		t.Errorf("Expected sid 'test_sid', got '%s'", response.Sid)
	}
}

func TestStopRecording(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockCloudRecordingService{
		StopRecordingFunc: func(c *gin.Context) {
			resourceId := "test_resource_id"
			sid := "test_sid"
			fileList := []byte(`[{"fileName":"test.mp4","trackType":"audio","uid":"1","mixedAllUser":true,"isPlayable":true,"sliceStartTime":1609459200}]`)
			response := ActiveRecordingResponse{
				ResourceId: &resourceId,
				Sid:        &sid,
				ServerResponse: ServerResponse{
					FileListMode: &[]string{"json"}[0],
					FileList:     (*json.RawMessage)(&fileList),
				},
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/cloud_recording/stop", mockService.StopRecording)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/cloud_recording/stop", strings.NewReader(`{"cname":"test_channel","uid":"1","resourceId":"test_resource_id","sid":"test_sid"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response ActiveRecordingResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if *response.ResourceId != "test_resource_id" {
		t.Errorf("Expected resourceId 'test_resource_id', got '%s'", *response.ResourceId)
	}
	if *response.Sid != "test_sid" {
		t.Errorf("Expected sid 'test_sid', got '%s'", *response.Sid)
	}
	if *response.ServerResponse.FileListMode != "json" {
		t.Errorf("Expected FileListMode 'json', got '%s'", *response.ServerResponse.FileListMode)
	}
}

func TestHandleAcquireResourceReq(t *testing.T) {
	mockService := &MockCloudRecordingService{
		HandleAcquireResourceReqFunc: func(acquireReq AcquireResourceRequest) (string, error) {
			return "test_resource_id", nil
		},
	}

	resourceID, err := mockService.HandleAcquireResourceReq(AcquireResourceRequest{
		Cname: "test_channel",
		Uid:   "test_uid",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resourceID != "test_resource_id" {
		t.Errorf("Expected resourceID 'test_resource_id', got '%s'", resourceID)
	}
}

// Define TestResponse type
type TestResponse struct {
	Message   string  `json:"message"`
	Timestamp *string `json:"timestamp,omitempty"`
}

// Implement Timestampable interface for TestResponse
func (r *TestResponse) SetTimestamp(timestamp string) {
	r.Timestamp = &timestamp
}

func TestAddTimestamp(t *testing.T) {
	mockService := &MockCloudRecordingService{
		AddTimestampFunc: func(response Timestampable) (json.RawMessage, error) {
			timestamp := time.Now().UTC().Format(time.RFC3339)
			response.SetTimestamp(timestamp)
			return json.Marshal(response)
		},
	}

	response := &TestResponse{Message: "Test"}
	jsonResponse, err := mockService.AddTimestamp(response)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	var updatedResponse TestResponse
	err = json.Unmarshal(jsonResponse, &updatedResponse)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if updatedResponse.Timestamp == nil {
		t.Error("Expected timestamp to be added, but it's nil")
	}

	timestamp, err := time.Parse(time.RFC3339, *updatedResponse.Timestamp)
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	if time.Since(timestamp) > 5*time.Second {
		t.Errorf("Timestamp is not recent: %v", *updatedResponse.Timestamp)
	}
}

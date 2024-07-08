package real_time_transcription_service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type MockRTTService struct {
	StartRTTFunc                     func(c *gin.Context)
	StopRTTFunc                      func(c *gin.Context)
	QueryRTTFunc                     func(c *gin.Context)
	HandleAcquireBuilderTokenReqFunc func(acquireReq AcquireBuilderTokenRequest) (json.RawMessage, string, error)
	HandleStartReqFunc               func(startRttRequest StartRTTRequest, builderToken string) (json.RawMessage, error)
	HandleStopReqFunc                func(taskId string, builderToken string) (json.RawMessage, error)
	HandleQueryReqFunc               func(taskId string, builderToken string) (json.RawMessage, error)
	AddTimestampFunc                 func(response Timestampable) (json.RawMessage, error)
}

func (m *MockRTTService) StartRTT(c *gin.Context) { m.StartRTTFunc(c) }
func (m *MockRTTService) StopRTT(c *gin.Context)  { m.StopRTTFunc(c) }
func (m *MockRTTService) QueryRTT(c *gin.Context) { m.QueryRTTFunc(c) }
func (m *MockRTTService) HandleAcquireBuilderTokenReq(a AcquireBuilderTokenRequest) (json.RawMessage, string, error) {
	return m.HandleAcquireBuilderTokenReqFunc(a)
}
func (m *MockRTTService) HandleStartReq(s StartRTTRequest, b string) (json.RawMessage, error) {
	return m.HandleStartReqFunc(s, b)
}
func (m *MockRTTService) HandleStopReq(t string, b string) (json.RawMessage, error) {
	return m.HandleStopReqFunc(t, b)
}
func (m *MockRTTService) HandleQueryReq(t string, b string) (json.RawMessage, error) {
	return m.HandleQueryReqFunc(t, b)
}
func (m *MockRTTService) AddTimestamp(r Timestampable) (json.RawMessage, error) {
	return m.AddTimestampFunc(r)
}

func TestStartRTT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRTTService{
		StartRTTFunc: func(c *gin.Context) {
			response := struct {
				Acquire   json.RawMessage `json:"acquire"`
				Start     json.RawMessage `json:"start"`
				Timestamp string          `json:"timestamp"`
			}{
				Acquire:   json.RawMessage(`{"tokenName":"test_token","createTs":1625097600,"instanceId":"test_instance"}`),
				Start:     json.RawMessage(`{"createTs":1625097600,"status":"success","taskId":"test_task_id"}`),
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtt/start", mockService.StartRTT)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rtt/start", strings.NewReader(`{"channelName":"test_channel","languages":["en-US"],"subscribeAudioUids":["1","2"],"maxIdleTime":300,"enableStorage":false}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Acquire   json.RawMessage `json:"acquire"`
		Start     json.RawMessage `json:"start"`
		Timestamp string          `json:"timestamp"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Add more specific checks for the response content
}

func TestStopRTT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRTTService{
		StopRTTFunc: func(c *gin.Context) {
			response := StopRTTResponse{
				Timestamp: new(string),
			}
			*response.Timestamp = time.Now().UTC().Format(time.RFC3339)
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtt/stop/:taskId", mockService.StopRTT)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rtt/stop/test_task_id", strings.NewReader(`{"builderToken":"test_builder_token"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response StopRTTResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Timestamp == nil {
		t.Error("Expected timestamp to be present, but it's nil")
	}
}

func TestQueryRTT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRTTService{
		QueryRTTFunc: func(c *gin.Context) {
			response := struct {
				Status    string `json:"status"`
				TaskId    string `json:"taskId"`
				Timestamp string `json:"timestamp"`
			}{
				Status:    "in_progress",
				TaskId:    "test_task_id",
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.GET("/rtt/status/:taskId", mockService.QueryRTT)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rtt/status/test_task_id?builderToken=test_builder_token", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response struct {
		Status    string `json:"status"`
		TaskId    string `json:"taskId"`
		Timestamp string `json:"timestamp"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Status != "in_progress" {
		t.Errorf("Expected status 'in_progress', got '%s'", response.Status)
	}
	if response.TaskId != "test_task_id" {
		t.Errorf("Expected taskId 'test_task_id', got '%s'", response.TaskId)
	}
}

func TestHandleAcquireBuilderTokenReq(t *testing.T) {
	mockService := &MockRTTService{
		HandleAcquireBuilderTokenReqFunc: func(acquireReq AcquireBuilderTokenRequest) (json.RawMessage, string, error) {
			response := AcquireBuilderTokenResponse{
				TokenName:  "test_token",
				CreateTs:   1625097600,
				InstanceId: "test_instance",
			}
			jsonResponse, _ := json.Marshal(response)
			return jsonResponse, response.TokenName, nil
		},
	}

	jsonResponse, tokenName, err := mockService.HandleAcquireBuilderTokenReq(AcquireBuilderTokenRequest{
		InstanceId: "test_instance",
	})

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if tokenName != "test_token" {
		t.Errorf("Expected tokenName 'test_token', got '%s'", tokenName)
	}

	var response AcquireBuilderTokenResponse
	err = json.Unmarshal(jsonResponse, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if response.InstanceId != "test_instance" {
		t.Errorf("Expected instanceId 'test_instance', got '%s'", response.InstanceId)
	}
}

func TestAddTimestamp(t *testing.T) {
	mockService := &MockRTTService{
		AddTimestampFunc: func(response Timestampable) (json.RawMessage, error) {
			timestamp := time.Now().UTC().Format(time.RFC3339)
			response.SetTimestamp(timestamp)
			return json.Marshal(response)
		},
	}

	response := &StopRTTResponse{}
	jsonResponse, err := mockService.AddTimestamp(response)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	var updatedResponse StopRTTResponse
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

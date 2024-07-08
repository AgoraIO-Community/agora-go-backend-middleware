package real_time_transcription_service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
			var req ClientStartRTTRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			if req.ChannelName == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "channelName is required"})
				return
			}
			// Happy path response
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

	// Test case 1: Valid request
	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/start", strings.NewReader(`{"channelName":"test_channel","languages":["en-US"],"subscribeAudioUids":["1","2"],"maxIdleTime":300,"enableStorage":false}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Add more assertions for the response body
	})

	// Test case 2: Missing required field
	t.Run("Missing Required Field", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/start", strings.NewReader(`{"languages":["en-US"],"subscribeAudioUids":["1","2"]}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "channelName is required")
	})

	// Test case 3: Malformed JSON
	t.Run("Malformed JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/start", strings.NewReader(`{"channelName":test_channel`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "invalid character")
	})
}

func TestStopRTT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRTTService{
		StopRTTFunc: func(c *gin.Context) {
			taskId := c.Param("taskId")
			if taskId == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "taskId is required"})
				return
			}
			var req struct {
				BuilderToken string `json:"builderToken" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			// Happy path response
			response := StopRTTResponse{
				Timestamp: new(string),
			}
			*response.Timestamp = time.Now().UTC().Format(time.RFC3339)
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtt/stop/:taskId", mockService.StopRTT)

	// Test case 1: Valid request
	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/stop/test_task_id", strings.NewReader(`{"builderToken":"test_builder_token"}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Add more assertions for the response body
	})

	// Test case 2: Missing taskId
	t.Run("Missing TaskId", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/stop/", strings.NewReader(`{"builderToken":"test_builder_token"}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test case 3: Missing builderToken
	t.Run("Missing BuilderToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtt/stop/test_task_id", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Error:Field validation for 'BuilderToken' failed on the 'required' tag")
	})
}

func TestQueryRTT(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRTTService{
		QueryRTTFunc: func(c *gin.Context) {
			taskId := c.Param("taskId")
			if taskId == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "taskId is required"})
				return
			}
			builderToken := c.Query("builderToken")
			if builderToken == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "builderToken is required"})
				return
			}
			// Happy path response
			response := struct {
				Status    string `json:"status"`
				TaskId    string `json:"taskId"`
				Timestamp string `json:"timestamp"`
			}{
				Status:    "in_progress",
				TaskId:    taskId,
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.GET("/rtt/status/:taskId", mockService.QueryRTT)

	// Test case 1: Valid request
	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rtt/status/test_task_id?builderToken=test_builder_token", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		// Add more assertions for the response body
	})

	// Test case 2: Missing taskId
	t.Run("Missing TaskId", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rtt/status/?builderToken=test_builder_token", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	// Test case 3: Missing builderToken
	t.Run("Missing BuilderToken", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rtt/status/test_task_id", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "builderToken is required")
	})
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

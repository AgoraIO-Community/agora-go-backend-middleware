package rtmp_service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockRtmpService struct {
	StartPushFunc           func(c *gin.Context)
	StopPushFunc            func(c *gin.Context)
	StartPullFunc           func(c *gin.Context)
	StopPullFunc            func(c *gin.Context)
	UpdateConverterFunc     func(c *gin.Context)
	UpdatePlayerFunc        func(c *gin.Context)
	HandleStartPushReqFunc  func(startReq RtmpPushRequest, region string, regionHintIp *string, requestID string) (json.RawMessage, error)
	HandleStopPushReqFunc   func(converterId string, region string, requestID string) (json.RawMessage, error)
	HandleStartPullReqFunc  func(startReq CloudPlayerStartRequest, region string, streamOriginIp *string, requestID string) (json.RawMessage, error)
	HandleStopPullReqFunc   func(playerId string, region string, requestID string) (json.RawMessage, error)
	HandleUpdatePushReqFunc func(updateReq RtmpPushRequest, converterId string, region string, requestID string, sequenceId *int) (json.RawMessage, error)
	HandleUpdatePullReqFunc func(updateReq CloudPlayerStartRequest, playerId string, region string, requestID string, sequenceId *int) (json.RawMessage, error)
	AddTimestampFunc        func(response Timestampable) (json.RawMessage, error)
}

func (m *MockRtmpService) StartPush(c *gin.Context) {
	var req ClientStartRtmpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Add region validation
	if !m.ValidateRegion(req.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified"})
		return
	}
	response := StartRtmpResponse{
		Converter: ConverterResponse{
			ConverterId: "test_converter_id",
			CreateTs:    time.Now().Unix(),
			UpdateTs:    time.Now().Unix(),
			State:       "active",
		},
		Fields: "test_fields",
	}
	c.JSON(http.StatusOK, response)
}

func (m *MockRtmpService) StopPush(c *gin.Context) { m.StopPushFunc(c) }

func (m *MockRtmpService) UpdateConverter(c *gin.Context) {
	var req ClientUpdateRtmpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Add region validation
	if !m.ValidateRegion(req.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified"})
		return
	}
	response := StartRtmpResponse{
		Converter: ConverterResponse{
			ConverterId: req.ConverterId,
			CreateTs:    time.Now().Unix(),
			UpdateTs:    time.Now().Unix(),
			State:       "active",
		},
		Fields: "updated_fields",
	}
	c.JSON(http.StatusOK, response)
}

func (m *MockRtmpService) StartPull(c *gin.Context) { m.StartPullFunc(c) }
func (m *MockRtmpService) StopPull(c *gin.Context)  { m.StopPullFunc(c) }

func (m *MockRtmpService) UpdatePlayer(c *gin.Context) {
	var req ClientUpdatePullRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Add region validation
	if !m.ValidateRegion(req.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified"})
		return
	}
	response := CloudPlayerUpdateResponse{
		Status: "success",
	}
	c.JSON(http.StatusOK, response)
}

// Add this method to your MockRtmpService
func (m *MockRtmpService) ValidateRegion(region string) bool {
	validRegions := []string{"na", "eu", "ap", "cn"}
	for _, r := range validRegions {
		if r == region {
			return true
		}
	}
	return false
}

func (m *MockRtmpService) HandleStartPushReq(s RtmpPushRequest, r string, rh *string, rid string) (json.RawMessage, error) {
	return m.HandleStartPushReqFunc(s, r, rh, rid)
}
func (m *MockRtmpService) HandleStopPushReq(c string, r string, rid string) (json.RawMessage, error) {
	return m.HandleStopPushReqFunc(c, r, rid)
}
func (m *MockRtmpService) HandleStartPullReq(s CloudPlayerStartRequest, r string, so *string, rid string) (json.RawMessage, error) {
	return m.HandleStartPullReqFunc(s, r, so, rid)
}
func (m *MockRtmpService) HandleStopPullReq(p string, r string, rid string) (json.RawMessage, error) {
	return m.HandleStopPullReqFunc(p, r, rid)
}
func (m *MockRtmpService) HandleUpdatePushReq(u RtmpPushRequest, c string, r string, rid string, s *int) (json.RawMessage, error) {
	return m.HandleUpdatePushReqFunc(u, c, r, rid, s)
}
func (m *MockRtmpService) HandleUpdatePullReq(u CloudPlayerStartRequest, p string, r string, rid string, s *int) (json.RawMessage, error) {
	return m.HandleUpdatePullReqFunc(u, p, r, rid, s)
}
func (m *MockRtmpService) AddTimestamp(r Timestampable) (json.RawMessage, error) {
	return m.AddTimestampFunc(r)
}

func TestStartPush(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		StartPushFunc: func(c *gin.Context) {
			var req ClientStartRtmpRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := StartRtmpResponse{
				Converter: ConverterResponse{
					ConverterId: "test_converter_id",
					CreateTs:    time.Now().Unix(),
					UpdateTs:    time.Now().Unix(),
					State:       "active",
				},
				Fields: "test_fields",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/push/start", mockService.StartPush)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/push/start", strings.NewReader(`{"rtcChannel":"test_channel","streamUrl":"rtmp://test.com/live","streamKey":"test_key","region":"na","useTranscoding":true}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response StartRtmpResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Converter.ConverterId)
	})

	t.Run("Invalid Region", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/push/start", strings.NewReader(`{"rtcChannel":"test_channel","streamUrl":"rtmp://test.com/live","streamKey":"test_key","region":"invalid","useTranscoding":true}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid region specified")
	})
}

func TestStopPush(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		StopPushFunc: func(c *gin.Context) {
			var req ClientStopRtmpRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := StopRtmpResponse{
				Status: "success",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/push/stop", mockService.StopPush)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/push/stop", strings.NewReader(`{"converterId":"test_converter_id","region":"na"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response StopRtmpResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
	})
}

func TestStartPull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		StartPullFunc: func(c *gin.Context) {
			var req ClientStartCloudPlayerRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := StartCloudPlayerResponse{
				Player: PlayerResponse{
					PlayerId: "test_player_id",
					CreateTs: time.Now().Unix(),
				},
				Fields: "test_fields",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/pull/start", mockService.StartPull)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/pull/start", strings.NewReader(`{"channelName":"test_channel","streamUrl":"rtmp://test.com/live","region":"na"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response StartCloudPlayerResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Player.PlayerId)
	})
}

func TestStopPull(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		StopPullFunc: func(c *gin.Context) {
			var req ClientStopPullRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := CloudPlayerUpdateResponse{
				Status: "success",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/pull/stop", mockService.StopPull)

	t.Run("Valid Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/pull/stop", strings.NewReader(`{"playerId":"test_player_id","region":"na"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response CloudPlayerUpdateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
	})
}

func TestAddTimestamp(t *testing.T) {
	mockService := &MockRtmpService{
		AddTimestampFunc: func(response Timestampable) (json.RawMessage, error) {
			timestamp := time.Now().UTC().Format(time.RFC3339)
			response.SetTimestamp(timestamp)
			return json.Marshal(response)
		},
	}

	response := &StartRtmpResponse{
		Converter: ConverterResponse{
			ConverterId: "test_converter_id",
			CreateTs:    time.Now().Unix(),
			UpdateTs:    time.Now().Unix(),
			State:       "active",
		},
		Fields: "test_fields",
	}
	jsonResponse, err := mockService.AddTimestamp(response)
	if err != nil {
		t.Fatalf("AddTimestamp failed: %v", err)
	}

	var updatedResponse StartRtmpResponse
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

func TestUpdateConverter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		UpdateConverterFunc: func(c *gin.Context) {
			var req ClientUpdateRtmpRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := StartRtmpResponse{
				Converter: ConverterResponse{
					ConverterId: req.ConverterId,
					CreateTs:    time.Now().Unix(),
					UpdateTs:    time.Now().Unix(),
					State:       "active",
				},
				Fields: "updated_fields",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/push/update", mockService.UpdateConverter)

	t.Run("Valid Update Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/push/update", strings.NewReader(`{"converterId":"test_converter_id","region":"na","rtcChannel":"updated_channel","videoOptions":{"bitrate":2000}}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response StartRtmpResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "test_converter_id", response.Converter.ConverterId)
		assert.Equal(t, "updated_fields", response.Fields)
	})

	t.Run("Invalid Region", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/push/update", strings.NewReader(`{"converterId":"test_converter_id","region":"invalid","rtcChannel":"updated_channel"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid region specified")
	})
}

func TestUpdatePlayer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &MockRtmpService{
		UpdatePlayerFunc: func(c *gin.Context) {
			var req ClientUpdatePullRequest
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			response := CloudPlayerUpdateResponse{
				Status: "success",
			}
			c.JSON(http.StatusOK, response)
		},
	}

	router.POST("/rtmp/pull/update", mockService.UpdatePlayer)

	t.Run("Valid Update Request", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/pull/update", strings.NewReader(`{"playerId":"test_player_id","region":"na","streamUrl":"rtmp://updated.com/live"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response CloudPlayerUpdateResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
	})

	t.Run("Invalid Region", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/rtmp/pull/update", strings.NewReader(`{"playerId":"test_player_id","region":"invalid","streamUrl":"rtmp://updated.com/live"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Request-ID", "test-request-id")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid region specified")
	})
}

func TestValidateRegion(t *testing.T) {
	service := &RtmpService{}

	testCases := []struct {
		name     string
		region   string
		expected bool
	}{
		{"Valid region na", "na", true},
		{"Valid region eu", "eu", true},
		{"Valid region ap", "ap", true},
		{"Valid region cn", "cn", true},
		{"Invalid region", "invalid", false},
		{"Empty region", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.ValidateRegion(tc.region)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGenerateUID(t *testing.T) {
	service := &RtmpService{}

	for i := 0; i < 100; i++ {
		uid := service.GenerateUID()
		uidInt, err := strconv.Atoi(uid)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, uidInt, 1)
		assert.LessOrEqual(t, uidInt, 4294967294)
	}
}

func TestIsValidIPv4(t *testing.T) {
	service := &RtmpService{}

	testCases := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"Valid IPv4", "192.168.1.1", true},
		{"Invalid IPv4", "256.0.0.1", false},
		{"IPv6 address", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", false},
		{"Non-IP string", "not-an-ip", false},
		{"Empty string", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.isValidIPv4(tc.ip)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateIdleTimeOut(t *testing.T) {
	service := &RtmpService{}

	testCases := []struct {
		name     string
		input    int
		expected int
	}{
		{"Valid timeout", 100, 100},
		{"Minimum timeout", 5, 5},
		{"Maximum timeout", 600, 600},
		{"Below minimum", 4, 300},
		{"Above maximum", 601, 300},
		{"Negative value", -1, 300},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := service.ValidateIdleTimeOut(&tc.input)
			assert.Equal(t, tc.expected, *result)
		})
	}
}

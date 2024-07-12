package rtmp_service

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

type MockRtmpService struct {
	*RtmpService
	MockHandleStartPushReq  func(startReq RtmpPushRequest, region string, regionHintIp *string, requestID string) (json.RawMessage, error)
	MockHandleStopPushReq   func(converterId string, region string, requestID string) (json.RawMessage, error)
	MockHandleUpdatePushReq func(updateReq RtmpPushRequest, converterId string, region string, requestID string) (json.RawMessage, error)
}

func (m *MockRtmpService) StartPush(c *gin.Context) {
	// Extract the necessary information from the context
	var clientStartReq ClientStartRtmpRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if clientStartReq.RtcChannel == "" || clientStartReq.StreamUrl == "" || clientStartReq.StreamKey == "" || clientStartReq.Region == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	// Validate region
	if !m.ValidateRegion(clientStartReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region"})
		return
	}

	if m.MockHandleStartPushReq == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MockHandleStartPushReq not implemented"})
		return
	}

	// Call the mock function
	response, err := m.MockHandleStartPushReq(RtmpPushRequest{}, clientStartReq.Region, clientStartReq.RegionHintIp, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func (m *MockRtmpService) StopPush(c *gin.Context) {
	var clientStopReq ClientStopRtmpRequest
	if err := c.ShouldBindJSON(&clientStopReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := m.MockHandleStopPushReq(clientStopReq.ConverterId, clientStopReq.Region, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func (m *MockRtmpService) UpdateConverter(c *gin.Context) {
	var clientUpdateReq ClientUpdateRtmpRequest
	if err := c.ShouldBindJSON(&clientUpdateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := m.MockHandleUpdatePushReq(RtmpPushRequest{}, clientUpdateReq.ConverterId, clientUpdateReq.Region, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", response)
}

func TestNewRtmpService(t *testing.T) {
	appID := "testAppID"
	baseURL := "https://api.agora.io/"
	rtmpURL := "v1/projects/{appId}/rtmp-converters"
	basicAuth := "Basic dGVzdDp0ZXN0"

	service := NewRtmpService(appID, baseURL, rtmpURL, basicAuth)

	assert.NotNil(t, service)
	assert.Equal(t, appID, service.appID)
	assert.Equal(t, baseURL, service.baseURL)
	assert.Equal(t, rtmpURL, service.rtmpURL)
	assert.Equal(t, basicAuth, service.basicAuth)
}

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	service.RegisterRoutes(router)

	routes := router.Routes()
	expectedRoutes := []string{
		"POST /rtmp/push/start",
		"POST /rtmp/push/stop",
		"GET /rtmp/push/status",
		"POST /rtmp/push/update",
	}

	for _, route := range expectedRoutes {
		found := false
		for _, r := range routes {
			if r.Method+" "+r.Path == route {
				found = true
				break
			}
		}
		assert.True(t, found, "Route %s not found", route)
	}
}

func TestStartPush(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockRtmpService{
		RtmpService: NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0"),
	}
	mockService.MockHandleStartPushReq = func(startReq RtmpPushRequest, region string, regionHintIp *string, requestID string) (json.RawMessage, error) {
		response := StartRtmpResponse{
			Converter: ConverterResponse{
				ConverterId: "test-converter-id",
				CreateTs:    time.Now().Unix(),
				UpdateTs:    time.Now().Unix(),
				State:       "ACTIVE",
			},
			Fields: "test-fields",
		}
		return json.Marshal(response)
	}
	router.POST("/rtmp/push/start", mockService.StartPush)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rtmp/push/start", strings.NewReader(`{
        "converterName": "test-converter",
        "rtcChannel": "test-channel",
        "streamUrl": "rtmp://test.com/live",
        "streamKey": "test-key",
        "region": "na",
        "useTranscoding": true
    }`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-request-id")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response StartRtmpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "test-converter-id", response.Converter.ConverterId)
	assert.Equal(t, "ACTIVE", response.Converter.State)
}

func TestStopPush(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockRtmpService{
		RtmpService: NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0"),
	}
	mockService.MockHandleStopPushReq = func(converterId string, region string, requestID string) (json.RawMessage, error) {
		response := StopRtmpResponse{
			Status: "Success",
		}
		return json.Marshal(response)
	}
	router.POST("/rtmp/push/stop", mockService.StopPush)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rtmp/push/stop", strings.NewReader(`{
        "converterId": "test-converter-id",
        "region": "na"
    }`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-request-id")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response StopRtmpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Success", response.Status)
}

func TestUpdateConverter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockRtmpService{
		RtmpService: NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0"),
	}
	mockService.MockHandleUpdatePushReq = func(updateReq RtmpPushRequest, converterId string, region string, requestID string) (json.RawMessage, error) {
		response := StopRtmpResponse{
			Status: "Success",
		}
		return json.Marshal(response)
	}
	router.POST("/rtmp/push/update", mockService.UpdateConverter)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rtmp/push/update", strings.NewReader(`{
        "converterId": "test-converter-id",
        "region": "na",
        "rtcChannel": "updated-channel",
        "streamUrl": "rtmp://updated.com/live",
        "streamKey": "updated-key"
    }`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", "test-request-id")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response StopRtmpResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Success", response.Status)
}

func TestValidateRegion(t *testing.T) {
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	testCases := []struct {
		region   string
		expected bool
	}{
		{"na", true},
		{"eu", true},
		{"ap", true},
		{"cn", true},
		{"invalid", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := service.ValidateRegion(tc.region)
		assert.Equal(t, tc.expected, result, "Region: %s", tc.region)
	}
}

func TestAddTimestamp(t *testing.T) {
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	response := &StartRtmpResponse{
		Converter: ConverterResponse{
			ConverterId: "test-converter-id",
			CreateTs:    time.Now().Unix(),
			UpdateTs:    time.Now().Unix(),
			State:       "ACTIVE",
		},
		Fields: "test-fields",
	}

	timestampedResponse, err := service.AddTimestamp(response)
	assert.NoError(t, err)

	var updatedResponse StartRtmpResponse
	err = json.Unmarshal(timestampedResponse, &updatedResponse)
	assert.NoError(t, err)

	assert.NotNil(t, updatedResponse.Timestamp)
	timestamp, err := time.Parse(time.RFC3339, *updatedResponse.Timestamp)
	assert.NoError(t, err)
	assert.True(t, time.Since(timestamp) < 5*time.Second)
}

func TestIsValidIPv4(t *testing.T) {
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	testCases := []struct {
		ip       string
		expected bool
	}{
		{"192.168.0.1", true},
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"256.1.2.3", false},
		{"1.2.3.4.5", false},
		{"::1", false},
		{"2001:db8::1", false},
		{"not an ip", false},
		{"", false},
	}

	for _, tc := range testCases {
		result := service.isValidIPv4(tc.ip)
		assert.Equal(t, tc.expected, result, "IP: %s", tc.ip)
	}
}

func TestMakeRequest(t *testing.T) {
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Basic dGVzdDp0ZXN0", r.Header.Get("Authorization"))
		assert.Equal(t, "test-request-id", r.Header.Get("X-Request-ID"))

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-ID", "test-request-id")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer testServer.Close()

	body := map[string]string{"key": "value"}
	response, err := service.makeRequest("POST", testServer.URL, body, "test-request-id")

	assert.NoError(t, err)
	assert.Equal(t, []byte(`{"status":"success"}`), response)
}

func TestStartPushErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := &MockRtmpService{
		RtmpService: NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0"),
	}
	router.POST("/rtmp/push/start", mockService.StartPush)

	testCases := []struct {
		name           string
		body           string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid JSON",
			body:           `{"invalid json"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "invalid character '}' after object key",
		},
		{
			name:           "Missing required fields",
			body:           `{"converterName": "test-converter"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Missing required fields",
		},
		{
			name:           "Invalid region",
			body:           `{"rtcChannel": "test-channel", "streamUrl": "rtmp://test.com/live", "streamKey": "test-key", "region": "invalid", "useTranscoding": true}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid region",
		},
		{
			name:           "MockHandleStartPushReq not implemented",
			body:           `{"rtcChannel": "test-channel", "streamUrl": "rtmp://test.com/live", "streamKey": "test-key", "region": "na", "useTranscoding": true}`,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "MockHandleStartPushReq not implemented",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/rtmp/push/start", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Request-ID", "test-request-id")
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"], tc.expectedError)
		})
	}
}

func TestMakeRequestErrors(t *testing.T) {
	service := NewRtmpService("testAppID", "https://api.agora.io/", "v1/projects/{appId}/rtmp-converters", "Basic dGVzdDp0ZXN0")

	// Test non-existent server
	_, err := service.makeRequest("POST", "http://non-existent-server.com", nil, "test-request-id")
	assert.Error(t, err)

	// Test timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	_, err = service.makeRequest("POST", server.URL, nil, "test-request-id")
	assert.Error(t, err)
}

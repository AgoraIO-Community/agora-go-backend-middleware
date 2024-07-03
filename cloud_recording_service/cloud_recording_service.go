package cloud_recording_service

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AgoraIO-Community/agora-backend-service/token_service"
	"github.com/gin-gonic/gin"
)

// CloudRecordingService represents the cloud recording service.
// It holds the necessary configurations and dependencies for managing cloud recordings.
type CloudRecordingService struct {
	appID         string                      // The Agora app ID
	baseURL       string                      // The base URL for the Agora cloud recording API
	basicAuth     string                      // Middleware for handling requests
	tokenService  *token_service.TokenService // Token service for generating tokens
	storageConfig StorageConfig
}

// NewCloudRecordingService returns a CloudRecordingService pointer with all configurations set.
// This function initializes a new CloudRecordingService with specified configurations. It ensures all provided parameters are valid and logs a fatal error if any required configurations are missing.
//
// Parameters:
//   - tokenService: *token_service.TokenService - The token service for generating tokens.
//
// Returns:
//   - *CloudRecordingService: The initialized CloudRecordingService struct.
//
// Behavior:
//   - Initializes and returns a CloudRecordingService struct with the given configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewCloudRecordingService(appID string, baseURL string, basicAuth string, tokenService *token_service.TokenService, storageConfig StorageConfig) *CloudRecordingService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &CloudRecordingService{
		appID:         appID,         // The Agora app ID used to identify the application within Agora services.
		baseURL:       baseURL,       // The base URL for the Agora cloud recording API where all API requests are sent.
		basicAuth:     basicAuth,     // Basic authentication credentials required for interacting with the Agora API.
		tokenService:  tokenService,  // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
		storageConfig: storageConfig, // Configuration for storage options including directory structure and file naming.
	}
}

// RegisterRoutes registers the routes for the CloudRecordingService.
// It sets up the API endpoints and applies necessary middleware for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for cloud recording routes.
//   - Applies middleware for NoCache and CORS.
//   - Registers routes for ping, acquireResource, startRecording, stopRecording, getStatus, update subscriber list, and update layout.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *CloudRecordingService) RegisterRoutes(r *gin.Engine) {
	// group route
	api := r.Group("/cloud_recording")
	// routes
	api.POST("/start", s.StartRecording)
	api.POST("/stop", s.StopRecording)
	api.GET("/status", s.GetStatus)
	// "update" group route
	updateAPI := api.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList)
	updateAPI.POST("/layout", s.UpdateLayout)
}

func (s *CloudRecordingService) StartRecording(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStartReq ClientStartRecordingRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sceneMode := 0
	sceneModes := map[string]int{"realtime": 0, "web": 1, "postponed": 2}
	if clientStartReq.SceneMode != nil {
		if mode, ok := sceneModes[*clientStartReq.SceneMode]; ok {
			sceneMode = mode
		}
	}

	// Default RecordingMode to "composite" if nil
	recordingMode := "mix"
	if clientStartReq.RecordingMode != nil {
		recordingMode = *clientStartReq.RecordingMode
	}

	// Validate recording mode against a set list
	validRecordingModes := []string{"individual", "mix", "web"}
	if !Contains(validRecordingModes, recordingMode) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recording mode."})
		return
	}

	// Generate a unique UID for this recording session
	uid := generateUID()

	// Generate token for recording using token_service
	tokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   clientStartReq.ChannelName,
		Uid:       uid,
	}
	token, err := s.tokenService.GenRtcToken(tokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Add dynamic directory structure ChannelName/YYYYMMDD/HHMMSS
	currentTimeUTC := time.Now().UTC()
	dateStr := currentTimeUTC.Format("20060102")
	hrsMinSecStr := currentTimeUTC.Format("150405")
	s.storageConfig.FileNamePrefix = &[]string{strings.ReplaceAll(clientStartReq.ChannelName, "-", ""), dateStr, hrsMinSecStr}

	// Check if RecordingConfig is nil, if so, create a default one
	if clientStartReq.RecordingConfig == nil {
		// create default recording config
		channelType := 0
		streamTypes := 2
		videoStreamType := 0
		maxIdleTime := 120
		subscribeUidGroup := 0
		streamMode := "standard" // Default to "standard", ensure to modify based on your actual logic
		subscribeAudioUids := []string{"#allstream#"}
		subscribeVideoUids := []string{"#allstream#"}
		// create default recording config
		clientStartReq.RecordingConfig = &RecordingConfig{
			ChannelType:        channelType,
			StreamTypes:        &streamTypes,
			VideoStreamType:    &videoStreamType,
			StreamMode:         &streamMode,
			MaxIdleTime:        &maxIdleTime,
			SubscribeAudioUids: &subscribeAudioUids,
			SubscribeVideoUids: &subscribeVideoUids,
			SubscribeUidGroup:  &subscribeUidGroup,
		}
	}

	// Assemble recording client request
	recClientReq := AquireClientRequest{
		Scene:               sceneMode,
		ResourceExpiredHour: 24, // Assuming 24 hours, adjust as needed
		StartParameter: ClientRequest{
			Token:           token,
			StorageConfig:   s.storageConfig,
			RecordingConfig: *clientStartReq.RecordingConfig,
		},
		ExcludeResourceIds: clientStartReq.ExcludeResourceIds,
	}

	// Acquire Resource
	acquireReq := AcquireResourceRequest{
		Cname:         clientStartReq.ChannelName,
		Uid:           uid,
		ClientRequest: &recClientReq, // Initialize as an empty map
	}
	resourceID, err := s.HandleAcquireResourceReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	log.Println("resourceID:", (resourceID))

	// Build the full StartRecordingRequest
	startReq := StartRecordingRequest{
		Cname: clientStartReq.ChannelName,
		Uid:   uid,
		ClientRequest: ClientRequest{
			Token:           token,
			StorageConfig:   s.storageConfig,
			RecordingConfig: *clientStartReq.RecordingConfig,
		},
	}

	// Start Recording
	response, err := s.HandleStartRecordingReq(startReq, resourceID, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start recording: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// StopRecording
func (s *CloudRecordingService) StopRecording(c *gin.Context) {
	var respWriter = c.Writer
	var clientStopReq ClientStopRecordingRequest
	err := json.NewDecoder(c.Request.Body).Decode(&clientStopReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// build the stop request from user request
	stopReq := StopRecordingRequest{
		Cname: clientStopReq.Cname,
		Uid:   clientStopReq.Uid,
		ClientRequest: StopClientRequest{
			AsyncStop: clientStopReq.AsyncStop,
		},
	}
	recordingMode := "mix"
	if clientStopReq.RecordingMode != nil {
		recordingMode = *clientStopReq.RecordingMode
	}
	// Send Stop Recording Request to Agora
	response, err := s.HandleStopRecording(stopReq, clientStopReq.ResourceId, clientStopReq.Sid, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// GetStatus
func (s *CloudRecordingService) GetStatus(c *gin.Context) {
	// s.HandleGetStatus(c.Writer, c.Request)
}

// UpdateSubscriptionList
func (s *CloudRecordingService) UpdateSubscriptionList(c *gin.Context) {
	var respWriter = c.Writer
	var clientUpdateReq ClientUpdateSubscriptionRequest
	err := json.NewDecoder(c.Request.Body).Decode(&clientUpdateReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate the update request
	if !clientUpdateReq.UpdateConfig.IsValid() {
		http.Error(respWriter, "Invalid update configuration", http.StatusBadRequest)
		return
	}
	// build the stop request from user request
	updateReq := UpdateSubscriptionRequest{
		Cname:         clientUpdateReq.Cname,
		Uid:           clientUpdateReq.Uid,
		ClientRequest: clientUpdateReq.UpdateConfig,
	}

	recordingMode := "mix"
	if clientUpdateReq.RecordingMode != nil {
		recordingMode = *clientUpdateReq.RecordingMode
	}

	// Send Stop Recording Request to Agora
	response, err := s.HandleUpdateSubscriptionList(updateReq, clientUpdateReq.ResourceId, clientUpdateReq.Sid, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// UpdateLayout
func (s *CloudRecordingService) UpdateLayout(c *gin.Context) {
	var respWriter = c.Writer
	var clientUpdateReq ClientUpdateLayoutRequest
	err := json.NewDecoder(c.Request.Body).Decode(&clientUpdateReq)
	if err != nil {
		// invalid request
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// build the stop request from user request
	updateReq := UpdateLayoutRequest{
		Cname:         clientUpdateReq.Cname,
		Uid:           clientUpdateReq.Uid,
		ClientRequest: clientUpdateReq.UpdateConfig,
	}

	recordingMode := "mix"
	if clientUpdateReq.RecordingMode != nil {
		recordingMode = *clientUpdateReq.RecordingMode
	}

	// Send Stop Recording Request to Agora
	response, err := s.HandleUpdateLayout(updateReq, clientUpdateReq.ResourceId, clientUpdateReq.Sid, recordingMode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

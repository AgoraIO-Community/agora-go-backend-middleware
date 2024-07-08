package real_time_transcription_service

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/cloud_recording_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
	"github.com/gin-gonic/gin"
)

// RTTService struct holds all the necessary configurations and dependencies
// required for managing real-time transcription services.
type RTTService struct {
	appID         string                                // Agora application ID to identify the application within Agora services.
	baseURL       string                                // Base URL for the Agora cloud recording API where all API requests are sent.
	basicAuth     string                                // Basic authentication credentials required for interacting with the Agora API.
	tokenService  *token_service.TokenService           // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
	storageConfig cloud_recording_service.StorageConfig // Configuration for storage options including directory structure and file naming.
}

// NewRTTService initializes a new instance of RTTService with the provided configurations.
// It seeds the random number generator to ensure varied operational behavior.
//
// Parameters:
//   - appID: The Agora application ID.
//   - baseURL: Base URL for the API interactions.
//   - basicAuth: Basic authentication credentials for the API.
//   - tokenService: Token service instance for generating tokens.
//   - storageConfig: Storage configuration detailing file and directory naming conventions.
//
// Returns:
//   - A pointer to the newly created RTTService.
func NewRTTService(appID string, baseURL string, basicAuth string, tokenService *token_service.TokenService, storageConfig cloud_recording_service.StorageConfig) *RTTService {
	rand.Seed(time.Now().UnixNano()) // Ensure varied randomness in the application operations.
	return &RTTService{
		appID:         appID,
		baseURL:       baseURL,
		basicAuth:     basicAuth,
		tokenService:  tokenService,
		storageConfig: storageConfig,
	}
}

// RegisterRoutes sets up the API endpoints related to the real-time transcription service.
// It creates a route group and registers individual routes for starting, stopping, and querying the transcription status.
//
// Parameters:
//   - r: *gin.Engine - Gin engine instance to register routes.
func (s *RTTService) RegisterRoutes(r *gin.Engine) {
	// group route
	api := r.Group("/rtt")
	// routes
	api.POST("/start", s.StartRTT)
	api.POST("/stop", s.StopRTT)
	api.GET("/status", s.QueryRTT)
}

// StartRTT handles the starting of the real-time transcription by binding JSON data from client requests,
// validating and setting default values, acquiring necessary tokens, and making the start request.
//
// Parameters:
//   - c: *gin.Context - Context instance containing HTTP request and response objects.
func (s *RTTService) StartRTT(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStartReq ClientStartRTTRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// s.ValidateAndSetDefaults(&clientStartReq) // Validate client request and set default values.

	// Acquire Builder Token
	acquireReq := AcquireBuilderTokenRequest{
		InstanceId: clientStartReq.ChannelName,
	}
	acquireResponse, builderToken, err := s.HandleAcquireBuilderTokenReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	subscriberBotUid := s.GenerateUID() // Generate a unique identifier for the audio subscriber bot.
	publisherBotUid := s.GenerateUID()  // Generate a unique identifier for the output bot.

	// Generate subscriber token for rtt using token_service
	subscriberBotTokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   clientStartReq.ChannelName,
		Uid:       subscriberBotUid,
	}
	subscriberBotToken, err := s.tokenService.GenRtcToken(subscriberBotTokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate publisher token for rtt using token_service
	publisherBotTokenRequest := token_service.TokenRequest{
		TokenType: "rtc",
		Channel:   clientStartReq.ChannelName,
		Uid:       subscriberBotUid,
	}
	publisherBotToken, err := s.tokenService.GenRtcToken(publisherBotTokenRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construct the start request
	startRttRequest := StartRTTRequest{
		Languages:   []string{},
		MaxIdleTime: *clientStartReq.MaxIdleTime,
		RTCConfig: RTCConfig{
			ChannelName:        clientStartReq.ChannelName,
			SubBotUID:          subscriberBotUid,
			SubBotToken:        &subscriberBotToken,
			PubBotUID:          publisherBotUid,
			PubBotToken:        &publisherBotToken,
			SubscribeAudioUIDs: clientStartReq.SubscribeAudioUIDs,
			CryptionMode:       clientStartReq.CryptionMode,
			Secret:             clientStartReq.Secret,
			Salt:               clientStartReq.Salt,
		},
		TranslateConfig: clientStartReq.TranslateConfig,
	}

	// If storage is in destinations list, add storage config
	if clientStartReq.EnableStorage != nil && *clientStartReq.EnableStorage {
		// Add dynamic directory structure ChannelName/YYYYMMDD/HHMMSS
		currentTimeUTC := time.Now().UTC()
		dateStr := currentTimeUTC.Format("20060102")
		hrsMinSecStr := currentTimeUTC.Format("150405")
		s.storageConfig.FileNamePrefix = &[]string{strings.ReplaceAll(clientStartReq.ChannelName, "-", ""), dateStr, hrsMinSecStr}

		// Enable subtitle sync
		if clientStartReq.EnableNTPtimestamp != nil && *clientStartReq.EnableNTPtimestamp {
			s.storageConfig.ExtensionParams.EnableNTPtimestamp = clientStartReq.EnableNTPtimestamp
		}
		// set cloud storage in request
		startRttRequest.CaptionConfig.Storage = s.storageConfig
	}

	// Make the Start Request to Agora Endpoint
	startResponse, err := s.HandleStartReq(startRttRequest, builderToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transcription: " + err.Error()})
		return
	}

	// Return acquire and start responses
	c.JSON(http.StatusOK, gin.H{
		"acquire":   acquireResponse,
		"start":     startResponse,
		"timestamp": time.Now().UTC(),
	})
}

func (s *RTTService) StopRTT(c *gin.Context)  {}
func (s *RTTService) QueryRTT(c *gin.Context) {}

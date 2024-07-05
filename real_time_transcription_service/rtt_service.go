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

type RTTService struct {
	appID         string                      // The Agora app ID
	baseURL       string                      // The base URL for the Agora cloud recording API
	basicAuth     string                      // Middleware for handling requests
	tokenService  *token_service.TokenService // Token service for generating tokens
	storageConfig cloud_recording_service.StorageConfig
}

func NewRTTService(appID string, baseURL string, basicAuth string, tokenService *token_service.TokenService, storageConfig cloud_recording_service.StorageConfig) *RTTService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &RTTService{
		appID:         appID,         // The Agora app ID used to identify the application within Agora services.
		baseURL:       baseURL,       // The base URL for the Agora cloud recording API where all API requests are sent.
		basicAuth:     basicAuth,     // Basic authentication credentials required for interacting with the Agora API.
		tokenService:  tokenService,  // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
		storageConfig: storageConfig, // Configuration for storage options including directory structure and file naming.
	}
}

func (s *RTTService) RegisterRoutes(r *gin.Engine) {
	// group route
	api := r.Group("/rtt")
	// routes
	api.POST("/start", s.StartRTT)
	api.POST("/stop", s.StopRTT)
	api.GET("/status", s.QueryRTT)
}

func (s *RTTService) StartRTT(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStartReq ClientStartRTTRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s.ValidateAndSetDefaults(&clientStartReq)

	// Acquire Builder Token
	acquireReq := AcquireBuilderTokenRequest{
		InstanceId: clientStartReq.ChannelName,
	}
	acquireResponse, builderToken, err := s.HandleAcquireBuilderTokenReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	// Generate a unique UID for this recording session
	uid := s.GenerateUID()

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

	// Construct the start request
	startRttRequest := StartRTTRequest{
		Audio: Audio{
			SubscribeSource: "AGORARTC",
			AgoraRtcConfig: AgoraRtcConfig{
				ChannelName:     clientStartReq.ChannelName,
				UID:             uid,
				Token:           token,
				ChannelType:     "LIVE_TYPE",
				SubscribeConfig: SubscribeConfig{},
				MaxIdleTime:     *clientStartReq.MaxIdleTime,
			},
		},
		Config: Config{
			Features: []string{"RECOGNIZE"},
			RecognizeConfig: RecognizeConfig{
				Language:        "",
				Model:           "Model",
				ProfanityFilter: clientStartReq.ProfanityFilter,
				Output: Output{
					Destinations:       *clientStartReq.Destinations,
					AgoraRTCDataStream: AgoraRTCDataStream{},
				},
			},
		},
	}

	// If storage is in destinations list, add storage config
	if s.Contains(clientStartReq.Destinations, "Storage") {
		// Add dynamic directory structure ChannelName/YYYYMMDD/HHMMSS
		currentTimeUTC := time.Now().UTC()
		dateStr := currentTimeUTC.Format("20060102")
		hrsMinSecStr := currentTimeUTC.Format("150405")
		s.storageConfig.FileNamePrefix = &[]string{strings.ReplaceAll(clientStartReq.ChannelName, "-", ""), dateStr, hrsMinSecStr}
		// set cloud storage in request
		startRttRequest.Config.RecognizeConfig.Output.CloudStorage = &[]CloudStorage{
			{Format: "HLS", StorageConfig: s.storageConfig},
		}
	}

	// Enable for subtitle sync
	if clientStartReq.EnableNTPtimestamp != nil {
		startRttRequest.PrivateParams.EnableNTPtimestamp = *clientStartReq.EnableNTPtimestamp
	}

	// Make the Start Request to Agora Endpoint
	startResponse, err := s.HandleStartReq(startRttRequest, builderToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}

	// Return Resource ID and Recording ID
	c.JSON(http.StatusOK, gin.H{
		"acquire":   acquireResponse,
		"start":     startResponse,
		"timestamp": time.Now().UTC(),
	})

}
func (s *RTTService) StopRTT(c *gin.Context)  {}
func (s *RTTService) QueryRTT(c *gin.Context) {}

package real_time_transcription_service

import (
	"math/rand"
	"net/http"
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
	// Acquire Builder Token
	acquireReq := AcquireBuilderTokenRequest{
		InstanceId: clientStartReq.ChannelName,
	}
	acquireResponse, err := s.HandleAcquireResourceReq(acquireReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to acquire resource: " + err.Error()})
		return
	}
	// Build the full StartRTTRequest

	// Return Resource ID and Recording ID
	c.JSON(http.StatusOK, gin.H{
		"acquire": acquireResponse,
		// "start":     startResponse,
		"timestamp": time.Now().UTC(),
	})

}
func (s *RTTService) StopRTT(c *gin.Context)  {}
func (s *RTTService) QueryRTT(c *gin.Context) {}

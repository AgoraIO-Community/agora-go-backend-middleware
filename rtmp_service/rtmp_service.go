package rtmp_service

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

// RtmpService represents the media push/pull service.
// It holds the necessary configurations and dependencies for managing Media Push and Pull with Agora Channels.
type RtmpService struct {
	appID     string // The Agora app ID
	baseURL   string // The base URL for the Agora API
	rtmpURL   string // The URL path for the rtmp converter
	basicAuth string // Middleware for handling requests
	// tokenService  *token_service.TokenService // Token service for generating tokens
}

// NewRtmpService returns a RtmpService pointer with all configurations set.
// This function initializes a new RtmpService with specified configurations. It ensures all provided parameters are valid and logs a fatal error if any required configurations are missing.
//
// Parameters:
//   - tokenService: *token_service.TokenService - The token service for generating tokens.
//
// Returns:
//   - *RtmpService: The initialized RtmpService struct.
//
// Behavior:
//   - Initializes and returns a RtmpService struct with the given configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewRtmpService(appID string, baseURL string, rtmpURL string, basicAuth string) *RtmpService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &RtmpService{
		appID:     appID,     // The Agora app ID used to identify the application within Agora services.
		baseURL:   baseURL,   // The base URL for the Agora API where all API requests are sent.
		rtmpURL:   rtmpURL,   // The URL path for the Agora rtmp converter endpoint
		basicAuth: basicAuth, // Basic authentication credentials required for interacting with the Agora API.
		// tokenService:  tokenService,  // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
	}
}

// RegisterRoutes registers the routes for the RtmpService.
// It sets up the API endpoints for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for rtmp routes.
//   - Registers routes for ping, acquireResource, startRecording, stopRecording, getStatus, update subscriber list, and update layout.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *RtmpService) RegisterRoutes(r *gin.Engine) {
	// group route
	api := r.Group("/rtmp")
	// routes
	api.POST("/start", s.StartPush)
	api.POST("/stop", s.StopPush)
	api.GET("/status", s.GetStatus)
	// "update" group route
	updateAPI := api.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList)
	updateAPI.POST("/layout", s.UpdateLayout)
}

func (s *RtmpService) StartPush(c *gin.Context) {}
func (s *RtmpService) StopPush(c *gin.Context)  {}
func (s *RtmpService) GetStatus(c *gin.Context) {}

func (s *RtmpService) UpdateSubscriptionList(c *gin.Context) {}
func (s *RtmpService) UpdateLayout(c *gin.Context)           {}

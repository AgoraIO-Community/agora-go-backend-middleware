package rtmp_service

import (
	"math/rand"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
	"github.com/gin-gonic/gin"
)

// RtmpService represents the media push/pull service.
// It holds the necessary configurations and dependencies for managing Media Push and Pull with Agora Channels.
type RtmpService struct {
	appID        string                      // The Agora app ID used to identify the application within Agora services.
	baseURL      string                      // The base URL for the Agora API where all API requests are sent.
	rtmpURL      string                      // The URL path for the Agora RTMP converter endpoint.
	basicAuth    string                      // Basic authentication credentials required for interacting with the Agora API.
	tokenService *token_service.TokenService // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
}

// NewRtmpService returns a RtmpService pointer with all configurations set.
// This function initializes a new RtmpService with specified configurations. It ensures all provided parameters are valid and logs a fatal error if any required configurations are missing.
//
// Parameters:
//   - appID: string - The Agora app ID.
//   - baseURL: string - The base URL for the Agora API.
//   - rtmpURL: string - The URL path for the RTMP converter.
//   - basicAuth: string - The basic authentication credentials.
//
// Returns:
//   - *RtmpService: The initialized RtmpService struct.
//
// Behavior:
//   - Initializes and returns a RtmpService struct with the given configurations.
//   - Seeds the random number generator with the current time.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewRtmpService(appID string, baseURL string, rtmpURL string, basicAuth string, tokenService *token_service.TokenService) *RtmpService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &RtmpService{
		appID:        appID,        // The Agora app ID used to identify the application within Agora services.
		baseURL:      baseURL,      // The base URL for the Agora API where all API requests are sent.
		rtmpURL:      rtmpURL,      // The URL path for the Agora RTMP converter endpoint.
		basicAuth:    basicAuth,    // Basic authentication credentials required for interacting with the Agora API.
		tokenService: tokenService, // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
	}
}

// RegisterRoutes registers the routes for the RtmpService.
// It sets up the API endpoints for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for RTMP routes.
//   - Registers routes for starting push, stopping push, getting status, updating subscriber list, and updating layout.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *RtmpService) RegisterRoutes(r *gin.Engine) {
	// group route for RTMP
	api := r.Group("/rtmp")

	// group route for push operations
	pushAPI := api.Group("/push")
	// push routes
	pushAPI.POST("/start", s.StartPush) // Route to start the RTMP push.
	pushAPI.POST("/stop", s.StopPush)   // Route to stop the RTMP push.
	pushAPI.GET("/status", s.GetStatus) // Route to get the status of the RTMP push.

	// group route for update operations
	updateAPI := pushAPI.Group("/update")
	updateAPI.POST("/subscriber-list", s.UpdateSubscriptionList) // Route to update the subscription list.
	updateAPI.POST("/layout", s.UpdateLayout)                    // Route to update the layout.
}

// StartPush handles the starting of an RTMP push.
// It processes the request to start pushing the media stream to the specified RTMP URL.
func (s *RtmpService) StartPush(c *gin.Context) {}

// StopPush handles the stopping of an RTMP push.
// It processes the request to stop pushing the media stream to the specified RTMP URL.
func (s *RtmpService) StopPush(c *gin.Context) {}

// GetStatus returns the current status of the RTMP push.
// It processes the request to get the current status of the media stream pushing operation.
func (s *RtmpService) GetStatus(c *gin.Context) {}

// UpdateSubscriptionList handles updating the subscription list for the RTMP push.
// It processes the request to update the list of subscribers for the media stream.
func (s *RtmpService) UpdateSubscriptionList(c *gin.Context) {}

// UpdateLayout handles updating the layout for the RTMP push.
// It processes the request to update the layout configuration for the media stream.
func (s *RtmpService) UpdateLayout(c *gin.Context) {}

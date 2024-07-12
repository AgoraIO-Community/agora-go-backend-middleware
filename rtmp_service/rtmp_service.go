package rtmp_service

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RtmpService represents the media push/pull service.
// It holds the necessary configurations and dependencies for managing Media Push and Pull with Agora Channels.
type RtmpService struct {
	appID     string // The Agora app ID used to identify the application within Agora services.
	baseURL   string // The base URL for the Agora API where all API requests are sent.
	rtmpURL   string // The URL path for the Agora RTMP converter endpoint.
	basicAuth string // Basic authentication credentials required for interacting with the Agora API.
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
func NewRtmpService(appID string, baseURL string, rtmpURL string, basicAuth string) *RtmpService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &RtmpService{
		appID:     appID,     // The Agora app ID used to identify the application within Agora services.
		baseURL:   baseURL,   // The base URL for the Agora API where all API requests are sent.
		rtmpURL:   rtmpURL,   // The URL path for the Agora RTMP converter endpoint.
		basicAuth: basicAuth, // Basic authentication credentials required for interacting with the Agora API.
	}
}

// Middleware to verify X-Request-ID header
func verifyRequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Request-ID") == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "X-Request-ID header is required"})
			c.Abort()
			return
		}
		c.Next()
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
	// make sure each request containers the X-Request-ID
	api.Use(verifyRequestID())
	// group route for push operations
	pushAPI := api.Group("/push")
	// push routes
	pushAPI.POST("/start", s.StartPush)        // Route to start the RTMP push.
	pushAPI.POST("/stop", s.StopPush)          // Route to stop the RTMP push.
	pushAPI.GET("/status", s.GetStatus)        // Route to get the status of the RTMP push.
	pushAPI.POST("/update", s.UpdateConverter) // Route to update the converter.

}

// StartPush handles the starting of an RTMP push.
// It processes the request to start pushing the media stream to the specified RTMP URL.
func (s *RtmpService) StartPush(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStartReq ClientStartRtmpRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate region
	if !s.ValidateRegion(clientStartReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified."})
		return
	}

	// Assemble rtmp client request
	rtmpClientReq := RtmpPushRequest{
		Converter: Converter{
			Name:               clientStartReq.ConverterName,
			RtmpUrl:            clientStartReq.StreamUrl + clientStartReq.StreamKey,
			IdleTimeOut:        clientStartReq.IdleTimeOut,
			JitterBufferSizeMs: clientStartReq.JitterBufferSizeMs,
		},
	}

	if clientStartReq.UseTranscoding {
		// set rtmp request to use TranscodeOptions
		rtmpClientReq.Converter.TranscodeOptions = &TranscodeOptions{
			RtcChannel:   clientStartReq.RtcChannel,
			AudioOptions: clientStartReq.AudioOptions,
			VideoOptions: clientStartReq.VideoOptions,
		}
	} else {
		// set rtmp request to use RawOptions
		rtmpClientReq.Converter.RawOptions = &RawOptions{
			RtcChannel:   clientStartReq.RtcChannel,
			RtcStreamUid: *clientStartReq.RtcStreamUid,
		}
	}

	// Start RTMP
	response, err := s.HandleStartPushReq(rtmpClientReq, clientStartReq.Region, clientStartReq.RegionHintIp, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start RTMP converter: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)

}

// StopPush handles the stopping of an RTMP push.
// It processes the request to stop pushing the media stream to the specified RTMP URL.
func (s *RtmpService) StopPush(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStopReq ClientStopRtmpRequest
	if err := c.ShouldBindJSON(&clientStopReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate region
	if !s.ValidateRegion(clientStopReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified."})
		return
	}

	// Stop RTMP
	response, err := s.HandleStopPushReq(clientStopReq.ConverterId, clientStopReq.Region, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop RTMP converter: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)

}

// GetStatus returns the current status of the RTMP push.
// It processes the request to get the current status of the media stream pushing operation.
func (s *RtmpService) GetStatus(c *gin.Context) {}

// UpdateConverter handles updating the transcoding options for the RTMP push.
// It processes the request to update the transcoding configuration for the media stream.
func (s *RtmpService) UpdateConverter(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientUpdateReq ClientUpdateRtmpRequest
	if err := c.ShouldBindJSON(&clientUpdateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if clientUpdateReq.VideoOptions != nil {
		// Update doesnt support changes to Codec or CodecProfile
		clientUpdateReq.VideoOptions.Codec = nil
		clientUpdateReq.VideoOptions.CodecProfile = nil
	}

	// Validate region
	if !s.ValidateRegion(clientUpdateReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified."})
		return
	}

	// Update RTMP
	response, err := s.HandleStopPushReq(clientUpdateReq.ConverterId, clientUpdateReq.Region, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop RTMP converter: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

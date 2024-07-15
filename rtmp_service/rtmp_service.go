package rtmp_service

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
	"github.com/gin-gonic/gin"
)

// RtmpService represents the media push and pull services.
// It holds the necessary configurations and dependencies for managing Media Push & Pull with Agora Channels.
type RtmpService struct {
	appID          string                      // The Agora app ID used to identify the application within Agora services.
	baseURL        string                      // The base URL for the Agora API where all API requests are sent.
	rtmpURL        string                      // The URL path for the Agora RTMP converter endpoint.
	cloudPlayerURL string                      // The URL path for the Agora Clpoud Player endpoint.
	basicAuth      string                      // Basic authentication credentials required for interacting with the Agora API.
	tokenService   *token_service.TokenService // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
}

// NewRtmpService returns a RtmpService pointer with all configurations set.
// This function initializes a new RtmpService with specified configurations. It ensures all provided parameters are valid and logs a fatal error if any required configurations are missing.
//
// Parameters:
//   - appID: string - The Agora app ID.
//   - baseURL: string - The base URL for the Agora API.
//   - rtmpURL: string - The URL path for Agora's RTMP converter service.
//   - cloudPlayerURL: string - The URL path for Agora's Cloud Player service.
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
func NewRtmpService(appID string, baseURL string, rtmpURL string, cloudPlayerURL string, basicAuth string, tokenService *token_service.TokenService) *RtmpService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &RtmpService{
		appID:          appID,          // The Agora app ID used to identify the application within Agora services.
		baseURL:        baseURL,        // The base URL for the Agora API where all API requests are sent.
		rtmpURL:        rtmpURL,        // The URL path for the Agora RTMP converter endpoint.
		cloudPlayerURL: cloudPlayerURL, // The URL path for the Agora Clpoud Player endpoint.
		basicAuth:      basicAuth,      // Basic authentication credentials required for interacting with the Agora API.
		tokenService:   tokenService,   // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.

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

	// setup push routes (rtmp converter)
	if s.rtmpURL != "" {
		// group route for push operations
		pushAPI := api.Group("/push")
		// push routes
		pushAPI.POST("/start", s.StartPush)        // Route to start the RTMP push.
		pushAPI.POST("/stop", s.StopPush)          // Route to stop the RTMP push.
		pushAPI.GET("/list", s.GetPushList)        // Route to get the list of RTMP converters.
		pushAPI.POST("/update", s.UpdateConverter) // Route to update the converter.
	}

	// setup pull routes (cloud player)
	if s.cloudPlayerURL != "" {
		// group route for pull operations
		pullAPI := api.Group("/pull")
		// pull routes
		pullAPI.POST("/start", s.StartPull)     // Route to start the RTMP push.
		pullAPI.POST("/stop", s.StopPull)       // Route to stop the RTMP push.
		pullAPI.POST("/update", s.UpdatePlayer) // Route to update the converter.
		pullAPI.GET("/list", s.GetPullList)     // RRoute to get the list of cloud players
	}
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
	rtmpPushURL := clientStartReq.StreamUrl + clientStartReq.StreamKey
	rtmpClientReq := RtmpPushRequest{
		Converter: Converter{
			Name:               clientStartReq.ConverterName,
			RtmpUrl:            &rtmpPushURL,
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

// GetPushList returns a list of the current RTMP converters.
// It processes the request to get the current list of the media stream pushing operations.
func (s *RtmpService) GetPushList(c *gin.Context) {
	// TODO: add logic to check context for api group
}

// UpdateConverter handles updating the transcoding options for the RTMP push.
// It processes the request to update the transcoding configuration for the media stream.
func (s *RtmpService) UpdateConverter(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientUpdateReq ClientUpdateRtmpRequest
	if err := c.ShouldBindJSON(&clientUpdateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate region
	if !s.ValidateRegion(clientUpdateReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified."})
		return
	}

	if clientUpdateReq.VideoOptions != nil {
		// Update doesnt support changes to Codec or CodecProfile
		clientUpdateReq.VideoOptions.Codec = nil
		clientUpdateReq.VideoOptions.CodecProfile = nil
	}

	// Assemble rtmp client request
	rtmpClientReq := RtmpPushRequest{
		Converter: Converter{
			JitterBufferSizeMs: clientUpdateReq.JitterBufferSizeMs,
			TranscodeOptions: &TranscodeOptions{
				RtcChannel:   clientUpdateReq.RtcChannel,
				VideoOptions: clientUpdateReq.VideoOptions,
			},
		},
	}

	// update rtmp url if defined
	if clientUpdateReq.StreamUrl != nil && clientUpdateReq.StreamKey != nil {
		rtmpPushURL := *clientUpdateReq.StreamUrl + *clientUpdateReq.StreamKey
		rtmpClientReq.Converter.RtmpUrl = &rtmpPushURL
	}

	// Update RTMP
	response, err := s.HandleUpdatePushReq(rtmpClientReq, clientUpdateReq.ConverterId, clientUpdateReq.Region, c.GetHeader("X-Request-ID"), clientUpdateReq.SequenceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update RTMP converter: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// StartPull handles the starting of an RTMP pull.
// It processes the request to start pulling the media stream from the specified RTMP URL into the given channel.
func (s *RtmpService) StartPull(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStartReq ClientStartCloudPlayerRequest
	if err := c.ShouldBindJSON(&clientStartReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate region against a set list
	if !s.ValidateRegion(clientStartReq.Region) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid region specified."})
		return
	}

	// Assign uid
	var uid string
	// check & asign uid from client start request
	if clientStartReq.Uid != nil {
		uid = *clientStartReq.Uid
	} else {
		// Generate a unique UID for this cloud player
		uid = s.GenerateUID()
	}

	if clientStartReq.IdleTimeOut != nil {
		clientStartReq.IdleTimeOut = s.ValidateIdleTimeOut(clientStartReq.IdleTimeOut)
	}

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

	cloudPlayerClientReq := CloudPlayerStartRequest{
		Player: Player{
			StreamUrl:   clientStartReq.StreamUrl,
			ChannelName: clientStartReq.ChannelName,
			Token:       token,
			Uid:         uid,
			IdleTimeOut: clientStartReq.IdleTimeOut,
			PlayTs:      clientStartReq.PlayTs,
			EncryptMode: clientStartReq.EncryptMode,
			PlayerName:  clientStartReq.PlayerName,
		},
	}

	// check and add transcoding config
	if clientStartReq.VideoOptions != nil {
		// check and add audio options
		var audioOptions PullAudioOptions
		if clientStartReq.AudioOptions != nil {
			audioOptions = *clientStartReq.AudioOptions
		} else {
			audioOptions = PullAudioOptions{
				Profile: 0,
			}
		}
		cloudPlayerClientReq.Player.AudioOptions = &audioOptions
		cloudPlayerClientReq.Player.VideoOptions = clientStartReq.VideoOptions
	}

	// Start Cloud Player
	response, err := s.HandleStartPullReq(cloudPlayerClientReq, clientStartReq.Region, clientStartReq.StreamOriginIp, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start Cloud Player: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// StopPull handles the stopping of an RTMP push.
// It processes the request to stop pushing the media stream to the specified RTMP URL.
func (s *RtmpService) StopPull(c *gin.Context) {
	// Verify the client's request. If binding fails, returns an HTTP 400 error with the specific binding error message.
	var clientStopReq ClientStopPullRequest
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
	response, err := s.HandleStopPullReq(clientStopReq.PlayerId, clientStopReq.Region, c.GetHeader("X-Request-ID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stop Cloud Player: " + err.Error()})
		return
	}

	// Return the wrapped Agora response
	c.Data(http.StatusOK, "application/json", response)
}

// UpdateConverter handles updating the transcoding options for the RTMP push.
// It processes the request to update the transcoding configuration for the media stream.
func (s *RtmpService) UpdatePlayer(c *gin.Context) {}

// GetPullList returns a list of the current cloud players.
// It processes the request to get the current list of the media stream pull operations.
func (s *RtmpService) GetPullList(c *gin.Context) {}

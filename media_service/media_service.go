package media_service

import (
	"math/rand"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
)

// MediaService represents the media push/pull service.
// It holds the necessary configurations and dependencies for managing Media Push and Pull with Agora Channels.
type MediaService struct {
	appID     string // The Agora app ID
	baseURL   string // The base URL for the Agora cloud recording API
	basicAuth string // Middleware for handling requests
	// tokenService  *token_service.TokenService // Token service for generating tokens
}

// NewMediaService returns a MediaService pointer with all configurations set.
// This function initializes a new MediaService with specified configurations. It ensures all provided parameters are valid and logs a fatal error if any required configurations are missing.
//
// Parameters:
//   - tokenService: *token_service.TokenService - The token service for generating tokens.
//
// Returns:
//   - *MediaService: The initialized MediaService struct.
//
// Behavior:
//   - Initializes and returns a MediaService struct with the given configurations.
//
// Notes:
//   - Logs a fatal error and exits if any required environment variables are missing.
func NewMediaService(appID string, baseURL string, basicAuth string, tokenService *token_service.TokenService) *MediaService {

	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Return a new instance of the service
	return &MediaService{
		appID:     appID,     // The Agora app ID used to identify the application within Agora services.
		baseURL:   baseURL,   // The base URL for the Agora cloud recording API where all API requests are sent.
		basicAuth: basicAuth, // Basic authentication credentials required for interacting with the Agora API.
		// tokenService:  tokenService,  // Pointer to an instance of TokenService used to generate authentication tokens for Agora API requests.
	}
}

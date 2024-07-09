package token_service

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// TokenService represents the main application token service.
// It holds the necessary configurations and dependencies for managing tokens.
type TokenService struct {
	Server         *http.Server   // The HTTP server for the application
	Sigint         chan os.Signal // Channel to handle OS signals, such as Ctrl+C
	appID          string         // The Agora app ID
	appCertificate string         // The Agora app certificate
}

// TokenRequest is a struct representing the JSON payload structure for token generation requests.
// It contains fields necessary for generating different types of tokens (RTC, RTM, or chat) based on the "TokenType".
// The "Channel", "RtcRole", "Uid", and "ExpirationSeconds" fields are used for specific token types.
//
// TokenType options: "rtc" for RTC token, "rtm" for RTM token, and "chat" for chat token.
type TokenRequest struct {
	TokenType         string `json:"tokenType"`         // The token type: "rtc", "rtm", or "chat"
	Channel           string `json:"channel,omitempty"` // The channel name (used for RTC and RTM tokens)
	RtcRole           string `json:"role,omitempty"`    // The role of the user for RTC tokens (publisher or subscriber)
	Uid               string `json:"uid,omitempty"`     // The user ID or account (used for RTC, RTM, and some chat tokens)
	ExpirationSeconds int    `json:"expire,omitempty"`  // The token expiration time in seconds (used for all token types)
}

// NewTokenService initializes and returns a TokenService pointer with all configurations set.
// It loads environment variables, validates their presence, and initializes the TokenService struct.
//
// Parameters:
//   - appIDEnv: string - The environment variable for the Agora app ID.
//   - appCertEnv: string - The environment variable for the Agora app certificate.
//   - corsAllowOrigin: string - The allowed origin for CORS.
//
// Returns:
//   - *TokenService: The initialized TokenService struct.
//
// Behavior:
//   - Initializes and returns a TokenService struct with the provided configurations.
//
// Notes:
//   - The necessary environment variables should be set before initializing the TokenService.
func NewTokenService(appIDEnv string, appCertEnv string) *TokenService {
	return &TokenService{
		appID:          appIDEnv,
		appCertificate: appCertEnv,
	}
}

// RegisterRoutes registers the routes for the TokenService.
// It sets up the API endpoints and applies necessary middleware for request handling.
//
// Parameters:
//   - r: *gin.Engine - The Gin engine instance to register the routes with.
//
// Behavior:
//   - Creates an API group for token routes.
//   - Applies middleware for NoCache and CORS.
//   - Registers routes for getting a new token.
//
// Notes:
//   - This function organizes the API routes and ensures that requests are handled with appropriate middleware.
func (s *TokenService) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/token")
	api.POST("/getNew", s.GetToken)
}

// GetToken handles the HTTP request to generate a token based on the provided TokenRequest.
// It forwards the HTTP response writer and request from the provided *gin.Context
// to the HandleGetToken method for token generation and response sending.
//
// Parameters:
//   - c: *gin.Context - The Gin context representing the HTTP request and response.
//
// Behavior:
//   - Parses the request body into a TokenRequest struct.
//   - Forwards the HTTP response writer and request to the HandleGetToken method for processing.
//
// Notes:
//   - This function acts as an intermediary to invoke the HandleGetToken method.
//   - It handles parsing the request body and passing it to the token generation logic.
//
// Example usage:
//
//	router.POST("/getNew", TokenService.GetToken)
func (s *TokenService) GetToken(c *gin.Context) {
	var req = c.Request
	var respWriter = c.Writer
	var tokenReq TokenRequest
	// Parse the request body into a TokenRequest struct
	err := json.NewDecoder(req.Body).Decode(&tokenReq)
	if err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}
	s.HandleGetToken(tokenReq, respWriter)
}

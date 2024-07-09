package token_service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// NewTestTokenService creates a TokenService instance for testing
func NewTestTokenService() *TokenService {
	// Mock credentials for testing
	return &TokenService{
		appID:          "6ce46dd303d54056a52f9a34c13c547e",
		appCertificate: "77be7e16f7482cef9fe796205b85831e",
	}
}

func TestGenRtcToken(t *testing.T) {
	service := NewTestTokenService()

	tests := []struct {
		name    string
		request TokenRequest
		wantErr bool
	}{
		{
			name: "Valid RTC token request",
			request: TokenRequest{
				TokenType: "rtc",
				Channel:   "test-channel",
				Uid:       "1234",
				RtcRole:   "publisher",
			},
			wantErr: false,
		},
		{
			name: "Missing channel",
			request: TokenRequest{
				TokenType: "rtc",
				Uid:       "1234",
				RtcRole:   "publisher",
			},
			wantErr: true,
		},
		{
			name: "Missing UID",
			request: TokenRequest{
				TokenType: "rtc",
				Channel:   "test-channel",
				RtcRole:   "publisher",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenRtcToken(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenRtcToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("GenRtcToken() returned empty token")
			}
			t.Logf("Generated token: %s", token)
		})
	}
}

func TestGenRtmToken(t *testing.T) {
	service := NewTestTokenService()

	tests := []struct {
		name    string
		request TokenRequest
		wantErr bool
	}{
		{
			name: "Valid RTM token request",
			request: TokenRequest{
				TokenType: "rtm",
				Uid:       "test-user",
			},
			wantErr: false,
		},
		{
			name: "Missing UID",
			request: TokenRequest{
				TokenType: "rtm",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenRtmToken(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenRtmToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("GenRtmToken() returned empty token")
			}
		})
	}
}

func TestGenChatToken(t *testing.T) {
	service := NewTestTokenService()

	tests := []struct {
		name    string
		request TokenRequest
		wantErr bool
	}{
		{
			name: "Valid chat app token request",
			request: TokenRequest{
				TokenType: "chat",
			},
			wantErr: false,
		},
		{
			name: "Valid chat user token request",
			request: TokenRequest{
				TokenType: "chat",
				Uid:       "test-user",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.GenChatToken(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenChatToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Errorf("GenChatToken() returned empty token")
			}
			t.Logf("Generated token: %s", token)
		})
	}
}

func TestGetToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := NewTestTokenService()

	tests := []struct {
		name           string
		requestBody    string
		wantStatusCode int
	}{
		{
			name:           "Valid RTC token request",
			requestBody:    `{"tokenType": "rtc", "channel": "test-channel", "uid": "1234", "role": "publisher"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Valid RTM token request",
			requestBody:    `{"tokenType": "rtm", "uid": "test-user"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Valid chat token request",
			requestBody:    `{"tokenType": "chat", "uid": "test-user"}`,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "Invalid token type",
			requestBody:    `{"tokenType": "invalid", "channel": "test-channel", "uid": "1234"}`,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Missing required fields",
			requestBody:    `{"tokenType": "rtc"}`,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/token/getNew", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(rr)
			c.Request = req

			service.GetToken(c)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.wantStatusCode)
			}

			if tt.wantStatusCode == http.StatusOK {
				var response struct {
					Token string `json:"token"`
				}
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				if err != nil {
					t.Errorf("Error unmarshaling response: %v", err)
				}
				if response.Token == "" {
					t.Errorf("Expected non-empty token in response")
				}
			}
		})
	}
}

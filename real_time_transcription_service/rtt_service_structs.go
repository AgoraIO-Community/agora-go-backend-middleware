package real_time_transcription_service

import "github.com/AgoraIO-Community/agora-go-backend-middleware/cloud_recording_service"

// ClientStartRTTRequest represents the JSON payload structure sent by the client to start real time transcription.
// It includes the instance ID ,
type ClientStartRTTRequest struct {
	ChannelName        string           `json:"channelName"`            // The name of the channel to transcribe
	Languages          []string         `json:"languages"`              // The language(s) to transcribe
	SubscribeAudioUIDs []string         `json:"subscribeAudioUids"`     // A list of UID's to subscribe to in the channel. Max 3
	CryptionMode       *string          `json:"cryptionMode,omitempty"` // Cryption mode (Optional, if need cryption for audio and caption text)
	Secret             *string          `json:"secret,omitempty"`       // Cryption secret (Optional, if need decryption for audio and caption text)
	Salt               *string          `json:"salt,omitempty"`         // Cryption salt (Optional, if need decryption for audio and caption text)forceTranslateInterval.languages
	MaxIdleTime        *int             `json:"maxIdleTime,omitempty"`  // The default is 30 seconds. The unit is seconds, Range 5 seconds - 2592000 seconds (30 days)
	TranslateConfig    *TranslateConfig `json:"translateConfig,omitempty"`
	EnableStorage      *bool            `json:"enableStorage,omitempty"`      // Use to enable storage of captions
	EnableNTPtimestamp *bool            `json:"enableNTPtimestamp,omitempty"` // Use to enable subtitle sync
}

type ClientStartRTTV1Request struct {
	ChannelName        string    `json:"channelName"`                  // The name of the channel to transcribe
	ProfanityFilter    *bool     `json:"profanityFilter,omitempty"`    // Optional Text filter, default is false
	Destinations       *[]string `json:"destinations,omitempty"`       // List of output destination, if empty defaults to ["AgoraRTCDataStream"]
	MaxIdleTime        *int      `json:"maxIdleTime,omitempty"`        // The default is 30 seconds. The unit is seconds, Range 5 seconds - 2592000 seconds (30 days)
	EnableNTPtimestamp *bool     `json:"enableNTPtimestamp,omitempty"` // Use to enable subtitle sync
}

// AcquireBuilderTokenRequest defines the structure for a request to acquire a builder token for real time transcription
// It includes the instance ID set by the developer. Best practice is to use the channel name.
type AcquireBuilderTokenRequest struct {
	InstanceId string `json:"instanceId"`
}

type StartRTTRequest struct {
	Languages       []string         `json:"languages"`                 // The language(s) to transcribe
	MaxIdleTime     int              `json:"maxIdleTime"`               // If there is no audio stream in the channel for more than this time, the RTT Task will stop automatically.
	RTCConfig       RTCConfig        `json:"rtcConfig"`                 // The RTC settings for the audio and data bots
	CaptionConfig   *CaptionConfig   `json:"captionConfig,omitempty"`   // The cloud recording configuration
	TranslateConfig *TranslateConfig `json:"translateConfig,omitempty"` // The settings for real-time translation
}

type RTCConfig struct {
	ChannelName        string   `json:"channelName"`            // The name of the channel to transcribe
	SubBotUID          string   `json:"subBotUid"`              // The Uid used by the audio streaming bot to join the channel.
	SubBotToken        *string  `json:"subBotToken"`            // RTC token for the audio streaming bot
	PubBotUID          string   `json:"pubBotUid"`              // The uid used for Data streaming bot, used to stream text content after conversion.
	PubBotToken        *string  `json:"pubBotToken"`            // RTC token for the audio streaming bot
	SubscribeAudioUIDs []string `json:"subscribeAudioUids"`     // A list of UID's to subscribe to in the channel. Max 3
	CryptionMode       *string  `json:"cryptionMode,omitempty"` // Cryption mode (Optional, if need cryption for audio and caption text)
	Secret             *string  `json:"secret,omitempty"`       // Cryption secret (Optional, if need decryption for audio and caption text)
	Salt               *string  `json:"salt,omitempty"`         // Cryption salt (Optional, if need decryption for audio and caption text)forceTranslateInterval.languages
}

type CaptionConfig struct {
	Storage cloud_recording_service.StorageConfig `json:"storage"`
}

type TranslateConfig struct {
	ForceTranslateInterval int        `json:"forceTranslateInterval"`
	Languages              []Language `json:"languages"`
}

type Language struct {
	Source string   `json:"source"`
	Target []string `json:"target"`
}

// Timestampable is an interface that allows struct types to receive a timestamp.
// Implementing this interface ensures that a timestamp can be set on the object, primarily for auditing or tracking purposes.
type Timestampable interface {
	SetTimestamp(timestamp string)
}

// StartRTTResponse represents the response received from the Agora server after successfully starting a recording.
// It includes the identifiers of the recording session along with an optional timestamp.
type AcquireBuilderTokenResponse struct {
	TokenName  string  `json:"tokenName"`           // The value of the dynamic key builderToken
	CreateTs   int     `json:"createTs"`            // The Unix timestamp (seconds) when the builderToken was generated.
	InstanceId string  `json:"instanceId"`          // The instance ID set in the request body.
	Timestamp  *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *AcquireBuilderTokenResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// AgpraRTTResponse represents the response received from the Agora server after successfully starting a recording.
// It includes the identifiers of the recording session along with an optional timestamp.
type AgpraRTTResponse struct {
	CreateTs  int     `json:"createTs"`            // The Unix timestamp (seconds) when the builderToken was generated.
	Status    string  `json:"status"`              // The channel name for the recording session.
	TaskId    string  `json:"taskId"`              // a UUID (Universal Unique Identifier) generated by the Agora server to identify the real-time transcription task that has been created.
	Timestamp *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *AgpraRTTResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// StopRTTResponse represents the response received from the Agora server after successfully starting a recording.
// It includes the identifiers of the recording session along with an optional timestamp.
type StopRTTResponse struct {
	Timestamp *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *StopRTTResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

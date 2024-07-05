package real_time_transcription_service

import "github.com/AgoraIO-Community/agora-go-backend-middleware/cloud_recording_service"

// ClientStartRTTRequest represents the JSON payload structure sent by the client to start real time transcription.
// It includes the instance ID ,
type ClientStartRTTRequest struct {
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
	Audio         Audio          `json:"audio"`                   // Audio streaming bot settings
	Config        Config         `json:"config"`                  // Data Stream streaming Bot settings
	PrivateParams *PrivateParams `json:"privateParams,omitempty"` // Use to enable subtitle sync
}

type Audio struct {
	SubscribeSource string         `json:"subscribeSource"` // The current value is fixed: "AGORARTC"
	AgoraRtcConfig  AgoraRtcConfig `json:"agoraRtcConfig"`  // RTC Config for audio bot
}

type AgoraRtcConfig struct {
	ChannelName     string          `json:"channelName"`     // The name of the channel to transcribe
	UID             string          `json:"uid"`             // The Uid used by the audio streaming bot
	Token           string          `json:"token"`           // RTC token for the audio streaming bot
	ChannelType     string          `json:"channelType"`     // The current value is fixed: "LIVE_TYPE"
	SubscribeConfig SubscribeConfig `json:"subscribeConfig"` // Subscription settings
	MaxIdleTime     int             `json:"maxIdleTime"`     // If there is no audio stream in the channel for more than this time, the RTT Task will stop automatically.
}

type SubscribeConfig struct {
	SubscribeMode string `json:"subscribeMode"` // The current value is fixed: "CHANNEL_MODE"
}

type Config struct {
	Features        []string        `json:"features"` // The current value is fixed: ["RECOGNIZE"]
	RecognizeConfig RecognizeConfig `json:"recognizeConfig"`
}

type RecognizeConfig struct {
	Language        string `json:"language"`                  // The current version supports at most two, separated by commas.
	Model           string `json:"model"`                     // The current value is fixed: "Model"
	ProfanityFilter *bool  `json:"profanityFilter,omitempty"` // Optional. Text filter,
	Output          Output `json:"output"`                    // Output settings for conversion result
}

type Output struct {
	Destinations       []string           `json:"destinations"`       // Output target, AgoraRTCDataStream is required but Storage is optional.
	AgoraRTCDataStream AgoraRTCDataStream `json:"agoraRTCDataStream"` // Output data stream RTC settings
	CloudStorage       *[]CloudStorage    `json:"cloudStorage"`       // Optional storage setting, write transcriptions to cloud storage
}

type AgoraRTCDataStream struct {
	ChannelName string `json:"channelName"` // The target channel name of the conversion result output. (must be same as audio channel)
	UID         string `json:"uid"`         //  the uid used for streaming text content after conversion.
	Token       string `json:"token"`
}

type CloudStorage struct {
	Format        string                                `json:"format"` // The current value is fixed: "HLS"
	StorageConfig cloud_recording_service.StorageConfig `json:"storageConfig"`
}

type PrivateParams struct {
	EnableNTPtimestamp bool `json:"enableNTPtimestamp"`
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

// StartRTTResponse represents the response received from the Agora server after successfully starting a recording.
// It includes the identifiers of the recording session along with an optional timestamp.
type StartRTTResponse struct {
	CreateTs  int     `json:"createTs"`            // The Unix timestamp (seconds) when the builderToken was generated.
	Status    string  `json:"status"`              // The channel name for the recording session.
	TaskId    string  `json:"taskId"`              // a UUID (Universal Unique Identifier) generated by the Agora server to identify the real-time transcription task that has been created.
	Timestamp *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *StartRTTResponse) SetTimestamp(timestamp string) {
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

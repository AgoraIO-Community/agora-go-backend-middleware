package rtmp_service

// ClientStartRtmpRequest represents the JSON payload structure sent by the client to start an RTMP push.
// It includes configuration details for the RTMP converter, stream settings, and (Optional) transcoding options.
type ClientStartRtmpRequest struct {
	ConverterName      *string           `json:"converterName,omitempty"`      // (Optional) name for the RTMP converter
	RtcChannel         string            `json:"rtcChannel"`                   // The RTC channel name to push from
	StreamUrl          string            `json:"streamUrl"`                    // The RTMP server URL to push to
	StreamKey          string            `json:"streamKey"`                    // The stream key for the RTMP server
	Region             string            `json:"region"`                       // The region for the RTMP push service
	RegionHintIp       *string           `json:"regionHintIp"`                 // (Optional) IP hint for region selection
	UseTranscoding     bool              `json:"useTranscoding"`               // Whether to use transcoding for the push
	RtcStreamUid       *string           `json:"rtcStreamUid,omitempty"`       // (Optional) RTC stream UID to push
	AudioOptions       *PushAudioOptions `json:"audioOptions,omitempty"`       // (Optional) audio transcoding options
	VideoOptions       *PushVideoOptions `json:"videoOptions,omitempty"`       // (Optional) video transcoding options
	IdleTimeOut        *int              `json:"idleTimeOut,omitempty"`        // (Optional) idle timeout in seconds
	JitterBufferSizeMs *int              `json:"jitterBufferSizeMs,omitempty"` // (Optional) jitter buffer size in milliseconds
}

// ClientStartCloudPlayerRequest represents the JSON payload structure sent by the client to start an Cloud Player instance.
// It includes configuration details for the RTMP converter, stream settings, and (Optional) transcoding options.
type ClientStartCloudPlayerRequest struct {
	ChannelName    string            `json:"channelName"`              // The RTC channel name to push into
	StreamUrl      string            `json:"streamUrl"`                // The CDN/RTMP URL to pull from
	Region         string            `json:"region"`                   // The region for the cloud player service
	Uid            *string           `json:"uid,omitempty"`            // (Optional) RTC stream UID to use for the cloud player instance
	PlayerName     *string           `json:"name,omitempty"`           // (Optional) name for the cloud player instance
	StreamOriginIp *string           `json:"streamOriginIp,omitempty"` // (Optional) The IP address of the media stream's origin server.
	AudioOptions   *PullAudioOptions `json:"audioOptions,omitempty"`   // (Optional) audio transcoding configuration
	VideoOptions   *PullVideoOptions `json:"videoOptions,omitempty"`   // (Optional) video transcoding options
	IdleTimeOut    *int              `json:"idleTimeOut,omitempty"`    // (Optional) idle timeout in seconds
	PlayTs         *int              `json:"playTs,omitempty"`         // (Optional) Unix timestamp (in seconds) when the cloud player starts playing the online media stream.
	EncryptMode    *string           `json:"encryptMode,omitempty"`    // (Optional) Encryption mode
}

// ClientStopRtmpRequest represents the JSON payload structure for stopping an RTMP push.
// It contains the necessary identifiers to locate and terminate a specific RTMP push.
type ClientStopRtmpRequest struct {
	ConverterId string `json:"converterId"` // The ID of the RTMP converter to stop
	Region      string `json:"region"`      // The region where the RTMP push is running
}

// ClientStopPullRequest represents the JSON payload structure for stopping an RTMP push.
// It contains the necessary identifiers to locate and terminate a specific RTMP push.
type ClientStopPullRequest struct {
	PlayerId string `json:"playerId"` // The ID of the RTMP converter to stop
	Region   string `json:"region"`   // The region where the RTMP push is running
}

// ClientUpdateRtmpRequest represents the JSON payload structure for updating an ongoing RTMP push.
// It allows for modifications to certain parameters of an active RTMP push.
type ClientUpdateRtmpRequest struct {
	ConverterId        string            `json:"converterId"`                  // The ID of the RTMP converter to update
	Region             string            `json:"region"`                       // The region where the RTMP push is running
	StreamUrl          *string           `json:"streamUrl"`                    // The RTMP server URL to push to
	StreamKey          *string           `json:"streamKey"`                    // The stream key for the RTMP server
	RtcChannel         string            `json:"rtcChannel"`                   // The RTC channel name (in case of change)
	VideoOptions       *PushVideoOptions `json:"videoOptions,omitempty"`       // (Optional) updated video options
	JitterBufferSizeMs *int              `json:"jitterBufferSizeMs,omitempty"` // (Optional) updated jitter buffer size
	SequenceId         *int              `json:"sequenceId,omitempty"`         // (Optional) Agora server updates cloud player according to the latest sequence id
}

// RtmpPushRequest defines the structure for a request to start or update an RTMP push to the Agora service.
// It encapsulates the converter configuration for the RTMP push.
type RtmpPushRequest struct {
	Converter Converter `json:"converter"` // The converter configuration for the RTMP push instance
}

type CloudPlayerStartRequest struct {
	Player Player `json:"player"` // The player configuration for the clout player instance
}

// Converter represents the configuration for an RTMP converter.
// It includes settings for both transcoded and raw RTMP pushes.
type Converter struct {
	Name               *string           `json:"name,omitempty"`               // (Optional) name for the converter
	TranscodeOptions   *TranscodeOptions `json:"transcodeOptions,omitempty"`   // Options for transcoded push
	RawOptions         *RawOptions       `json:"rawOptions,omitempty"`         // Options for raw (non-transcoded) push
	RtmpUrl            *string           `json:"rtmpUrl,omitempty"`            // The RTMP URL to push to
	IdleTimeOut        *int              `json:"idleTimeOut,omitempty"`        // (Optional) idle timeout in seconds
	JitterBufferSizeMs *int              `json:"jitterBufferSizeMs,omitempty"` // (Optional) jitter buffer size in milliseconds
}

type Player struct {
	AudioOptions *PullAudioOptions `json:"audioOptions,omitempty"` // (Optional) audio transcoding configuration
	VideoOptions *PullVideoOptions `json:"videoOptions,omitempty"` // (Optional) video transcoding options
	StreamUrl    string            `json:"streamUrl"`              // The CDN/RTMP URL to pull from
	ChannelName  string            `json:"channelName"`            // The RTC channel name to push into
	Token        string            `json:"token"`                  // The access token
	Uid          string            `json:"uid"`                    // (Optional) RTC stream UID to use for the cloud player instance
	IdleTimeOut  *int              `json:"idleTimeOut,omitempty"`  // (Optional) idle timeout in seconds
	PlayTs       *int              `json:"playTs,omitempty"`       // (Optional) Unix timestamp (in seconds) when the cloud player starts playing the online media stream.
	EncryptMode  *string           `json:"encryptMode,omitempty"`  // (Optional) Encryption mode
	PlayerName   *string           `json:"name,omitempty"`         // (Optional) name for the cloud player instance
}

// RawOptions defines the parameters for a raw (non-transcoded) RTMP push.
type RawOptions struct {
	RtcChannel   string `json:"rtcChannel"`   // The RTC channel to push from
	RtcStreamUid string `json:"rtcStreamUid"` // The UID of the RTC stream to push
}

// TranscodeOptions defines the parameters for a transcoded RTMP push.
// It includes audio and video transcoding options.
type TranscodeOptions struct {
	RtcChannel   string            `json:"rtcChannel"`             // The RTC channel to push from
	AudioOptions *PushAudioOptions `json:"audioOptions,omitempty"` // Audio transcoding options
	VideoOptions *PushVideoOptions `json:"videoOptions,omitempty"` // Video transcoding options
}

// PushAudioOptions specifies the audio transcoding settings for an RTMP push.
type PushAudioOptions struct {
	CodecProfile  string `json:"codecProfile"`  // Audio codec profile
	SampleRate    int    `json:"sampleRate"`    // Audio sample rate in Hz
	Bitrate       int    `json:"bitrate"`       // Audio bitrate in Kbps
	AudioChannels int    `json:"audioChannels"` // Number of audio channels
}

// PushAudioOptions specifies the audio transcoding settings for an RTMP push.
type PullAudioOptions struct {
	Profile int `json:"profile"` // Sets the audio profile sample rate, bitrate, encoding mode, and the number of channels (0-5)
}

// PushVideoOptions specifies the video transcoding settings for an RTMP push.
// It includes layout, codec, and quality settings.
type PushVideoOptions struct {
	Canvas                     Canvas      `json:"canvas"`                               // Canvas dimensions for the video
	Layout                     []Layout    `json:"layout"`                               // Layout configuration for multiple streams
	Vertical                   *Vertical   `json:"vertical,omitempty"`                   // (Optional) vertical layout settings
	DefaultPlaceholderImageUrl *string     `json:"defaultPlaceholderImageUrl,omitempty"` // (Optional) default placeholder image URL
	Codec                      *string     `json:"codec,omitempty"`                      // (Optional) video codec
	CodecProfile               *string     `json:"codecProfile,omitempty"`               // (Optional) codec profile
	FrameRate                  *int        `json:"frameRate,omitempty"`                  // (Optional) frame rate
	Gop                        *int        `json:"gop,omitempty"`                        // (Optional) Group of Pictures (GOP) size
	Bitrate                    int         `json:"bitrate"`                              // Video bitrate in Kbps
	SeiOptions                 *SeiOptions `json:"seiOptions,omitempty"`                 // (Optional) Supplemental Enhancement Information (SEI) options
}

// PullVideoOptions specifies the video transcoding settings for the cloud player.
// It includes layout, codec, and quality settings.
type PullVideoOptions struct {
	Width               int     `json:"width"`               // Width of the canvas in pixels
	Height              int     `json:"height"`              // Height of the canvas in pixels
	WidthHeightAdaption *bool   `json:"widthHeightAdaption"` // Whether to enable horizontal and vertical screen adaptive mode (default: false)
	FrameRate           *int    `json:"frameRate,omitempty"` // (Optional) frame rate
	Bitrate             int     `json:"bitrate"`             // Video bitrate in Kbps
	Codec               string  `json:"codec"`               // Video codec (H264 or VP9)
	FillMode            *string `json:"fillMode,omitempty"`  // (Optional) The fill mode of the output video (fit/fill)
	Gop                 *int    `json:"gop,omitempty"`       // (Optional) Group of Pictures (GOP) size
}

// Canvas defines the dimensions of the video canvas for transcoding.
type Canvas struct {
	Width  int `json:"width"`  // Width of the canvas in pixels
	Height int `json:"height"` // Height of the canvas in pixels
}

// Layout defines the position and size of a single stream within the video canvas.
type Layout struct {
	RtcStreamUid        string `json:"rtcStreamUid"`                  // UID of the RTC stream
	Region              Region `json:"region"`                        // Position and size within the canvas
	FillMode            string `json:"fillMode,omitempty"`            // (Optional) fill mode for the video
	PlaceholderImageUrl string `json:"placeholderImageUrl,omitempty"` // (Optional) placeholder image URL
}

// Vertical defines settings for vertical video layout.
type Vertical struct {
	MaxResolutionUid int    `json:"maxResolutionUid"` // UID of the stream with max resolution
	FillMode         string `json:"fillMode"`         // Fill mode for vertical layout
}

// Region defines the position and size of a video stream within the canvas.
type Region struct {
	XPos   int `json:"xPos"`   // X-axis position
	YPos   int `json:"yPos"`   // Y-axis position
	ZIndex int `json:"zIndex"` // Z-index for layering
	Width  int `json:"width"`  // Width of the region
	Height int `json:"height"` // Height of the region
}

// SeiOptions defines options for Supplemental Enhancement Information (SEI) in the video stream.
type SeiOptions struct {
	Source Source `json:"source"` // SEI source options
	Sink   Sink   `json:"sink"`   // SEI sink options
}

// Source defines the SEI source configuration.
type Source struct {
	Metadata   bool        `json:"metadata"`             // Whether to include metadata
	Datastream bool        `json:"datastream"`           // Whether to include datastream
	Customized *Customized `json:"customized,omitempty"` // (Optional) custom SEI data
}

// Customized defines custom SEI data.
type Customized struct {
	PrefixForAgoraSei string `json:"prefixForAgoraSei"` // Prefix for Agora SEI
	Payload           string `json:"payload"`           // Custom SEI payload
}

// Sink defines the SEI sink configuration.
type Sink struct {
	Type int `json:"type"` // Type of SEI sink
}

// Timestampable is an interface that allows struct types to receive a timestamp.
// Implementing this interface ensures that a timestamp can be set on the object, primarily for auditing or tracking purposes.
type Timestampable interface {
	SetTimestamp(timestamp string)
}

// StartRtmpResponse represents the response received from the Agora server after successfully starting an RTMP push.
// It includes details about the created converter and an (Optional) timestamp.
type StartRtmpResponse struct {
	Converter ConverterResponse `json:"converter"`           // Details of the newly created RTMP converter
	Fields    string            `json:"fields"`              // Additional fields returned by the server
	Timestamp *string           `json:"timestamp,omitempty"` // (Optional) timestamp for when the RTMP push was started
}

// ConverterResponse contains the details of an RTMP converter returned by the Agora server.
type ConverterResponse struct {
	ConverterId string `json:"id"`       // Unique identifier for the converter
	CreateTs    int64  `json:"createTs"` // Timestamp of converter creation
	UpdateTs    int64  `json:"updateTs"` // Timestamp of last update
	State       string `json:"state"`    // Current state of the converter
}

// SetTimestamp implements the Timestampable interface for StartRtmpResponse.
func (s *StartRtmpResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// StopRtmpResponse represents the response received from the Agora server after stopping an RTMP push.
// It includes a status message and an (Optional) timestamp.
type StopRtmpResponse struct {
	Status    string  `json:"status"`              // Status of the stop operation
	Timestamp *string `json:"timestamp,omitempty"` // (Optional) timestamp for when the RTMP push was stopped
}

// SetTimestamp implements the Timestampable interface for StopRtmpResponse.
func (s *StopRtmpResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// StartCloudPlayerResponse represents the response received from the Agora server after successfully starting an RTMP push.
// It includes details about the created converter and an (Optional) timestamp.
type StartCloudPlayerResponse struct {
	Player    PlayerResponse `json:"player"`              // Details of the newly created cloud player
	Fields    string         `json:"fields"`              // Additional fields returned by the server
	Timestamp *string        `json:"timestamp,omitempty"` // (Optional) timestamp for when the RTMP push was started
}

// ConverterResponse contains the details of an RTMP converter returned by the Agora server.
type PlayerResponse struct {
	PlayerId string  `json:"id"`            // Unique identifier for the cloud player
	CreateTs int64   `json:"createTs"`      // Timestamp of converter creation
	Uid      *string `json:"uid,omitempty"` // (Optional) RTC stream UID to use for the cloud player instance
}

// SetTimestamp implements the Timestampable interface for StartCloudPlayerResponse.
func (s *StartCloudPlayerResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// StopCloudPlayerResponse represents the response received from the Agora server after stopping an RTMP push.
// It includes a status message and an (Optional) timestamp.
type StopCloudPlayerResponse struct {
	Status    string  `json:"status"`              // Status of the stop request
	Timestamp *string `json:"timestamp,omitempty"` // (Optional) timestamp for when the RTMP push was stopped
}

// SetTimestamp implements the Timestampable interface for StopCloudPlayerResponse.
func (s *StopCloudPlayerResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

package rtmp_service

type ClientStartRequest struct {
}

// Agora Media Push Request structs
type MediaPushRequest struct {
	Converter Converter `json:"converter"`
}

type Converter struct {
	Name               *string           `json:"name,omitempty"`
	TranscodeOptions   *TranscodeOptions `json:"transcodeOptions,omitempty"`
	RawOptions         *RawOptions       `json:"rawOptions,omitempty"`
	RtmpUrl            string            `json:"rtmpUrl"`
	IdleTimeOut        *int              `json:"idleTimeOut,omitempty"`
	JitterBufferSizeMs *int              `json:"jitterBufferSizeMs,omitempty"`
}

type RawOptions struct {
	RtcChannel   string `json:"rtcChannel"`
	Token        string `json:"token"`
	RtcStreamUid string `json:"rtcStreamUid"`
}

type TranscodeOptions struct {
	RtcChannel   string        `json:"rtcChannel"`
	Token        string        `json:"token"`
	AudioOptions *AudioOptions `json:"audioOptions,omitempty"`
	VideoOptions *VideoOptions `json:"videoOptions,omitempty"`
}

type AudioOptions struct {
	CodecProfile  string `json:"codecProfile"`
	SampleRate    int    `json:"sampleRate"`
	Bitrate       int    `json:"bitrate"`
	AudioChannels int    `json:"audioChannels"`
}

type VideoOptions struct {
	Canvas                     Canvas      `json:"canvas"`
	Layout                     []Layout    `json:"layout"`
	Vertical                   *Vertical   `json:"vertical,omitempty"`
	DefaultPlaceholderImageUrl *string     `json:"defaultPlaceholderImageUrl,omitempty"`
	Codec                      *string     `json:"codec,omitempty"`
	CodecProfile               *string     `json:"codecProfile,omitempty"`
	FrameRate                  *int        `json:"frameRate,omitempty"`
	Gop                        *int        `json:"gop,omitempty"`
	Bitrate                    int         `json:"bitrate"`
	SeiOptions                 *SeiOptions `json:"seiOptions,omitempty"`
}

type Canvas struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Layout struct {
	RtcStreamUid        string `json:"rtcStreamUid"`
	Region              Region `json:"region"`
	FillMode            string `json:"fillMode,omitempty"`
	PlaceholderImageUrl string `json:"placeholderImageUrl,omitempty"`
}

type Vertical struct {
	MaxResolutionUid int    `json:"maxResolutionUid"`
	FillMode         string `json:"fillMode"`
}

type Region struct {
	XPos   int `json:"xPos"`
	YPos   int `json:"yPos"`
	ZIndex int `json:"zIndex"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type SeiOptions struct {
	Source Source `json:"source"`
	Sink   Sink   `json:"sink"`
}

type Source struct {
	Metadata   bool        `json:"metadata"`
	Datastream bool        `json:"datastream"`
	Customized *Customized `json:"customized,omitempty"`
}

type Customized struct {
	PrefixForAgoraSei string `json:"prefixForAgoraSei"`
	Payload           string `json:"payload"`
}

type Sink struct {
	Type int `json:"type"`
}

// Agora Media Push Response structs

type MediaPushResponse struct {
	Converter ConverterResponse `json:"converter"`
	Fields    string            `json:"fields"`
}

type ConverterResponse struct {
	ID       string `json:"id"`
	CreateTs int64  `json:"createTs"`
	UpdateTs int64  `json:"updateTs"`
	State    string `json:"state"`
}

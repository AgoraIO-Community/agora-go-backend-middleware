package cloud_recording_service

import "encoding/json"

// ClientStartRecordingRequest represents the JSON payload structure sent by the client to start a cloud recording.
type ClientStartRecordingRequest struct {
	ChannelName        string           `json:"channelName"`
	SceneMode          *string          `json:"sceneMode,omitempty"`
	RecordingMode      *string          `json:"recordingMode,omitempty"`
	ExcludeResourceIds *[]string        `json:"excludeResourceIds,omitempty"`
	RecordingConfig    *RecordingConfig `json:"recordingConfig,omitempty"`
}

type ClientStopRecordingRequest struct {
	Cname         string  `json:"cname"` // The channel name for the cloud recording
	Uid           string  `json:"uid"`   // The UID for the cloud recording session
	ResourceId    string  `json:"resourceId"`
	RecordingId   string  `json:"recordingId"`
	RecordingMode *string `json:"recordingMode,omitempty"`
	AsyncStop     *bool   `json:"async_stop,omitempty"`
}

type ClientUpdateRecordingRequest struct {
	Cname         string          `json:"cname"`         // The channel name for the cloud recording
	Uid           string          `json:"uid"`           // The UID for the cloud recording session
	ClientRequest RecordingConfig `json:"clientRequest"` // The updated recording configuration request for the cloud recording session
}

// AcquireResourceRequest represents the JSON payload structure for acquiring a cloud recording resource.
// It contains the channel name and UID necessary for resource acquisition.
type AcquireResourceRequest struct {
	Cname         string               `json:"cname"`         // The channel name for the cloud recording
	Uid           string               `json:"uid"`           // The UID for the cloud recording session
	ClientRequest *AquireClientRequest `json:"clientRequest"` // The client request, an empty object
}

// StartRecordingRequest represents the JSON payload structure for starting a cloud recording.
// It includes the channel name, UID, and the client request configuration.
type StartRecordingRequest struct {
	Cname         string        `json:"cname"`         // The channel name for the cloud recording
	Uid           string        `json:"uid"`           // The UID for the cloud recording session
	ClientRequest ClientRequest `json:"clientRequest"` // The client request configuration for the cloud recording
}

type StopRecordingRequest struct {
	Cname         string            `json:"cname"` // The channel name for the cloud recording
	Uid           string            `json:"uid"`   // The UID for the cloud recording session
	ResourceId    string            `json:"resourceId"`
	ClientRequest StopClientRequest `json:"clientRequest"` // The client request to stop the cloud recording
}

type StopClientRequest struct {
	AsyncStop *bool `json:"async_stop,omitempty"`
}

// AquireClientRequest represents the client request configuration for starting a cloud recording.
type AquireClientRequest struct {
	Scene               int           `json:"scene,omitempty"`
	ResourceExpiredHour int           `json:"resourceExpiredHour,omitempty"`
	StartParameter      ClientRequest `json:"startParameter,omitempty"`
	ExcludeResourceIds  *[]string     `json:"excludeResourceIds,omitempty"`
}

// Server Response from Agora after successful Start
type StartRecordingResponse struct {
	Cname      string `json:"cname"`
	Uid        string `json:"uid"`
	ResourceId string `json:"resourceId"`
	Sid        string `json:"sid"`
}

// StopRecordingResponse main struct for the recording response
type StopRecordingResponse struct {
	ResourceId     *string        `json:"resourceId"`
	Sid            *string        `json:"sid"`
	ServerResponse ServerResponse `json:"serverResponse,omitempty"` // Use RawMessage to defer unmarshaling
	Cname          *string        `json:"cname"`
	Uid            *string        `json:"uid"`
	Timestamp      *string        `json:"timestamp,omitempty"`
}

// ServerResponse all responsese
type ServerResponse struct {
	ExtensionServiceState   *ExtensionServiceState `json:"extensionServiceState,omitempty"`
	UploadingStatusResponse *string                `json:"uploadingStatus,omitempty"`
	FileListMode            *string                `json:"fileListMode,omitempty"` // values: string or json
	FileList                *json.RawMessage       `json:"fileList,omitempty"`     // []FileDetail | []FileListEntry
}

// ServerResponse: Web page recording (scenario 1)
type ExtensionServiceState struct {
	PlayloadStop *PlayloadStop `json:"playload"`
	ServiceName  string        `json:"serviceName"`
}

type PlayloadStop struct {
	UploadingStatus *string       `json:"uploadingStatus,omitempty"`
	FileList        *[]FileDetail `json:"fileList,omitempty"`
	OnHold          *bool         `json:"onhold"`
	State           *string       `json:"state"`
}

// Use as part of two ServerResponses
// - Web page recording (part of Playload)
// - []FileDetail (individual recording)
type FileDetail struct {
	Filename       string `json:"filename"`
	SliceStartTime int64  `json:"sliceStartTime"`
}

// []FileListEntry for handling fileList-json
type FileListEntry struct {
	FileName       string `json:"fileName"`
	TrackType      string `json:"trackType"`
	Uid            string `json:"uid"`
	MixedAllUser   bool   `json:"mixedAllUser"`
	IsPlayable     bool   `json:"isPlayable"`
	SliceStartTime int64  `json:"sliceStartTime"`
}

// ClientRequestcontains the detailed parameters for starting or updating a recording.
// It includes the token, storage configuration, and recording configuration.
type ClientRequest struct {
	Token                  string                  `json:"token,omitempty"` // The token for the cloud recording session
	StorageConfig          StorageConfig           `json:"storageConfig"`   // The storage configuration for the cloud recording
	RecordingConfig        RecordingConfig         `json:"recordingConfig"` // The recording configuration for the cloud recording
	RecordingFileConfig    *RecordingFileConfig    `json:"recordingFileConfig,omitempty"`
	SnapshotConfig         *SnapshotConfig         `json:"snapshotConfig,omitempty"` // Snapshot configuration
	ExtensionServiceConfig *ExtensionServiceConfig `json:"extensionServiceConfig,omitempty"`
	AppsCollection         *AppsCollection         `json:"appsCollection,omitempty"`
	TranscodeOptions       *TranscodeOptions       `json:"transcodeOptions,omitempty"`
}

// StorageConfig represents the storage configuration for cloud recording.
// It includes the secret key, vendor, region, bucket, and access key for storage.
type StorageConfig struct {
	Vendor          int              `json:"vendor"`                   // The storage vendor identifier
	Region          int              `json:"region"`                   // The storage region identifier
	Bucket          string           `json:"bucket"`                   // The storage bucket name
	AccessKey       string           `json:"accessKey"`                // The access key for storage authentication
	SecretKey       string           `json:"secretKey"`                // The secret key for storage authentication
	FileNamePrefix  *[]string        `json:"fileNamePrefix,omitempty"` // Array of folder names ["directory1","directory2"] => "directory1/directory2/" => directory1/directory2/xxx.m3u8
	ExtensionParams *ExtensionParams `json:"extensionParams,omitempty"`
}

// ExtensionParams represents additional parameters for storage configuration.
type ExtensionParams struct {
	SSE *string `json:"sse,omitempty"`
	Tag *string `json:"tag,omitempty"`
}

// RecordingConfig represents the recording configuration for cloud recording.
type RecordingConfig struct {
	ChannelType          int                `json:"channelType"`
	DecryptionMode       *int               `json:"decryptionMode,omitempty"`
	Secret               *string            `json:"secret,omitempty"`
	Salt                 *string            `json:"salt,omitempty"`
	MaxIdleTime          *int               `json:"maxIdleTime,omitempty"`
	StreamTypes          *int               `json:"streamTypes,omitempty"`
	VideoStreamType      *int               `json:"videoStreamType,omitempty"`
	SubscribeAudioUids   *[]string          `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUids *[]string          `json:"unsubscribeAudioUids,omitempty"`
	SubscribeVideoUids   *[]string          `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUids *[]string          `json:"unsubscribeVideoUids,omitempty"`
	SubscribeUidGroup    *int               `json:"subscribeUidGroup,omitempty"`
	StreamMode           *string            `json:"streamMode,omitempty"` // "individual", "composite", or "web"
	AudioProfile         *int               `json:"audioProfile,omitempty"`
	TranscodingConfig    *TranscodingConfig `json:"transcodingConfig,omitempty"`
}

// TranscodingConfig represents the transcoding configuration for cloud recording.
type TranscodingConfig struct {
	Width                      *int                `json:"width,omitempty"`
	Height                     *int                `json:"height,omitempty"`
	Fps                        *int                `json:"fps,omitempty"`
	Bitrate                    *int                `json:"bitrate,omitempty"`
	MaxResolutionUid           *string             `json:"maxResolutionUid,omitempty"`
	MixedVideoLayout           *int                `json:"mixedVideoLayout,omitempty"`
	BackgroundColor            *string             `json:"backgroundColor,omitempty"`
	BackgroundImage            *string             `json:"backgroundImage,omitempty"`
	DefaultUserBackgroundImage *string             `json:"defaultUserBackgroundImage,omitempty"`
	LayoutConfig               *[]LayoutConfig     `json:"layoutConfig,omitempty"`
	BackgroundConfig           *[]BackgroundConfig `json:"backgroundConfig,omitempty"`
}

// LayoutConfig represents the layout configuration for transcoding.
type LayoutConfig struct {
	Uid        string `json:"uid"`
	XAxis      int    `json:"x_axis"`
	YAxis      int    `json:"y_axis"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Alpha      int    `json:"alpha"`
	RenderMode int    `json:"render_mode"`
}

// BackgroundConfig represents the background configuration for transcoding.
type BackgroundConfig struct {
	Uid        string `json:"uid"`
	ImageURL   string `json:"image_url"`
	RenderMode int    `json:"render_mode"`
}

// RecordingFileConfig represents the recording file configuration.
type RecordingFileConfig struct {
	AVFileType []string `json:"avFileType,omitempty"`
}

// SnapshotConfig represents the snapshot configuration.
type SnapshotConfig struct {
	CaptureInterval int      `json:"captureInterval,omitempty"`
	FileType        []string `json:"fileType,omitempty"`
}

// ExtensionServiceConfig represents the extension service configuration.
type ExtensionServiceConfig struct {
	ErrorHandlePolicy string             `json:"errorHandlePolicy,omitempty"`
	ExtensionServices []ExtensionService `json:"extensionServices,omitempty"`
}

// ExtensionService represents a single extension service.
type ExtensionService struct {
	ServiceName       string       `json:"serviceName"`
	ErrorHandlePolicy *string      `json:"errorHandlePolicy,omitempty"`
	ServiceParam      ServiceParam `json:"serviceParam"`
}

// ServiceParam represents the parameters for an extension service.
type ServiceParam struct {
	URL              string `json:"url"`
	AudioProfile     *int   `json:"audioProfile,omitempty"`
	VideoWidth       *int   `json:"videoWidth,omitempty"`
	VideoHeight      *int   `json:"videoHeight,omitempty"`
	MaxRecordingHour *int   `json:"maxRecordingHour,omitempty"`
	VideoBitrate     *int   `json:"videoBitrate,omitempty"`
	VideoFps         *int   `json:"videoFps,omitempty"`
	Mobile           *bool  `json:"mobile,omitempty"`
	MaxVideoDuration *int   `json:"maxVideoDuration,omitempty"`
	OnHold           *bool  `json:"onhold,omitempty"`
	ReadyTimeout     *int   `json:"readyTimeout,omitempty"`
}

// AppsCollection represents the collection of apps.
type AppsCollection struct {
	CombinationPolicy *string `json:"combinationPolicy,omitempty"`
}

// TranscodeOptions represents the transcode options.
type TranscodeOptions struct {
	TransConfig *TransConfig `json:"transConfig,omitempty"`
	Container   *Container   `json:"container,omitempty"`
	Audio       *Audio       `json:"audio,omitempty"`
}

// TransConfig represents the transcode configuration.
type TransConfig struct {
	TransMode *string `json:"transMode,omitempty"`
}

// Container represents the container configuration.
type Container struct {
	Format *string `json:"format,omitempty"`
}

// Audio represents the audio configuration.
type Audio struct {
	SampleRate *string `json:"sampleRate,omitempty"`
	Bitrate    *string `json:"bitrate,omitempty"`
	Channels   *string `json:"channels,omitempty"`
}

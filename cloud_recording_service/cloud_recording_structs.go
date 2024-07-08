package cloud_recording_service

import "encoding/json"

// ClientStartRecordingRequest represents the JSON payload structure sent by the client to start a cloud recording.
// It includes channel name, optional scene mode, recording mode, excluded resource IDs, and recording configuration details.
type ClientStartRecordingRequest struct {
	ChannelName        string           `json:"channelName"`                  // The name of the channel to record.
	SceneMode          *string          `json:"sceneMode,omitempty"`          // The recording scene type.
	RecordingMode      *string          `json:"recordingMode,omitempty"`      // The recording mode (indvidual, mix, web).
	ExcludeResourceIds *[]string        `json:"excludeResourceIds,omitempty"` // UID's to other recording or rtt services in the channel.
	RecordingConfig    *RecordingConfig `json:"recordingConfig,omitempty"`    // The configuration to use for the new cloud recording session.
}

// ClientUpdateSubscriptionRequest represents the JSON payload structure sent by the client to update a cloud recording's subscription.
// It includes identifiers and a nested update configuration specific to the recording session.
type ClientUpdateSubscriptionRequest struct {
	Cname         string                          `json:"cname"`                   // The name of the channel being recorded.
	Uid           string                          `json:"uid"`                     // The UID for the existing cloud recording session.
	ResourceId    string                          `json:"resourceId"`              // The ResourceId for the existing cloud recording session.
	Sid           string                          `json:"sid"`                     // The Sid for the existing cloud recording session.
	RecordingMode *string                         `json:"recordingMode,omitempty"` // The recording mode (indvidual, mix, web).
	UpdateConfig  UpdateSubscriptionClientRequest `json:"recordingConfig"`         // The updated recording configuration for the given cloud recording session.
}

// ClientUpdateLayoutRequest represents the JSON payload for updating the layout of a cloud recording session.
// It includes channel and session identifiers and layout update configurations.
type ClientUpdateLayoutRequest struct {
	Cname         string                    `json:"cname"`                   // The name of the channel being recorded.
	Uid           string                    `json:"uid"`                     // The UID for the existing cloud recording session.
	ResourceId    string                    `json:"resourceId"`              // The ResourceId for the existing cloud recording session.
	Sid           string                    `json:"sid"`                     // The Sid for the existing cloud recording session.
	RecordingMode *string                   `json:"recordingMode,omitempty"` // The recording mode (indvidual, mix, web).
	UpdateConfig  UpdateLayoutClientRequest `json:"recordingConfig"`         // The updated layout configuration for the given cloud recording session.
}

// ClientStopRecordingRequest represents the JSON payload structure for requesting the stop of a cloud recording.
// It contains identifiers necessary to identify the recording session to be stopped.
type ClientStopRecordingRequest struct {
	Cname         string  `json:"cname"`                   // The name of the channel being recorded.
	Uid           string  `json:"uid"`                     // The UID for an existing cloud recording session.
	ResourceId    string  `json:"resourceId"`              // The ResourceId for the existing cloud recording session.
	Sid           string  `json:"sid"`                     // The Sid for the existing cloud recording session.
	RecordingMode *string `json:"recordingMode,omitempty"` // The recording mode (indvidual, mix, web).
	AsyncStop     *bool   `json:"async_stop,omitempty"`    // Stop immediately or asynchronously.
}

// AcquireResourceRequest defines the structure for a request to acquire resources for cloud recording.
// It includes the channel name, UID, and an optional client request with additional parameters.
type AcquireResourceRequest struct {
	Cname         string               `json:"cname"`         // The channel name for the cloud recording.
	Uid           string               `json:"uid"`           // The UID for the cloud recording session.
	ClientRequest *AquireClientRequest `json:"clientRequest"` // Additional client-specific request parameters.
}

// StartRecordingRequest defines the structure for a request to start a cloud recording.
// It encapsulates the channel and session identifiers along with client-specific recording configurations.
type StartRecordingRequest struct {
	Cname         string        `json:"cname"`         // The channel name for the cloud recording.
	Uid           string        `json:"uid"`           // The UID for the cloud recording session.
	ClientRequest ClientRequest `json:"clientRequest"` // Configuration parameters for the recording.
}

// StopRecordingRequest defines the structure for a request to stop a cloud recording.
// It includes the channel and session identifiers and client-specific parameters for stopping the recording.
type StopRecordingRequest struct {
	Cname         string            `json:"cname"`         // The channel name for the cloud recording.
	Uid           string            `json:"uid"`           // The UID for the cloud recording session.
	ResourceId    string            `json:"resourceId"`    // The ResourceId for the existing cloud recording session.
	ClientRequest StopClientRequest `json:"clientRequest"` // Client-specific parameters to control the stop process.
}

// StopClientRequest defines additional parameters that can be specified when stopping a cloud recording.
type StopClientRequest struct {
	AsyncStop *bool `json:"async_stop,omitempty"` // Indicates whether the stop should be performed asynchronously.
}

// UpdateSubscriptionRequest defines the structure for a request to update a subscription for a cloud recording.
// It includes the channel and session identifiers and client-specific parameters for the subscription update.
type UpdateSubscriptionRequest struct {
	Cname         string                          `json:"cname"`         // The channel name for the cloud recording.
	Uid           string                          `json:"uid"`           // The UID for the cloud recording session.
	ClientRequest UpdateSubscriptionClientRequest `json:"clientRequest"` // Client-specific parameters for updating the subscription.
}

// UpdateSubscriptionClientRequest encapsulates the various configurations that can be modified in a subscription update request.
type UpdateSubscriptionClientRequest struct {
	StreamSubscribe    *StreamSubscribe    `json:"streamSubscribe,omitempty"`
	WebRecordingConfig *WebRecordingConfig `json:"webRecordingConfig,omitempty"`
	RTMPPublishConfig  *RTMPPublishConfig  `json:"rtmpPublishConfig,omitempty"`
}

// IsValid checks if the UpdateSubscriptionClientRequest is valid by ensuring only one of its fields is non-nil.
func (ucr *UpdateSubscriptionClientRequest) IsValid() bool {
	count := 0
	if ucr.StreamSubscribe != nil {
		count++
	}
	if ucr.WebRecordingConfig != nil {
		count++
	}
	if ucr.RTMPPublishConfig != nil {
		count++
	}
	return count == 1 // Ensures exactly one of the fields should be non-nil.
}

// StreamSubscribe defines the UID lists for subscribing to or unsubscribing from audio and video streams.
type StreamSubscribe struct {
	AudioUidList *AudioUidList `json:"audioUidList,omitempty"`
	VideoUidList *VideoUidList `json:"videoUidList,omitempty"`
}

// AudioUidList defines the UIDs for subscribing to or unsubscribing from audio streams.
type AudioUidList struct {
	SubscribeAudioUids   *[]string `json:"subscribeAudioUids,omitempty"`
	UnsubscribeAudioUids *[]string `json:"unsubscribeAudioUids,omitempty"`
}

// VideoUidList defines the UIDs for subscribing to or unsubscribing from video streams.
type VideoUidList struct {
	SubscribeVideoUids   *[]string `json:"subscribeVideoUids,omitempty"`
	UnsubscribeVideoUids *[]string `json:"unsunscribeVideoUids,omitempty"`
}

// WebRecordingConfig specifies the on-hold status of a web recording.
type WebRecordingConfig struct {
	Onhold bool `json:"onhold"`
}

// RTMPPublishConfig defines the configuration for RTMP publishing, including output URLs.
type RTMPPublishConfig struct {
	Outputs []Output `json:"outputs"` // List of output configurations for RTMP streams.
}

// Output defines the URL and additional parameters for an RTMP stream.
type Output struct {
	RTMPUrl string `json:"rtmpUrl"` // The URL for the RTMP stream.
}

// UpdateLayoutRequest defines the structure for a request to update the layout of a cloud recording.
// It includes the channel and session identifiers and layout-specific configurations.
type UpdateLayoutRequest struct {
	Cname         string                    `json:"cname"`         // The channel name for the cloud recording.
	Uid           string                    `json:"uid"`           // The UID for the cloud recording session.
	ClientRequest UpdateLayoutClientRequest `json:"clientRequest"` // Configuration parameters for the layout update.
}

// UpdateLayoutClientRequest encapsulates the parameters for configuring the video layout in a cloud recording.
type UpdateLayoutClientRequest struct {
	MaxResolutionUid           *string             `json:"maxResolutionUid,omitempty"`           // UID of the participant with the maximum resolution video.
	MixedVideoLayout           *int                `json:"mixedVideoLayout,omitempty"`           // The layout type for mixed video.
	BackgroundColor            *string             `json:"backgroundColor,omitempty"`            // Background color for the layout.
	BackgroundImage            *string             `json:"backgroundImage,omitempty"`            // URL of the background image.
	DefaultUserBackgroundImage *string             `json:"defaultUserBackgroundImage,omitempty"` // Default background image for users.
	LayoutConfig               *[]LayoutConfig     `json:"layoutConfig,omitempty"`               // Array of individual layout configurations.
	BackgroundConfig           *[]BackgroundConfig `json:"backgroundConfig,omitempty"`           // Array of background configurations.
}

// AquireClientRequest defines the parameters for acquiring a cloud recording resource.
// It includes the scene type, resource expiry, and initial recording parameters.
type AquireClientRequest struct {
	Scene               int           `json:"scene,omitempty"`               // The recording scene type.
	ResourceExpiredHour int           `json:"resourceExpiredHour,omitempty"` // The hour after which the resource expires.
	StartParameter      ClientRequest `json:"startParameter,omitempty"`      // Initial parameters for the recording.
	ExcludeResourceIds  *[]string     `json:"excludeResourceIds,omitempty"`  // List of resource IDs to exclude from recording.
}

// Timestampable is an interface that allows struct types to receive a timestamp.
// Implementing this interface ensures that a timestamp can be set on the object, primarily for auditing or tracking purposes.
type Timestampable interface {
	SetTimestamp(timestamp string)
}

// StartRecordingResponse represents the response received from the Agora server after successfully starting a recording.
// It includes the identifiers of the recording session along with an optional timestamp.
type StartRecordingResponse struct {
	Cname      string  `json:"cname"`               // The channel name for the recording session.
	Uid        string  `json:"uid"`                 // The UID for the recording session.
	ResourceId string  `json:"resourceId"`          // The unique resource identifier for the recording session.
	Sid        string  `json:"sid"`                 // The session identifier for the recording.
	Timestamp  *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *StartRecordingResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// ActiveRecordingResponse represents the current state of a recording session.
// It contains details about the recording session and includes a server response that can be parsed for more specific details.
type ActiveRecordingResponse struct {
	ResourceId     *string        `json:"resourceId"`               // The unique resource identifier for the recording session.
	Sid            *string        `json:"sid"`                      // The session identifier for the recording.
	ServerResponse ServerResponse `json:"serverResponse,omitempty"` // Detailed server response, parsed according to the specific recording scenario.
	Cname          *string        `json:"cname"`                    // The channel name for the recording session.
	Uid            *string        `json:"uid"`                      // The UID for the recording session.
	Timestamp      *string        `json:"timestamp,omitempty"`      // Optional timestamp for the current state of the recording.
}

func (s *ActiveRecordingResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// UpdateRecordingResponse represents the response from updating a recording session's settings.
// It includes session identifiers and an optional timestamp.
type UpdateRecordingResponse struct {
	Cname      *string `json:"cname"`               // The channel name for the recording session.
	Uid        *string `json:"uid"`                 // The UID for the recording session.
	ResourceId *string `json:"resourceId"`          // The unique resource identifier for the recording session.
	Sid        string  `json:"sid"`                 // The session identifier for the recording.
	Timestamp  *string `json:"timestamp,omitempty"` // Optional timestamp for when the recording was started.
}

func (s *UpdateRecordingResponse) SetTimestamp(timestamp string) {
	s.Timestamp = &timestamp
}

// ServerResponse encapsulates various possible states and details returned by the Agora server in response to recording commands.
// It is flexible enough to contain different types of data depending on the operation performed.
type ServerResponse struct {
	ExtensionServiceState   *ExtensionServiceState `json:"extensionServiceState,omitempty"`
	UploadingStatusResponse *string                `json:"uploadingStatus,omitempty"`
	FileListMode            *string                `json:"fileListMode,omitempty"` // Specifies how the file list is presented, e.g., as a string or JSON.
	FileList                *json.RawMessage       `json:"fileList,omitempty"`     // Contains details about the recorded files, format dependent on fileListMode.
}

// ExtensionServiceState details the state of any extension services used during the recording, such as additional data streaming or processing services.
type ExtensionServiceState struct {
	PlayloadStop *PlayloadStop `json:"playload"`    // Specific details about the stop condition or status of the service.
	ServiceName  string        `json:"serviceName"` // Name of the service provided.
}

// PlayloadStop contains details about the stopping condition of an extension service, including file list and recording state.
type PlayloadStop struct {
	UploadingStatus *string       `json:"uploadingStatus,omitempty"`
	FileList        *[]FileDetail `json:"fileList,omitempty"` // Details about each file created during the recording.
	OnHold          *bool         `json:"onhold"`
	State           *string       `json:"state"` // Current state of the service, such as "active" or "stopped".
}

// FileDetail provides specific details about a recorded file, including its name and the start time of the recording slice it contains.
type FileDetail struct {
	Filename       string `json:"filename"`       // Name of the file.
	SliceStartTime int64  `json:"sliceStartTime"` // UNIX timestamp marking the start of the recording slice contained in the file.
}

// FileListEntry details a file in a list of recorded files, including track type and whether it includes all users' mixed audio/video.
type FileListEntry struct {
	FileName       string `json:"fileName"`       // Name of the file.
	TrackType      string `json:"trackType"`      // Type of track, such as "audio" or "video".
	Uid            string `json:"uid"`            // UID associated with the track, if applicable.
	MixedAllUser   bool   `json:"mixedAllUser"`   // Indicates whether the file contains mixed content from all users.
	IsPlayable     bool   `json:"isPlayable"`     // Indicates whether the file is ready to be played.
	SliceStartTime int64  `json:"sliceStartTime"` // UNIX timestamp marking the start of the slice in the file.
}

// ClientRequest contains the detailed parameters for starting or updating a recording session.
// It encapsulates the authentication token, configurations for storage, recording, file management, snapshot options,
// extension services, application behaviors, and transcoding specifications.
type ClientRequest struct {
	Token                  string                  `json:"token,omitempty"`                  // Authentication token for the cloud recording session.
	StorageConfig          StorageConfig           `json:"storageConfig"`                    // Configuration parameters for storage during recording.
	RecordingConfig        RecordingConfig         `json:"recordingConfig"`                  // Settings related to the recording process.
	RecordingFileConfig    *RecordingFileConfig    `json:"recordingFileConfig,omitempty"`    // Optional configurations for recorded files.
	SnapshotConfig         *SnapshotConfig         `json:"snapshotConfig,omitempty"`         // Optional configurations for snapshots during recording.
	ExtensionServiceConfig *ExtensionServiceConfig `json:"extensionServiceConfig,omitempty"` // Optional configurations for any extension services used.
	AppsCollection         *AppsCollection         `json:"appsCollection,omitempty"`         // Collection of applications used during recording.
	TranscodeOptions       *TranscodeOptions       `json:"transcodeOptions,omitempty"`       // Optional transcoding settings.
}

// StorageConfig represents the details necessary for storing recorded content.
// It specifies the storage service provider details and authentication credentials.
type StorageConfig struct {
	Vendor          int              `json:"vendor"`                    // Identifier for the storage service provider.
	Region          int              `json:"region"`                    // Geographical region of the storage.
	Bucket          string           `json:"bucket"`                    // Storage bucket name.
	AccessKey       string           `json:"accessKey"`                 // Access key for storage authentication.
	SecretKey       string           `json:"secretKey"`                 // Secret key for storage authentication.
	FileNamePrefix  *[]string        `json:"fileNamePrefix,omitempty"`  // Optional prefixes for generating file paths within the bucket.
	ExtensionParams *ExtensionParams `json:"extensionParams,omitempty"` // Additional parameters for storage configuration.
}

// ExtensionParams adds further customization to the storage configuration, supporting specific features of the storage provider.
type ExtensionParams struct {
	SSE                *string `json:"sse,omitempty"`      // Server-side encryption option.
	Tag                *string `json:"tag,omitempty"`      // Custom tag for identifying storage settings or operations.
	EnableNTPtimestamp *bool   `json:"enableNTPtimestamp"` // Private Param for enabling subtitle sync in RTT
}

// RecordingConfig encapsulates all settings related to the recording process itself.
// This includes options for channel type, decryption, stream selection, and video/audio quality.
type RecordingConfig struct {
	ChannelType          int                `json:"channelType"`                    // Type of channel being recorded.
	DecryptionMode       *int               `json:"decryptionMode,omitempty"`       // Optional mode for decryption of incoming streams.
	Secret               *string            `json:"secret,omitempty"`               // Optional secret key for decryption.
	Salt                 *string            `json:"salt,omitempty"`                 // Optional salt used in conjunction with the secret key.
	MaxIdleTime          *int               `json:"maxIdleTime,omitempty"`          // Maximum time in seconds the recording can remain idle.
	StreamTypes          *int               `json:"streamTypes,omitempty"`          // Types of streams to be included in the recording.
	VideoStreamType      *int               `json:"videoStreamType,omitempty"`      // Specific type of video stream to record.
	SubscribeAudioUids   *[]string          `json:"subscribeAudioUids,omitempty"`   // List of audio UIDs to subscribe.
	UnsubscribeAudioUids *[]string          `json:"unsubscribeAudioUids,omitempty"` // List of audio UIDs to unsubscribe.
	SubscribeVideoUids   *[]string          `json:"subscribeVideoUids,omitempty"`   // List of video UIDs to subscribe.
	UnsubscribeVideoUids *[]string          `json:"unsubscribeVideoUids,omitempty"` // List of video UIDs to unsubscribe.
	SubscribeUidGroup    *int               `json:"subscribeUidGroup,omitempty"`    // Group of UIDs to subscribe collectively.
	StreamMode           *string            `json:"streamMode,omitempty"`           // Recording mode, such as individual, composite, or web.
	AudioProfile         *int               `json:"audioProfile,omitempty"`         // Audio quality profile.
	TranscodingConfig    *TranscodingConfig `json:"transcodingConfig,omitempty"`    // Optional transcoding settings.
}

// TranscodingConfig specifies the parameters for video transcoding during the recording process.
// It includes settings for resolution, frame rate, bitrate, and layout.
type TranscodingConfig struct {
	Width                      *int                `json:"width,omitempty"`                      // Width of the transcoded video.
	Height                     *int                `json:"height,omitempty"`                     // Height of the transcoded video.
	Fps                        *int                `json:"fps,omitempty"`                        // Frames per second of the transcoded video.
	Bitrate                    *int                `json:"bitrate,omitempty"`                    // Bitrate of the transcoded video.
	MaxResolutionUid           *string             `json:"maxResolutionUid,omitempty"`           // UID with the highest resolution video.
	MixedVideoLayout           *int                `json:"mixedVideoLayout,omitempty"`           // Layout of the mixed video.
	BackgroundColor            *string             `json:"backgroundColor,omitempty"`            // Background color for the video layout.
	BackgroundImage            *string             `json:"backgroundImage,omitempty"`            // Image used as a background.
	DefaultUserBackgroundImage *string             `json:"defaultUserBackgroundImage,omitempty"` // Default background image for users.
	LayoutConfig               *[]LayoutConfig     `json:"layoutConfig,omitempty"`               // Individual layout configurations.
	BackgroundConfig           *[]BackgroundConfig `json:"backgroundConfig,omitempty"`           // Background settings for individual layouts.
}

// LayoutConfig defines individual video layout positions and dimensions for participants in a recorded session.
type LayoutConfig struct {
	Uid        string `json:"uid"`         // User identifier for the layout configuration.
	XAxis      int    `json:"x_axis"`      // X-axis position for the layout.
	YAxis      int    `json:"y_axis"`      // Y-axis position for the layout.
	Width      int    `json:"width"`       // Width of the video in the layout.
	Height     int    `json:"height"`      // Height of the video in the layout.
	Alpha      int    `json:"alpha"`       // Transparency level of the video in the layout.
	RenderMode int    `json:"render_mode"` // Rendering mode for the video.
}

// BackgroundConfig specifies the background settings for individual participants or the entire session.
type BackgroundConfig struct {
	Uid        string `json:"uid"`         // User identifier for which the background is configured.
	ImageURL   string `json:"image_url"`   // URL of the image used as a background.
	RenderMode int    `json:"render_mode"` // Rendering mode for the background image.
}

// RecordingFileConfig represents the configuration for recorded files, specifically the types of audio/video files.
type RecordingFileConfig struct {
	AVFileType []string `json:"avFileType,omitempty"` // List of audio/video file types to be generated.
}

// SnapshotConfig specifies the settings for taking snapshots during a cloud recording session.
// It includes the interval between snapshots and the types of files to generate.
type SnapshotConfig struct {
	CaptureInterval int      `json:"captureInterval,omitempty"` // Time interval between snapshots, in seconds.
	FileType        []string `json:"fileType,omitempty"`        // Types of snapshot files to generate.
}

// ExtensionServiceConfig holds configuration for any extension services used during recording.
// This includes error handling policies and a list of specific extension services.
type ExtensionServiceConfig struct {
	ErrorHandlePolicy string             `json:"errorHandlePolicy,omitempty"` // Policy for handling errors in extension services.
	ExtensionServices []ExtensionService `json:"extensionServices,omitempty"` // List of extension services configured.
}

// ExtensionService defines a single service that can extend the functionality of cloud recording,
// such as additional processing or streaming capabilities.
type ExtensionService struct {
	ServiceName       string       `json:"serviceName"`                 // Name of the service.
	ErrorHandlePolicy *string      `json:"errorHandlePolicy,omitempty"` // Optional specific error handling policy for this service.
	ServiceParam      ServiceParam `json:"serviceParam"`                // Parameters specific to this service.
}

// ServiceParam encapsulates the parameters for an extension service, which can include video and audio settings.
type ServiceParam struct {
	URL              string `json:"url"`                        // URL of the extension service.
	AudioProfile     *int   `json:"audioProfile,omitempty"`     // Audio profile setting, if applicable.
	VideoWidth       *int   `json:"videoWidth,omitempty"`       // Width of the video stream.
	VideoHeight      *int   `json:"videoHeight,omitempty"`      // Height of the video stream.
	MaxRecordingHour *int   `json:"maxRecordingHour,omitempty"` // Maximum duration of the recording in hours.
	VideoBitrate     *int   `json:"videoBitrate,omitempty"`     // Bitrate of the video stream.
	VideoFps         *int   `json:"videoFps,omitempty"`         // Frames per second of the video stream.
	Mobile           *bool  `json:"mobile,omitempty"`           // Indicates if the service is used on mobile devices.
	MaxVideoDuration *int   `json:"maxVideoDuration,omitempty"` // Maximum duration of a single video file.
	OnHold           *bool  `json:"onhold,omitempty"`           // Indicates if the recording is on hold.
	ReadyTimeout     *int   `json:"readyTimeout,omitempty"`     // Timeout for the service to be ready.
}

// AppsCollection represents a collection of application settings used during the recording.
// It could include policies on how multiple apps interact or are used together.
type AppsCollection struct {
	CombinationPolicy *string `json:"combinationPolicy,omitempty"` // Policy defining how multiple apps are combined or used.
}

// TranscodeOptions encapsulates options for transcoding the recording, including video and audio configuration.
type TranscodeOptions struct {
	TransConfig *TransConfig `json:"transConfig,omitempty"` // Transcoding configurations like mode of transcoding.
	Container   *Container   `json:"container,omitempty"`   // Container settings for the transcoded files.
	Audio       *Audio       `json:"audio,omitempty"`       // Audio settings for transcoding.
}

// TransConfig defines the transcoding mode used for converting media from one format to another.
type TransConfig struct {
	TransMode *string `json:"transMode,omitempty"` // Transcoding mode, defines how the media is processed.
}

// Container represents the container format for the recorded media files.
type Container struct {
	Format *string `json:"format,omitempty"` // Media container format, e.g., mp4, mkv.
}

// Audio specifies the audio settings used in transcoding, affecting quality and compatibility.
type Audio struct {
	SampleRate *string `json:"sampleRate,omitempty"` // Sampling rate of the audio.
	Bitrate    *string `json:"bitrate,omitempty"`    // Bitrate of the audio.
	Channels   *string `json:"channels,omitempty"`   // Number of audio channels.
}

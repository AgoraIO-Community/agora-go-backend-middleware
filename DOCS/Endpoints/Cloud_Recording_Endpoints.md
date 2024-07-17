# Cloud Recording API

This document provides details about the Cloud Recording API endpoints with curl examples for testing.

## Start Recording

Starts a cloud recording session.

### Endpoint

**POST:** `/cloud_recording/start`

### Request Body

```json
{
  "channelName": "string",
  "uid": "string",
  "recordingConfig": {
    // RecordingConfig fields
  },
  "storageConfig": {
    // StorageConfig fields
  }
}
```

### Response

```json
{
  "resourceId": "string",
  "sid": "string",
  "timestamp": "string"
}
```

## Stop Recording

Stops an ongoing cloud recording session.

### Endpoint

**POST:** `/cloud_recording/stop`

### Request Body

```json
{
  "cname": "string",
  "uid": "string",
  "resourceId": "string",
  "sid": "string",
  "recordingMode": "string",
  "async_stop": boolean
}
```

### Response

```json
{
  "resourceId": "string",
  "sid": "string",
  "serverResponse": {
    "fileListMode": "string",
    "fileList": [
      {
        "fileName": "string",
        "trackType": "string",
        "uid": "string",
        "mixedAllUser": boolean,
        "isPlayable": boolean,
        "sliceStartTime": number
      }
    ]
  },
  "timestamp": "string"
}
```

## Get Recording Status

Retrieves the status of a cloud recording session.

### Endpoint

**GET:** `/cloud_recording/status`

### Query Parameters

- `resourceId`: string
- `sid`: string
- `mode`: string

### Response

```json
{
  "resourceId": "string",
  "sid": "string",
  "serverResponse": {
    "fileListMode": "string",
    "fileList": [
      {
        "fileName": "string",
        "trackType": "string",
        "uid": "string",
        "mixedAllUser": boolean,
        "isPlayable": boolean,
        "sliceStartTime": number
      }
    ]
  },
  "timestamp": "string"
}
```

## Update Subscriber List

Updates the subscriber list for a cloud recording session.

### Endpoint

**POST:** `/cloud_recording/update/subscriber-list`

### Request Body

```json
{
  "cname": "string",
  "uid": "string",
  "resourceId": "string",
  "sid": "string",
  "recordingMode": "string",
  "recordingConfig": {
    // UpdateSubscriptionClientRequest fields
  }
}
```

### Response

```json
{
  "cname": "string",
  "uid": "string",
  "resourceId": "string",
  "sid": "string",
  "timestamp": "string"
}
```

## Update Layout

Updates the layout of a cloud recording session.

### Endpoint

**POST:** `/cloud_recording/update/layout`

### Request Body

```json
{
  "cname": "string",
  "uid": "string",
  "resourceId": "string",
  "sid": "string",
  "recordingMode": "string",
  "recordingConfig": {
    // UpdateLayoutClientRequest fields
  }
}
```

### Response

```json
{
  "cname": "string",
  "uid": "string",
  "resourceId": "string",
  "sid": "string",
  "timestamp": "string"
}
```

Replace `localhost:8080` with your server's address if different.

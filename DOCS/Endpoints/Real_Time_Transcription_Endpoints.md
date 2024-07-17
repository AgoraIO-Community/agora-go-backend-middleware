# Real-Time Transcription (RTT) API

This document provides details about the Real-Time Transcription API endpoints with curl examples for testing.

## Start RTT

Starts a real-time transcription session.

### Endpoint

**POST:** `/rtt/start`

### Request Body

```json
{
  "channelName": "string",
  "languages": ["string"],
  "subscribeAudioUIDs": ["string"],
  "cryptionMode": "string",
  "secret": "string",
  "salt": "string",
  "maxIdleTime": int,
  "translateConfig": {
    "forceTranslateInterval": int,
    "languages": [
      {
        "source": "string",
        "target": ["string"]
      }
    ]
  },
  "enableStorage": boolean,
  "enableNTPtimestamp": boolean
}
```

### Response

```json
{
  "acquire": {
    "tokenName": "string",
    "createTs": number,
    "instanceId": "string",
    "timestamp": "string"
  },
  "start": {
    "createTs": number,
    "status": "string",
    "taskId": "string",
    "timestamp": "string"
  },
  "timestamp": "string"
}
```

## Stop RTT

Stops an ongoing real-time transcription session.

### Endpoint

**DELETE:** `/rtt/stop/:taskId`

### Request Body

```json
{
  "builderToken": "string"
}
```

### Response

```json
{
  "stop": {
    "timestamp": "string"
  },
  "timestamp": "string"
}
```

## Get RTT Status

Retrieves the status of a real-time transcription session.

### Endpoint

**GET:** `/rtt/status/:taskId`

### Query Parameters

- `builderToken`: string

### Response

```json
{
  "createTs": number,
  "status": "string",
  "taskId": "string",
  "timestamp": "string"
}
```

Replace `localhost:8080` with your server's address if different.

Note: All responses include a `timestamp` field for auditing purposes.

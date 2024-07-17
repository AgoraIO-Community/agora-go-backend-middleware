# RTMP API

This document provides details about the RTMP API endpoints.

## Start RTMP Push

Starts an RTMP push session.

### Endpoint

**POST:** `/rtmp/push/start`

### Request Body

```json
{
  "converterName": "string",
  "rtcChannel": "string",
  "streamUrl": "string",
  "streamKey": "string",
  "region": "string",
  "regionHintIp": "string",
  "useTranscoding": boolean,
  "rtcStreamUid": "string",
  "audioOptions": {
    // PushAudioOptions fields
  },
  "videoOptions": {
    // PushVideoOptions fields
  },
  "idleTimeOut": int,
  "jitterBufferSizeMs": int
}
```

### Response

```json
{
  "converter": {
    "id": "string",
    "createTs": number,
    "updateTs": number,
    "state": "string"
  },
  "fields": "string",
  "timestamp": "string"
}
```

## Stop RTMP Push

Stops an ongoing RTMP push session.

### Endpoint

**POST:** `/rtmp/push/stop`

### Request Body

```json
{
  "converterId": "string",
  "region": "string"
}
```

### Response

```json
{
  "status": "string",
  "timestamp": "string"
}
```

## Update RTMP Push

Updates an ongoing RTMP push session.

### Endpoint

**POST:** `/rtmp/push/update`

### Request Body

```json
{
  "converterId": "string",
  "region": "string",
  "streamUrl": "string",
  "streamKey": "string",
  "rtcChannel": "string",
  "videoOptions": {
    // PushVideoOptions fields
  },
  "jitterBufferSizeMs": int,
  "sequenceId": int
}
```

### Response

```json
{
  "converter": {
    "id": "string",
    "createTs": number,
    "updateTs": number,
    "state": "string"
  },
  "fields": "string",
  "timestamp": "string"
}
```

## Start Cloud Player (RTMP Pull)

Starts a Cloud Player session.

### Endpoint

**POST:** `/rtmp/pull/start`

### Request Body

```json
{
  "channelName": "string",
  "streamUrl": "string",
  "region": "string",
  "uid": "string",
  "playerName": "string",
  "streamOriginIp": "string",
  "audioOptions": {
    // PullAudioOptions fields
  },
  "videoOptions": {
    // PullVideoOptions fields
  },
  "idleTimeOut": int,
  "playTs": int,
  "encryptMode": "string"
}
```

### Response

```json
{
  "player": {
    "id": "string",
    "createTs": number,
    "uid": "string"
  },
  "fields": "string",
  "timestamp": "string"
}
```

## Stop Cloud Player (RTMP Pull)

Stops an ongoing Cloud Player session.

### Endpoint

**POST:** `/rtmp/pull/stop`

### Request Body

```json
{
  "playerId": "string",
  "region": "string"
}
```

### Response

```json
{
  "status": "string",
  "timestamp": "string"
}
```

## Update Cloud Player (RTMP Pull)

Updates an ongoing Cloud Player session.

### Endpoint

**POST:** `/rtmp/pull/update`

### Request Body

```json
{
  "playerId": "string",
  "region": "string",
  "streamUrl": "string",
  "audioOptions": {
    // PullAudioOptions fields
  },
  "isPause": boolean,
  "seekPosition": int,
  "sequenceId": int
}
```

### Response

```json
{
  "status": "string",
  "timestamp": "string"
}
```

Replace `localhost:8080` with your server's address if different.

# RTMP API

This document provides curl examples for testing the backend's RTMP API endpoints.

## Start RTMP Push

Starts an RTMP push session.

```
POST /rtmp/push/start
```

```bash
curl -X POST http://localhost:8080/rtmp/push/start \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "converterName": "test_converter",
    "rtcChannel": "test_channel",
    "streamUrl": "rtmp://example.com/live",
    "streamKey": "your-stream-key",
    "region": "na",
    "useTranscoding": true,
    "rtcStreamUid": "12345",
    "audioOptions": {
      "codecProfile": "lc-aac",
      "sampleRate": 44100,
      "bitrate": 128,
      "audioChannels": 2
    },
    "videoOptions": {
      "codec": "h264",
      "width": 1280,
      "height": 720,
      "frameRate": 30,
      "bitrate": 2000
    },
    "idleTimeOut": 300
  }'
```

## Stop RTMP Push

Stops an ongoing RTMP push session.

```
POST /rtmp/push/stop
```

```bash
curl -X POST http://localhost:8080/rtmp/push/stop \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "converterId": "your-converter-id",
    "region": "na"
  }'
```

## Update RTMP Push

Updates an ongoing RTMP push session.

```
POST /rtmp/push/update
```

```bash
curl -X POST http://localhost:8080/rtmp/push/update \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "converterId": "your-converter-id",
    "region": "na",
    "streamUrl": "rtmp://example.com/live",
    "streamKey": "your-new-stream-key",
    "rtcChannel": "new_test_channel",
    "videoOptions": {
      "width": 1920,
      "height": 1080,
      "frameRate": 60,
      "bitrate": 4000
    },
    "jitterBufferSizeMs": 1000
  }'
```

## Start Cloud Player (RTMP Pull)

Starts a Cloud Player session.

```
POST /rtmp/pull/start
```

```bash
curl -X POST http://localhost:8080/rtmp/pull/start \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "channelName": "test_channel",
    "streamUrl": "rtmp://example.com/live/stream",
    "region": "na",
    "playerName": "test_player",
    "audioOptions": {
      "profile": 0
    },
    "videoOptions": {
      "width": 1280,
      "height": 720,
      "frameRate": 30,
      "bitrate": 2000,
      "codec": "H264"
    },
    "idleTimeOut": 300
  }'
```

## Stop Cloud Player (RTMP Pull)

Stops an ongoing Cloud Player session.

```
POST /rtmp/pull/stop
```

```bash
curl -X POST http://localhost:8080/rtmp/pull/stop \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "playerId": "your-player-id",
    "region": "na"
  }'
```

## Update Cloud Player (RTMP Pull)

Updates an ongoing Cloud Player session.

```
POST /rtmp/pull/update
```

```bash
curl -X POST http://localhost:8080/rtmp/pull/update \
  -H "Content-Type: application/json" \
  -H "X-Request-ID: unique-request-id" \
  -d '{
    "playerId": "your-player-id",
    "region": "na",
    "streamUrl": "rtmp://example.com/live/new-stream",
    "audioOptions": {
      "profile": 1
    },
    "isPause": false,
    "seekPosition": 30
  }'
```

Replace `localhost:8080` with your server's address if different.

> Note: All requests require the `X-Request-ID` header for request tracing.

# Cloud Recording API

This document provides curl examples for testing the backend's Cloud Recording API endpoints.

## Start Recording

Starts a cloud recording session.

`POST /cloud_recording/start`

(simple)

```bash
curl -X POST http://localhost:8080/cloud_recording/start \
  -H "Content-Type: application/json" \
  -d '{
    "channelName": "test_channel",
    "sceneMode": "realtime",
    "recordingMode": "mix",
    "excludeResourceIds": []
  }'

```

(full)

```bash
curl -X POST http://localhost:8080/cloud_recording/start \
  -H "Content-Type: application/json" \
  -d '{
    "channelName": "testChannel",
    "sceneMode": "realtime",
    "recordingMode": "mix",
    "excludeResourceIds": [],
    "recordingConfig": {
      "channelType": 0,
      "decryptionMode": 1,
      "secret": "your_secret",
      "salt": "your_salt",
      "maxIdleTime": 120,
      "streamTypes": 2,
      "videoStreamType": 0,
      "subscribeAudioUids": ["#allstream#"],
      "unsubscribeAudioUids": [],
      "subscribeVideoUids": ["#allstream#"],
      "unsubscribeVideoUids": [],
      "subscribeUidGroup": 0,
      "streamMode": "individual",
      "audioProfile": 1,
      "transcodingConfig": {
        "width": 640,
        "height": 360,
        "fps": 15,
        "bitrate": 500,
        "maxResolutionUid": "1",
        "layoutConfig": [
          {
            "x_axis": 0,
            "y_axis": 0,
            "width": 640,
            "height": 360,
            "alpha": 1,
            "render_mode": 1
          }
        ]
      }
    }
  }'
```

## Stop Recording

Stops an ongoing cloud recording session.

`POST /cloud_recording/stop`

```bash
curl -X POST http://localhost:8080/cloud_recording/stop \
  -H "Content-Type: application/json" \
  -d '{
    "cname": "test_channel",
    "uid": "uid-from-start-response",
    "resourceId": "resource-id-from-start-response",
    "sid": "sid-from-start-response",
    "recordingMode": "mix",
    "async_stop": false
  }'
```

## Get Recording Status

Retrieves the status of a cloud recording session.

`GET /cloud_recording/status`

```bash
curl -X GET "http://localhost:8080/cloud_recording/status?resourceId=your-resource-id&sid=your-sid&mode=mix"
```

## Update Subscriber List

Updates the subscriber list for a cloud recording session.

`POST /cloud_recording/update/subscriber-list`

```bash
curl -X POST http://localhost:8080/cloud_recording/update/subscriber-list \
  -H "Content-Type: application/json" \
  -d '{
    "cname": "test_channel",
    "uid": "uid-from-start-response",
    "resourceId": "your-resource-id",
    "sid": "your-sid",
    "recordingMode": "mix",
    "recordingConfig": {
      "streamSubscribe": {
        "audioUidList": {
          "subscribeAudioUids": ["2345", "3456"]
        },
        "videoUidList": {
          "subscribeVideoUids": ["2345", "3456"]
        }
      }
    }
  }'
```

## Update Layout

Updates the layout of a cloud recording session.

`POST /cloud_recording/update/layout`

```bash
curl -X POST http://localhost:8080/cloud_recording/update/layout \
  -H "Content-Type: application/json" \
  -d '{
    "cname": "test_channel",
    "uid": "uid-from-start-response",
    "resourceId": "your-resource-id",
    "sid": "your-sid",
    "recordingMode": "mix",
    "recordingConfig": {
      "mixedVideoLayout": 1,
      "backgroundColor": "#000000",
      "layoutConfig": [
        {
          "uid": "2345",
          "x_axis": 0,
          "y_axis": 0,
          "width": 360,
          "height": 640,
          "alpha": 1,
          "render_mode": 1
        }
      ]
    }
  }'
```

Replace `localhost:8080` with your server's address if different.

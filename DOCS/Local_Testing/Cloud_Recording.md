# Cloud Recording API

This document provides curl examples for testing the backend's Cloud Recording API endpoints.

## Start Recording

Starts a cloud recording session.

`POST /cloud_recording/start`

```bash
curl -X POST http://localhost:8080/cloud_recording/start \
  -H "Content-Type: application/json" \
  -d '{
    "channelName": "test_channel",
    "uid": "1234",
    "recordingConfig": {
      "maxIdleTime": 30,
      "streamTypes": 2,
      "channelType": 1,
      "subscribeUidGroup": 0
    },
    "storageConfig": {
      "vendor": 1,
      "region": 0,
      "bucket": "your-bucket-name",
      "accessKey": "your-access-key",
      "secretKey": "your-secret-key"
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
    "uid": "1234",
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
    "uid": "1234",
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
    "uid": "1234",
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

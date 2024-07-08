# Real-Time Transcription (RTT) API

This document provides details about the Real-Time Transcription API endpoints with curl examples for testing.

## Start RTT

Starts a real-time transcription session.

`POST /rtt/start`

(simple)

```bash
curl -X POST http://localhost:8080/rtt/start \
  -H "Content-Type: application/json" \
  -d '{
    "channelName": "test_channel",
    "languages": ["en-US"],
    "subscribeAudioUids": ["410431250", "802976399"],
    "maxIdleTime": 300,
    "enableStorage": false
  }'
```

(full)

```bash
curl -X POST http://localhost:8080/rtt/start \
  -H "Content-Type: application/json" \
  -d '{
    "channelName": "test_channel",
    "languages": ["en-US"],
    "subscribeAudioUIDs": ["1234", "5678"],
    "maxIdleTime": 300,
    "translateConfig": {
      "forceTranslateInterval": 60,
      "languages": [
        {
          "source": "en-US",
          "target": ["es-ES", "fr-FR"]
        }
      ]
    },
    "enableStorage": false,
    "enableNTPtimestamp": true
  }'
```

## Stop RTT

Stops an ongoing real-time transcription session.

`DELETE /rtt/stop/:taskId`

```bash
curl -X DELETE http://localhost:8080/rtt/stop/your-task-id \
  -H "Content-Type: application/json" \
  -d '{
    "builderToken": "your-builder-token"
  }'
```

## Get RTT Status

Retrieves the status of a real-time transcription session.

`GET /rtt/status/:taskId`

```bash
curl -X GET "http://localhost:8080/rtt/status/your-task-id?builderToken=your-builder-token"
```

Replace `localhost:8080` with your server's address if different.

Note: All responses include a `timestamp` field for auditing purposes.

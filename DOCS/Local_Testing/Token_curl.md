# Token Generation API

This document provides curl examples for testing the backend's Token Generation API endpoints.

## Generate Token

Generates a token based on the provided parameters.

**POST:** `/token/getNew`

### RTC

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "rtc",
  "channel": "testChannel",
  "role": "publisher",
  "uid": "12345",
  "expire": 3600
}'
```

### RTM

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "rtm",
  "uid": "12345",
  "expire": 3600
}'
```

### CHAT

```bash
curl -X POST http://localhost:8080/token/getNew \
-H "Content-Type: application/json" \
-d '{
  "tokenType": "chat",
  "uid": "12345",
  "expire": 3600
}'
```

Replace `localhost:8080` with your server's address if different.

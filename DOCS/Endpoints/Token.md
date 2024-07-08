# Token Generation API

This document provides details about the Token Generation API endpoints with curl examples for testing.

## Generate Token

Generates a token based on the provided parameters.

### Endpoint

`POST /token/getNew`

### Request Body

```json
{
  "tokenType": "rtc|rtm|chat",
  "channel": "string",
  "uid": "string",
  "role": "publisher|subscriber",
  "expire": int
}
```

### Response

```json
{
  "token": "string"
}
```

Replace `localhost:8080` with your server's address if different.

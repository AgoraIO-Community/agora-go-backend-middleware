# Agora Backend Webservice

This project is a fork of the [AgoraIO Community Token Service](https://github.com/AgoraIO-Community/agora-token-service/), that adds support for Cloud Recording, Real-Time Transcription, and Media Services.

Written in Golang, using [Gin framework](https://github.com/gin-gonic/gin). A RESTful webservice for interacting with [Agora.io](https://www.agora.io).

## How to Run

Set the APP_ID, APP_CERTIFICATE and CORS_ALLOW_ORIGIN env variables.

```bash
cp .env.example .env
```

```bash
go run cmd/main.go
```

## Endpoints

- [GET /ping]()
- [POST /token/getNew]()
- [POST /cloudrecording/start]()
- [POST /cloudrecording/stop]()
- [GET /cloudrecording/status]()
- [POST /cloudrecording/update/subscriber-list]()
- [POST /cloudrecording/update/layout]()

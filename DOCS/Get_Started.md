# Get Started with the Agora Middleware Service

Welcome to the comprehensive guide for getting started with the Agora Middleware Service! This powerful tool is designed to streamline your integration with Agora's cloud services, providing a robust backend solution for token generation, cloud recording, and real-time transcription. Whether you're a seasoned developer or just getting started with Agora, this guide will walk you through the setup process and help you make the most of the middleware service.

## Table of Contents

1. [Introduction](#introduction)
2. [Prerequisites](#prerequisites)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [Running the Service](#running-the-service)
6. [API Overview](#api-overview)
   - [Token Generation](#token-generation)
   - [Cloud Recording](#cloud-recording)
   - [Real-Time Transcription](#real-time-transcription)
7. [Testing the APIs](#testing-the-apis)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [Getting Help](#getting-help)

## Introduction

The Agora Middleware Service is a Go-based backend solution that simplifies the process of integrating Agora's cloud services into your applications. It provides a set of RESTful APIs for token generation, cloud recording management, and real-time transcription. By using this middleware, you can offload complex server-side operations and focus on building great user experiences in your applications.

Key features of the middleware service include:

- Token generation for RTC, RTM, and Chat services
- Cloud recording management (start, stop, update layout, etc.)
- Real-time transcription control
- Configurable storage options for recordings and transcriptions
- Easy-to-use RESTful API endpoints

Let's dive into setting up and using this powerful tool!

## Prerequisites

Before you begin, make sure you have the following:

1. Go (version 1.16 or later) installed on your system
1. An Agora developer account (if you don't have one, sign up at [https://www.agora.io/](https://www.agora.io/))
1. Your Agora App ID and App Certificate
1. Your Agora Customer ID and Customer Secret
1. A cloud storage provider: Amazon S3, Alibaba Cloud, Tencent Cloud, Microsoft Azure, Google Cloud, Huawei Cloud, Baidu IntelligentCloud
1. Basic knowledge of RESTful APIs and JSON
1. (Optional) A tool for making HTTP requests, such as Bruno or Postman. You can also use cURL.

## Installation

To get started with the Agora Middleware Service, follow these steps:

1. Clone the repository:

   ```
   git clone https://github.com/AgoraIO-Community/agora-go-backend-middleware.git
   ```

2. Navigate to the project directory:

   ```
   cd agora-go-backend-middleware
   ```

3. Install the required dependencies:
   ```
   go mod download
   ```

## Configuration

Proper configuration is crucial for the middleware service to function correctly. Follow these steps to set up your environment:

1. Copy the example environment file:

   ```
   cp .env.example .env
   ```

2. Open the `.env` file in your favorite text editor and fill in the required values:

   ```
   APP_ID=your_agora_app_id
   APP_CERTIFICATE=your_agora_app_certificate
   CUSTOMER_ID=your_customer_id
   CUSTOMER_SECRET=your_customer_secret
   CORS_ALLOW_ORIGIN=*
   SERVER_PORT=8080
   AGORA_BASE_URL=https://api.agora.io/
   AGORA_CLOUD_RECORDING_URL=v1/apps/{appId}/cloud_recording
   AGORA_RTT_URL=v1/projects/{appId}/rtsc/speech-to-text
   AGORA_RTMP_URL=v1/projects/{appId}/rtmp-converters
   AGORA_CLOUD_PLAYER_URL=v1/projects/{appId}/cloud-player
   STORAGE_VENDOR=
   STORAGE_REGION=
   STORAGE_BUCKET=
   STORAGE_BUCKET_ACCESS_KEY=
   STORAGE_BUCKET_SECRET_KEY=
   ```

   Make sure to replace the placeholder values with your actual Agora credentials and desired configuration options.

3. If you're using cloud storage for recordings or transcriptions, fill in the appropriate storage configuration values.

## Running the Service

Now that you've configured the service, it's time to run it. We have a few options:

1. Using Go:

   ```bash
   go run cmd/main.go
   ```

2. Using Docker and Make:

   ```bash
   make build
   make run
   ```

If everything is set up correctly, you should see output indicating that the server is running on `http://localhost:8080` (unless you've specified a different port in the `.env` file).

## API Overview

The middleware service exposes Agora APIs for: Token Generation, Cloud Recording, Real-Time Transcription, and RTMP (Push & Pull). Let's explore each set of endpoints and how to use them.

### Token Generation

The Token Generation API allows you to create tokens for Agora services. These tokens are essential for authenticating users and maintaining the security of your application.

Endpoint: `POST /token/getNew`

You can generate tokens for:

- RTC (Real-Time Communication)
- RTM (Real-Time Messaging)
- Chat

Example request body:

```json
{
  "tokenType": "rtc",
  "channel": "myChannel",
  "uid": "12345",
  "role": "publisher",
  "expire": 3600
}
```

### Cloud Recording

The Cloud Recording API provides endpoints for managing cloud recording sessions. This includes starting and stopping recordings, updating layouts, and managing recording resources.

Key endpoints:

- Start Recording: `POST /cloud_recording/start`
- Stop Recording: `POST /cloud_recording/stop`
- Get Recording Status: `GET /cloud_recording/status`
- Update Layout: `POST /cloud_recording/update/layout`

Example start recording request body:

```json
{
  "channelName": "myChannel",
  "uid": "12345",
  "recordingConfig": {
    "maxIdleTime": 30,
    "streamTypes": 2,
    "channelType": 1,
    "videoStreamType": 0,
    "subscribeUidGroup": 0
  }
}
```

### Real-Time Transcription

The Real-Time Transcription (RTT) API allows you to control transcription sessions, enabling you to convert speech to text in real-time during your Agora calls.

Key endpoints:

- Start RTT: `POST /rtt/start`
- Stop RTT: `POST /rtt/stop/:taskId`
- Query RTT Status: `GET /rtt/status/:taskId`

Example start RTT request body:

```json
{
  "channelName": "myChannel",
  "languages": ["en-US"],
  "subscribeAudioUIDs": ["12345", "67890"],
  "maxIdleTime": 300,
  "enableStorage": true
}
```

## Testing the APIs

To ensure your middleware service is working correctly, you can test the APIs using cURL or a tool like Postman. Here are some example cURL commands for each main API category:

1. Generate an RTC token:

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

2. Start a cloud recording:

   ```bash
   curl -X POST http://localhost:8080/cloud_recording/start \
   -H "Content-Type: application/json" \
   -d '{
     "channelName": "testChannel",
     "uid": "12345",
     "recordingConfig": {
       "maxIdleTime": 30,
       "streamTypes": 2,
       "channelType": 1,
       "videoStreamType": 0,
       "subscribeUidGroup": 0
     }
   }'
   ```

3. Start a real-time transcription session:
   ```bash
   curl -X POST http://localhost:8080/rtt/start \
   -H "Content-Type: application/json" \
   -d '{
     "channelName": "testChannel",
     "languages": ["en-US"],
     "subscribeAudioUIDs": ["12345", "67890"],
     "maxIdleTime": 300,
     "enableStorage": false
   }'
   ```

Make sure to replace `localhost:8080` with the appropriate address if your service is running on a different host or port.

## Best Practices

To make the most of the Agora Middleware Service, consider the following best practices:

1. **Security**: Always use HTTPS in production environments to encrypt data in transit.
2. **Error Handling**: Implement proper error handling in your client applications to gracefully manage API responses.
3. **Rate Limiting**: Be mindful of Agora's API rate limits and implement appropriate throttling mechanisms if necessary.
4. **Logging**: Enable and monitor logs to track usage and troubleshoot issues.
5. **Regular Updates**: Keep the middleware service updated to benefit from the latest features and security improvements.

## Troubleshooting

If you encounter issues while using the middleware service, try the following:

1. Double-check your `.env` configuration to ensure all values are correct.
2. Verify that your Agora account is active and has the necessary permissions.
3. Check the console output for any error messages or logs.
4. Ensure your firewall or network settings are not blocking the service.
5. Verify that the required ports are open and accessible.

## Getting Help

If you need further assistance or have questions about the Agora Middleware Service, you have several options:

1. Check the [official Agora documentation](https://docs.agora.io/) for detailed information about Agora's services.
2. Visit the [Agora Developer Center](https://www.agora.io/en/developer-center/) for additional resources and support.
3. Join the [Agora Developer Slack community](https://www.agora.io/en/join-slack/) to connect with other developers and Agora experts.
4. Open an issue on the [GitHub repository](https://github.com/AgoraIO-Community/agora-go-backend-middleware) if you believe you've found a bug or have a feature request.

Remember, the Agora team is here to help you (our community) succeed in building amazing real-time communication applications!

With this guide, you should now have a solid understanding of how to set up, configure, and use the Agora Middleware Service. Happy coding, and we can't wait to see what you'll build with Agora!

# Start from a Debian image with the latest version of Go installed
FROM golang:1.21-alpine

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Declare build-time ARGs
ARG APP_ID
ARG APP_CERTIFICATE
ARG CUSTOMER_ID
ARG CUSTOMER_SECRET
ARG SERVER_PORT
ARG CORS_ALLOW_ORIGIN
ARG AGORA_BASE_URL
ARG AGORA_CLOUD_RECORDING_URL
ARG AGORA_RTT_URL
ARG STORAGE_VENDOR
ARG STORAGE_REGION
ARG STORAGE_BUCKET
ARG STORAGE_BUCKET_ACCESS_KEY
ARG STORAGE_BUCKET_SECRET_KEY

# Set them as persistent ENV variables
ENV APP_ID=$APP_ID \
    APP_CERTIFICATE=$APP_CERTIFICATE \
    CUSTOMER_ID=$CUSTOMER_ID \
    CUSTOMER_SECRET=$CUSTOMER_SECRET \
    SERVER_PORT=$SERVER_PORT \
    CORS_ALLOW_ORIGIN=$CORS_ALLOW_ORIGIN \
    AGORA_BASE_URL=$AGORA_BASE_URL \
    AGORA_CLOUD_RECORDING_URL=$AGORA_CLOUD_RECORDING_URL \
    AGORA_RTT_URL=$AGORA_RTT_URL \
    STORAGE_VENDOR=$STORAGE_VENDOR \
    STORAGE_REGION=$STORAGE_REGION \
    STORAGE_BUCKET=$STORAGE_BUCKET \
    STORAGE_BUCKET_ACCESS_KEY=$STORAGE_BUCKET_ACCESS_KEY \
    STORAGE_BUCKET_SECRET_KEY=$STORAGE_BUCKET_SECRET_KEY

# Build the application
RUN go build -v -o agora-backend-middleware ./cmd/main.go

# Run the application
CMD ["./agora-backend-middleware"]

# Document that the service listens on the specified port
EXPOSE $SERVER_PORT
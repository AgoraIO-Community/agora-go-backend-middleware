package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/AgoraIO-Community/agora-go-backend-middleware/cloud_recording_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/http_headers"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/real_time_transcription_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/rtmp_service"
	"github.com/AgoraIO-Community/agora-go-backend-middleware/token_service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from a .env file, logging an error if the file cannot be loaded.
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Error loading .env file. Using existing environment variables.")
	}

	// Retrieve all configuration values from environment variables.
	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")
	customerIDEnv, customerIDExists := os.LookupEnv("CUSTOMER_ID")
	customerSecretEnv, customerSecretExists := os.LookupEnv("CUSTOMER_SECRET")
	corsAllowOrigin, _ := os.LookupEnv("CORS_ALLOW_ORIGIN")
	baseURLEnv, baseURLExists := os.LookupEnv("AGORA_BASE_URL")
	cloudRecordingURLEnv, cloudRecordingURLExists := os.LookupEnv("AGORA_CLOUD_RECORDING_URL")
	realTimeTranscriptionURLEnv, realTimeTranscriptionURLExists := os.LookupEnv("AGORA_RTT_URL")
	rtmpURLEnv, rtmpURLExists := os.LookupEnv("AGORA_RTMP_URL")
	storageVendorEnv, vendorExists := os.LookupEnv("STORAGE_VENDOR")
	storageRegionEnv, regionExists := os.LookupEnv("STORAGE_REGION")
	storageBucketEnv, bucketExists := os.LookupEnv("STORAGE_BUCKET")
	storageAccessKeyEnv, accessKeyExists := os.LookupEnv("STORAGE_BUCKET_ACCESS_KEY")
	storageSecretKeyEnv, secretKeyExists := os.LookupEnv("STORAGE_BUCKET_SECRET_KEY")

	// Check for for the presence of core environment variables
	if !appIDExists || !appCertExists {
		log.Fatal("FATAL ERROR: ENV not properly configured, APP ID and APP CERTIFICATE are required.")
	}

	// Set up the Gin HTTP router with headers for CORS, caching, and timestamp.
	router := gin.Default()
	var httpHeaders = http_headers.NewHttpHeaders(corsAllowOrigin)
	router.Use(httpHeaders.NoCache())
	router.Use(httpHeaders.CORShttpHeaders())
	router.Use(httpHeaders.Timestamp())

	// Initialize services & register routes.
	tokenService := token_service.NewTokenService(appIDEnv, appCertEnv)
	tokenService.RegisterRoutes(router)

	if baseURLExists {
		// Check for Basic Auth settings if baseURL is provided
		if !customerIDExists || !customerSecretExists {
			log.Fatal("FATAL ERROR: ENV not properly configured for Basic Auth, check .env file for all required variables")
		}
		// get basicAuth key
		basicAuthKey := getBasicAuth(customerIDEnv, customerSecretEnv)

		if cloudRecordingURLExists || realTimeTranscriptionURLExists {
			if !vendorExists || !regionExists || !bucketExists || !accessKeyExists || !secretKeyExists {
				log.Fatal("FATAL ERROR: ENV not properly configured for cloud storage, check .env file for all required variables")
			}

			// Convert storage vendor and region environment variables to integers.
			storageVendorInt, storageVendorErr := strconv.Atoi(storageVendorEnv)
			storageRegionInt, storageRegionErr := strconv.Atoi(storageRegionEnv)

			if storageVendorErr != nil || storageRegionErr != nil {
				log.Fatal("FATAL ERROR: Invalid STORAGE_VENDOR / STORAGE_REGION not properly configured")
			}

			// Configure storage settings based on environment variables.
			storageConfig := cloud_recording_service.StorageConfig{
				Vendor:    storageVendorInt,
				Region:    storageRegionInt,
				Bucket:    storageBucketEnv,
				AccessKey: storageAccessKeyEnv,
				SecretKey: storageSecretKeyEnv,
			}

			if cloudRecordingURLExists {
				// Init Cloud Recording Service
				cloudRecordingUrl := baseURLEnv + strings.Replace(cloudRecordingURLEnv, "{appId}", appIDEnv, 1) // replace the place-holder value with appID from Env
				cloudRecordingService := cloud_recording_service.NewCloudRecordingService(appIDEnv, cloudRecordingUrl, basicAuthKey, tokenService, storageConfig)
				cloudRecordingService.RegisterRoutes(router)
			}

			if realTimeTranscriptionURLExists {
				// Init Real Time Transcription Service
				realTimeTranscriptionUrl := baseURLEnv + strings.Replace(realTimeTranscriptionURLEnv, "{appId}", appIDEnv, 1) //replace the place-holder value with appID
				realTimeTranscriptionService := real_time_transcription_service.NewRTTService(appIDEnv, realTimeTranscriptionUrl, basicAuthKey, tokenService, storageConfig)
				realTimeTranscriptionService.RegisterRoutes(router)
			}
		}

		if rtmpURLExists {
			// Init RTMP Service
			rtmpURL := strings.Replace(rtmpURLEnv, "{appId}", appIDEnv, 1) // replace the place-holder value with appID from Env
			rtmpService := rtmp_service.NewRtmpService(appIDEnv, baseURLEnv, rtmpURL, basicAuthKey)
			rtmpService.RegisterRoutes(router)
		}
	} else {
		log.Print("WARNING: baseURLEnv Not Found - SKIPPING Cloud Recording, RTT and RTMP services ")
	}

	// Register healthcheck route
	router.GET("/ping", Ping)

	// Retrieve server port from environment variables or default to 8080.
	serverPort, exists := os.LookupEnv("SERVER_PORT")
	if !exists {
		serverPort = "8080"
	}

	// Configure and start the HTTP server.
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}

	// Start the server in a separate goroutine to handle graceful shutdown.
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		log.Printf("Server starting on port %s", serverPort)
	}()

	// Prepare to handle graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Wait for a shutdown signal.
	<-quit
	log.Println("Shutting down server...")

	// Attempt to gracefully shutdown the server with a timeout of 5 seconds.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

// Ping is a handler function for the "/ping" route. It serves as a basic health check endpoint.
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// getBasicAuth generates a basic authentication string from a customer ID and secret.
func getBasicAuth(customerID string, customerSecret string) string {
	auth := fmt.Sprintf("%s:%s", customerID, customerSecret)
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

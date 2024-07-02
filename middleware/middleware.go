package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	AllowOrigin string
}

func NewMiddleware(allowOrigin string) *Middleware {
	return &Middleware{AllowOrigin: allowOrigin}
}

func (m *Middleware) NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		// set headers
		c.Header("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Header("Expires", "-1")
		c.Header("Pragma", "no-cache")
	}
}

// Add CORSMiddleware to handle CORS requests and set the necessary headers
func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if !m.isOriginAllowed(origin) {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Origin not allowed",
			})
			c.Abort()
			return
		}
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

func (m *Middleware) isOriginAllowed(origin string) bool {
	if m.AllowOrigin == "*" {
		return true
	}

	allowedOrigins := strings.Split(m.AllowOrigin, ",")
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}

	return false
}

func (m *Middleware) TimestampMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed to the next middleware/handler
		c.Next()

		// Add the current timestamp to the response header
		timestamp := time.Now().Format(time.RFC3339)
		c.Writer.Header().Set("X-Timestamp", timestamp)
	}
}

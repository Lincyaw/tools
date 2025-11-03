package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lincyaw/tools/services/shortcode/internal/service"
)

// NewRouter creates router
func NewRouter(service service.ShortCodeService) *gin.Engine {
	// Set to release mode to improve performance
	// gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Middleware chain
	router.Use(gin.Recovery())              // Recovery middleware
	router.Use(errorHandlerMiddleware())    // Error handling middleware
	router.Use(requestIDMiddleware())       // Request ID middleware
	router.Use(loggerMiddleware())          // Logger middleware
	router.Use(securityHeadersMiddleware()) // Security headers middleware
	router.Use(corsMiddleware())            // CORS middleware

	// Create rate limiter: 100 requests per minute
	limiter := NewRateLimiter(100, time.Minute)
	router.Use(rateLimitMiddleware(limiter))        // Rate limiting middleware
	router.Use(timeoutMiddleware(30 * time.Second)) // Request timeout

	handler := NewHandler(service)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/shorten", handler.CreateShortCode)
		v1.GET("/stats/:code", handler.GetStats)
		v1.DELETE("/shorten/:code", handler.DeleteShortCode) // New delete functionality
	}

	// Health check
	router.GET("/health", handler.Health)
	router.GET("/metrics", handler.Metrics) // New metrics endpoint

	// Short link redirection (placed last to avoid conflicts)
	router.GET("/:code", handler.RedirectToOriginal)

	return router
}

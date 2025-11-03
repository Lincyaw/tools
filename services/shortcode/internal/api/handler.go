package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lincyaw/tools/services/shortcode/internal/model"
	"github.com/lincyaw/tools/services/shortcode/internal/service"
)

type Handler struct {
	service service.ShortCodeService
}

func NewHandler(service service.ShortCodeService) *Handler {
	return &Handler{
		service: service,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CreateShortCode create short link
// @Summary Create short link
// @Description Create a new short link
// @Tags shortcode
// @Accept json
// @Produce json
// @Param request body model.CreateShortCodeRequest true "Create short link request"
// @Success 201 {object} model.CreateShortCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/shorten [post]
func (h *Handler) CreateShortCode(c *gin.Context) {
	var req model.CreateShortCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	resp, err := h.service.CreateShortCode(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_url",
				Message: "The provided URL is not valid",
			})
		case errors.Is(err, service.ErrCodeExists):
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "code_exists",
				Message: "The custom code already exists",
			})
		case errors.Is(err, service.ErrInvalidCode):
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_code",
				Message: "The code format is invalid (4-50 alphanumeric characters)",
			})
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to create short code",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// RedirectToOriginal redirect to original URL
// @Summary Redirect to original URL
// @Description Redirect to original URL based on short code
// @Tags shortcode
// @Param code path string true "Short code"
// @Success 302 "Redirect to original URL"
// @Failure 404 {object} ErrorResponse
// @Router /{code} [get]
func (h *Handler) RedirectToOriginal(c *gin.Context) {
	code := c.Param("code")

	originalURL, err := h.service.GetOriginalURL(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Short code not found or expired",
		})
		return
	}

	// Asynchronously record click
	// Use context.Background() to avoid context cancellation after redirect
	go func() {
		ipAddress := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		referer := c.GetHeader("Referer")
		// Create a new context with timeout to avoid goroutine leak
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = h.service.RecordClick(ctx, code, ipAddress, userAgent, referer)
	}()

	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetStats get short link statistics
// @Summary Get statistics
// @Description Get short link statistics
// @Tags shortcode
// @Produce json
// @Param code path string true "Short code"
// @Success 200 {object} model.ShortCodeStats
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/stats/{code} [get]
func (h *Handler) GetStats(c *gin.Context) {
	code := c.Param("code")

	stats, err := h.service.GetStats(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Short code not found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Health health check
// @Summary Health check
// @Description Check service health status
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "shortcode",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// Metrics get service metrics
// @Summary Get metrics
// @Description Get service running metrics
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func (h *Handler) Metrics(c *gin.Context) {
	metrics, err := h.service.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to fetch metrics",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// DeleteShortCode delete short link
// @Summary Delete short link
// @Description Delete the specified short link
// @Tags shortcode
// @Produce json
// @Param code path string true "Short code"
// @Success 200 {object} map[string]string
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/shorten/{code} [delete]
func (h *Handler) DeleteShortCode(c *gin.Context) {
	code := c.Param("code")

	if err := h.service.DeleteShortCode(c.Request.Context(), code); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Short code not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Short code deleted successfully",
	})
}

// GetDetailedStats get detailed statistics with hourly buckets
// @Summary Get detailed statistics
// @Description Get detailed statistics including hourly access data and location information
// @Tags shortcode
// @Produce json
// @Param code path string true "Short code"
// @Param hours query int false "Number of hours to look back (default: all time)"
// @Success 200 {object} model.DetailedStats
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/stats/{code}/detailed [get]
func (h *Handler) GetDetailedStats(c *gin.Context) {
	code := c.Param("code")

	// Get hours parameter (default to 0 = all time)
	hours := 0
	if hoursParam := c.Query("hours"); hoursParam != "" {
		if h, err := time.ParseDuration(hoursParam + "h"); err == nil {
			hours = int(h.Hours())
		}
	}

	stats, err := h.service.GetDetailedStats(c.Request.Context(), code, hours)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Short code not found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

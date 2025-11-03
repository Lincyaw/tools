package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/lincyaw/tools/services/shortcode/internal/model"
	"github.com/lincyaw/tools/services/shortcode/internal/repository"
)

var (
	// ErrInvalidURL invalid URL
	ErrInvalidURL = errors.New("invalid URL")
	// ErrCodeExists code already exists
	ErrCodeExists = errors.New("code already exists")
	// ErrCodeNotFound code not found
	ErrCodeNotFound = errors.New("code not found")
	// ErrInvalidCode invalid code format
	ErrInvalidCode = errors.New("invalid code format")
)

const (
	defaultCodeLength = 6
	charset           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	maxRetries        = 5
)

// ShortCodeService short link service interface
type ShortCodeService interface {
	CreateShortCode(ctx context.Context, req *model.CreateShortCodeRequest) (*model.CreateShortCodeResponse, error)
	GetOriginalURL(ctx context.Context, code string) (string, error)
	GetStats(ctx context.Context, code string) (*model.ShortCodeStats, error)
	RecordClick(ctx context.Context, code, ipAddress, userAgent, referer string) error
	DeleteShortCode(ctx context.Context, code string) error
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
	GetDetailedStats(ctx context.Context, code string, hours int) (*model.DetailedStats, error)
}

type shortCodeService struct {
	repo    repository.ShortCodeRepository
	baseURL string
}

// NewShortCodeService creates short link service instance
func NewShortCodeService(repo repository.ShortCodeRepository, baseURL string) ShortCodeService {
	return &shortCodeService{
		repo:    repo,
		baseURL: baseURL,
	}
}

// CreateShortCode creates short link
func (s *shortCodeService) CreateShortCode(ctx context.Context, req *model.CreateShortCodeRequest) (*model.CreateShortCodeResponse, error) {
	// Validate URL
	if !isValidURL(req.URL) {
		return nil, ErrInvalidURL
	}

	var code string
	var err error

	// If custom code is provided
	if req.CustomCode != "" {
		if !isValidCode(req.CustomCode) {
			return nil, ErrInvalidCode
		}

		exists, err := s.repo.CodeExists(ctx, req.CustomCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check code existence: %w", err)
		}
		if exists {
			return nil, ErrCodeExists
		}
		code = req.CustomCode
	} else {
		// Generate random code
		code, err = s.generateUniqueCode(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to generate unique code: %w", err)
		}
	}

	// Calculate expiration time
	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		expiry := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		expiresAt = &expiry
	}

	// Create short code
	shortCode := &model.ShortCode{
		Code:        code,
		OriginalURL: req.URL,
		ExpiresAt:   expiresAt,
	}

	if err := s.repo.Create(ctx, shortCode); err != nil {
		return nil, fmt.Errorf("failed to create short code: %w", err)
	}

	return &model.CreateShortCodeResponse{
		ShortCode:   code,
		ShortURL:    fmt.Sprintf("%s/%s", s.baseURL, code),
		OriginalURL: req.URL,
		CreatedAt:   shortCode.CreatedAt,
		ExpiresAt:   expiresAt,
	}, nil
}

// GetOriginalURL gets original URL
func (s *shortCodeService) GetOriginalURL(ctx context.Context, code string) (string, error) {
	shortCode, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return "", ErrCodeNotFound
	}

	return shortCode.OriginalURL, nil
}

// GetStats gets statistics
func (s *shortCodeService) GetStats(ctx context.Context, code string) (*model.ShortCodeStats, error) {
	stats, err := s.repo.GetStats(ctx, code)
	if err != nil {
		return nil, ErrCodeNotFound
	}
	return stats, nil
}

// RecordClick records click
func (s *shortCodeService) RecordClick(ctx context.Context, code, ipAddress, userAgent, referer string) error {
	shortCode, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to get short code: %w", err)
	}

	// Update click count
	if err := s.repo.UpdateClickCount(ctx, shortCode.ID); err != nil {
		return fmt.Errorf("failed to update click count: %w", err)
	}

	// Record click log
	log := &model.ClickLog{
		ShortCodeID: shortCode.ID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Referer:     referer,
	}

	if err := s.repo.LogClick(ctx, log); err != nil {
		return fmt.Errorf("failed to log click: %w", err)
	}

	// Get IP location information
	location := s.getIPLocation(ipAddress)

	// Record access statistics with hourly bucket
	hourBucket := time.Now().Truncate(time.Hour)
	stats := &model.AccessStatistics{
		ShortCodeID: shortCode.ID,
		IPAddress:   ipAddress,
		Country:     location.Country,
		Region:      location.Region,
		City:        location.City,
		HourBucket:  hourBucket,
	}

	if err := s.repo.RecordAccessStats(ctx, stats); err != nil {
		return fmt.Errorf("failed to record access stats: %w", err)
	}

	return nil
}

// DeleteShortCode deletes short link
func (s *shortCodeService) DeleteShortCode(ctx context.Context, code string) error {
	return s.repo.Delete(ctx, code)
}

// GetMetrics gets service metrics
func (s *shortCodeService) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	return s.repo.GetMetrics(ctx)
}

// generateUniqueCode generates unique code
func (s *shortCodeService) generateUniqueCode(ctx context.Context) (string, error) {
	for i := 0; i < maxRetries; i++ {
		code := generateRandomCode(defaultCodeLength)

		exists, err := s.repo.CodeExists(ctx, code)
		if err != nil {
			return "", err
		}

		if !exists {
			return code, nil
		}
	}

	return "", errors.New("failed to generate unique code after max retries")
}

// generateRandomCode generates random code
func generateRandomCode(length int) string {
	code := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		randomIndex, _ := rand.Int(rand.Reader, charsetLen)
		code[i] = charset[randomIndex.Int64()]
	}

	return string(code)
}

// isValidURL validates if URL is valid
func isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

// isValidCode validates if code format is valid
func isValidCode(code string) bool {
	// Code can only contain letters and numbers, length between 4-50
	re := regexp.MustCompile(`^[a-zA-Z0-9]{4,50}$`)
	return re.MatchString(code)
}

// GetDetailedStats gets detailed statistics
func (s *shortCodeService) GetDetailedStats(ctx context.Context, code string, hours int) (*model.DetailedStats, error) {
	stats, err := s.repo.GetDetailedStats(ctx, code, hours)
	if err != nil {
		return nil, ErrCodeNotFound
	}
	return stats, nil
}

// getIPLocation gets IP location information
func (s *shortCodeService) getIPLocation(ipAddress string) model.IPLocation {
	// Default location
	location := model.IPLocation{
		Country: "Unknown",
		Region:  "Unknown",
		City:    "Unknown",
	}

	// Skip for local/private IPs
	if isPrivateIP(ipAddress) {
		location.Country = "Private"
		location.Region = "Local"
		location.City = "Local"
		return location
	}

	// Use ip-api.com free API (limited to 45 requests per minute)
	// In production, consider using a paid service or caching
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,regionName,city", ipAddress)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return location
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return location
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return location
	}

	var result struct {
		Status     string `json:"status"`
		Country    string `json:"country"`
		RegionName string `json:"regionName"`
		City       string `json:"city"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return location
	}

	if result.Status == "success" {
		if result.Country != "" {
			location.Country = result.Country
		}
		if result.RegionName != "" {
			location.Region = result.RegionName
		}
		if result.City != "" {
			location.City = result.City
		}
	}

	return location
}

// isPrivateIP checks if IP is private/local
func isPrivateIP(ip string) bool {
	// Simple check for common private IP ranges and localhost
	if ip == "" || ip == "::1" || ip == "localhost" {
		return true
	}

	// Check for private IPv4 ranges
	privateRanges := []string{
		"10.", "172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
		"172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
		"172.30.", "172.31.", "192.168.", "127.",
	}

	for _, prefix := range privateRanges {
		if len(ip) >= len(prefix) && ip[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

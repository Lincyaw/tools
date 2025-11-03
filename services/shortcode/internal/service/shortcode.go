package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
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

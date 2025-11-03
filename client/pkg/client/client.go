package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
			CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
				// Do not automatically follow redirects to test redirect functionality
				return http.ErrUseLastResponse
			},
		},
	}
}

// CreateShortCodeRequest create short link request
type CreateShortCodeRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
	ExpiresIn  int    `json:"expires_in,omitempty"`
}

// CreateShortCodeResponse create short link response
type CreateShortCodeResponse struct {
	ShortCode   string     `json:"short_code"`
	ShortURL    string     `json:"short_url"`
	OriginalURL string     `json:"original_url"`
	CreatedAt   time.Time  `json:"created_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// ShortCodeStats short link statistics
type ShortCodeStats struct {
	Code           string     `json:"code"`
	OriginalURL    string     `json:"original_url"`
	ClickCount     int64      `json:"click_count"`
	CreatedAt      time.Time  `json:"created_at"`
	LastAccessedAt *time.Time `json:"last_accessed_at,omitempty"`
}

// DetailedStats detailed statistics response
type DetailedStats struct {
	Code           string             `json:"code"`
	OriginalURL    string             `json:"original_url"`
	TotalClicks    int64              `json:"total_clicks"`
	UniqueIPs      int64              `json:"unique_ips"`
	CreatedAt      time.Time          `json:"created_at"`
	LastAccessedAt *time.Time         `json:"last_accessed_at,omitempty"`
	HourlyStats    []HourlyStatItem   `json:"hourly_stats"`
	LocationStats  []LocationStatItem `json:"location_stats"`
	RecentAccesses []RecentAccessItem `json:"recent_accesses"`
}

// HourlyStatItem hourly statistics item
type HourlyStatItem struct {
	HourBucket  time.Time `json:"hour_bucket"`
	AccessCount int64     `json:"access_count"`
	UniqueIPs   int64     `json:"unique_ips"`
}

// LocationStatItem location statistics item
type LocationStatItem struct {
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	AccessCount int64  `json:"access_count"`
}

// RecentAccessItem recent access item
type RecentAccessItem struct {
	IPAddress  string    `json:"ip_address"`
	Country    string    `json:"country"`
	Region     string    `json:"region"`
	City       string    `json:"city"`
	AccessTime time.Time `json:"access_time"`
	UserAgent  string    `json:"user_agent"`
}

// ErrorResponse error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// RedirectInfo redirect information
type RedirectInfo struct {
	StatusCode  int
	Location    string
	OriginalURL string
}

// CreateShortCode create short link
func (c *Client) CreateShortCode(req CreateShortCodeRequest) (*CreateShortCodeResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/v1/shorten",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Message)
	}

	var result CreateShortCodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &result, nil
}

// GetStats get short link statistics
func (c *Client) GetStats(code string) (*ShortCodeStats, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/v1/stats/" + code)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Message)
	}

	var stats ShortCodeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &stats, nil
}

// GetDetailedStats get detailed short link statistics
func (c *Client) GetDetailedStats(code string, hours int) (*DetailedStats, error) {
	url := fmt.Sprintf("%s/api/v1/stats/%s/detailed", c.BaseURL, code)
	if hours > 0 {
		url = fmt.Sprintf("%s?hours=%d", url, hours)
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Message)
	}

	var stats DetailedStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &stats, nil
}

// TestRedirect test short link redirect
func (c *Client) TestRedirect(code string) (*RedirectInfo, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/" + code)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 300 || resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("expected redirect, got status %d: %s", resp.StatusCode, string(body))
	}

	location := resp.Header.Get("Location")
	if location == "" {
		return nil, fmt.Errorf("redirect response missing Location header")
	}

	return &RedirectInfo{
		StatusCode:  resp.StatusCode,
		Location:    location,
		OriginalURL: location,
	}, nil
}

// DeleteShortCode delete short link
func (c *Client) DeleteShortCode(code string) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+"/api/v1/shorten/"+code, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		var errResp ErrorResponse
		if err := json.Unmarshal(body, &errResp); err != nil {
			return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}
		return fmt.Errorf("API error (%d): %s - %s", resp.StatusCode, errResp.Error, errResp.Message)
	}

	return nil
}

// HealthCheck health check
func (c *Client) HealthCheck() error {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/")
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return fmt.Errorf("service unavailable: status %d", resp.StatusCode)
	}

	return nil
}

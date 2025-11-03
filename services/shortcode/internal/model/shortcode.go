package model

import (
	"time"

	"gorm.io/gorm"
)

// ShortCode short link model
type ShortCode struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"uniqueIndex;size:50;not null" json:"code"`
	OriginalURL    string         `gorm:"type:text;not null" json:"original_url"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	ExpiresAt      *time.Time     `gorm:"index" json:"expires_at,omitempty"`
	ClickCount     int64          `gorm:"default:0" json:"click_count"`
	LastAccessedAt *time.Time     `json:"last_accessed_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specify table name
func (ShortCode) TableName() string {
	return "short_codes"
}

// ClickLog click log model
type ClickLog struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ShortCodeID uint      `gorm:"index;not null" json:"short_code_id"`
	ShortCode   ShortCode `gorm:"foreignKey:ShortCodeID" json:"-"`
	IPAddress   string    `gorm:"size:45" json:"ip_address"`
	UserAgent   string    `gorm:"type:text" json:"user_agent"`
	Referer     string    `gorm:"type:text" json:"referer"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specify table name
func (ClickLog) TableName() string {
	return "click_logs"
}

// CreateShortCodeRequest create short link request
type CreateShortCodeRequest struct {
	URL        string `json:"url" binding:"required,url"`
	CustomCode string `json:"custom_code,omitempty" binding:"omitempty,min=4,max=50,alphanum"`
	ExpiresIn  int    `json:"expires_in,omitempty" binding:"omitempty,min=1"` // Expiration time (hours)
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

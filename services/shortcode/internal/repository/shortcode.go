package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/lincyaw/tools/services/shortcode/internal/config"
	"github.com/lincyaw/tools/services/shortcode/internal/model"
)

// ShortCodeRepository short link repository interface
type ShortCodeRepository interface {
	Create(ctx context.Context, shortCode *model.ShortCode) error
	GetByCode(ctx context.Context, code string) (*model.ShortCode, error)
	UpdateClickCount(ctx context.Context, id uint) error
	GetStats(ctx context.Context, code string) (*model.ShortCodeStats, error)
	LogClick(ctx context.Context, log *model.ClickLog) error
	CodeExists(ctx context.Context, code string) (bool, error)
	Delete(ctx context.Context, code string) error
	InvalidateCache(ctx context.Context, code string) error
	GetMetrics(ctx context.Context) (map[string]interface{}, error)
}

type shortCodeRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

// NewPostgresDB create PostgreSQL database connection
func NewPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get the underlying *sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Auto migrate
	if err := db.AutoMigrate(&model.ShortCode{}, &model.ClickLog{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// NewRedisClient create Redis client
func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return client
}

// NewShortCodeRepository create short link repository instance
func NewShortCodeRepository(db *gorm.DB, redisClient *redis.Client) ShortCodeRepository {
	return &shortCodeRepository{
		db:          db,
		redisClient: redisClient,
	}
}

// Create create short link
func (r *shortCodeRepository) Create(ctx context.Context, shortCode *model.ShortCode) error {
	return r.db.WithContext(ctx).Create(shortCode).Error
}

// GetByCode get short link by code
func (r *shortCodeRepository) GetByCode(ctx context.Context, code string) (*model.ShortCode, error) {
	// First try to get from cache
	cacheKey := fmt.Sprintf("shortcode:%s", code)
	cached, err := r.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var shortCode model.ShortCode
		if err := json.Unmarshal([]byte(cached), &shortCode); err == nil {
			return &shortCode, nil
		}
	}

	// Get from database
	var shortCode model.ShortCode
	err = r.db.WithContext(ctx).
		Where("code = ?", code).
		Where("expires_at IS NULL OR expires_at > ?", time.Now()).
		First(&shortCode).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("short code not found")
		}
		return nil, err
	}

	// Cache to Redis (24 hours)
	if data, err := json.Marshal(shortCode); err == nil {
		r.redisClient.Set(ctx, cacheKey, data, 24*time.Hour)
	}

	return &shortCode, nil
}

// UpdateClickCount update click count
func (r *shortCodeRepository) UpdateClickCount(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.ShortCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"click_count":      gorm.Expr("click_count + ?", 1),
			"last_accessed_at": now,
		}).Error
}

// GetStats get statistics
func (r *shortCodeRepository) GetStats(ctx context.Context, code string) (*model.ShortCodeStats, error) {
	var shortCode model.ShortCode
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&shortCode).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("short code not found")
		}
		return nil, err
	}

	stats := &model.ShortCodeStats{
		Code:           shortCode.Code,
		OriginalURL:    shortCode.OriginalURL,
		ClickCount:     shortCode.ClickCount,
		CreatedAt:      shortCode.CreatedAt,
		LastAccessedAt: shortCode.LastAccessedAt,
	}

	return stats, nil
}

// LogClick log click
func (r *shortCodeRepository) LogClick(ctx context.Context, log *model.ClickLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// CodeExists check if code exists
func (r *shortCodeRepository) CodeExists(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.ShortCode{}).
		Where("code = ?", code).
		Count(&count).Error

	return count > 0, err
}

// Delete delete short link
func (r *shortCodeRepository) Delete(ctx context.Context, code string) error {
	// Delete cache
	if err := r.InvalidateCache(ctx, code); err != nil {
		// Log cache invalidation error but continue with deletion
		log.Printf("Warning: Failed to invalidate cache for code %s: %v", code, err)
	}

	// Delete from database
	result := r.db.WithContext(ctx).
		Where("code = ?", code).
		Delete(&model.ShortCode{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("short code not found")
	}

	return nil
}

// InvalidateCache invalidate cache
func (r *shortCodeRepository) InvalidateCache(ctx context.Context, code string) error {
	cacheKey := fmt.Sprintf("shortcode:%s", code)
	return r.redisClient.Del(ctx, cacheKey).Err()
}

// GetMetrics get system metrics
func (r *shortCodeRepository) GetMetrics(ctx context.Context) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Count total short codes
	var totalCodes int64
	if err := r.db.WithContext(ctx).Model(&model.ShortCode{}).Count(&totalCodes).Error; err != nil {
		return nil, err
	}
	metrics["total_codes"] = totalCodes

	// Count total clicks
	var totalClicks int64
	if err := r.db.WithContext(ctx).Model(&model.ShortCode{}).Select("COALESCE(SUM(click_count), 0)").Scan(&totalClicks).Error; err != nil {
		return nil, err
	}
	metrics["total_clicks"] = totalClicks

	// Count clicks in the past 24 hours
	var clicks24h int64
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	if err := r.db.WithContext(ctx).Model(&model.ClickLog{}).
		Where("created_at > ?", oneDayAgo).
		Count(&clicks24h).Error; err != nil {
		return nil, err
	}
	metrics["clicks_24h"] = clicks24h

	// Count active short codes (with clicks)
	var activeCodes int64
	if err := r.db.WithContext(ctx).Model(&model.ShortCode{}).
		Where("click_count > 0").
		Count(&activeCodes).Error; err != nil {
		return nil, err
	}
	metrics["active_codes"] = activeCodes

	return metrics, nil
}

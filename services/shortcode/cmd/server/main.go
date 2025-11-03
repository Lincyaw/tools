package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lincyaw/tools/services/shortcode/internal/api"
	"github.com/lincyaw/tools/services/shortcode/internal/config"
	"github.com/lincyaw/tools/services/shortcode/internal/repository"
	"github.com/lincyaw/tools/services/shortcode/internal/service"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	log.Println("Database connection established")

	// Initialize Redis
	redisClient := repository.NewRedisClient(cfg.Redis)
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("Redis connection established")
	}

	// Initialize repository layer
	repo := repository.NewShortCodeRepository(db, redisClient)

	// Initialize service layer
	svc := service.NewShortCodeService(repo, cfg.BaseURL)

	// Initialize HTTP server
	router := api.NewRouter(svc)

	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server
	go func() {
		log.Printf("Starting server on port %s (environment: %s)", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)

	if err := srv.Shutdown(shutdownCtx); err != nil {
		shutdownCancel()
		log.Printf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}
	shutdownCancel()

	log.Println("Server exited gracefully")
}

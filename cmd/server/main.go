package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/config"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/handlers"
	"github.com/quckapp/media-service/internal/services"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize Redis
	redisClient := database.NewRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, 7)
	defer redisClient.Close()

	// Initialize S3 storage
	s3Storage, err := services.NewS3Storage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}

	// Initialize services
	mediaService := services.NewMediaService(mongoDB, redisClient, s3Storage)

	// Initialize handlers
	mediaHandler := handlers.NewMediaHandler(mediaService)
	healthHandler := handlers.NewHealthHandler(mongoDB, redisClient)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Health endpoints
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Ready)

	// API routes
	api := router.Group("/api/v1/media")
	api.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		api.POST("/upload", mediaHandler.Upload)
		api.POST("/upload/presigned", mediaHandler.GetPresignedURL)
		api.GET("/:id", mediaHandler.Get)
		api.DELETE("/:id", mediaHandler.Delete)
		api.GET("/user/:userId", mediaHandler.GetUserMedia)
		api.POST("/:id/thumbnail", mediaHandler.GenerateThumbnail)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Media service starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Media service stopped")
}

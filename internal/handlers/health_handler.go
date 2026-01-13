package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quckchat/media-service/internal/database"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	db    *database.MongoDB
	redis *redis.Client
}

func NewHealthHandler(db *database.MongoDB, redis *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, redis: redis}
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "media-service",
		"version": "1.0.0",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	mongoOK := h.db.Ping() == nil
	redisOK := h.redis.Ping(ctx).Err() == nil

	allHealthy := mongoOK && redisOK
	statusCode := http.StatusOK
	if !allHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"ready": allHealthy,
		"checks": gin.H{
			"mongo": mongoOK,
			"redis": redisOK,
		},
	})
}

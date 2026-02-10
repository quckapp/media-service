package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/services"
)

type AnalyticsHandler struct {
	service *services.AnalyticsService
}

func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetUploadTrends(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trends, err := h.service.GetUploadTrends(c.Request.Context(), workspaceID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": trends})
}

func (h *AnalyticsHandler) GetStorageTrends(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))

	trends, err := h.service.GetStorageTrends(c.Request.Context(), workspaceID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": trends})
}

func (h *AnalyticsHandler) GetFileTypeDistribution(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	dist, err := h.service.GetFileTypeDistribution(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": dist})
}

func (h *AnalyticsHandler) GetUserUploadStats(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	stats, err := h.service.GetUserUploadStats(c.Request.Context(), workspaceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

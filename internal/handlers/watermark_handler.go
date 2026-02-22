package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type WatermarkHandler struct {
	service *services.WatermarkService
}

func NewWatermarkHandler(service *services.WatermarkService) *WatermarkHandler {
	return &WatermarkHandler{service: service}
}

func (h *WatermarkHandler) Upload(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.UploadWatermarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	watermark, err := h.service.Upload(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": watermark})
}

func (h *WatermarkHandler) List(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	watermarks, err := h.service.List(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": watermarks})
}

func (h *WatermarkHandler) Apply(c *gin.Context) {
	var req models.ApplyWatermarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.Apply(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Watermark applied"})
}

func (h *WatermarkHandler) Remove(c *gin.Context) {
	mediaID := c.Param("mediaId")

	if err := h.service.Remove(c.Request.Context(), mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Watermark removed"})
}

func (h *WatermarkHandler) GetSettings(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	settings, err := h.service.GetSettings(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Settings not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": settings})
}

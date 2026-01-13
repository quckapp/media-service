package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckchat/media-service/internal/models"
	"github.com/quckchat/media-service/internal/services"
)

type MediaHandler struct {
	service *services.MediaService
}

func NewMediaHandler(service *services.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

func (h *MediaHandler) Upload(c *gin.Context) {
	userID := c.GetString("userID")
	
	var req models.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	media, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *MediaHandler) GetPresignedURL(c *gin.Context) {
	userID := c.GetString("userID")
	
	var req models.UploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	resp, err := h.service.GetPresignedUploadURL(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

func (h *MediaHandler) Get(c *gin.Context) {
	mediaID := c.Param("id")
	
	media, err := h.service.Get(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Media not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *MediaHandler) Delete(c *gin.Context) {
	mediaID := c.Param("id")
	userID := c.GetString("userID")

	if err := h.service.Delete(c.Request.Context(), mediaID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Media deleted"})
}

func (h *MediaHandler) GetUserMedia(c *gin.Context) {
	userID := c.Param("userId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	media, err := h.service.GetUserMedia(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *MediaHandler) GenerateThumbnail(c *gin.Context) {
	// Placeholder for thumbnail generation
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Thumbnail generation queued"})
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
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

func (h *MediaHandler) BulkDelete(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result := h.service.BulkDelete(c.Request.Context(), req.MediaIDs, userID)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *MediaHandler) GetChannelMedia(c *gin.Context) {
	channelID := c.Param("channelId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	media, err := h.service.GetMediaByChannel(c.Request.Context(), channelID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *MediaHandler) GetWorkspaceMedia(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	media, err := h.service.GetMediaByWorkspace(c.Request.Context(), workspaceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *MediaHandler) GetUserStats(c *gin.Context) {
	userID := c.Param("userId")

	stats, err := h.service.GetUserStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *MediaHandler) UpdateMetadata(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var metadata map[string]string
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.UpdateMetadata(c.Request.Context(), mediaID, userID, metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Metadata updated"})
}

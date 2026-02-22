package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type SearchHandler struct {
	service     *services.SearchService
	mediaSvc    *services.MediaService
}

func NewSearchHandler(service *services.SearchService, mediaSvc *services.MediaService) *SearchHandler {
	return &SearchHandler{service: service, mediaSvc: mediaSvc}
}

func (h *SearchHandler) Search(c *gin.Context) {
	userID := c.GetString("userID")

	var params models.MediaSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	media, total, err := h.service.Search(c.Request.Context(), userID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    media,
		"total":   total,
		"page":    params.Page,
		"limit":   params.Limit,
	})
}

func (h *SearchHandler) GetWorkspaceStats(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	stats, err := h.service.GetWorkspaceStats(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

func (h *SearchHandler) GetDuplicates(c *gin.Context) {
	userID := c.GetString("userID")

	media, err := h.service.GetDuplicates(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *SearchHandler) Rename(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.RenameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.mediaSvc.Rename(c.Request.Context(), mediaID, userID, req.Filename); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Media renamed"})
}

func (h *SearchHandler) CopyMedia(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.CopyMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	newMedia, err := h.mediaSvc.CopyMedia(c.Request.Context(), mediaID, userID, req.TargetWorkspaceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": newMedia})
}

func (h *SearchHandler) MoveMedia(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.MoveMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.mediaSvc.MoveMedia(c.Request.Context(), mediaID, userID, req.TargetWorkspaceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Media moved"})
}

func (h *SearchHandler) BulkMove(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.BulkMoveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	result := h.mediaSvc.BulkMove(c.Request.Context(), req.MediaIDs, userID, req.TargetWorkspaceID)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (h *SearchHandler) GetByType(c *gin.Context) {
	userID := c.Param("userId")
	mediaType := c.Param("type")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	media, err := h.mediaSvc.GetMediaByType(c.Request.Context(), userID, mediaType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

func (h *SearchHandler) GetDownloadURL(c *gin.Context) {
	mediaID := c.Param("id")

	url, err := h.mediaSvc.GetDownloadURL(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Media not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"downloadUrl": url}})
}

func (h *SearchHandler) GetRecent(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	media, err := h.mediaSvc.GetRecentMedia(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": media})
}

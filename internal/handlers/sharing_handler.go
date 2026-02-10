package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type SharingHandler struct {
	service *services.SharingService
}

func NewSharingHandler(service *services.SharingService) *SharingHandler {
	return &SharingHandler{service: service}
}

func (h *SharingHandler) ShareWithUser(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.ShareMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	share, err := h.service.ShareWithUser(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": share})
}

func (h *SharingHandler) GetSharedWithMe(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	shares, err := h.service.GetSharedWithUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": shares})
}

func (h *SharingHandler) GetSharedByMe(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	shares, err := h.service.GetSharedByUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": shares})
}

func (h *SharingHandler) RevokeShare(c *gin.Context) {
	userID := c.GetString("userID")
	shareID := c.Param("shareId")

	if err := h.service.RevokeShare(c.Request.Context(), shareID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Share revoked"})
}

func (h *SharingHandler) CreateShareLink(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	link, err := h.service.CreateShareLink(c.Request.Context(), mediaID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": link})
}

func (h *SharingHandler) GetShareLink(c *gin.Context) {
	token := c.Param("token")

	link, err := h.service.GetShareLink(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Share link not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": link})
}

func (h *SharingHandler) DeactivateShareLink(c *gin.Context) {
	userID := c.GetString("userID")
	linkID := c.Param("linkId")

	if err := h.service.DeactivateShareLink(c.Request.Context(), linkID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Share link deactivated"})
}

func (h *SharingHandler) GetShareLinks(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	links, err := h.service.GetShareLinks(c.Request.Context(), mediaID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": links})
}

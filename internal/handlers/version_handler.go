package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type VersionHandler struct {
	service *services.VersionService
}

func NewVersionHandler(service *services.VersionService) *VersionHandler {
	return &VersionHandler{service: service}
}

func (h *VersionHandler) CreateVersion(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.CreateVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	version, err := h.service.CreateVersion(c.Request.Context(), mediaID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": version})
}

func (h *VersionHandler) GetVersions(c *gin.Context) {
	mediaID := c.Param("id")

	versions, err := h.service.GetVersions(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": versions})
}

func (h *VersionHandler) GetVersion(c *gin.Context) {
	versionID := c.Param("versionId")

	version, err := h.service.GetVersion(c.Request.Context(), versionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Version not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": version})
}

func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	userID := c.GetString("userID")
	versionID := c.Param("versionId")

	if err := h.service.DeleteVersion(c.Request.Context(), versionID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Version deleted"})
}

func (h *VersionHandler) RestoreVersion(c *gin.Context) {
	userID := c.GetString("userID")
	versionID := c.Param("versionId")

	if err := h.service.RestoreVersion(c.Request.Context(), versionID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Version restored"})
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type ScanningHandler struct {
	service *services.ScanningService
}

func NewScanningHandler(service *services.ScanningService) *ScanningHandler {
	return &ScanningHandler{service: service}
}

func (h *ScanningHandler) ScanMedia(c *gin.Context) {
	var req models.ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	scan, err := h.service.ScanMedia(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": scan})
}

func (h *ScanningHandler) GetResults(c *gin.Context) {
	mediaID := c.Param("mediaId")

	scans, err := h.service.GetScanResults(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": scans})
}

func (h *ScanningHandler) ListFlagged(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	scans, err := h.service.ListFlagged(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": scans})
}

func (h *ScanningHandler) UpdateStatus(c *gin.Context) {
	scanID := c.Param("scanId")
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "status is required"})
		return
	}

	if err := h.service.UpdateStatus(c.Request.Context(), scanID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Scan status updated"})
}

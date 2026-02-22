package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type QuotaHandler struct {
	service *services.QuotaService
}

func NewQuotaHandler(service *services.QuotaService) *QuotaHandler {
	return &QuotaHandler{service: service}
}

func (h *QuotaHandler) GetQuota(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	quota, err := h.service.GetByWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Quota not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": quota})
}

func (h *QuotaHandler) SetQuota(c *gin.Context) {
	var req models.SetQuotaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	quota, err := h.service.SetQuota(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": quota})
}

func (h *QuotaHandler) GetUsage(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	quota, err := h.service.GetUsage(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Usage not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": quota})
}

func (h *QuotaHandler) ListOverQuota(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	quotas, err := h.service.ListOverQuota(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": quotas})
}

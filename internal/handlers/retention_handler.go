package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type RetentionHandler struct {
	service *services.RetentionService
}

func NewRetentionHandler(service *services.RetentionService) *RetentionHandler {
	return &RetentionHandler{service: service}
}

func (h *RetentionHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.CreateRetentionPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	policy, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": policy})
}

func (h *RetentionHandler) Get(c *gin.Context) {
	policyID := c.Param("policyId")

	policy, err := h.service.GetByID(c.Request.Context(), policyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Policy not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": policy})
}

func (h *RetentionHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	policyID := c.Param("policyId")

	var req models.UpdateRetentionPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	policy, err := h.service.Update(c.Request.Context(), policyID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": policy})
}

func (h *RetentionHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	policyID := c.Param("policyId")

	if err := h.service.Delete(c.Request.Context(), policyID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Policy deleted"})
}

func (h *RetentionHandler) GetByWorkspace(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	policies, err := h.service.GetByWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": policies})
}

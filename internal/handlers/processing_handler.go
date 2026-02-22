package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type ProcessingHandler struct {
	service *services.ProcessingService
}

func NewProcessingHandler(service *services.ProcessingService) *ProcessingHandler {
	return &ProcessingHandler{service: service}
}

func (h *ProcessingHandler) CreateJob(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.CreateProcessingJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	job, err := h.service.CreateJob(c.Request.Context(), mediaID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": job})
}

func (h *ProcessingHandler) GetJob(c *gin.Context) {
	jobID := c.Param("jobId")

	job, err := h.service.GetJob(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": job})
}

func (h *ProcessingHandler) GetJobsByMedia(c *gin.Context) {
	mediaID := c.Param("id")

	jobs, err := h.service.GetJobsByMedia(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": jobs})
}

func (h *ProcessingHandler) GetUserJobs(c *gin.Context) {
	userID := c.GetString("userID")
	status := c.Query("status")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	jobs, err := h.service.GetUserJobs(c.Request.Context(), userID, status, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": jobs})
}

func (h *ProcessingHandler) CancelJob(c *gin.Context) {
	userID := c.GetString("userID")
	jobID := c.Param("jobId")

	if err := h.service.CancelJob(c.Request.Context(), jobID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Job cancelled"})
}

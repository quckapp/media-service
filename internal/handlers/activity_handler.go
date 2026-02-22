package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/services"
)

type ActivityHandler struct {
	service *services.ActivityService
}

func NewActivityHandler(service *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) GetByMedia(c *gin.Context) {
	mediaID := c.Param("id")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	activities, err := h.service.GetByMedia(c.Request.Context(), mediaID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": activities})
}

func (h *ActivityHandler) GetByUser(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	activities, err := h.service.GetByUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": activities})
}

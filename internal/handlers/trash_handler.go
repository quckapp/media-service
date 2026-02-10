package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/services"
)

type TrashHandler struct {
	service *services.TrashService
}

func NewTrashHandler(service *services.TrashService) *TrashHandler {
	return &TrashHandler{service: service}
}

func (h *TrashHandler) MoveToTrash(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	if err := h.service.MoveToTrash(c.Request.Context(), mediaID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Moved to trash"})
}

func (h *TrashHandler) RestoreFromTrash(c *gin.Context) {
	userID := c.GetString("userID")
	trashID := c.Param("trashId")

	if err := h.service.RestoreFromTrash(c.Request.Context(), trashID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Restored from trash"})
}

func (h *TrashHandler) GetTrash(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	trashed, err := h.service.GetTrash(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": trashed})
}

func (h *TrashHandler) PermanentDelete(c *gin.Context) {
	userID := c.GetString("userID")
	trashID := c.Param("trashId")

	if err := h.service.PermanentDelete(c.Request.Context(), trashID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Permanently deleted"})
}

func (h *TrashHandler) EmptyTrash(c *gin.Context) {
	userID := c.GetString("userID")

	count, err := h.service.EmptyTrash(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Trash emptied", "deleted": count})
}

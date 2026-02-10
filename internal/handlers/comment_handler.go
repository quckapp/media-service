package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	comment, err := h.service.Create(c.Request.Context(), mediaID, userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": comment})
}

func (h *CommentHandler) GetByMedia(c *gin.Context) {
	mediaID := c.Param("id")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	comments, err := h.service.GetByMedia(c.Request.Context(), mediaID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": comments})
}

func (h *CommentHandler) GetReplies(c *gin.Context) {
	commentID := c.Param("commentId")

	replies, err := h.service.GetReplies(c.Request.Context(), commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": replies})
}

func (h *CommentHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	commentID := c.Param("commentId")

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	comment, err := h.service.Update(c.Request.Context(), commentID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": comment})
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	commentID := c.Param("commentId")

	if err := h.service.Delete(c.Request.Context(), commentID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Comment deleted"})
}

func (h *CommentHandler) CountByMedia(c *gin.Context) {
	mediaID := c.Param("id")

	count, err := h.service.CountByMedia(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
}

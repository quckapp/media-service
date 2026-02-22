package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type TagHandler struct {
	service *services.TagService
}

func NewTagHandler(service *services.TagService) *TagHandler {
	return &TagHandler{service: service}
}

func (h *TagHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	tag, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": tag})
}

func (h *TagHandler) GetByUser(c *gin.Context) {
	userID := c.Param("userId")

	tags, err := h.service.GetByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *TagHandler) GetByWorkspace(c *gin.Context) {
	workspaceID := c.Param("workspaceId")

	tags, err := h.service.GetByWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *TagHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	tagID := c.Param("tagId")

	var req models.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	tag, err := h.service.Update(c.Request.Context(), tagID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tag})
}

func (h *TagHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	tagID := c.Param("tagId")

	if err := h.service.Delete(c.Request.Context(), tagID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tag deleted"})
}

func (h *TagHandler) TagMedia(c *gin.Context) {
	mediaID := c.Param("id")

	var req models.TagMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.TagMedia(c.Request.Context(), mediaID, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tags added"})
}

func (h *TagHandler) UntagMedia(c *gin.Context) {
	mediaID := c.Param("id")
	tagID := c.Param("tagId")

	if err := h.service.UntagMedia(c.Request.Context(), mediaID, tagID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Tag removed"})
}

func (h *TagHandler) GetMediaTags(c *gin.Context) {
	mediaID := c.Param("id")

	tags, err := h.service.GetMediaTags(c.Request.Context(), mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": tags})
}

func (h *TagHandler) GetMediaByTag(c *gin.Context) {
	tagID := c.Param("tagId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	mediaIDs, err := h.service.GetMediaByTag(c.Request.Context(), tagID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": mediaIDs})
}

func (h *TagHandler) BulkTag(c *gin.Context) {
	var req models.BulkTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.BulkTag(c.Request.Context(), req.MediaIDs, req.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Bulk tag applied"})
}

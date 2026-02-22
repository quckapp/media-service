package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type GalleryHandler struct {
	service *services.GalleryService
}

func NewGalleryHandler(service *services.GalleryService) *GalleryHandler {
	return &GalleryHandler{service: service}
}

func (h *GalleryHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.CreateGalleryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	gallery, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": gallery})
}

func (h *GalleryHandler) List(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	galleries, err := h.service.List(c.Request.Context(), workspaceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": galleries})
}

func (h *GalleryHandler) Get(c *gin.Context) {
	galleryID := c.Param("galleryId")

	gallery, err := h.service.GetByID(c.Request.Context(), galleryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Gallery not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gallery})
}

func (h *GalleryHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	galleryID := c.Param("galleryId")

	var req models.UpdateGalleryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	gallery, err := h.service.Update(c.Request.Context(), galleryID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gallery})
}

func (h *GalleryHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	galleryID := c.Param("galleryId")

	if err := h.service.Delete(c.Request.Context(), galleryID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Gallery deleted"})
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/models"
	"github.com/quckapp/media-service/internal/services"
)

type AlbumHandler struct {
	service *services.AlbumService
}

func NewAlbumHandler(service *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: service}
}

func (h *AlbumHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")

	var req models.CreateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	album, err := h.service.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": album})
}

func (h *AlbumHandler) GetByID(c *gin.Context) {
	albumID := c.Param("albumId")

	album, err := h.service.GetByID(c.Request.Context(), albumID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": album})
}

func (h *AlbumHandler) GetByUser(c *gin.Context) {
	userID := c.Param("userId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	albums, err := h.service.GetByUser(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": albums})
}

func (h *AlbumHandler) GetByWorkspace(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	albums, err := h.service.GetByWorkspace(c.Request.Context(), workspaceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": albums})
}

func (h *AlbumHandler) GetPublicByWorkspace(c *gin.Context) {
	workspaceID := c.Param("workspaceId")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	albums, err := h.service.GetPublicByWorkspace(c.Request.Context(), workspaceID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": albums})
}

func (h *AlbumHandler) Update(c *gin.Context) {
	userID := c.GetString("userID")
	albumID := c.Param("albumId")

	var req models.UpdateAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	album, err := h.service.Update(c.Request.Context(), albumID, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": album})
}

func (h *AlbumHandler) Delete(c *gin.Context) {
	userID := c.GetString("userID")
	albumID := c.Param("albumId")

	if err := h.service.Delete(c.Request.Context(), albumID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Album deleted"})
}

func (h *AlbumHandler) AddMedia(c *gin.Context) {
	userID := c.GetString("userID")
	albumID := c.Param("albumId")

	var req models.AddToAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.AddMedia(c.Request.Context(), albumID, userID, req.MediaIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Media added to album"})
}

func (h *AlbumHandler) RemoveMedia(c *gin.Context) {
	userID := c.GetString("userID")
	albumID := c.Param("albumId")

	var req models.RemoveFromAlbumRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := h.service.RemoveMedia(c.Request.Context(), albumID, userID, req.MediaIDs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Media removed from album"})
}

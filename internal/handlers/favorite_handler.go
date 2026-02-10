package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/services"
)

type FavoriteHandler struct {
	service *services.FavoriteService
}

func NewFavoriteHandler(service *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	if err := h.service.AddFavorite(c.Request.Context(), userID, mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Added to favorites"})
}

func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	if err := h.service.RemoveFavorite(c.Request.Context(), userID, mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Removed from favorites"})
}

func (h *FavoriteHandler) IsFavorite(c *gin.Context) {
	userID := c.GetString("userID")
	mediaID := c.Param("id")

	isFav, err := h.service.IsFavorite(c.Request.Context(), userID, mediaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"isFavorite": isFav}})
}

func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := c.GetString("userID")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)

	favs, err := h.service.GetFavorites(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": favs})
}

func (h *FavoriteHandler) GetFavoriteCount(c *gin.Context) {
	userID := c.GetString("userID")

	count, err := h.service.GetFavoriteCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"count": count}})
}

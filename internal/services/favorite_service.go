package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FavoriteService struct {
	db *database.MongoDB
}

func NewFavoriteService(db *database.MongoDB) *FavoriteService {
	return &FavoriteService{db: db}
}

func (s *FavoriteService) AddFavorite(ctx context.Context, userID, mediaID string) error {
	fav := &models.MediaFavorite{
		ID:        uuid.New().String(),
		UserID:    userID,
		MediaID:   mediaID,
		CreatedAt: time.Now(),
	}

	_, err := s.db.Collection("media_favorites").InsertOne(ctx, fav)
	return err
}

func (s *FavoriteService) RemoveFavorite(ctx context.Context, userID, mediaID string) error {
	_, err := s.db.Collection("media_favorites").DeleteOne(ctx,
		bson.M{"userId": userID, "mediaId": mediaID},
	)
	return err
}

func (s *FavoriteService) IsFavorite(ctx context.Context, userID, mediaID string) (bool, error) {
	count, err := s.db.Collection("media_favorites").CountDocuments(ctx,
		bson.M{"userId": userID, "mediaId": mediaID},
	)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *FavoriteService) GetFavorites(ctx context.Context, userID string, limit int64) ([]models.MediaFavorite, error) {
	cursor, err := s.db.Collection("media_favorites").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var favs []models.MediaFavorite
	if err := cursor.All(ctx, &favs); err != nil {
		return nil, err
	}
	return favs, nil
}

func (s *FavoriteService) GetFavoriteCount(ctx context.Context, userID string) (int64, error) {
	return s.db.Collection("media_favorites").CountDocuments(ctx, bson.M{"userId": userID})
}

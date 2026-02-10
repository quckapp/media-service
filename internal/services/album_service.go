package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AlbumService struct {
	db *database.MongoDB
}

func NewAlbumService(db *database.MongoDB) *AlbumService {
	return &AlbumService{db: db}
}

func (s *AlbumService) Create(ctx context.Context, userID string, req *models.CreateAlbumRequest) (*models.MediaAlbum, error) {
	album := &models.MediaAlbum{
		ID:          uuid.New().String(),
		UserID:      userID,
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		MediaIDs:    []string{},
		IsPublic:    req.IsPublic,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Collection("media_albums").InsertOne(ctx, album)
	if err != nil {
		return nil, err
	}
	return album, nil
}

func (s *AlbumService) GetByID(ctx context.Context, albumID string) (*models.MediaAlbum, error) {
	var album models.MediaAlbum
	err := s.db.Collection("media_albums").FindOne(ctx, bson.M{"_id": albumID}).Decode(&album)
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (s *AlbumService) GetByUser(ctx context.Context, userID string, limit int64) ([]models.MediaAlbum, error) {
	cursor, err := s.db.Collection("media_albums").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []models.MediaAlbum
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AlbumService) GetByWorkspace(ctx context.Context, workspaceID string, limit int64) ([]models.MediaAlbum, error) {
	cursor, err := s.db.Collection("media_albums").Find(ctx,
		bson.M{"workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []models.MediaAlbum
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

func (s *AlbumService) Update(ctx context.Context, albumID, userID string, req *models.UpdateAlbumRequest) (*models.MediaAlbum, error) {
	album, err := s.GetByID(ctx, albumID)
	if err != nil {
		return nil, err
	}
	if album.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	update := bson.M{"updatedAt": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Description != "" {
		update["description"] = req.Description
	}
	if req.CoverURL != "" {
		update["coverUrl"] = req.CoverURL
	}
	if req.IsPublic != nil {
		update["isPublic"] = *req.IsPublic
	}

	_, err = s.db.Collection("media_albums").UpdateOne(ctx,
		bson.M{"_id": albumID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}

	return s.GetByID(ctx, albumID)
}

func (s *AlbumService) Delete(ctx context.Context, albumID, userID string) error {
	album, err := s.GetByID(ctx, albumID)
	if err != nil {
		return err
	}
	if album.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media_albums").DeleteOne(ctx, bson.M{"_id": albumID})
	return err
}

func (s *AlbumService) AddMedia(ctx context.Context, albumID, userID string, mediaIDs []string) error {
	album, err := s.GetByID(ctx, albumID)
	if err != nil {
		return err
	}
	if album.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media_albums").UpdateOne(ctx,
		bson.M{"_id": albumID},
		bson.M{
			"$addToSet": bson.M{"mediaIds": bson.M{"$each": mediaIDs}},
			"$set":      bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (s *AlbumService) RemoveMedia(ctx context.Context, albumID, userID string, mediaIDs []string) error {
	album, err := s.GetByID(ctx, albumID)
	if err != nil {
		return err
	}
	if album.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media_albums").UpdateOne(ctx,
		bson.M{"_id": albumID},
		bson.M{
			"$pullAll": bson.M{"mediaIds": mediaIDs},
			"$set":     bson.M{"updatedAt": time.Now()},
		},
	)
	return err
}

func (s *AlbumService) GetPublicByWorkspace(ctx context.Context, workspaceID string, limit int64) ([]models.MediaAlbum, error) {
	cursor, err := s.db.Collection("media_albums").Find(ctx,
		bson.M{"workspaceId": workspaceID, "isPublic": true},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var albums []models.MediaAlbum
	if err := cursor.All(ctx, &albums); err != nil {
		return nil, err
	}
	return albums, nil
}

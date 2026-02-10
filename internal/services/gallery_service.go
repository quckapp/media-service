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

type GalleryService struct {
	db *database.MongoDB
}

func NewGalleryService(db *database.MongoDB) *GalleryService {
	return &GalleryService{db: db}
}

func (s *GalleryService) Create(ctx context.Context, userID string, req *models.CreateGalleryRequest) (*models.MediaGallery, error) {
	layout := req.Layout
	if layout == "" {
		layout = "grid"
	}

	gallery := &models.MediaGallery{
		ID:          uuid.New().String(),
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Description: req.Description,
		MediaIDs:    []string{},
		Layout:      layout,
		IsPublic:    req.IsPublic,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.db.Collection("media_galleries").InsertOne(ctx, gallery)
	if err != nil {
		return nil, err
	}
	return gallery, nil
}

func (s *GalleryService) List(ctx context.Context, workspaceID string, limit int64) ([]models.MediaGallery, error) {
	cursor, err := s.db.Collection("media_galleries").Find(ctx,
		bson.M{"workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var galleries []models.MediaGallery
	if err := cursor.All(ctx, &galleries); err != nil {
		return nil, err
	}
	return galleries, nil
}

func (s *GalleryService) GetByID(ctx context.Context, galleryID string) (*models.MediaGallery, error) {
	var gallery models.MediaGallery
	err := s.db.Collection("media_galleries").FindOne(ctx, bson.M{"_id": galleryID}).Decode(&gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, nil
}

func (s *GalleryService) Update(ctx context.Context, galleryID, userID string, req *models.UpdateGalleryRequest) (*models.MediaGallery, error) {
	gallery, err := s.GetByID(ctx, galleryID)
	if err != nil {
		return nil, err
	}
	if gallery.CreatedBy != userID {
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
	if req.Layout != "" {
		update["layout"] = req.Layout
	}
	if req.IsPublic != nil {
		update["isPublic"] = *req.IsPublic
	}
	if req.MediaIDs != nil {
		update["mediaIds"] = req.MediaIDs
	}

	_, err = s.db.Collection("media_galleries").UpdateOne(ctx,
		bson.M{"_id": galleryID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, galleryID)
}

func (s *GalleryService) Delete(ctx context.Context, galleryID, userID string) error {
	gallery, err := s.GetByID(ctx, galleryID)
	if err != nil {
		return err
	}
	if gallery.CreatedBy != userID {
		return fmt.Errorf("unauthorized")
	}
	_, err = s.db.Collection("media_galleries").DeleteOne(ctx, bson.M{"_id": galleryID})
	return err
}

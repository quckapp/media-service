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

type WatermarkService struct {
	db *database.MongoDB
}

func NewWatermarkService(db *database.MongoDB) *WatermarkService {
	return &WatermarkService{db: db}
}

func (s *WatermarkService) Upload(ctx context.Context, userID string, req *models.UploadWatermarkRequest) (*models.Watermark, error) {
	watermark := &models.Watermark{
		ID:          uuid.New().String(),
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		ImageURL:    req.ImageURL,
		Position:    req.Position,
		Opacity:     req.Opacity,
		Size:        req.Size,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
	}

	_, err := s.db.Collection("watermarks").InsertOne(ctx, watermark)
	if err != nil {
		return nil, err
	}
	return watermark, nil
}

func (s *WatermarkService) List(ctx context.Context, workspaceID string) ([]models.Watermark, error) {
	cursor, err := s.db.Collection("watermarks").Find(ctx,
		bson.M{"workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var watermarks []models.Watermark
	if err := cursor.All(ctx, &watermarks); err != nil {
		return nil, err
	}
	return watermarks, nil
}

func (s *WatermarkService) Apply(ctx context.Context, req *models.ApplyWatermarkRequest) error {
	_, err := s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": req.MediaID},
		bson.M{"$set": bson.M{
			"metadata.watermarkId":       req.WatermarkID,
			"metadata.watermarkPosition": req.Position,
			"updatedAt": time.Now(),
		}},
	)
	return err
}

func (s *WatermarkService) Remove(ctx context.Context, mediaID string) error {
	_, err := s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": mediaID},
		bson.M{"$unset": bson.M{
			"metadata.watermarkId":       "",
			"metadata.watermarkPosition": "",
		}, "$set": bson.M{"updatedAt": time.Now()}},
	)
	return err
}

func (s *WatermarkService) GetSettings(ctx context.Context, workspaceID string) (*models.WatermarkSettings, error) {
	var settings models.WatermarkSettings
	err := s.db.Collection("watermark_settings").FindOne(ctx, bson.M{"workspaceId": workspaceID}).Decode(&settings)
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

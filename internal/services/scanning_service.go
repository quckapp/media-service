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

type ScanningService struct {
	db *database.MongoDB
}

func NewScanningService(db *database.MongoDB) *ScanningService {
	return &ScanningService{db: db}
}

func (s *ScanningService) ScanMedia(ctx context.Context, req *models.ScanRequest) (*models.MediaScan, error) {
	scan := &models.MediaScan{
		ID:         uuid.New().String(),
		MediaID:    req.MediaID,
		ScanType:   req.ScanType,
		Status:     "pending",
		Confidence: 0,
		ScannedAt:  time.Now(),
	}

	_, err := s.db.Collection("media_scans").InsertOne(ctx, scan)
	if err != nil {
		return nil, err
	}
	return scan, nil
}

func (s *ScanningService) GetScanResults(ctx context.Context, mediaID string) ([]models.MediaScan, error) {
	cursor, err := s.db.Collection("media_scans").Find(ctx,
		bson.M{"mediaId": mediaID},
		options.Find().SetSort(bson.M{"scannedAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scans []models.MediaScan
	if err := cursor.All(ctx, &scans); err != nil {
		return nil, err
	}
	return scans, nil
}

func (s *ScanningService) ListFlagged(ctx context.Context, limit int64) ([]models.MediaScan, error) {
	cursor, err := s.db.Collection("media_scans").Find(ctx,
		bson.M{"status": "flagged"},
		options.Find().SetSort(bson.M{"scannedAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var scans []models.MediaScan
	if err := cursor.All(ctx, &scans); err != nil {
		return nil, err
	}
	return scans, nil
}

func (s *ScanningService) UpdateStatus(ctx context.Context, scanID, status string) error {
	_, err := s.db.Collection("media_scans").UpdateOne(ctx,
		bson.M{"_id": scanID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

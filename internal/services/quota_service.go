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

type QuotaService struct {
	db *database.MongoDB
}

func NewQuotaService(db *database.MongoDB) *QuotaService {
	return &QuotaService{db: db}
}

func (s *QuotaService) GetByWorkspace(ctx context.Context, workspaceID string) (*models.StorageQuota, error) {
	var quota models.StorageQuota
	err := s.db.Collection("storage_quotas").FindOne(ctx, bson.M{"workspaceId": workspaceID}).Decode(&quota)
	if err != nil {
		return nil, err
	}
	return &quota, nil
}

func (s *QuotaService) SetQuota(ctx context.Context, req *models.SetQuotaRequest) (*models.StorageQuota, error) {
	now := time.Now()
	filter := bson.M{"workspaceId": req.WorkspaceID}
	update := bson.M{
		"$set": bson.M{
			"maxStorageMB": req.MaxStorageMB,
			"maxFileCount": req.MaxFileCount,
			"updatedAt":    now,
		},
		"$setOnInsert": bson.M{
			"_id":              uuid.New().String(),
			"workspaceId":      req.WorkspaceID,
			"usedStorageMB":    int64(0),
			"currentFileCount": int64(0),
			"createdAt":        now,
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err := s.db.Collection("storage_quotas").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}
	return s.GetByWorkspace(ctx, req.WorkspaceID)
}

func (s *QuotaService) GetUsage(ctx context.Context, workspaceID string) (*models.StorageQuota, error) {
	return s.GetByWorkspace(ctx, workspaceID)
}

func (s *QuotaService) ListOverQuota(ctx context.Context, limit int64) ([]models.StorageQuota, error) {
	cursor, err := s.db.Collection("storage_quotas").Find(ctx,
		bson.M{"$expr": bson.M{"$gt": []string{"$usedStorageMB", "$maxStorageMB"}}},
		options.Find().SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var quotas []models.StorageQuota
	if err := cursor.All(ctx, &quotas); err != nil {
		return nil, err
	}
	return quotas, nil
}

package services

import (
	"context"
	"time"

	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsService struct {
	db *database.MongoDB
}

func NewAnalyticsService(db *database.MongoDB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) GetUploadTrends(ctx context.Context, workspaceID string, days int) ([]models.UploadTrend, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"createdAt": bson.M{"$gte": startDate},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":       bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
			"count":     bson.M{"$sum": 1},
			"totalSize": bson.M{"$sum": "$size"},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Date      string `bson:"_id"`
		Count     int64  `bson:"count"`
		TotalSize int64  `bson:"totalSize"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	trends := make([]models.UploadTrend, len(results))
	for i, r := range results {
		trends[i] = models.UploadTrend{Date: r.Date, Count: r.Count, TotalSize: r.TotalSize}
	}
	return trends, nil
}

func (s *AnalyticsService) GetStorageTrends(ctx context.Context, workspaceID string, days int) ([]models.StorageTrend, error) {
	startDate := time.Now().AddDate(0, 0, -days)
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"createdAt": bson.M{"$gte": startDate},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id":    bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
			"usedMB": bson.M{"$sum": bson.M{"$divide": []interface{}{"$size", 1048576}}},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Date   string `bson:"_id"`
		UsedMB int64  `bson:"usedMB"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	trends := make([]models.StorageTrend, len(results))
	for i, r := range results {
		trends[i] = models.StorageTrend{Date: r.Date, UsedMB: r.UsedMB}
	}
	return trends, nil
}

func (s *AnalyticsService) GetFileTypeDistribution(ctx context.Context, workspaceID string) ([]models.FileTypeDistribution, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":       "$contentType",
			"count":     bson.M{"$sum": 1},
			"totalSize": bson.M{"$sum": "$size"},
		}}},
		{{Key: "$sort", Value: bson.M{"count": -1}}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Type      string `bson:"_id"`
		Count     int64  `bson:"count"`
		TotalSize int64  `bson:"totalSize"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	dist := make([]models.FileTypeDistribution, len(results))
	for i, r := range results {
		dist[i] = models.FileTypeDistribution{Type: r.Type, Count: r.Count, TotalSize: r.TotalSize}
	}
	return dist, nil
}

func (s *AnalyticsService) GetUserUploadStats(ctx context.Context, workspaceID string, limit int64) ([]models.UserUploadStats, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":        "$uploadedBy",
			"fileCount":  bson.M{"$sum": 1},
			"totalSize":  bson.M{"$sum": "$size"},
			"lastUpload": bson.M{"$max": "$createdAt"},
		}}},
		{{Key: "$sort", Value: bson.M{"totalSize": -1}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		UserID     string    `bson:"_id"`
		FileCount  int64     `bson:"fileCount"`
		TotalSize  int64     `bson:"totalSize"`
		LastUpload time.Time `bson:"lastUpload"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	stats := make([]models.UserUploadStats, len(results))
	for i, r := range results {
		stats[i] = models.UserUploadStats{UserID: r.UserID, FileCount: r.FileCount, TotalSize: r.TotalSize, LastUpload: r.LastUpload}
	}
	return stats, nil
}

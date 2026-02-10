package services

import (
	"context"
	"time"

	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchService struct {
	db      *database.MongoDB
	storage *S3Storage
}

func NewSearchService(db *database.MongoDB, storage *S3Storage) *SearchService {
	return &SearchService{db: db, storage: storage}
}

func (s *SearchService) Search(ctx context.Context, userID string, params models.MediaSearchParams) ([]models.Media, int64, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Page < 0 {
		params.Page = 0
	}

	filter := bson.M{"userId": userID}

	if params.Type != "" {
		filter["type"] = params.Type
	}
	if params.WorkspaceID != "" {
		filter["metadata.workspaceId"] = params.WorkspaceID
	}
	if params.ChannelID != "" {
		filter["metadata.channelId"] = params.ChannelID
	}
	if params.Query != "" {
		filter["filename"] = bson.M{"$regex": primitive.Regex{Pattern: params.Query, Options: "i"}}
	}
	if params.MinSize > 0 {
		filter["size"] = bson.M{"$gte": params.MinSize}
	}
	if params.MaxSize > 0 {
		if existing, ok := filter["size"]; ok {
			existing.(bson.M)["$lte"] = params.MaxSize
		} else {
			filter["size"] = bson.M{"$lte": params.MaxSize}
		}
	}
	if params.DateFrom != "" {
		if t, err := time.Parse("2006-01-02", params.DateFrom); err == nil {
			filter["createdAt"] = bson.M{"$gte": t}
		}
	}
	if params.DateTo != "" {
		if t, err := time.Parse("2006-01-02", params.DateTo); err == nil {
			if existing, ok := filter["createdAt"]; ok {
				existing.(bson.M)["$lte"] = t
			} else {
				filter["createdAt"] = bson.M{"$lte": t}
			}
		}
	}

	// Sorting
	sortField := "createdAt"
	sortOrder := -1
	if params.SortBy != "" {
		sortField = params.SortBy
	}
	if params.SortOrder == "asc" {
		sortOrder = 1
	}

	total, err := s.db.Collection("media").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := int64(params.Page * params.Limit)
	cursor, err := s.db.Collection("media").Find(ctx, filter,
		options.Find().
			SetSort(bson.M{sortField: sortOrder}).
			SetSkip(skip).
			SetLimit(int64(params.Limit)),
	)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var media []models.Media
	if err := cursor.All(ctx, &media); err != nil {
		return nil, 0, err
	}

	for i := range media {
		media[i].URL, _ = s.storage.GetPresignedDownloadURL(media[i].S3Key, time.Hour)
	}

	return media, total, nil
}

func (s *SearchService) GetWorkspaceStats(ctx context.Context, workspaceID string) (*models.WorkspaceMediaStats, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{"metadata.workspaceId": workspaceID}},
		bson.M{"$group": bson.M{
			"_id":       "$type",
			"count":     bson.M{"$sum": 1},
			"totalSize": bson.M{"$sum": "$size"},
		}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	stats := &models.WorkspaceMediaStats{
		ByType: make(map[string]int64),
		ByUser: make(map[string]int64),
	}

	for cursor.Next(ctx) {
		var result struct {
			Type      string `bson:"_id"`
			Count     int64  `bson:"count"`
			TotalSize int64  `bson:"totalSize"`
		}
		if err := cursor.Decode(&result); err == nil {
			stats.ByType[result.Type] = result.Count
			stats.TotalFiles += result.Count
			stats.TotalSize += result.TotalSize
		}
	}

	stats.StorageUsed = stats.TotalSize

	// By user breakdown
	userPipeline := bson.A{
		bson.M{"$match": bson.M{"metadata.workspaceId": workspaceID}},
		bson.M{"$group": bson.M{
			"_id":   "$userId",
			"count": bson.M{"$sum": 1},
		}},
	}

	userCursor, err := s.db.Collection("media").Aggregate(ctx, userPipeline)
	if err == nil {
		defer userCursor.Close(ctx)
		for userCursor.Next(ctx) {
			var result struct {
				UserID string `bson:"_id"`
				Count  int64  `bson:"count"`
			}
			if err := userCursor.Decode(&result); err == nil {
				stats.ByUser[result.UserID] = result.Count
			}
		}
	}

	// Recent uploads (last 24h)
	dayAgo := time.Now().Add(-24 * time.Hour)
	recentCount, _ := s.db.Collection("media").CountDocuments(ctx,
		bson.M{"metadata.workspaceId": workspaceID, "createdAt": bson.M{"$gte": dayAgo}},
	)
	stats.RecentUploads = recentCount

	return stats, nil
}

func (s *SearchService) GetDuplicates(ctx context.Context, userID string) ([]models.Media, error) {
	// Find media with same filename and size
	pipeline := bson.A{
		bson.M{"$match": bson.M{"userId": userID}},
		bson.M{"$group": bson.M{
			"_id":   bson.M{"filename": "$filename", "size": "$size"},
			"count": bson.M{"$sum": 1},
			"ids":   bson.M{"$push": "$_id"},
		}},
		bson.M{"$match": bson.M{"count": bson.M{"$gt": 1}}},
	}

	cursor, err := s.db.Collection("media").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dupIDs []string
	for cursor.Next(ctx) {
		var result struct {
			IDs []string `bson:"ids"`
		}
		if err := cursor.Decode(&result); err == nil {
			dupIDs = append(dupIDs, result.IDs...)
		}
	}

	if len(dupIDs) == 0 {
		return []models.Media{}, nil
	}

	mediaCursor, err := s.db.Collection("media").Find(ctx,
		bson.M{"_id": bson.M{"$in": dupIDs}},
		options.Find().SetSort(bson.M{"filename": 1, "createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer mediaCursor.Close(ctx)

	var media []models.Media
	if err := mediaCursor.All(ctx, &media); err != nil {
		return nil, err
	}
	return media, nil
}

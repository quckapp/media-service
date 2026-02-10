package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MediaService struct {
	db      *database.MongoDB
	redis   *redis.Client
	storage *S3Storage
}

func NewMediaService(db *database.MongoDB, redis *redis.Client, storage *S3Storage) *MediaService {
	return &MediaService{db: db, redis: redis, storage: storage}
}

func (s *MediaService) Create(ctx context.Context, userID string, req *models.UploadRequest) (*models.Media, error) {
	mediaID := uuid.New().String()
	s3Key := fmt.Sprintf("media/%s/%s/%s", userID, mediaID, req.Filename)

	media := &models.Media{
		ID:        mediaID,
		UserID:    userID,
		Type:      getMediaType(req.MimeType),
		Filename:  req.Filename,
		MimeType:  req.MimeType,
		Size:      req.Size,
		S3Key:     s3Key,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Collection("media").InsertOne(ctx, media)
	if err != nil {
		return nil, err
	}

	return media, nil
}

func (s *MediaService) GetPresignedUploadURL(ctx context.Context, userID string, req *models.UploadRequest) (*models.PresignedURLResponse, error) {
	media, err := s.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	uploadURL, err := s.storage.GetPresignedUploadURL(media.S3Key, req.MimeType, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	return &models.PresignedURLResponse{
		UploadURL: uploadURL,
		MediaID:   media.ID,
		S3Key:     media.S3Key,
		ExpiresAt: time.Now().Add(15 * time.Minute).Format(time.RFC3339),
	}, nil
}

func (s *MediaService) Get(ctx context.Context, mediaID string) (*models.Media, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("media:%s", mediaID)
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var media models.Media
		if json.Unmarshal([]byte(cached), &media) == nil {
			return &media, nil
		}
	}

	// Fetch from DB
	var media models.Media
	err = s.db.Collection("media").FindOne(ctx, bson.M{"_id": mediaID}).Decode(&media)
	if err != nil {
		return nil, err
	}

	// Generate signed URL
	media.URL, _ = s.storage.GetPresignedDownloadURL(media.S3Key, time.Hour)

	// Cache result
	if data, err := json.Marshal(media); err == nil {
		s.redis.Set(ctx, cacheKey, data, time.Hour)
	}

	return &media, nil
}

func (s *MediaService) Delete(ctx context.Context, mediaID, userID string) error {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return err
	}

	if media.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Delete from S3
	if err := s.storage.Delete(media.S3Key); err != nil {
		return err
	}

	// Delete from DB
	_, err = s.db.Collection("media").DeleteOne(ctx, bson.M{"_id": mediaID})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.redis.Del(ctx, fmt.Sprintf("media:%s", mediaID))

	return nil
}

func (s *MediaService) GetUserMedia(ctx context.Context, userID string, limit int64) ([]models.Media, error) {
	cursor, err := s.db.Collection("media").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var media []models.Media
	if err := cursor.All(ctx, &media); err != nil {
		return nil, err
	}

	// Generate signed URLs
	for i := range media {
		media[i].URL, _ = s.storage.GetPresignedDownloadURL(media[i].S3Key, time.Hour)
	}

	return media, nil
}

func (s *MediaService) SetURL(ctx context.Context, mediaID, url string) error {
	_, err := s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": mediaID},
		bson.M{"$set": bson.M{"url": url, "updatedAt": time.Now()}},
	)
	return err
}

func (s *MediaService) BulkDelete(ctx context.Context, mediaIDs []string, userID string) *models.BulkDeleteResponse {
	resp := &models.BulkDeleteResponse{}
	for _, id := range mediaIDs {
		if err := s.Delete(ctx, id, userID); err != nil {
			resp.Failed = append(resp.Failed, id)
		} else {
			resp.Deleted = append(resp.Deleted, id)
		}
	}
	return resp
}

func (s *MediaService) GetMediaByChannel(ctx context.Context, channelID string, limit int64) ([]models.Media, error) {
	cursor, err := s.db.Collection("media").Find(ctx,
		bson.M{"metadata.channelId": channelID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var media []models.Media
	if err := cursor.All(ctx, &media); err != nil {
		return nil, err
	}

	for i := range media {
		media[i].URL, _ = s.storage.GetPresignedDownloadURL(media[i].S3Key, time.Hour)
	}

	return media, nil
}

func (s *MediaService) GetMediaByWorkspace(ctx context.Context, workspaceID string, limit int64) ([]models.Media, error) {
	cursor, err := s.db.Collection("media").Find(ctx,
		bson.M{"metadata.workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var media []models.Media
	if err := cursor.All(ctx, &media); err != nil {
		return nil, err
	}

	for i := range media {
		media[i].URL, _ = s.storage.GetPresignedDownloadURL(media[i].S3Key, time.Hour)
	}

	return media, nil
}

func (s *MediaService) GetUserStats(ctx context.Context, userID string) (*models.MediaStatsResponse, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{"userId": userID}},
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

	stats := &models.MediaStatsResponse{
		ByType: make(map[string]int64),
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

	return stats, nil
}

func (s *MediaService) UpdateMetadata(ctx context.Context, mediaID, userID string, metadata map[string]string) error {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return err
	}
	if media.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": mediaID},
		bson.M{"$set": bson.M{"metadata": metadata, "updatedAt": time.Now()}},
	)
	if err != nil {
		return err
	}

	s.redis.Del(ctx, fmt.Sprintf("media:%s", mediaID))
	return nil
}

func (s *MediaService) Rename(ctx context.Context, mediaID, userID, newFilename string) error {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return err
	}
	if media.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": mediaID},
		bson.M{"$set": bson.M{"filename": newFilename, "updatedAt": time.Now()}},
	)
	if err != nil {
		return err
	}

	s.redis.Del(ctx, fmt.Sprintf("media:%s", mediaID))
	return nil
}

func (s *MediaService) CopyMedia(ctx context.Context, mediaID, userID, targetWorkspaceID string) (*models.Media, error) {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return nil, err
	}
	if media.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	newID := uuid.New().String()
	newMedia := &models.Media{
		ID:        newID,
		UserID:    userID,
		Type:      media.Type,
		Filename:  media.Filename,
		MimeType:  media.MimeType,
		Size:      media.Size,
		S3Key:     media.S3Key, // shares same S3 key
		Metadata:  map[string]string{"workspaceId": targetWorkspaceID},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.Collection("media").InsertOne(ctx, newMedia)
	if err != nil {
		return nil, err
	}
	return newMedia, nil
}

func (s *MediaService) MoveMedia(ctx context.Context, mediaID, userID, targetWorkspaceID string) error {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return err
	}
	if media.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	_, err = s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": mediaID},
		bson.M{"$set": bson.M{
			"metadata.workspaceId": targetWorkspaceID,
			"updatedAt":            time.Now(),
		}},
	)
	if err != nil {
		return err
	}

	s.redis.Del(ctx, fmt.Sprintf("media:%s", mediaID))
	return nil
}

func (s *MediaService) BulkMove(ctx context.Context, mediaIDs []string, userID, targetWorkspaceID string) *models.BulkDeleteResponse {
	resp := &models.BulkDeleteResponse{}
	for _, id := range mediaIDs {
		if err := s.MoveMedia(context.Background(), id, userID, targetWorkspaceID); err != nil {
			resp.Failed = append(resp.Failed, id)
		} else {
			resp.Deleted = append(resp.Deleted, id)
		}
	}
	return resp
}

func (s *MediaService) GetMediaByType(ctx context.Context, userID, mediaType string, limit int64) ([]models.Media, error) {
	cursor, err := s.db.Collection("media").Find(ctx,
		bson.M{"userId": userID, "type": mediaType},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var media []models.Media
	if err := cursor.All(ctx, &media); err != nil {
		return nil, err
	}

	for i := range media {
		media[i].URL, _ = s.storage.GetPresignedDownloadURL(media[i].S3Key, time.Hour)
	}

	return media, nil
}

func (s *MediaService) GetDownloadURL(ctx context.Context, mediaID string) (string, error) {
	media, err := s.Get(ctx, mediaID)
	if err != nil {
		return "", err
	}

	url, err := s.storage.GetPresignedDownloadURL(media.S3Key, time.Hour)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (s *MediaService) GetRecentMedia(ctx context.Context, userID string, limit int64) ([]models.Media, error) {
	return s.GetUserMedia(ctx, userID, limit)
}

func getMediaType(mimeType string) string {
	if len(mimeType) < 5 {
		return "document"
	}
	switch {
	case mimeType[:5] == "image":
		return "image"
	case mimeType[:5] == "video":
		return "video"
	case mimeType[:5] == "audio":
		return "audio"
	default:
		return "document"
	}
}

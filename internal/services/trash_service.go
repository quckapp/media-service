package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TrashService struct {
	db      *database.MongoDB
	redis   *redis.Client
	storage *S3Storage
}

func NewTrashService(db *database.MongoDB, redis *redis.Client, storage *S3Storage) *TrashService {
	return &TrashService{db: db, redis: redis, storage: storage}
}

func (s *TrashService) MoveToTrash(ctx context.Context, mediaID, userID string) error {
	// Get the media document
	var media models.Media
	err := s.db.Collection("media").FindOne(ctx, bson.M{"_id": mediaID}).Decode(&media)
	if err != nil {
		return err
	}
	if media.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Create trash entry
	trashed := &models.TrashedMedia{
		ID:          uuid.New().String(),
		MediaID:     mediaID,
		UserID:      userID,
		OriginalDoc: media,
		TrashedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(30 * 24 * time.Hour), // 30 days
	}

	_, err = s.db.Collection("media_trash").InsertOne(ctx, trashed)
	if err != nil {
		return err
	}

	// Remove from main collection
	_, err = s.db.Collection("media").DeleteOne(ctx, bson.M{"_id": mediaID})
	if err != nil {
		return err
	}

	// Invalidate cache
	s.redis.Del(ctx, fmt.Sprintf("media:%s", mediaID))
	return nil
}

func (s *TrashService) RestoreFromTrash(ctx context.Context, trashID, userID string) error {
	var trashed models.TrashedMedia
	err := s.db.Collection("media_trash").FindOne(ctx,
		bson.M{"_id": trashID, "userId": userID},
	).Decode(&trashed)
	if err != nil {
		return err
	}

	// Restore to main collection
	_, err = s.db.Collection("media").InsertOne(ctx, trashed.OriginalDoc)
	if err != nil {
		return err
	}

	// Remove from trash
	_, err = s.db.Collection("media_trash").DeleteOne(ctx, bson.M{"_id": trashID})
	return err
}

func (s *TrashService) GetTrash(ctx context.Context, userID string, limit int64) ([]models.TrashedMedia, error) {
	cursor, err := s.db.Collection("media_trash").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"trashedAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trashed []models.TrashedMedia
	if err := cursor.All(ctx, &trashed); err != nil {
		return nil, err
	}
	return trashed, nil
}

func (s *TrashService) PermanentDelete(ctx context.Context, trashID, userID string) error {
	var trashed models.TrashedMedia
	err := s.db.Collection("media_trash").FindOne(ctx,
		bson.M{"_id": trashID, "userId": userID},
	).Decode(&trashed)
	if err != nil {
		return err
	}

	// Delete from S3
	_ = s.storage.Delete(trashed.OriginalDoc.S3Key)

	// Remove from trash
	_, err = s.db.Collection("media_trash").DeleteOne(ctx, bson.M{"_id": trashID})
	return err
}

func (s *TrashService) EmptyTrash(ctx context.Context, userID string) (int64, error) {
	// Get all trashed items
	cursor, err := s.db.Collection("media_trash").Find(ctx, bson.M{"userId": userID})
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var trashed []models.TrashedMedia
	if err := cursor.All(ctx, &trashed); err != nil {
		return 0, err
	}

	// Delete from S3
	for _, t := range trashed {
		_ = s.storage.Delete(t.OriginalDoc.S3Key)
	}

	// Remove all from trash
	result, err := s.db.Collection("media_trash").DeleteMany(ctx, bson.M{"userId": userID})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

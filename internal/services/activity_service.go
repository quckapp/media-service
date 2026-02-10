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

type ActivityService struct {
	db *database.MongoDB
}

func NewActivityService(db *database.MongoDB) *ActivityService {
	return &ActivityService{db: db}
}

func (s *ActivityService) LogActivity(ctx context.Context, mediaID, userID, action, details string) error {
	activity := &models.MediaActivity{
		ID:        uuid.New().String(),
		MediaID:   mediaID,
		UserID:    userID,
		Action:    action,
		Details:   details,
		CreatedAt: time.Now(),
	}

	_, err := s.db.Collection("media_activity").InsertOne(ctx, activity)
	return err
}

func (s *ActivityService) GetByMedia(ctx context.Context, mediaID string, limit int64) ([]models.MediaActivity, error) {
	cursor, err := s.db.Collection("media_activity").Find(ctx,
		bson.M{"mediaId": mediaID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []models.MediaActivity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) GetByUser(ctx context.Context, userID string, limit int64) ([]models.MediaActivity, error) {
	cursor, err := s.db.Collection("media_activity").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var activities []models.MediaActivity
	if err := cursor.All(ctx, &activities); err != nil {
		return nil, err
	}
	return activities, nil
}

func (s *ActivityService) GetRecent(ctx context.Context, userID string, limit int64) ([]models.MediaActivity, error) {
	return s.GetByUser(ctx, userID, limit)
}

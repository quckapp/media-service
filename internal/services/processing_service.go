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

type ProcessingService struct {
	db *database.MongoDB
}

func NewProcessingService(db *database.MongoDB) *ProcessingService {
	return &ProcessingService{db: db}
}

func (s *ProcessingService) CreateJob(ctx context.Context, mediaID, userID string, req *models.CreateProcessingJobRequest) (*models.ProcessingJob, error) {
	job := &models.ProcessingJob{
		ID:        uuid.New().String(),
		MediaID:   mediaID,
		UserID:    userID,
		Type:      req.Type,
		Status:    "pending",
		Params:    req.Params,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Collection("media_processing_jobs").InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (s *ProcessingService) GetJob(ctx context.Context, jobID string) (*models.ProcessingJob, error) {
	var job models.ProcessingJob
	err := s.db.Collection("media_processing_jobs").FindOne(ctx, bson.M{"_id": jobID}).Decode(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (s *ProcessingService) GetJobsByMedia(ctx context.Context, mediaID string) ([]models.ProcessingJob, error) {
	cursor, err := s.db.Collection("media_processing_jobs").Find(ctx,
		bson.M{"mediaId": mediaID},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.ProcessingJob
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *ProcessingService) GetUserJobs(ctx context.Context, userID string, status string, limit int64) ([]models.ProcessingJob, error) {
	filter := bson.M{"userId": userID}
	if status != "" {
		filter["status"] = status
	}

	cursor, err := s.db.Collection("media_processing_jobs").Find(ctx,
		filter,
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var jobs []models.ProcessingJob
	if err := cursor.All(ctx, &jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (s *ProcessingService) UpdateJobStatus(ctx context.Context, jobID, status string, result map[string]interface{}, jobError string) error {
	update := bson.M{
		"status":    status,
		"updatedAt": time.Now(),
	}
	if result != nil {
		update["result"] = result
	}
	if jobError != "" {
		update["error"] = jobError
	}

	_, err := s.db.Collection("media_processing_jobs").UpdateOne(ctx,
		bson.M{"_id": jobID},
		bson.M{"$set": update},
	)
	return err
}

func (s *ProcessingService) CancelJob(ctx context.Context, jobID, userID string) error {
	_, err := s.db.Collection("media_processing_jobs").UpdateOne(ctx,
		bson.M{"_id": jobID, "userId": userID, "status": "pending"},
		bson.M{"$set": bson.M{"status": "cancelled", "updatedAt": time.Now()}},
	)
	return err
}

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

type RetentionService struct {
	db *database.MongoDB
}

func NewRetentionService(db *database.MongoDB) *RetentionService {
	return &RetentionService{db: db}
}

func (s *RetentionService) Create(ctx context.Context, userID string, req *models.CreateRetentionPolicyRequest) (*models.RetentionPolicy, error) {
	policy := &models.RetentionPolicy{
		ID:            uuid.New().String(),
		WorkspaceID:   req.WorkspaceID,
		Name:          req.Name,
		RetentionDays: req.RetentionDays,
		ApplyTo:       req.ApplyTo,
		AutoDelete:    req.AutoDelete,
		CreatedBy:     userID,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := s.db.Collection("retention_policies").InsertOne(ctx, policy)
	if err != nil {
		return nil, err
	}
	return policy, nil
}

func (s *RetentionService) GetByID(ctx context.Context, policyID string) (*models.RetentionPolicy, error) {
	var policy models.RetentionPolicy
	err := s.db.Collection("retention_policies").FindOne(ctx, bson.M{"_id": policyID}).Decode(&policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

func (s *RetentionService) Update(ctx context.Context, policyID, userID string, req *models.UpdateRetentionPolicyRequest) (*models.RetentionPolicy, error) {
	policy, err := s.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if policy.CreatedBy != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	update := bson.M{"updatedAt": time.Now()}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.RetentionDays != nil {
		update["retentionDays"] = *req.RetentionDays
	}
	if req.ApplyTo != "" {
		update["applyTo"] = req.ApplyTo
	}
	if req.AutoDelete != nil {
		update["autoDelete"] = *req.AutoDelete
	}

	_, err = s.db.Collection("retention_policies").UpdateOne(ctx,
		bson.M{"_id": policyID},
		bson.M{"$set": update},
	)
	if err != nil {
		return nil, err
	}
	return s.GetByID(ctx, policyID)
}

func (s *RetentionService) Delete(ctx context.Context, policyID, userID string) error {
	policy, err := s.GetByID(ctx, policyID)
	if err != nil {
		return err
	}
	if policy.CreatedBy != userID {
		return fmt.Errorf("unauthorized")
	}
	_, err = s.db.Collection("retention_policies").DeleteOne(ctx, bson.M{"_id": policyID})
	return err
}

func (s *RetentionService) GetByWorkspace(ctx context.Context, workspaceID string) ([]models.RetentionPolicy, error) {
	cursor, err := s.db.Collection("retention_policies").Find(ctx,
		bson.M{"workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var policies []models.RetentionPolicy
	if err := cursor.All(ctx, &policies); err != nil {
		return nil, err
	}
	return policies, nil
}

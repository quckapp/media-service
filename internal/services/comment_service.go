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

type CommentService struct {
	db *database.MongoDB
}

func NewCommentService(db *database.MongoDB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) Create(ctx context.Context, mediaID, userID string, req *models.CreateCommentRequest) (*models.MediaComment, error) {
	comment := &models.MediaComment{
		ID:        uuid.New().String(),
		MediaID:   mediaID,
		UserID:    userID,
		Content:   req.Content,
		ParentID:  req.ParentID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.db.Collection("media_comments").InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *CommentService) GetByMedia(ctx context.Context, mediaID string, limit int64) ([]models.MediaComment, error) {
	cursor, err := s.db.Collection("media_comments").Find(ctx,
		bson.M{"mediaId": mediaID, "parentId": bson.M{"$in": []interface{}{nil, ""}}},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.MediaComment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) GetReplies(ctx context.Context, parentID string) ([]models.MediaComment, error) {
	cursor, err := s.db.Collection("media_comments").Find(ctx,
		bson.M{"parentId": parentID},
		options.Find().SetSort(bson.M{"createdAt": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []models.MediaComment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentService) Update(ctx context.Context, commentID, userID string, req *models.UpdateCommentRequest) (*models.MediaComment, error) {
	result, err := s.db.Collection("media_comments").UpdateOne(ctx,
		bson.M{"_id": commentID, "userId": userID},
		bson.M{"$set": bson.M{"content": req.Content, "updatedAt": time.Now()}},
	)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, fmt.Errorf("comment not found or unauthorized")
	}

	var comment models.MediaComment
	err = s.db.Collection("media_comments").FindOne(ctx, bson.M{"_id": commentID}).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *CommentService) Delete(ctx context.Context, commentID, userID string) error {
	// Delete comment and its replies
	_, _ = s.db.Collection("media_comments").DeleteMany(ctx, bson.M{"parentId": commentID})
	_, err := s.db.Collection("media_comments").DeleteOne(ctx,
		bson.M{"_id": commentID, "userId": userID},
	)
	return err
}

func (s *CommentService) CountByMedia(ctx context.Context, mediaID string) (int64, error) {
	return s.db.Collection("media_comments").CountDocuments(ctx, bson.M{"mediaId": mediaID})
}

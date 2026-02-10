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

type TagService struct {
	db *database.MongoDB
}

func NewTagService(db *database.MongoDB) *TagService {
	return &TagService{db: db}
}

func (s *TagService) Create(ctx context.Context, userID string, req *models.CreateTagRequest) (*models.MediaTag, error) {
	tag := &models.MediaTag{
		ID:          uuid.New().String(),
		UserID:      userID,
		WorkspaceID: req.WorkspaceID,
		Name:        req.Name,
		Color:       req.Color,
		MediaCount:  0,
		CreatedAt:   time.Now(),
	}

	_, err := s.db.Collection("media_tags").InsertOne(ctx, tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *TagService) GetByID(ctx context.Context, tagID string) (*models.MediaTag, error) {
	var tag models.MediaTag
	err := s.db.Collection("media_tags").FindOne(ctx, bson.M{"_id": tagID}).Decode(&tag)
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *TagService) GetByUser(ctx context.Context, userID string) ([]models.MediaTag, error) {
	cursor, err := s.db.Collection("media_tags").Find(ctx,
		bson.M{"userId": userID},
		options.Find().SetSort(bson.M{"name": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []models.MediaTag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *TagService) GetByWorkspace(ctx context.Context, workspaceID string) ([]models.MediaTag, error) {
	cursor, err := s.db.Collection("media_tags").Find(ctx,
		bson.M{"workspaceId": workspaceID},
		options.Find().SetSort(bson.M{"name": 1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tags []models.MediaTag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *TagService) Update(ctx context.Context, tagID, userID string, req *models.UpdateTagRequest) (*models.MediaTag, error) {
	tag, err := s.GetByID(ctx, tagID)
	if err != nil {
		return nil, err
	}
	if tag.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	update := bson.M{}
	if req.Name != "" {
		update["name"] = req.Name
	}
	if req.Color != "" {
		update["color"] = req.Color
	}

	if len(update) > 0 {
		_, err = s.db.Collection("media_tags").UpdateOne(ctx,
			bson.M{"_id": tagID},
			bson.M{"$set": update},
		)
		if err != nil {
			return nil, err
		}
	}

	return s.GetByID(ctx, tagID)
}

func (s *TagService) Delete(ctx context.Context, tagID, userID string) error {
	tag, err := s.GetByID(ctx, tagID)
	if err != nil {
		return err
	}
	if tag.UserID != userID {
		return fmt.Errorf("unauthorized")
	}

	// Remove all mappings
	_, _ = s.db.Collection("media_tag_mappings").DeleteMany(ctx, bson.M{"tagId": tagID})

	_, err = s.db.Collection("media_tags").DeleteOne(ctx, bson.M{"_id": tagID})
	return err
}

func (s *TagService) TagMedia(ctx context.Context, mediaID string, tagIDs []string) error {
	for _, tagID := range tagIDs {
		mapping := models.MediaTagMapping{
			ID:      uuid.New().String(),
			MediaID: mediaID,
			TagID:   tagID,
			AddedAt: time.Now(),
		}
		_, _ = s.db.Collection("media_tag_mappings").InsertOne(ctx, mapping)

		// Increment count
		_, _ = s.db.Collection("media_tags").UpdateOne(ctx,
			bson.M{"_id": tagID},
			bson.M{"$inc": bson.M{"mediaCount": 1}},
		)
	}
	return nil
}

func (s *TagService) UntagMedia(ctx context.Context, mediaID, tagID string) error {
	result, err := s.db.Collection("media_tag_mappings").DeleteOne(ctx,
		bson.M{"mediaId": mediaID, "tagId": tagID},
	)
	if err != nil {
		return err
	}
	if result.DeletedCount > 0 {
		_, _ = s.db.Collection("media_tags").UpdateOne(ctx,
			bson.M{"_id": tagID},
			bson.M{"$inc": bson.M{"mediaCount": -1}},
		)
	}
	return nil
}

func (s *TagService) GetMediaTags(ctx context.Context, mediaID string) ([]models.MediaTag, error) {
	// Get tag IDs from mappings
	cursor, err := s.db.Collection("media_tag_mappings").Find(ctx, bson.M{"mediaId": mediaID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mappings []models.MediaTagMapping
	if err := cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}

	if len(mappings) == 0 {
		return []models.MediaTag{}, nil
	}

	tagIDs := make([]string, len(mappings))
	for i, m := range mappings {
		tagIDs[i] = m.TagID
	}

	tagCursor, err := s.db.Collection("media_tags").Find(ctx,
		bson.M{"_id": bson.M{"$in": tagIDs}},
	)
	if err != nil {
		return nil, err
	}
	defer tagCursor.Close(ctx)

	var tags []models.MediaTag
	if err := tagCursor.All(ctx, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (s *TagService) GetMediaByTag(ctx context.Context, tagID string, limit int64) ([]string, error) {
	cursor, err := s.db.Collection("media_tag_mappings").Find(ctx,
		bson.M{"tagId": tagID},
		options.Find().SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mappings []models.MediaTagMapping
	if err := cursor.All(ctx, &mappings); err != nil {
		return nil, err
	}

	mediaIDs := make([]string, len(mappings))
	for i, m := range mappings {
		mediaIDs[i] = m.MediaID
	}
	return mediaIDs, nil
}

func (s *TagService) BulkTag(ctx context.Context, mediaIDs, tagIDs []string) error {
	for _, mediaID := range mediaIDs {
		for _, tagID := range tagIDs {
			mapping := models.MediaTagMapping{
				ID:      uuid.New().String(),
				MediaID: mediaID,
				TagID:   tagID,
				AddedAt: time.Now(),
			}
			_, _ = s.db.Collection("media_tag_mappings").InsertOne(ctx, mapping)
		}
	}
	// Update counts
	for _, tagID := range tagIDs {
		_, _ = s.db.Collection("media_tags").UpdateOne(ctx,
			bson.M{"_id": tagID},
			bson.M{"$inc": bson.M{"mediaCount": len(mediaIDs)}},
		)
	}
	return nil
}

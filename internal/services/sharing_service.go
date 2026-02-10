package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SharingService struct {
	db *database.MongoDB
}

func NewSharingService(db *database.MongoDB) *SharingService {
	return &SharingService{db: db}
}

func (s *SharingService) ShareWithUser(ctx context.Context, userID string, req *models.ShareMediaRequest) (*models.MediaShare, error) {
	share := &models.MediaShare{
		ID:         uuid.New().String(),
		MediaID:    req.MediaID,
		SharedBy:   userID,
		SharedWith: req.SharedWith,
		Permission: req.Permission,
		CreatedAt:  time.Now(),
	}

	if req.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		share.ExpiresAt = &exp
	}

	_, err := s.db.Collection("media_shares").InsertOne(ctx, share)
	if err != nil {
		return nil, err
	}
	return share, nil
}

func (s *SharingService) GetSharedWithUser(ctx context.Context, userID string, limit int64) ([]models.MediaShare, error) {
	cursor, err := s.db.Collection("media_shares").Find(ctx,
		bson.M{"sharedWith": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shares []models.MediaShare
	if err := cursor.All(ctx, &shares); err != nil {
		return nil, err
	}
	return shares, nil
}

func (s *SharingService) GetSharedByUser(ctx context.Context, userID string, limit int64) ([]models.MediaShare, error) {
	cursor, err := s.db.Collection("media_shares").Find(ctx,
		bson.M{"sharedBy": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(limit),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var shares []models.MediaShare
	if err := cursor.All(ctx, &shares); err != nil {
		return nil, err
	}
	return shares, nil
}

func (s *SharingService) RevokeShare(ctx context.Context, shareID, userID string) error {
	_, err := s.db.Collection("media_shares").DeleteOne(ctx,
		bson.M{"_id": shareID, "sharedBy": userID},
	)
	return err
}

func (s *SharingService) CreateShareLink(ctx context.Context, mediaID, userID string, req *models.CreateShareLinkRequest) (*models.MediaShareLink, error) {
	token := generateToken(32)

	link := &models.MediaShareLink{
		ID:        uuid.New().String(),
		MediaID:   mediaID,
		CreatedBy: userID,
		Token:     token,
		IsActive:  true,
		ViewCount: 0,
		CreatedAt: time.Now(),
	}

	if req.ExpiresIn > 0 {
		exp := time.Now().Add(time.Duration(req.ExpiresIn) * time.Hour)
		link.ExpiresAt = &exp
	}

	_, err := s.db.Collection("media_share_links").InsertOne(ctx, link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

func (s *SharingService) GetShareLink(ctx context.Context, token string) (*models.MediaShareLink, error) {
	var link models.MediaShareLink
	err := s.db.Collection("media_share_links").FindOne(ctx,
		bson.M{"token": token, "isActive": true},
	).Decode(&link)
	if err != nil {
		return nil, err
	}

	// Increment view count
	_, _ = s.db.Collection("media_share_links").UpdateOne(ctx,
		bson.M{"_id": link.ID},
		bson.M{"$inc": bson.M{"viewCount": 1}},
	)

	return &link, nil
}

func (s *SharingService) DeactivateShareLink(ctx context.Context, linkID, userID string) error {
	_, err := s.db.Collection("media_share_links").UpdateOne(ctx,
		bson.M{"_id": linkID, "createdBy": userID},
		bson.M{"$set": bson.M{"isActive": false}},
	)
	return err
}

func (s *SharingService) GetShareLinks(ctx context.Context, mediaID, userID string) ([]models.MediaShareLink, error) {
	cursor, err := s.db.Collection("media_share_links").Find(ctx,
		bson.M{"mediaId": mediaID, "createdBy": userID},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var links []models.MediaShareLink
	if err := cursor.All(ctx, &links); err != nil {
		return nil, err
	}
	return links, nil
}

func generateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

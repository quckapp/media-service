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

type VersionService struct {
	db      *database.MongoDB
	storage *S3Storage
}

func NewVersionService(db *database.MongoDB, storage *S3Storage) *VersionService {
	return &VersionService{db: db, storage: storage}
}

func (s *VersionService) CreateVersion(ctx context.Context, mediaID, userID string, req *models.CreateVersionRequest) (*models.MediaVersion, error) {
	// Get current max version
	var latest models.MediaVersion
	err := s.db.Collection("media_versions").FindOne(ctx,
		bson.M{"mediaId": mediaID},
		options.FindOne().SetSort(bson.M{"version": -1}),
	).Decode(&latest)

	nextVersion := 1
	if err == nil {
		nextVersion = latest.Version + 1
	}

	s3Key := fmt.Sprintf("media/%s/%s/v%d/%s", userID, mediaID, nextVersion, req.Filename)

	version := &models.MediaVersion{
		ID:         uuid.New().String(),
		MediaID:    mediaID,
		Version:    nextVersion,
		Filename:   req.Filename,
		MimeType:   req.MimeType,
		Size:       req.Size,
		S3Key:      s3Key,
		UploadedBy: userID,
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}

	_, err = s.db.Collection("media_versions").InsertOne(ctx, version)
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (s *VersionService) GetVersions(ctx context.Context, mediaID string) ([]models.MediaVersion, error) {
	cursor, err := s.db.Collection("media_versions").Find(ctx,
		bson.M{"mediaId": mediaID},
		options.Find().SetSort(bson.M{"version": -1}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []models.MediaVersion
	if err := cursor.All(ctx, &versions); err != nil {
		return nil, err
	}

	// Generate signed URLs
	for i := range versions {
		versions[i].URL, _ = s.storage.GetPresignedDownloadURL(versions[i].S3Key, time.Hour)
	}

	return versions, nil
}

func (s *VersionService) GetVersion(ctx context.Context, versionID string) (*models.MediaVersion, error) {
	var version models.MediaVersion
	err := s.db.Collection("media_versions").FindOne(ctx, bson.M{"_id": versionID}).Decode(&version)
	if err != nil {
		return nil, err
	}

	version.URL, _ = s.storage.GetPresignedDownloadURL(version.S3Key, time.Hour)
	return &version, nil
}

func (s *VersionService) DeleteVersion(ctx context.Context, versionID, userID string) error {
	version, err := s.GetVersion(ctx, versionID)
	if err != nil {
		return err
	}
	if version.UploadedBy != userID {
		return fmt.Errorf("unauthorized")
	}

	// Delete from S3
	_ = s.storage.Delete(version.S3Key)

	_, err = s.db.Collection("media_versions").DeleteOne(ctx, bson.M{"_id": versionID})
	return err
}

func (s *VersionService) RestoreVersion(ctx context.Context, versionID, userID string) error {
	version, err := s.GetVersion(ctx, versionID)
	if err != nil {
		return err
	}

	// Update the main media record to point to this version
	_, err = s.db.Collection("media").UpdateOne(ctx,
		bson.M{"_id": version.MediaID},
		bson.M{"$set": bson.M{
			"filename":  version.Filename,
			"mimeType":  version.MimeType,
			"size":      version.Size,
			"s3Key":     version.S3Key,
			"updatedAt": time.Now(),
		}},
	)
	return err
}

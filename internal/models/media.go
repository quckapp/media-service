package models

import "time"

type Media struct {
	ID          string            `json:"id" bson:"_id"`
	UserID      string            `json:"userId" bson:"userId"`
	Type        string            `json:"type" bson:"type"` // image, video, audio, document
	Filename    string            `json:"filename" bson:"filename"`
	MimeType    string            `json:"mimeType" bson:"mimeType"`
	Size        int64             `json:"size" bson:"size"`
	URL         string            `json:"url" bson:"url"`
	ThumbnailURL string           `json:"thumbnailUrl,omitempty" bson:"thumbnailUrl,omitempty"`
	S3Key       string            `json:"s3Key" bson:"s3Key"`
	Metadata    map[string]string `json:"metadata,omitempty" bson:"metadata,omitempty"`
	CreatedAt   time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type UploadRequest struct {
	Filename string `json:"filename" binding:"required"`
	MimeType string `json:"mimeType" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
}

type PresignedURLResponse struct {
	UploadURL string `json:"uploadUrl"`
	MediaID   string `json:"mediaId"`
	S3Key     string `json:"s3Key"`
	ExpiresAt string `json:"expiresAt"`
}

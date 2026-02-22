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

type BulkDeleteRequest struct {
	MediaIDs []string `json:"media_ids" binding:"required,min=1,max=50"`
}

type BulkDeleteResponse struct {
	Deleted []string `json:"deleted"`
	Failed  []string `json:"failed"`
}

type MediaStatsResponse struct {
	TotalFiles int64            `json:"totalFiles"`
	TotalSize  int64            `json:"totalSize"`
	ByType     map[string]int64 `json:"byType"`
}

// ── Albums ──

type MediaAlbum struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"userId" bson:"userId"`
	WorkspaceID string    `json:"workspaceId" bson:"workspaceId"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description,omitempty" bson:"description,omitempty"`
	CoverURL    string    `json:"coverUrl,omitempty" bson:"coverUrl,omitempty"`
	MediaIDs    []string  `json:"mediaIds" bson:"mediaIds"`
	IsPublic    bool      `json:"isPublic" bson:"isPublic"`
	Position    int       `json:"position" bson:"position"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CreateAlbumRequest struct {
	Name        string `json:"name" binding:"required"`
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"isPublic"`
}

type UpdateAlbumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	IsPublic    *bool  `json:"isPublic"`
}

type AddToAlbumRequest struct {
	MediaIDs []string `json:"mediaIds" binding:"required,min=1"`
}

type RemoveFromAlbumRequest struct {
	MediaIDs []string `json:"mediaIds" binding:"required,min=1"`
}

// ── Tags ──

type MediaTag struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"userId" bson:"userId"`
	WorkspaceID string    `json:"workspaceId" bson:"workspaceId"`
	Name        string    `json:"name" bson:"name"`
	Color       string    `json:"color,omitempty" bson:"color,omitempty"`
	MediaCount  int       `json:"mediaCount" bson:"mediaCount"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Color       string `json:"color"`
}

type UpdateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TagMediaRequest struct {
	TagIDs []string `json:"tagIds" binding:"required,min=1"`
}

type MediaTagMapping struct {
	ID      string    `json:"id" bson:"_id"`
	MediaID string    `json:"mediaId" bson:"mediaId"`
	TagID   string    `json:"tagId" bson:"tagId"`
	AddedAt time.Time `json:"addedAt" bson:"addedAt"`
}

// ── Sharing ──

type MediaShare struct {
	ID         string     `json:"id" bson:"_id"`
	MediaID    string     `json:"mediaId" bson:"mediaId"`
	SharedBy   string     `json:"sharedBy" bson:"sharedBy"`
	SharedWith string     `json:"sharedWith" bson:"sharedWith"`
	Permission string     `json:"permission" bson:"permission"` // view, download, edit
	ExpiresAt  *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedAt  time.Time  `json:"createdAt" bson:"createdAt"`
}

type ShareMediaRequest struct {
	MediaID    string `json:"mediaId" binding:"required"`
	SharedWith string `json:"sharedWith" binding:"required"`
	Permission string `json:"permission" binding:"required"`
	ExpiresIn  int    `json:"expiresIn"` // hours, 0 = no expiry
}

type ShareLinkResponse struct {
	ShareID   string     `json:"shareId"`
	ShareLink string     `json:"shareLink"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type MediaShareLink struct {
	ID        string     `json:"id" bson:"_id"`
	MediaID   string     `json:"mediaId" bson:"mediaId"`
	CreatedBy string     `json:"createdBy" bson:"createdBy"`
	Token     string     `json:"token" bson:"token"`
	IsActive  bool       `json:"isActive" bson:"isActive"`
	ViewCount int        `json:"viewCount" bson:"viewCount"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty" bson:"expiresAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt" bson:"createdAt"`
}

type CreateShareLinkRequest struct {
	ExpiresIn int `json:"expiresIn"` // hours, 0 = no expiry
}

// ── Versions ──

type MediaVersion struct {
	ID         string    `json:"id" bson:"_id"`
	MediaID    string    `json:"mediaId" bson:"mediaId"`
	Version    int       `json:"version" bson:"version"`
	Filename   string    `json:"filename" bson:"filename"`
	MimeType   string    `json:"mimeType" bson:"mimeType"`
	Size       int64     `json:"size" bson:"size"`
	S3Key      string    `json:"s3Key" bson:"s3Key"`
	URL        string    `json:"url,omitempty" bson:"url,omitempty"`
	UploadedBy string    `json:"uploadedBy" bson:"uploadedBy"`
	Comment    string    `json:"comment,omitempty" bson:"comment,omitempty"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
}

type CreateVersionRequest struct {
	Filename string `json:"filename" binding:"required"`
	MimeType string `json:"mimeType" binding:"required"`
	Size     int64  `json:"size" binding:"required"`
	Comment  string `json:"comment"`
}

// ── Trash ──

type TrashedMedia struct {
	ID          string    `json:"id" bson:"_id"`
	MediaID     string    `json:"mediaId" bson:"mediaId"`
	UserID      string    `json:"userId" bson:"userId"`
	OriginalDoc Media     `json:"originalDoc" bson:"originalDoc"`
	TrashedAt   time.Time `json:"trashedAt" bson:"trashedAt"`
	ExpiresAt   time.Time `json:"expiresAt" bson:"expiresAt"` // auto-delete after 30 days
}

// ── Processing Jobs ──

type ProcessingJob struct {
	ID        string                 `json:"id" bson:"_id"`
	MediaID   string                 `json:"mediaId" bson:"mediaId"`
	UserID    string                 `json:"userId" bson:"userId"`
	Type      string                 `json:"type" bson:"type"`     // thumbnail, resize, compress, transcode
	Status    string                 `json:"status" bson:"status"` // pending, processing, completed, failed
	Params    map[string]interface{} `json:"params,omitempty" bson:"params,omitempty"`
	Result    map[string]interface{} `json:"result,omitempty" bson:"result,omitempty"`
	Error     string                 `json:"error,omitempty" bson:"error,omitempty"`
	CreatedAt time.Time              `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt" bson:"updatedAt"`
}

type CreateProcessingJobRequest struct {
	Type   string                 `json:"type" binding:"required"`
	Params map[string]interface{} `json:"params"`
}

// ── Favorites ──

type MediaFavorite struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"userId" bson:"userId"`
	MediaID   string    `json:"mediaId" bson:"mediaId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

// ── Comments ──

type MediaComment struct {
	ID        string    `json:"id" bson:"_id"`
	MediaID   string    `json:"mediaId" bson:"mediaId"`
	UserID    string    `json:"userId" bson:"userId"`
	Content   string    `json:"content" bson:"content"`
	ParentID  string    `json:"parentId,omitempty" bson:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CreateCommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parentId"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// ── Search & Filters ──

type MediaSearchParams struct {
	Query       string `form:"q"`
	Type        string `form:"type"`
	WorkspaceID string `form:"workspaceId"`
	ChannelID   string `form:"channelId"`
	AlbumID     string `form:"albumId"`
	TagID       string `form:"tagId"`
	SortBy      string `form:"sortBy"`    // createdAt, size, filename
	SortOrder   string `form:"sortOrder"` // asc, desc
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	MinSize     int64  `form:"minSize"`
	MaxSize     int64  `form:"maxSize"`
	DateFrom    string `form:"dateFrom"`
	DateTo      string `form:"dateTo"`
}

// ── Workspace Stats ──

type WorkspaceMediaStats struct {
	TotalFiles    int64            `json:"totalFiles"`
	TotalSize     int64            `json:"totalSize"`
	ByType        map[string]int64 `json:"byType"`
	ByUser        map[string]int64 `json:"byUser"`
	StorageUsed   int64            `json:"storageUsed"`
	StorageLimit  int64            `json:"storageLimit"`
	RecentUploads int64            `json:"recentUploads"`
}

// ── File Operations ──

type RenameRequest struct {
	Filename string `json:"filename" binding:"required"`
}

type CopyMediaRequest struct {
	TargetWorkspaceID string `json:"targetWorkspaceId" binding:"required"`
}

type MoveMediaRequest struct {
	TargetWorkspaceID string `json:"targetWorkspaceId" binding:"required"`
}

type BulkTagRequest struct {
	MediaIDs []string `json:"mediaIds" binding:"required,min=1"`
	TagIDs   []string `json:"tagIds" binding:"required,min=1"`
}

type BulkMoveRequest struct {
	MediaIDs          []string `json:"mediaIds" binding:"required,min=1"`
	TargetWorkspaceID string   `json:"targetWorkspaceId" binding:"required"`
}

type BulkAlbumRequest struct {
	MediaIDs []string `json:"mediaIds" binding:"required,min=1"`
	AlbumID  string   `json:"albumId" binding:"required"`
}

// ── Activity / Audit ──

type MediaActivity struct {
	ID        string    `json:"id" bson:"_id"`
	MediaID   string    `json:"mediaId" bson:"mediaId"`
	UserID    string    `json:"userId" bson:"userId"`
	Action    string    `json:"action" bson:"action"` // upload, delete, view, download, share, rename, move, tag
	Details   string    `json:"details,omitempty" bson:"details,omitempty"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

// Retention Policies

type RetentionPolicy struct {
	ID            string    `json:"id" bson:"_id"`
	WorkspaceID   string    `json:"workspaceId" bson:"workspaceId"`
	Name          string    `json:"name" bson:"name"`
	RetentionDays int       `json:"retentionDays" bson:"retentionDays"`
	ApplyTo       string    `json:"applyTo" bson:"applyTo"`
	AutoDelete    bool      `json:"autoDelete" bson:"autoDelete"`
	CreatedBy     string    `json:"createdBy" bson:"createdBy"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CreateRetentionPolicyRequest struct {
	WorkspaceID   string `json:"workspaceId" binding:"required"`
	Name          string `json:"name" binding:"required"`
	RetentionDays int    `json:"retentionDays" binding:"required"`
	ApplyTo       string `json:"applyTo" binding:"required"`
	AutoDelete    bool   `json:"autoDelete"`
}

type UpdateRetentionPolicyRequest struct {
	Name          string `json:"name"`
	RetentionDays *int   `json:"retentionDays"`
	ApplyTo       string `json:"applyTo"`
	AutoDelete    *bool  `json:"autoDelete"`
}

// Storage Quotas

type StorageQuota struct {
	ID               string    `json:"id" bson:"_id"`
	WorkspaceID      string    `json:"workspaceId" bson:"workspaceId"`
	MaxStorageMB     int64     `json:"maxStorageMB" bson:"maxStorageMB"`
	UsedStorageMB    int64     `json:"usedStorageMB" bson:"usedStorageMB"`
	MaxFileCount     int64     `json:"maxFileCount" bson:"maxFileCount"`
	CurrentFileCount int64     `json:"currentFileCount" bson:"currentFileCount"`
	CreatedAt        time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
}

type SetQuotaRequest struct {
	WorkspaceID  string `json:"workspaceId" binding:"required"`
	MaxStorageMB int64  `json:"maxStorageMB" binding:"required"`
	MaxFileCount int64  `json:"maxFileCount" binding:"required"`
}

// Watermarks

type Watermark struct {
	ID          string    `json:"id" bson:"_id"`
	WorkspaceID string    `json:"workspaceId" bson:"workspaceId"`
	Name        string    `json:"name" bson:"name"`
	ImageURL    string    `json:"imageUrl" bson:"imageUrl"`
	Position    string    `json:"position" bson:"position"`
	Opacity     float64   `json:"opacity" bson:"opacity"`
	Size        int       `json:"size" bson:"size"`
	CreatedBy   string    `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type WatermarkSettings struct {
	ID          string    `json:"id" bson:"_id"`
	WorkspaceID string    `json:"workspaceId" bson:"workspaceId"`
	AutoApply   bool      `json:"autoApply" bson:"autoApply"`
	ApplyTo     string    `json:"applyTo" bson:"applyTo"`
	Position    string    `json:"position" bson:"position"`
	Opacity     float64   `json:"opacity" bson:"opacity"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UploadWatermarkRequest struct {
	WorkspaceID string  `json:"workspaceId" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	ImageURL    string  `json:"imageUrl" binding:"required"`
	Position    string  `json:"position" binding:"required"`
	Opacity     float64 `json:"opacity"`
	Size        int     `json:"size"`
}

type ApplyWatermarkRequest struct {
	WatermarkID string  `json:"watermarkId" binding:"required"`
	MediaID     string  `json:"mediaId" binding:"required"`
	Position    string  `json:"position"`
	Opacity     float64 `json:"opacity"`
}

// Media Scanning / Moderation

type MediaScan struct {
	ID         string    `json:"id" bson:"_id"`
	MediaID    string    `json:"mediaId" bson:"mediaId"`
	ScanType   string    `json:"scanType" bson:"scanType"`
	Status     string    `json:"status" bson:"status"`
	Confidence float64   `json:"confidence" bson:"confidence"`
	Details    string    `json:"details" bson:"details"`
	ScannedAt  time.Time `json:"scannedAt" bson:"scannedAt"`
}

type ScanRequest struct {
	MediaID  string `json:"mediaId" binding:"required"`
	ScanType string `json:"scanType" binding:"required"`
}

// Media Analytics

type UploadTrend struct {
	Date      string `json:"date" bson:"_id"`
	Count     int64  `json:"count" bson:"count"`
	TotalSize int64  `json:"totalSize" bson:"totalSize"`
}

type StorageTrend struct {
	Date   string `json:"date" bson:"_id"`
	UsedMB int64  `json:"usedMB" bson:"usedMB"`
}

type FileTypeDistribution struct {
	Type      string `json:"type" bson:"_id"`
	Count     int64  `json:"count" bson:"count"`
	TotalSize int64  `json:"totalSize" bson:"totalSize"`
}

type UserUploadStats struct {
	UserID     string    `json:"userId" bson:"_id"`
	FileCount  int64     `json:"fileCount" bson:"fileCount"`
	TotalSize  int64     `json:"totalSize" bson:"totalSize"`
	LastUpload time.Time `json:"lastUpload" bson:"lastUpload"`
}

// Media Galleries / Collections

type MediaGallery struct {
	ID          string    `json:"id" bson:"_id"`
	WorkspaceID string    `json:"workspaceId" bson:"workspaceId"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	CoverURL    string    `json:"coverUrl" bson:"coverUrl"`
	MediaIDs    []string  `json:"mediaIds" bson:"mediaIds"`
	Layout      string    `json:"layout" bson:"layout"`
	IsPublic    bool      `json:"isPublic" bson:"isPublic"`
	CreatedBy   string    `json:"createdBy" bson:"createdBy"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CreateGalleryRequest struct {
	WorkspaceID string `json:"workspaceId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Layout      string `json:"layout"`
	IsPublic    bool   `json:"isPublic"`
}

type UpdateGalleryRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	CoverURL    string   `json:"coverUrl"`
	Layout      string   `json:"layout"`
	IsPublic    *bool    `json:"isPublic"`
	MediaIDs    []string `json:"mediaIds"`
}
// ── Paginated Response ──

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

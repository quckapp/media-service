package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quckapp/media-service/internal/config"
	"github.com/quckapp/media-service/internal/database"
	"github.com/quckapp/media-service/internal/handlers"
	"github.com/quckapp/media-service/internal/services"
)

func main() {
	cfg := config.Load()

	// Initialize MongoDB
	mongoDB, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	// Initialize Redis
	redisClient := database.NewRedis(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, 7)
	defer redisClient.Close()

	// Initialize S3 storage
	s3Storage, err := services.NewS3Storage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3: %v", err)
	}

	// ── Initialize Services ──
	mediaService := services.NewMediaService(mongoDB, redisClient, s3Storage)
	albumService := services.NewAlbumService(mongoDB)
	tagService := services.NewTagService(mongoDB)
	sharingService := services.NewSharingService(mongoDB)
	versionService := services.NewVersionService(mongoDB, s3Storage)
	trashService := services.NewTrashService(mongoDB, redisClient, s3Storage)
	processingService := services.NewProcessingService(mongoDB)
	favoriteService := services.NewFavoriteService(mongoDB)
	commentService := services.NewCommentService(mongoDB)
	activityService := services.NewActivityService(mongoDB)
	searchService := services.NewSearchService(mongoDB, s3Storage)
	retentionService := services.NewRetentionService(mongoDB)
	quotaService := services.NewQuotaService(mongoDB)
	watermarkService := services.NewWatermarkService(mongoDB)
	scanningService := services.NewScanningService(mongoDB)
	analyticsService := services.NewAnalyticsService(mongoDB)
	galleryService := services.NewGalleryService(mongoDB)

	// ── Initialize Handlers ──
	mediaHandler := handlers.NewMediaHandler(mediaService)
	albumHandler := handlers.NewAlbumHandler(albumService)
	tagHandler := handlers.NewTagHandler(tagService)
	sharingHandler := handlers.NewSharingHandler(sharingService)
	versionHandler := handlers.NewVersionHandler(versionService)
	trashHandler := handlers.NewTrashHandler(trashService)
	processingHandler := handlers.NewProcessingHandler(processingService)
	favoriteHandler := handlers.NewFavoriteHandler(favoriteService)
	commentHandler := handlers.NewCommentHandler(commentService)
	activityHandler := handlers.NewActivityHandler(activityService)
	searchHandler := handlers.NewSearchHandler(searchService, mediaService)
	healthHandler := handlers.NewHealthHandler(mongoDB, redisClient)
	retentionHandler := handlers.NewRetentionHandler(retentionService)
	quotaHandler := handlers.NewQuotaHandler(quotaService)
	watermarkHandler := handlers.NewWatermarkHandler(watermarkService)
	scanningHandler := handlers.NewScanningHandler(scanningService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
	galleryHandler := handlers.NewGalleryHandler(galleryService)

	// Setup router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Health endpoints
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Ready)

	// API routes
	api := router.Group("/api/v1/media")
	api.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		// ── Upload ──
		api.POST("/upload", mediaHandler.Upload)
		api.POST("/upload/presigned", mediaHandler.GetPresignedURL)

		// ── Bulk Operations ──
		api.POST("/bulk-delete", mediaHandler.BulkDelete)
		api.POST("/bulk-move", searchHandler.BulkMove)
		api.POST("/bulk-tag", tagHandler.BulkTag)

		// ── Search & Discovery ──
		api.GET("/search", searchHandler.Search)
		api.GET("/recent", searchHandler.GetRecent)
		api.GET("/duplicates", searchHandler.GetDuplicates)

		// ── Single Media Operations ──
		api.GET("/:id", mediaHandler.Get)
		api.DELETE("/:id", mediaHandler.Delete)
		api.PUT("/:id/metadata", mediaHandler.UpdateMetadata)
		api.POST("/:id/thumbnail", mediaHandler.GenerateThumbnail)
		api.PUT("/:id/rename", searchHandler.Rename)
		api.POST("/:id/copy", searchHandler.CopyMedia)
		api.POST("/:id/move", searchHandler.MoveMedia)
		api.GET("/:id/download-url", searchHandler.GetDownloadURL)

		// ── Tags on Media ──
		api.POST("/:id/tags", tagHandler.TagMedia)
		api.DELETE("/:id/tags/:tagId", tagHandler.UntagMedia)
		api.GET("/:id/tags", tagHandler.GetMediaTags)

		// ── Versions ──
		api.POST("/:id/versions", versionHandler.CreateVersion)
		api.GET("/:id/versions", versionHandler.GetVersions)
		api.GET("/:id/versions/:versionId", versionHandler.GetVersion)
		api.DELETE("/:id/versions/:versionId", versionHandler.DeleteVersion)
		api.POST("/:id/versions/:versionId/restore", versionHandler.RestoreVersion)

		// ── Trash (Soft Delete) ──
		api.POST("/:id/trash", trashHandler.MoveToTrash)

		// ── Processing Jobs ──
		api.POST("/:id/process", processingHandler.CreateJob)
		api.GET("/:id/jobs", processingHandler.GetJobsByMedia)

		// ── Favorites ──
		api.POST("/:id/favorite", favoriteHandler.AddFavorite)
		api.DELETE("/:id/favorite", favoriteHandler.RemoveFavorite)
		api.GET("/:id/favorite", favoriteHandler.IsFavorite)

		// ── Comments ──
		api.POST("/:id/comments", commentHandler.Create)
		api.GET("/:id/comments", commentHandler.GetByMedia)
		api.GET("/:id/comments/count", commentHandler.CountByMedia)
		api.GET("/:id/comments/:commentId/replies", commentHandler.GetReplies)
		api.PUT("/:id/comments/:commentId", commentHandler.Update)
		api.DELETE("/:id/comments/:commentId", commentHandler.Delete)

		// ── Activity / Audit ──
		api.GET("/:id/activity", activityHandler.GetByMedia)

		// ── Share Links ──
		api.POST("/:id/share-links", sharingHandler.CreateShareLink)
		api.GET("/:id/share-links", sharingHandler.GetShareLinks)
	}

	// ── Sharing (User-Scoped) ──
	shares := router.Group("/api/v1/media/shares")
	shares.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		shares.POST("", sharingHandler.ShareWithUser)
		shares.GET("/received", sharingHandler.GetSharedWithMe)
		shares.GET("/sent", sharingHandler.GetSharedByMe)
		shares.DELETE("/:shareId", sharingHandler.RevokeShare)
	}

	// ── Public Share Link Access ──
	router.GET("/api/v1/media/shared/:token", sharingHandler.GetShareLink)

	// ── Share Link Deactivation ──
	shareLinks := router.Group("/api/v1/media/share-links")
	shareLinks.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		shareLinks.DELETE("/:linkId", sharingHandler.DeactivateShareLink)
	}

	// ── Trash Management ──
	trash := router.Group("/api/v1/media/trash")
	trash.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		trash.GET("", trashHandler.GetTrash)
		trash.POST("/:trashId/restore", trashHandler.RestoreFromTrash)
		trash.DELETE("/:trashId", trashHandler.PermanentDelete)
		trash.DELETE("", trashHandler.EmptyTrash)
	}

	// ── Processing Jobs (User-Scoped) ──
	jobs := router.Group("/api/v1/media/jobs")
	jobs.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		jobs.GET("", processingHandler.GetUserJobs)
		jobs.GET("/:jobId", processingHandler.GetJob)
		jobs.POST("/:jobId/cancel", processingHandler.CancelJob)
	}

	// ── Favorites (User-Scoped) ──
	favorites := router.Group("/api/v1/media/favorites")
	favorites.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		favorites.GET("", favoriteHandler.GetFavorites)
		favorites.GET("/count", favoriteHandler.GetFavoriteCount)
	}

	// ── Albums ──
	albums := router.Group("/api/v1/media/albums")
	albums.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		albums.POST("", albumHandler.Create)
		albums.GET("/:albumId", albumHandler.GetByID)
		albums.PUT("/:albumId", albumHandler.Update)
		albums.DELETE("/:albumId", albumHandler.Delete)
		albums.POST("/:albumId/media", albumHandler.AddMedia)
		albums.DELETE("/:albumId/media", albumHandler.RemoveMedia)
	}

	// ── Tags (CRUD) ──
	tags := router.Group("/api/v1/media/tags")
	tags.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		tags.POST("", tagHandler.Create)
		tags.PUT("/:tagId", tagHandler.Update)
		tags.DELETE("/:tagId", tagHandler.Delete)
		tags.GET("/:tagId/media", tagHandler.GetMediaByTag)
	}

	// ── User Activity ──
	userActivity := router.Group("/api/v1/media/activity")
	userActivity.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		userActivity.GET("", activityHandler.GetByUser)
	}

	// ── Query Endpoints ──
	query := router.Group("/api/v1/media/user")
	query.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		query.GET("/:userId", mediaHandler.GetUserMedia)
		query.GET("/:userId/stats", mediaHandler.GetUserStats)
		query.GET("/:userId/albums", albumHandler.GetByUser)
		query.GET("/:userId/tags", tagHandler.GetByUser)
		query.GET("/:userId/type/:type", searchHandler.GetByType)
	}

	// ── Workspace Endpoints ──
	workspace := router.Group("/api/v1/media/workspace")
	workspace.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		workspace.GET("/:workspaceId", mediaHandler.GetWorkspaceMedia)
		workspace.GET("/:workspaceId/stats", searchHandler.GetWorkspaceStats)
		workspace.GET("/:workspaceId/albums", albumHandler.GetByWorkspace)
		workspace.GET("/:workspaceId/albums/public", albumHandler.GetPublicByWorkspace)
		workspace.GET("/:workspaceId/tags", tagHandler.GetByWorkspace)
	}

	// ── Channel Endpoints ──
	channel := router.Group("/api/v1/media/channel")
	channel.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		channel.GET("/:channelId", mediaHandler.GetChannelMedia)
	}

	// ── Retention Policies ──
	retention := router.Group("/api/v1/media/retention")
	retention.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		retention.POST("", retentionHandler.Create)
		retention.GET("/:policyId", retentionHandler.Get)
		retention.PUT("/:policyId", retentionHandler.Update)
		retention.DELETE("/:policyId", retentionHandler.Delete)
		retention.GET("/workspace/:workspaceId", retentionHandler.GetByWorkspace)
	}

	// ── Storage Quotas ──
	quotas := router.Group("/api/v1/media/quotas")
	quotas.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		quotas.GET("/:workspaceId", quotaHandler.GetQuota)
		quotas.POST("", quotaHandler.SetQuota)
		quotas.GET("/:workspaceId/usage", quotaHandler.GetUsage)
		quotas.GET("/over-quota", quotaHandler.ListOverQuota)
	}

	// ── Watermarks ──
	watermarks := router.Group("/api/v1/media/watermarks")
	watermarks.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		watermarks.POST("", watermarkHandler.Upload)
		watermarks.GET("/workspace/:workspaceId", watermarkHandler.List)
		watermarks.POST("/apply", watermarkHandler.Apply)
		watermarks.DELETE("/:mediaId", watermarkHandler.Remove)
		watermarks.GET("/settings/:workspaceId", watermarkHandler.GetSettings)
	}

	// ── Media Scanning / Moderation ──
	scanning := router.Group("/api/v1/media/scanning")
	scanning.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		scanning.POST("/scan", scanningHandler.ScanMedia)
		scanning.GET("/:mediaId/results", scanningHandler.GetResults)
		scanning.GET("/flagged", scanningHandler.ListFlagged)
		scanning.PUT("/:scanId/status", scanningHandler.UpdateStatus)
	}

	// ── Media Analytics ──
	analytics := router.Group("/api/v1/media/analytics")
	analytics.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		analytics.GET("/:workspaceId/upload-trends", analyticsHandler.GetUploadTrends)
		analytics.GET("/:workspaceId/storage-trends", analyticsHandler.GetStorageTrends)
		analytics.GET("/:workspaceId/file-types", analyticsHandler.GetFileTypeDistribution)
		analytics.GET("/:workspaceId/user-stats", analyticsHandler.GetUserUploadStats)
	}

	// ── Media Galleries ──
	galleries := router.Group("/api/v1/media/galleries")
	galleries.Use(handlers.AuthMiddleware(cfg.JWTSecret))
	{
		galleries.POST("", galleryHandler.Create)
		galleries.GET("/workspace/:workspaceId", galleryHandler.List)
		galleries.GET("/:galleryId", galleryHandler.Get)
		galleries.PUT("/:galleryId", galleryHandler.Update)
		galleries.DELETE("/:galleryId", galleryHandler.Delete)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Media service starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Media service stopped")
}

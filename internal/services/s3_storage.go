package services

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/quckapp/media-service/internal/config"
)

type S3Storage struct {
	client  *s3.Client
	bucket  string
	presign *s3.PresignClient
}

func NewS3Storage(cfg *config.Config) (*S3Storage, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWSAccessKey,
			cfg.AWSSecretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)
	presign := s3.NewPresignClient(client)

	return &S3Storage{
		client:  client,
		bucket:  cfg.S3Bucket,
		presign: presign,
	}, nil
}

func (s *S3Storage) GetPresignedUploadURL(key, contentType string, expiry time.Duration) (string, error) {
	req, err := s.presign.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *S3Storage) GetPresignedDownloadURL(key string, expiry time.Duration) (string, error) {
	req, err := s.presign.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

func (s *S3Storage) Delete(key string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

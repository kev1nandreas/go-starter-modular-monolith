package factory

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/kev1nandreas/go-rest-api-template/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage wraps MinIO. Internal client does bucket ops; public client signs
// presigned URLs against the public host (signature is bound to the host).
type Storage struct {
	client         *minio.Client
	publicClient   *minio.Client
	bucket         string
	UploadExpiry   time.Duration
	DownloadExpiry time.Duration
}

func NewStorage(ctx context.Context) (*Storage, error) {
	endpoint := config.GetEnvString("MINIO_ENDPOINT", "localhost:9000")
	accessKey := config.GetEnvString("MINIO_ACCESS_KEY", "")
	secretKey := config.GetEnvString("MINIO_SECRET_KEY", "")
	bucket := config.GetEnvString("MINIO_BUCKET", "")
	useSSL := config.GetEnvBool("MINIO_USE_SSL", false)

	publicEndpoint := config.GetEnvString("MINIO_PUBLIC_ENDPOINT", endpoint)
	publicUseSSL := config.GetEnvBool("MINIO_PUBLIC_USE_SSL", useSSL)

	newClient := func(ep string, ssl bool) (*minio.Client, error) {
		return minio.New(ep, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: ssl,
		})
	}

	client, err := newClient(endpoint, useSSL)
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}
	publicClient, err := newClient(publicEndpoint, publicUseSSL)
	if err != nil {
		return nil, fmt.Errorf("minio public client: %w", err)
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("minio bucket check: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("minio make bucket: %w", err)
		}
	}

	return &Storage{
		client:         client,
		publicClient:   publicClient,
		bucket:         bucket,
		UploadExpiry:   time.Duration(config.GetEnvInt("STORAGE_UPLOAD_URL_EXPIRY_MINUTES", 15)) * time.Minute,
		DownloadExpiry: time.Duration(config.GetEnvInt("STORAGE_DOWNLOAD_URL_EXPIRY_MINUTES", 60)) * time.Minute,
	}, nil
}

// Client returns the internal MinIO client for direct operations.
func (s *Storage) Client() *minio.Client { return s.client }

// Bucket returns the configured bucket name.
func (s *Storage) Bucket() string { return s.bucket }

// PresignedUploadURL returns a PUT URL valid for UploadExpiry.
func (s *Storage) PresignedUploadURL(ctx context.Context, objectKey string) (*url.URL, error) {
	return s.publicClient.PresignedPutObject(ctx, s.bucket, objectKey, s.UploadExpiry)
}

// PresignedDownloadURL returns a GET URL valid for DownloadExpiry.
func (s *Storage) PresignedDownloadURL(ctx context.Context, objectKey string) (*url.URL, error) {
	return s.publicClient.PresignedGetObject(ctx, s.bucket, objectKey, s.DownloadExpiry, nil)
}

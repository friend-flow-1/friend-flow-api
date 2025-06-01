package minio

import (
	"context"
	"fmt"

	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(cfg *config.Config) (*minio.Client, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MinIO: %w", err)
	}

	// Check bucket or create
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		return nil, fmt.Errorf("error checking MinIO bucket: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return client, nil
}

package minioinit

import (
	"context"
	"log"

	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio(cfg *config.Config) *minio.Client {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinioAccessKey, cfg.MinioSecretKey, ""),
		Secure: cfg.MinioUseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to MinIO: %v", err)
	}

	log.Println("✅ Connected to MinIO")

	// Check bucket or create
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.MinioBucket)
	if err != nil {
		log.Fatalf("❌ Error checking MinIO bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.MinioBucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		log.Printf("🪣 Created bucket: %s", cfg.MinioBucket)
	} else {
		log.Printf("🪣 Bucket already exists: %s", cfg.MinioBucket)
	}

	return client
}

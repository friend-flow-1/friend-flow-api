package internal

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
	chat_server "github.com/haxxu/friend-flow-api/internal/modules/chat/server"
	"github.com/haxxu/friend-flow-api/internal/modules/media_upload"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"github.com/haxxu/friend-flow-api/internal/rbac"
	"github.com/haxxu/friend-flow-api/internal/routes"
	minioinit "github.com/haxxu/friend-flow-api/pkg/minio"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type App struct {
	Router     *gin.Engine
	Session    *gocql.Session
	PostgresDB *gorm.DB
}

func InitializeApp(cfg *config.Config) (*App, error) {
	postgresDB, err := db.InitPostgres(cfg)
	if err != nil {
		return nil, err
	}

	if err := rbac.InitEnforcerPostgres(postgresDB); err != nil {
		return nil, err
	}

	// Minio Start ----------------------------------------------------------
	// Minio client initialization can be added here if needed
	minioClient := minioinit.InitMinio(cfg)
	ctx := context.Background()

	// Make sure bucket exists
	bucket := cfg.MinioBucket
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("❌ Failed to check bucket: %v", err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			log.Fatalf("❌ Failed to create bucket: %v", err)
		}
		log.Printf("✅ Bucket created: %s", bucket)
	} else {
		log.Printf("🪣 Bucket already exists: %s", bucket)
	}

	// Prepare file data
	objectName := "test-file.txt"
	content := []byte("This is just a test file 🧪")
	contentType := "text/plain"

	// Upload to MinIO
	_, err = minioClient.PutObject(ctx, bucket, objectName, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		log.Fatalf("❌ Upload failed: %v", err)
	}

	log.Printf("✅ Successfully uploaded %s to bucket %s", objectName, bucket)
	// Minio End ----------------------------------------------------------

	// Router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize User module
	userModule := user.InitModule(postgresDB)

	authModule := auth.InitModule(postgresDB, rbac.EnforcerPG, cfg, userModule.Service)

	chatServerModule := chat_server.InitModule(postgresDB)

	mediaUploadModule := media_upload.InitModule(postgresDB)

	routes.SetupRoutes(router, &routes.Handlers{
		AuthHandler:        authModule.Handler,
		UserHandler:        userModule.Handler,
		ChatServerHandler:  chatServerModule.Handler,
		MediaUploadHandler: mediaUploadModule.Handler,
	})

	return &App{
		Router:     router,
		PostgresDB: postgresDB,
	}, nil
}

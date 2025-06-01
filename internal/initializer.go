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
	"github.com/haxxu/friend-flow-api/internal/shared"
	minioinit "github.com/haxxu/friend-flow-api/pkg/minio"
	"github.com/minio/minio-go/v7"
)

type App struct {
	Router  *gin.Engine
	Session *gocql.Session
}

func InitializeApp(cfg *config.Config) (*App, error) {
	postgresDB, err := db.InitPostgres(cfg)
	if err != nil {
		return nil, err
	}

	EnforcerPG, err := rbac.InitEnforcerPostgres(postgresDB)
	if err != nil {
		return nil, err
	}

	// Minio Start ----------------------------------------------------------
	// Minio client initialization can be added here if needed
	minioClient, err := minioinit.InitMinio(cfg)
	if err != nil {
		return nil, err
	}
	if err := ensureBucket(minioClient, cfg); err != nil {
		return nil, err
	}

	deps := &shared.CoreDeps{
		Config:      cfg,
		PostgresDB:  postgresDB,
		MinioClient: minioClient,
		EnforcerPG:  EnforcerPG,
	}

	// Router
	router := setupRouter()

	// Initialize Modules and Routes
	setupModulesAndRoutes(router, deps)

	return &App{
		Router: router,
	}, nil
}

func ensureBucket(client *minio.Client, cfg *config.Config) error {
	ctx := context.Background()
	bucket := cfg.MinioBucket

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}

	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("✅ Bucket created: %s", bucket)
	} else {
		log.Printf("🪣 Bucket already exists: %s", bucket)
	}

	// Optional test upload
	objectName := "test-file.txt"
	content := []byte("This is just a test file 🧪")
	contentType := "text/plain"

	_, err = client.PutObject(ctx, bucket, objectName, bytes.NewReader(content), int64(len(content)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	log.Printf("✅ Successfully uploaded %s to bucket %s", objectName, bucket)
	return nil
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	return r
}

func setupModulesAndRoutes(r *gin.Engine, deps *shared.CoreDeps) {
	userModule := user.InitModule(deps.PostgresDB)
	authModule := auth.InitModule(deps.PostgresDB, deps.EnforcerPG, deps.Config, userModule.Service)
	chatServerModule := chat_server.InitModule(deps.PostgresDB)
	mediaUploadModule := media_upload.InitModule(deps.PostgresDB)

	routes.SetupRoutes(r, &routes.Handlers{
		AuthHandler:        authModule.Handler,
		UserHandler:        userModule.Handler,
		ChatServerHandler:  chatServerModule.Handler,
		MediaUploadHandler: mediaUploadModule.Handler,
	})
}

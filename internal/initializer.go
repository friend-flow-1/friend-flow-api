package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strconv"
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
	"github.com/haxxu/friend-flow-api/pkg/idgen"
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

	enforcerPG, err := rbac.InitEnforcerPostgres(postgresDB)
	if err != nil {
		return nil, err
	}

	// Minio client
	minioClient, err := minioinit.InitMinio(cfg)
	if err != nil {
		return nil, err
	}
	if err := ensureBucket(minioClient, cfg); err != nil {
		return nil, err
	}

	// Snowflake ID generator
	snowflakeNodeID, err := strconv.ParseInt(cfg.SnowflakeNodeID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid SNOWFLAKE_NODE_ID: %w", err)
	}
	snowflake, err := idgen.InitSnowflakeNode(snowflakeNodeID)
	if err != nil {
		return nil, err
	}

	deps := &shared.CoreDeps{
		Config:        cfg,
		PostgresDB:    postgresDB,
		MinioClient:   minioClient,
		EnforcerPG:    enforcerPG,
		SnowflakeNode: snowflake,
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

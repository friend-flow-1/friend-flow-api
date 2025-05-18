package internal

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
	chat_server "github.com/haxxu/friend-flow-api/internal/modules/chat/server"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"github.com/haxxu/friend-flow-api/internal/rbac"
	"github.com/haxxu/friend-flow-api/internal/routes"
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

	routes.SetupRoutes(router, &routes.Handlers{
		AuthHandler:       authModule.Handler,
		UserHandler:       userModule.Handler,
		ChatServerHandler: chatServerModule.Handler,
	})

	return &App{
		Router:     router,
		PostgresDB: postgresDB,
	}, nil
}

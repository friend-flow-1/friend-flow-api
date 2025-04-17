package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
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

	// Initialize User module
	userModule := user.InitModule(postgresDB)

	authModule := auth.InitModule(postgresDB, rbac.EnforcerPG, cfg.JWTSecret, userModule.Service)

	// Router
	router := gin.Default()
	routes.SetupRoutes(router, &routes.Handlers{
		AuthHandler: authModule.Handler,
		UserHandler: userModule.Handler,
	})

	return &App{
		Router:     router,
		PostgresDB: postgresDB,
	}, nil
}

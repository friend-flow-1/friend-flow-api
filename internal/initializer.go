package internal

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/modules/auth"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"github.com/haxxu/friend-flow-api/internal/rbac"
	"github.com/haxxu/friend-flow-api/internal/routes"
)

type App struct {
	Router  *mux.Router
	Session *gocql.Session
}

func InitializeApp(cfg *config.Config) (*App, error) {
	// Connect to ScyllaDB
	session, err := db.InitScylla(cfg)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(session)
	db.CreateIndexes(session)

	if err := rbac.InitEnforcer(session); err != nil {
		return nil, err
	}

	// Initialize User module
	userRepo := user.NewRepository(session)
	userService := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	// Initialize Auth module
	authService := auth.NewAuthService(userRepo, rbac.Enforcer, cfg.JWTSecret)
	authHandler := auth.NewAuthHandler(authService)

	// Router
	router := mux.NewRouter()
	routes.SetupRoutes(router, &routes.Handlers{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	})

	return &App{
		Router:  router,
		Session: session,
	}, nil
}

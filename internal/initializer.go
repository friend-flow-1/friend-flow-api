package internal

import (
	"log"
	"os"
	"path/filepath"

	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/haxxu/friend-flow-api/internal/auth"
	mycasbin "github.com/haxxu/friend-flow-api/internal/casbin"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/user"
)

func InitializeModules(cfg *config.Config) (*user.Handler, *auth.Handler, error) {
	// Connect to ScyllaDB
	session, err := db.InitScylla(cfg)
	if err != nil {
		return nil, nil, err // Return an error if DB initialization fails
	}
	db.AutoMigrate(session)

	// Initialize Casbin Enforcer
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current working directory: %v", err)
		return nil, nil, err
	}
	policyPath := filepath.Join(currentDir, "internal", "casbin", "policy.csv")
	adapter := fileadapter.NewAdapter(policyPath)
	if err := mycasbin.InitEnforcer(adapter); err != nil {
		log.Fatalf("Failed to initialize Casbin Enforcer: %v", err)
		return nil, nil, err
	}

	// Initialize User module
	userRepo := user.NewRepository(session)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	// Initialize Auth module
	authService := auth.NewAuthService(userRepo, mycasbin.Enforcer, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	db.CreateIndexes(session)

	return userHandler, authHandler, nil
}

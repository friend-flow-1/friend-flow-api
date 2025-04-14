package internal

import (
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/db"
	"github.com/haxxu/friend-flow-api/internal/user"
)

func InitializeModules(cfg *config.Config) (*user.Handler, error) {
	// Connect to ScyllaDB
	session, err := db.InitScylla(cfg)
	if err != nil {
		return nil, err // Return an error if DB initialization fails
	}
	db.AutoMigrate(session)

	// Initialize User module
	userRepo := user.NewRepository(session)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	return userHandler, nil
}

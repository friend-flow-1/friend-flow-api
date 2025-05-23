package chat_server

import (
	repositories "github.com/haxxu/friend-flow-api/internal/modules/chat/server/repositories"
	"gorm.io/gorm"
)

type Module struct {
	Repo    repositories.ServerRepository
	Service ServerService
	Handler *ServerHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := repositories.NewServerRepository(db)
	service := NewServerService(repo)
	handler := NewServerHandler(service)

	return &Module{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}

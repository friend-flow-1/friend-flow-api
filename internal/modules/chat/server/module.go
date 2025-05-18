package chat_server

import (
	chat_server_repositories "github.com/haxxu/friend-flow-api/internal/modules/chat/server/repositories"
	"gorm.io/gorm"
)

type Module struct {
	ServerRepo *chat_server_repositories.ServerRepository
	Service    *ServerService
	Handler    *ServerHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := chat_server_repositories.NewServerRepository(db)
	service := NewServerService(repo)
	handler := NewServerHandler(service)

	return &Module{
		ServerRepo: &repo,
		Service:    service,
		Handler:    handler,
	}
}

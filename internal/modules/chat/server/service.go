package chat_server

import (
	"errors"

	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/modules/chat/server/models"
	repositories "github.com/haxxu/friend-flow-api/internal/modules/chat/server/repositories"
)

type ServerService interface {
	CreateServer(server *models.Server) error
}

type serverService struct {
	repo repositories.ServerRepository
}

func NewServerService(repo repositories.ServerRepository) ServerService {
	return &serverService{repo: repo}
}

func (s *serverService) CreateServer(server *models.Server) error {
	if server.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}
	if server.Name == "" {
		return errors.New("server name is required")
	}
	return s.repo.Create(server)
}

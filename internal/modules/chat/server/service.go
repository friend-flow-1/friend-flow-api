package chat_server

import (
	"errors"

	"github.com/google/uuid"
	models "github.com/haxxu/friend-flow-api/internal/modules/chat/server/models"
	repositories "github.com/haxxu/friend-flow-api/internal/modules/chat/server/repositories"
)

type ServerService struct {
	ServerRepo repositories.ServerRepository
}

func NewServerService(serverRepo repositories.ServerRepository) *ServerService {
	return &ServerService{ServerRepo: serverRepo}
}

// CreateServer handles creating a new server.
func (s *ServerService) CreateServer(server *models.Server) error {
	if server.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}
	if server.Name == "" {
		return errors.New("server name is required")
	}
	return s.ServerRepo.Create(server)
}

// GetServerByID fetches a server by its ID.
func (s *ServerService) GetServerByID(serverID string) (*models.Server, error) {
	return s.ServerRepo.FindByID(serverID)
}

// GetServersByOwnerID returns all servers owned by a user.
func (s *ServerService) GetServersByOwnerID(ownerID string) ([]models.Server, error) {
	return s.ServerRepo.FindByOwnerID(ownerID)
}

// GetPublicServers returns all public servers.
func (s *ServerService) GetPublicServers() ([]models.Server, error) {
	return s.ServerRepo.FindPublicServers()
}

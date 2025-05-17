package chat_server

import (
	"errors"

	"github.com/google/uuid"
	chat_server "github.com/haxxu/friend-flow-api/internal/modules/server/models"
)

type ServerService struct {
	Repo Repository
}

func NewServerService(repo Repository) *ServerService {
	return &ServerService{Repo: repo}
}

// CreateServer handles creating a new server.
func (s *ServerService) CreateServer(server *chat_server.Server) error {
	if server.OwnerID == uuid.Nil {
		return errors.New("owner ID is required")
	}
	if server.Name == "" {
		return errors.New("server name is required")
	}
	return s.Repo.Create(server)
}

// GetServerByID fetches a server by its ID.
func (s *ServerService) GetServerByID(serverID string) (*chat_server.Server, error) {
	return s.Repo.FindByID(serverID)
}

// GetServersByOwnerID returns all servers owned by a user.
func (s *ServerService) GetServersByOwnerID(ownerID string) ([]chat_server.Server, error) {
	return s.Repo.FindByOwnerID(ownerID)
}

// GetPublicServers returns all public servers.
func (s *ServerService) GetPublicServers() ([]chat_server.Server, error) {
	return s.Repo.FindPublicServers()
}

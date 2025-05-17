package chat_server

import (
	chat_server "github.com/haxxu/friend-flow-api/internal/modules/server/models"

	"gorm.io/gorm"
)

type Repository interface {
	Create(server *chat_server.Server) error
	FindByID(serverID string) (*chat_server.Server, error)
	FindByOwnerID(ownerID string) ([]chat_server.Server, error)
	FindPublicServers() ([]chat_server.Server, error)
}

type repository struct {
	db *gorm.DB
}

// NewServerRepository initializes the server repository with the database connection.
func NewServerRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new server into the database.
func (r *repository) Create(server *chat_server.Server) error {
	return r.db.Create(server).Error
}

// FindByID fetches a server by its ID.
func (r *repository) FindByID(serverID string) (*chat_server.Server, error) {
	var server chat_server.Server
	if err := r.db.Preload("Members").
		Preload("Roles").
		Preload("Invites").
		Preload("Rules").
		First(&server, "id = ?", serverID).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

// FindByOwnerID fetches all servers owned by a specific user.
func (r *repository) FindByOwnerID(ownerID string) ([]chat_server.Server, error) {
	var servers []chat_server.Server
	if err := r.db.Where("owner_id = ?", ownerID).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

// FindPublicServers fetches all public servers.
func (r *repository) FindPublicServers() ([]chat_server.Server, error) {
	var servers []chat_server.Server
	if err := r.db.Where("is_private = false").Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

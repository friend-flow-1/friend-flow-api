package chat_server_repositories

import (
	models "github.com/haxxu/friend-flow-api/internal/modules/chat/server/models"

	"gorm.io/gorm"
)

type ServerRepository interface {
	Create(server *models.Server) error
	FindByID(serverID string) (*models.Server, error)
	FindByOwnerID(ownerID string) ([]models.Server, error)
	FindPublicServers() ([]models.Server, error)
}

type serverRepository struct {
	db *gorm.DB
}

// NewServerRepository initializes the server repository with the database connection.
func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{db: db}
}

// Create inserts a new server into the database.
func (r *serverRepository) Create(server *models.Server) error {
	return r.db.Create(server).Error
}

// FindByID fetches a server by its ID.
func (r *serverRepository) FindByID(serverID string) (*models.Server, error) {
	var server models.Server
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
func (r *serverRepository) FindByOwnerID(ownerID string) ([]models.Server, error) {
	var servers []models.Server
	if err := r.db.Where("owner_id = ?", ownerID).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

// FindPublicServers fetches all public servers.
func (r *serverRepository) FindPublicServers() ([]models.Server, error) {
	var servers []models.Server
	if err := r.db.Where("is_private = false").Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

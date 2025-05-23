package chat_server

import (
	"github.com/haxxu/friend-flow-api/internal/modules/chat/server/models"

	"gorm.io/gorm"
)

type ServerRepository interface {
	Create(server *models.Server) error
}

type serverRepository struct {
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{db: db}
}

func (r *serverRepository) Create(server *models.Server) error {
	return r.db.Create(server).Error
}

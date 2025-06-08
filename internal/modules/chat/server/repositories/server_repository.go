package chat_server

import (
	"github.com/haxxu/friend-flow-api/internal/modules/chat/server/models"
	"github.com/haxxu/friend-flow-api/pkg/repository"

	"gorm.io/gorm"
)

type ServerRepository interface {
	repository.BaseRepository[models.Server]
}

type serverRepository struct {
	repository.BaseRepository[models.Server]
	db *gorm.DB
}

func NewServerRepository(db *gorm.DB) ServerRepository {
	return &serverRepository{
		BaseRepository: repository.NewBaseRepository[models.Server](db),
		db:             db,
	}
}

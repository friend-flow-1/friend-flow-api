package chat_server

import (
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type ServerRole struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServerID    uuid.UUID `json:"server_id"`
	Name        string    `json:"name"`
	Permissions []string  `json:"permissions" gorm:"type:text[]"` // e.g. view_channel, manage_roles

	models.AuditFields `json:",inline"`
}

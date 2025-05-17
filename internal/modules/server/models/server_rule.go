package chat_server

import (
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type ServerRule struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServerID uuid.UUID `json:"server_id"`
	Content  string    `json:"content"`

	models.AuditFields `json:",inline"`
}

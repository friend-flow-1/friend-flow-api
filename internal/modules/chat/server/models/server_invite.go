package chat_server_models

import (
	"time"

	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type ServerInvite struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServerID  uuid.UUID  `json:"server_id"`
	Code      string     `json:"code" gorm:"uniqueIndex"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	models.AuditFields `json:",inline"`
}

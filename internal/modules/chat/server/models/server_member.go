package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
)

type ServerMember struct {
	ID       uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ServerID uuid.UUID  `json:"server_id"`
	UserID   uuid.UUID  `json:"user_id"` // FK to User
	RoleID   *uuid.UUID `json:"role_id,omitempty"`

	JoinedAt time.Time `json:"joined_at"`

	models.AuditFields `json:",inline"`
}

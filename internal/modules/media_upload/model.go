package media_upload

import (
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
)

type MediaUpload struct {
	ID                 uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Type               string    `json:"type"`
	OriginalName       string    `json:"original_name"`
	Path               string    `json:"path"`
	URL                string    `json:"url"`
	Size               int64     `json:"size"`
	Verified           bool      `json:"verified"`
	OwnerID            uuid.UUID `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner              user.User `gorm:"foreignKey:OwnerID;references:ID" json:"-"`
	models.AuditFields `json:",inline"`
}

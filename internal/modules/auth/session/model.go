package authsession

import (
	"time"

	"github.com/google/uuid"
)

type AuthSession struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	RefreshToken string    `gorm:"unique;not null"`
	UserAgent    string
	IPAddress    string
	ExpiresAt    time.Time
	Revoked      bool `gorm:"default:false"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

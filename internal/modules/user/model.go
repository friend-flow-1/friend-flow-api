package user

import (
	"time"

	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
	"gorm.io/gorm"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"
)

type User struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Phone      string    `json:"phone,omitempty"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Password   string    `json:"-"` // hide from JSON output
	Status     string    `json:"status" gorm:"default:'active'"`
	Background string    `json:"background"`
	Avatar     string    `json:"avatar,omitempty"`
	BirthDate  time.Time `json:"birth_date,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Role       Role      `json:"role" gorm:"type:varchar(20);not null"`

	// Embed the AuditFields struct to include the common audit fields
	models.AuditFields `json:",inline"`
}

// To be used with GORM for automatic timestamp handling.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	u.CreatedBy = "system" // This should be dynamically set based on your application
	u.UpdatedBy = "system" // Same as above
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	u.UpdatedBy = "system" // Set it to the current user or system dynamically
	return
}

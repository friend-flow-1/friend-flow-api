package user

import (
	"time"

	"github.com/haxxu/friend-flow-api/internal/models"
)

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleUser       Role = "user"
)

type User struct {
	ID         string    `json:"id"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone,omitempty"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Password   string    `json:"_"`      // hide from JSON output
	Status     string    `json:"status"` // "active", "suspended", "locked", ...
	Background string    `json:"background,omitempty"`
	Avatar     string    `json:"avatar,omitempty"`
	BirthDate  time.Time `json:"birth_date,omitempty"`
	Gender     string    `json:"gender,omitempty"`
	Role       Role      `json:"role"`

	models.AuditFields
}

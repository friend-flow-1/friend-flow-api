package auth

import "github.com/go-playground/validator/v10"

// RegisterDTO represents the payload for user registration
type RegisterDTO struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"required"`
	// BirthDate time.Time `json:"birth_date"`
}

// LoginDTO represents the payload for user login
type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Validate validates the fields of RegisterDTO
func (r *RegisterDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Validate validates the fields of LoginDTO
func (l *LoginDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

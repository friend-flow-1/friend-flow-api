package auth

import "github.com/go-playground/validator/v10"

type RegisterDTO struct {
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Gender    string `json:"gender" validate:"required"`
	// BirthDate time.Time `json:"birth_date"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenDTO struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
	UserAgent    string `json:"userAgent"`
	IP           string `json:"ip"`
}

func (r *RegisterDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func (l *LoginDTO) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}

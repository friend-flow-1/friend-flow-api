package user

import "time"

type User struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone,omitempty"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	FullName   string     `json:"full_name"`
	Password   string     `json:"-"`
	Status     string     `json:"status"`
	Background string     `json:"background,omitempty"`
	Avatar     string     `json:"avatar,omitempty"`
	BirthDate  time.Time  `json:"birth_date"`
	Gender     string     `json:"gender,omitempty"`
	CreatedAt  time.Time  `json:"created_at" gocql:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gocql:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty" gocql:"deleted_at"`
}

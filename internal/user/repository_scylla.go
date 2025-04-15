package user

import (
	"github.com/gocql/gocql"
)

type repository struct {
	Session *gocql.Session
}

func NewRepository(session *gocql.Session) Repository {
	return &repository{Session: session}
}

func (r *repository) Create(user *User) error {
	return r.Session.Query(`
		INSERT INTO users (
			id, email, phone, first_name, last_name, 
			password, status, background, avatar, birth_date, gender, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.Phone, user.FirstName, user.LastName,
		user.Password, user.Status, user.Background, user.Avatar, user.BirthDate, user.Gender, user.CreatedAt, user.UpdatedAt,
	).Exec()
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var u User
	query := `
		SELECT id, email, first_name, last_name, 
			password, status, gender, created_at, updated_at FROM users 
		WHERE email = ? LIMIT 1
	`
	if err := r.Session.Query(query, email).Consistency(gocql.One).Scan(
		&u.ID, &u.Email, &u.FirstName, &u.LastName,
		&u.Password, &u.Status, &u.Gender, &u.CreatedAt, &u.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

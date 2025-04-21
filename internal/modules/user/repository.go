package user

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindById(userId string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

// NewUserRepository initializes the user repository with the database connection.
func NewUserRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// Create implements Repository.
func (r *repository) Create(user *User) error {
	return r.db.Create(user).Error
}

// FindByEmail implements Repository.
func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindById(userId string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

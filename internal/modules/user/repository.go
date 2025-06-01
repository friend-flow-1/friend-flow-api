package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindById(userId string) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository initializes the user repository with the database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create implements Repository.
func (r *userRepository) Create(user *User) error {
	return r.db.Create(user).Error
}

// FindByEmail implements Repository.
func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindById(userId string) (*User, error) {
	var user User
	if err := r.db.First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

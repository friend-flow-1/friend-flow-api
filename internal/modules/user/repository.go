package user

import (
	"github.com/haxxu/friend-flow-api/pkg/repository"
	"gorm.io/gorm"
)

type UserRepository interface {
	repository.BaseRepository[User]
	FindByEmail(email string) (*User, error)
}

type userRepository struct {
	db *gorm.DB
	repository.BaseRepository[User]
}

// NewUserRepository initializes the user repository with the database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: repository.NewBaseRepository[User](db),
		db:             db,
	}
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

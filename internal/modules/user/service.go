package user

import "errors"

type UserService interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	GetUserById(userId string) (*User, error)
}

type userService struct {
	Repo Repository
}

func NewUserService(repo Repository) UserService {
	return &userService{Repo: repo}
}

// CreateUser handles creating a new user.
func (s *userService) CreateUser(user *User) error {
	// Check if email already exists
	existing, _ := s.Repo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Create the user in the repository
	return s.Repo.Create(user)
}

// FindByEmail fetches a user by email.
func (s *userService) FindByEmail(email string) (*User, error) {
	return s.Repo.FindByEmail(email)
}

func (s *userService) GetUserById(userId string) (*User, error) {
	user, err := s.Repo.FindById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

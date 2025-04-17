package user

import "errors"

type UserService struct {
	Repo Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{Repo: repo}
}

// CreateUser handles creating a new user.
func (s *UserService) CreateUser(user *User) error {
	// Check if email already exists
	existing, _ := s.Repo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	// Create the user in the repository
	return s.Repo.Create(user)
}

// FindByEmail fetches a user by email.
func (s *UserService) FindByEmail(email string) (*User, error) {
	return s.Repo.FindByEmail(email)
}

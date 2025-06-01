package user

import "errors"

type UserService interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	GetUserById(userId string) (*User, error)
}

type userService struct {
	UserRepo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{UserRepo: repo}
}

func (s *userService) CreateUser(user *User) error {
	// Check if email already exists
	existing, _ := s.UserRepo.FindByEmail(user.Email)
	if existing != nil {
		return errors.New("email already registered")
	}

	return s.UserRepo.Create(user)
}

func (s *userService) FindByEmail(email string) (*User, error) {
	return s.UserRepo.FindByEmail(email)
}

func (s *userService) GetUserById(userId string) (*User, error) {
	user, err := s.UserRepo.FindById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

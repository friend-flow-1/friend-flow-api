package user

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) RegisterUser(req RegisterRequest) (*User, error) {
	existing, _ := s.Repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	user := &User{
		ID:        uuid.New().String(),
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		FullName:  req.FirstName + " " + req.LastName,
		Status:    "active",
		Gender:    req.Gender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}
	user.Password = "" // Hide password
	return user, nil
}

func (s *Service) LoginUser(req LoginRequest) (*User, error) {
	user, err := s.Repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}
	user.Password = ""
	return user, nil
}

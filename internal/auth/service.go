package auth

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/haxxu/friend-flow-api/internal/models"
	"github.com/haxxu/friend-flow-api/internal/user"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo  user.Repository
	Enforcer  *casbin.Enforcer
	JWTSecret string
}

func NewAuthService(repo user.Repository, enforcer *casbin.Enforcer, secret string) *AuthService {
	return &AuthService{
		UserRepo:  repo,
		Enforcer:  enforcer,
		JWTSecret: secret,
	}
}

type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	BirthDate time.Time
	Gender    string
}

type LoginRequest struct {
	Email    string
	Password string
}

func (s *AuthService) Register(req RegisterRequest) (*user.User, error) {
	existing, _ := s.UserRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}

	newId := uuid.New().String()
	u := &user.User{
		ID:        newId,
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    "active",
		BirthDate: req.BirthDate,
		Gender:    req.Gender,
		Role:      user.RoleUser,
		AuditFields: models.AuditFields{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			CreatedBy: newId,
			UpdatedBy: newId,
		},
	}

	if err := s.UserRepo.Create(u); err != nil {
		return nil, err
	}

	// Assign casbin role
	s.Enforcer.AddRoleForUser(u.ID, string(u.Role))

	u.Password = ""
	return u, nil
}

func (s *AuthService) Login(req LoginRequest) (string, *user.User, error) {
	u, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println(err)
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		log.Println(err)
		return "", nil, errors.New("invalid credentials")
	}

	token, err := s.generateJWT(u)
	if err != nil {
		return "", nil, err
	}

	u.Password = ""
	return token, u, nil
}

func (s *AuthService) generateJWT(u *user.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":  u.ID,
		"role": u.Role,
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return t.SignedString([]byte(s.JWTSecret))
}

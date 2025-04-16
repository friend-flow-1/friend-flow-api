package auth

import (
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/haxxu/friend-flow-api/internal/models"
	"github.com/haxxu/friend-flow-api/internal/modules/user"

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

func (s *AuthService) Register(req RegisterDTO) (*user.User, error) {
	existing, _ := s.UserRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}
	log.Println("hello", string(hashed))

	newId := uuid.New().String()
	u := &user.User{
		ID:        newId,
		Email:     req.Email,
		Password:  string(hashed),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Status:    "active",
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
	// s.Enforcer.AddRoleForUser(u.ID, string(u.Role))

	u.Password = ""
	return u, nil
}

func (s *AuthService) Login(req LoginDTO) (string, *user.User, error) {
	u, err := s.UserRepo.FindByEmail(req.Email)
	if err != nil {
		log.Println("Error finding user:", err)
		return "", nil, errors.New("invalid credentials")
	}

	// Log request and user password for debugging (ensure no sensitive info is exposed)
	log.Println("Request Email:", req.Email)
	log.Println("User Stored Password:", u.Password)
	log.Println("Request Password:", req.Password)

	// Compare password hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		log.Println("Error comparing passwords:", err)
		return "", nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(u)
	if err != nil {
		log.Println("Error generating JWT token:", err)
		return "", nil, err
	}

	// Log generated token for debugging (ensure sensitive data is logged properly)
	log.Println("Generated Token:", token)

	// Clear sensitive information before returning user
	u.Password = ""
	return token, u, nil
}

func (s *AuthService) generateJWT(u *user.User) (string, error) {
	// Ensure JWTSecret is valid and logged for debugging
	if s.JWTSecret == "" {
		log.Println("JWT Secret is empty!")
		return "", errors.New("missing JWT secret")
	}
	log.Println("Using JWT Secret for signing")

	// Set JWT claims
	claims := jwt.MapClaims{
		"sub":  u.ID,                                  // Subject (user ID)
		"role": u.Role,                                // User role
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // Expiry time (3 days)
	}

	// Use HS256 method for signing (using a shared secret)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and return the JWT token
	token, err := t.SignedString([]byte(s.JWTSecret))
	if err != nil {
		log.Println("Error signing JWT token:", err)
		return "", err
	}

	return token, nil
}

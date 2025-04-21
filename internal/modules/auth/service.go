package auth

import (
	"errors"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/models"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService *user.UserService
	Enforcer    *casbin.Enforcer
	JWTSecret   string
}

func NewAuthService(userService *user.UserService, enforcer *casbin.Enforcer, secret string) *AuthService {
	return &AuthService{
		UserService: userService,
		Enforcer:    enforcer,
		JWTSecret:   secret,
	}
}

func (s *AuthService) Register(req RegisterDTO) (*user.User, error) {
	// Check if email already exists
	existing, _ := s.UserService.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash the password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, err
	}
	log.Println("hashed password:", string(hashed))

	// Generate new UUID for user ID
	newId := uuid.New()

	// Create a user DTO
	u := &user.User{
		ID:        newId, // Set the ID as uuid.UUID
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
			CreatedBy: newId.String(), // Here you can set newId.String() if you want to store it as string
			UpdatedBy: newId.String(), // Same for UpdatedBy
		},
	}

	// Call the Create method from UserService to handle user creation
	if err := s.UserService.CreateUser(u); err != nil {
		return nil, err
	}

	// Assign casbin role
	s.Enforcer.AddRoleForUser(u.ID.String(), string(u.Role))
	_, _ = s.Enforcer.AddGroupingPolicy(u.ID.String(), string(u.Role))

	// Clear the password field before returning user object
	u.Password = ""

	return u, nil
}

// Login authenticates a user and generates a JWT token
func (s *AuthService) Login(req LoginDTO) (string, *user.User, error) {
	// Find user by email
	u, err := s.UserService.FindByEmail(req.Email)
	if err != nil {
		log.Println("Error finding user:", err)
		return "", nil, errors.New("invalid credentials")
	}

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

	// Clear the password field before returning user
	u.Password = ""
	return token, u, nil
}

// generateJWT generates a JWT token for the user
func (s *AuthService) generateJWT(u *user.User) (string, error) {
	// Ensure JWTSecret is valid
	if s.JWTSecret == "" {
		log.Println("JWT Secret is empty!")
		return "", errors.New("missing JWT secret")
	}

	claims := models.JWTClaims{
		Sub: u.ID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(72 * time.Hour).Unix(),
		},
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

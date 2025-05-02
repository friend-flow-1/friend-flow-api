package auth

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/haxxu/friend-flow-api/internal/config"
	"github.com/haxxu/friend-flow-api/internal/models"
	authsession "github.com/haxxu/friend-flow-api/internal/modules/auth/session"
	authtoken "github.com/haxxu/friend-flow-api/internal/modules/auth/token"
	"github.com/haxxu/friend-flow-api/internal/modules/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService        *user.UserService
	AuthSessionService *authsession.AuthSessionService
	TokenService       *authtoken.TokenService
	Enforcer           *casbin.Enforcer
	Config             *config.Config
}

func NewAuthService(userService *user.UserService, authSessionService *authsession.AuthSessionService, tokenService *authtoken.TokenService, enforcer *casbin.Enforcer, cfg *config.Config) *AuthService {
	return &AuthService{
		UserService:        userService,
		AuthSessionService: authSessionService,
		TokenService:       tokenService,
		Enforcer:           enforcer,
		Config:             cfg,
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

func (s *AuthService) Login(dto LoginDTO, ua string, ip string) (string, string, *user.User, error) {
	// Validate user credentials (pseudo-code)
	user, err := s.validateUser(dto)
	if err != nil {
		return "", "", nil, err
	}

	// Generate tokens
	accessToken, _ := s.TokenService.GenerateAccessToken(user.ID)
	refreshToken, _ := s.TokenService.GenerateRefreshToken(user.ID)

	// Save session
	err = s.AuthSessionService.Create(user.ID, refreshToken, ua, ip)
	if err != nil {
		return "", "", nil, err
	}

	return accessToken, refreshToken, user, nil
}

func (s *AuthService) validateUser(dto LoginDTO) (*user.User, error) {
	user, err := s.UserService.FindByEmail(dto.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := comparePasswords(user.Password, dto.Password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func comparePasswords(hashedPwd string, plainPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
}

// generateJWT generates a JWT token for the user
func (s *AuthService) generateJWT(u *user.User) (string, error) {
	// Ensure JWTSecret is valid
	if s.Config.JWTSecret == "" {
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
	token, err := t.SignedString([]byte(s.Config.JWTSecret))
	if err != nil {
		log.Println("Error signing JWT token:", err)
		return "", err
	}

	return token, nil
}

func (s *AuthService) Logout(token string) error {
	return s.AuthSessionService.Revoke(token)
}

func (s *AuthService) LogoutAll(userID uuid.UUID) error {
	return s.AuthSessionService.RevokeAll(userID)
}

func (s *AuthService) GetActiveSessions(userID uuid.UUID) ([]authsession.AuthSession, error) {
	return s.AuthSessionService.List(userID)
}

func (s *AuthService) ValidateRefreshToken(refreshToken string) (uuid.UUID, error) {
	session, err := s.AuthSessionService.FindByToken(refreshToken)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session not found: %w", err)
	}

	if session.Revoked || session.ExpiresAt.Before(time.Now()) {
		return uuid.Nil, errors.New("refresh token expired or revoked")
	}

	return session.UserID, nil
}

package authsession

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AuthSessionService struct {
	repo AuthSessionRepository
}

func NewAuthSessionService(repo AuthSessionRepository) *AuthSessionService {
	return &AuthSessionService{repo: repo}
}

func (s *AuthSessionService) Create(userID uuid.UUID, token string, userAgent string, ip string) error {
	session := &AuthSession{
		UserID:       userID,
		RefreshToken: token,
		UserAgent:    userAgent,
		IPAddress:    ip,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
	}
	return s.repo.Create(session)
}

func (s *AuthSessionService) Validate(token string) (*AuthSession, error) {
	return s.repo.FindByToken(token)
}

func (s *AuthSessionService) Rotate(sessionID uuid.UUID, oldToken string, newToken string, userAgent string, ip string) error {
	err := s.repo.Revoke(oldToken)
	if err != nil {
		return err
	}
	return s.Create(sessionID, newToken, userAgent, ip)
}

func (s *AuthSessionService) Revoke(token string) error {
	return s.repo.Revoke(token)
}

func (s *AuthSessionService) RevokeAll(userID uuid.UUID) error {
	return s.repo.RevokeAll(userID)
}

func (s *AuthSessionService) List(userID uuid.UUID) ([]AuthSession, error) {
	return s.repo.FindByUserID(userID)
}

func (s *AuthSessionService) FindByToken(token string) (*AuthSession, error) {
	return s.repo.FindByToken(token)
}

func (s *AuthSessionService) Update(userID uuid.UUID, oldToken, newToken, userAgent, ip string) error {
	session, err := s.repo.FindByToken(oldToken)
	if err != nil {
		return fmt.Errorf("failed to find session: %w", err)
	}

	if session.UserID != userID {
		return errors.New("user mismatch")
	}

	session.RefreshToken = newToken
	session.UpdatedAt = time.Now()
	session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour)

	session.UserAgent = userAgent
	session.IPAddress = ip

	return s.repo.Update(session)
}

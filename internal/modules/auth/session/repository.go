package authsession

import "github.com/google/uuid"

type AuthSessionRepository interface {
	Create(session *AuthSession) error
	FindByToken(token string) (*AuthSession, error)
	FindByUserID(userID uuid.UUID) ([]AuthSession, error)
	Update(session *AuthSession) error
	Delete(session *AuthSession) error
	Revoke(token string) error
	RevokeAll(userID uuid.UUID) error
}

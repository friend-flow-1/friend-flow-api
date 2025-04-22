package authsession

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormAuthSessionRepository struct {
	db *gorm.DB
}

func NewGormAuthSessionRepository(db *gorm.DB) *GormAuthSessionRepository {
	return &GormAuthSessionRepository{
		db: db,
	}
}

func (r *GormAuthSessionRepository) Create(session *AuthSession) error {
	return r.db.Create(session).Error
}

func (r *GormAuthSessionRepository) FindByToken(token string) (*AuthSession, error) {
	var session AuthSession
	if err := r.db.Where("refresh_token = ? AND revoked = false AND expires_at > ?", token, time.Now()).First(&session).Error; err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}
	return &session, nil
}

func (r *GormAuthSessionRepository) FindByUserID(userID uuid.UUID) ([]AuthSession, error) {
	var sessions []AuthSession
	err := r.db.Where("user_id = ? AND revoked = false", userID).Find(&sessions).Error
	return sessions, err
}

func (r *GormAuthSessionRepository) Update(session *AuthSession) error {
	return r.db.Save(session).Error
}

func (r *GormAuthSessionRepository) Delete(session *AuthSession) error {
	return r.db.Delete(session).Error
}

func (r *GormAuthSessionRepository) Revoke(token string) error {
	return r.db.Model(&AuthSession{}).Where("refresh_token = ?", token).Update("revoked", true).Error
}

func (r *GormAuthSessionRepository) RevokeAll(userID uuid.UUID) error {
	return r.db.Model(&AuthSession{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

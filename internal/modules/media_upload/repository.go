package media_upload

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaUploadRepository interface {
	Create(mediaUpload *MediaUpload) error
	GetByID(id uuid.UUID) (*MediaUpload, error)
	GetByOwnerID(ownerID uuid.UUID) ([]MediaUpload, error)
	Update(mediaUpload *MediaUpload) error
	Delete(id uuid.UUID) error
	// ListByOwnerID(ownerID string, limit, offset int) ([]MediaUpload, error)
	// ListByTypeAndOwnerID(mediaType, ownerID string, limit, offset int) ([]MediaUpload, error)
	// ListByType(mediaType string, limit, offset int) ([]MediaUpload, error)
	// ListAll(limit, offset int) ([]MediaUpload, error)
	// CountByOwnerID(ownerID string) (int64, error)
	// CountByTypeAndOwnerID(mediaType, ownerID string) (int64, error)
	// CountByType(mediaType string) (int64, error)
	// CountAll() (int64, error)
}
type mediaUploadRepository struct {
	db *gorm.DB
}

func NewMediaUploadRepository(db *gorm.DB) MediaUploadRepository {
	return &mediaUploadRepository{db: db}
}
func (r *mediaUploadRepository) Create(mediaUpload *MediaUpload) error {
	return r.db.Create(mediaUpload).Error
}

func (r *mediaUploadRepository) GetByID(id uuid.UUID) (*MediaUpload, error) {
	var mediaUpload MediaUpload
	err := r.db.First(&mediaUpload, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &mediaUpload, nil
}
func (r *mediaUploadRepository) GetByOwnerID(ownerID uuid.UUID) ([]MediaUpload, error) {
	var mediaUploads []MediaUpload
	err := r.db.Where("owner_id = ?", ownerID).Find(&mediaUploads).Error
	if err != nil {
		return nil, err
	}
	return mediaUploads, nil
}
func (r *mediaUploadRepository) Update(mediaUpload *MediaUpload) error {
	return r.db.Save(mediaUpload).Error
}

func (r *mediaUploadRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&MediaUpload{}, "id = ?", id).Error
}

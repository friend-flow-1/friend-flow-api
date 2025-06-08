package media_upload

import (
	"github.com/haxxu/friend-flow-api/pkg/repository"
	"gorm.io/gorm"
)

type MediaUploadRepository interface {
	repository.BaseRepository[MediaUpload]
	GetByOwnerID(ownerID string) ([]MediaUpload, error)
	Update(mediaUpload *MediaUpload) error
	Delete(id string) error
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
	repository.BaseRepository[MediaUpload]
	db *gorm.DB
}

func NewMediaUploadRepository(db *gorm.DB) MediaUploadRepository {
	return &mediaUploadRepository{
		BaseRepository: repository.NewBaseRepository[MediaUpload](db),
		db:             db,
	}
}

func (r *mediaUploadRepository) GetByOwnerID(ownerID string) ([]MediaUpload, error) {
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

func (r *mediaUploadRepository) Delete(id string) error {
	return r.db.Delete(&MediaUpload{}, "id = ?", id).Error
}

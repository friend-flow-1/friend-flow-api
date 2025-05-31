package media_upload

import (
	"gorm.io/gorm"
)

type Module struct {
	Repo    MediaUploadRepository
	Service MediaUploadService
	Handler *MediaUploadHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewMediaUploadRepository(db)
	service := NewMediaUploadService(repo)
	handler := NewMediaUploadHandler(service)

	return &Module{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}

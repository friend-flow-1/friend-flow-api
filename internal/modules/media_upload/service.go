package media_upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type MediaUploadService interface {
	Create(mediaUpload *MediaUpload) error
	GetByID(id uuid.UUID) (*MediaUpload, error)
	GetByOwnerID(ownerID uuid.UUID) ([]MediaUpload, error)
	Update(mediaUpload *MediaUpload) error
	Delete(id uuid.UUID) error
	UploadFile(c *gin.Context)
	DeleteFile(c *gin.Context)
}
type mediaUploadService struct {
	repo        MediaUploadRepository
	minioClient *minio.Client
	ctx         context.Context
	bucketName  string
}

func NewMediaUploadService(repo MediaUploadRepository) MediaUploadService {
	return &mediaUploadService{repo: repo}
}
func (s *mediaUploadService) Create(mediaUpload *MediaUpload) error {
	return s.repo.Create(mediaUpload)
}

func (s *mediaUploadService) GetByID(id uuid.UUID) (*MediaUpload, error) {
	return s.repo.GetByID(id)
}
func (s *mediaUploadService) GetByOwnerID(ownerID uuid.UUID) ([]MediaUpload, error) {

	return s.repo.GetByOwnerID(ownerID)
}
func (s *mediaUploadService) Update(mediaUpload *MediaUpload) error {
	return s.repo.Update(mediaUpload)
}
func (s *mediaUploadService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *mediaUploadService) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	userID := c.GetString("user_id")
	ownerID, err := uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	mediaID := uuid.New()
	filePath := fmt.Sprintf("uploads/%s_%s", mediaID.String(), header.Filename)

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	_, err = s.minioClient.PutObject(s.ctx, s.bucketName, filePath, bytes.NewReader(buf.Bytes()), int64(buf.Len()), minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	url := fmt.Sprintf("http://localhost:9000/%s/%s", s.bucketName, filePath)
	media := &MediaUpload{
		ID:           mediaID,
		Type:         "image",
		OriginalName: header.Filename,
		Path:         filePath,
		URL:          url,
		Size:         header.Size,
		Verified:     false,
		OwnerID:      ownerID,
	}
	if err := s.repo.Create(media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database insert failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"media": media})
}

func (s *mediaUploadService) DeleteFile(c *gin.Context) {
	id := c.Query("id")
	mediaID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	media, err := s.repo.GetByID(mediaID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}
	if media.Verified {
		c.JSON(http.StatusOK, gin.H{"message": "File is verified and won't be deleted"})
		return
	}

	err = s.minioClient.RemoveObject(s.ctx, s.bucketName, media.Path, minio.RemoveObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}
	if err := s.repo.Delete(mediaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media from DB"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

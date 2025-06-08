package media_upload

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MediaUploadHandler struct {
	service MediaUploadService
}

func NewMediaUploadHandler(service MediaUploadService) *MediaUploadHandler {
	return &MediaUploadHandler{
		service: service,
	}
}

func (h *MediaUploadHandler) UploadFile(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	userID := c.GetString("user_id")
	_, err = uuid.Parse(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	h.service.UploadFile(c)
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func (h *MediaUploadHandler) DeleteFile(c *gin.Context) {
	id := c.Query("id")
	mediaID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	err = h.service.Delete(mediaID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

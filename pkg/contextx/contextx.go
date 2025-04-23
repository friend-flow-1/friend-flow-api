package contextx

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	KeyUserID = "user_id"
)

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userID := c.GetString(KeyUserID)
	if userID == "" {
		return uuid.Nil, errors.New("user ID not found in context")
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID format")
	}

	return uid, nil
}

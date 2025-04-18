package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditFields struct {
	CreatedAt time.Time      `json:"created_at"`
	CreatedBy string         `json:"created_by"`
	UpdatedAt time.Time      `json:"updated_at"`
	UpdatedBy string         `json:"updated_by"`
	DeletedAt gorm.DeletedAt `json:"_" gorm:"index"`
	DeletedBy string         `json:"_"`
}

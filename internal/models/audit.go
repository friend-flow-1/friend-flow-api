package models

import "time"

type AuditFields struct {
	CreatedAt time.Time  `json:"created_at" gocql:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" gocql:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gocql:"deleted_at"`
	CreatedBy string     `json:"created_by,omitempty" gocql:"created_by"`
	UpdatedBy string     `json:"updated_by,omitempty" gocql:"updated_by"`
	DeletedBy string     `json:"deleted_by,omitempty" gocql:"deleted_by"`
}

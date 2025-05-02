package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel defines the common fields for GORM models with JSON tags.
// It includes ID, CreatedAt, UpdatedAt, and DeletedAt.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"` // Use omitempty for null values
}

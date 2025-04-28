// Package models contains database models
package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminUser represents an administrator account in the database
type AdminUser struct {
	gorm.Model            // Includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
	Username   string     `gorm:"uniqueIndex;not null" json:"username"` // Unique username
	Password   string     `gorm:"not null" json:"-"`             // Hashed password - Exclude from JSON
	LastLogin  *time.Time `json:"lastLogin,omitempty"` // Pointer allows null value for last login timestamp
}

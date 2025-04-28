package models

import (
	"gorm.io/gorm"
)

// QuestionType defines the allowed types for questions (single or multiple choice)
type QuestionType string

const (
	SingleChoice QuestionType = "single"
	MultiChoice  QuestionType = "multi"
)

// Quiz represents a single quiz created by an admin
type Quiz struct {
	gorm.Model
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description,omitempty"`
	TimeLimit   *uint      `json:"timeLimit,omitempty"`
	Status      string     `gorm:"not null;default:'Draft'" json:"status"`
	AdminUserID uint       `gorm:"not null" json:"adminUserId"`
	AdminUser   AdminUser  `json:"-"` 
	Questions   []Question `json:"questions,omitempty"`
}

// Question represents a single question within a quiz
type Question struct {
	gorm.Model
	Text    string       `gorm:"not null"`
	Type    QuestionType `gorm:"type:varchar(10);not null"` // 'single' or 'multi'
	QuizID  uint         `gorm:"not null"`                  // Foreign key to Quiz
	Choices []Choice     // Has many relationship
}

// Choice represents a possible answer choice for a question
type Choice struct {
	gorm.Model
	Text       string `gorm:"not null"`
	IsCorrect  bool   `gorm:"not null;default:false"`
	QuestionID uint   `gorm:"not null"` // Foreign key to Question
}

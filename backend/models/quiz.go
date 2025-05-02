package models

import (
)

// QuestionType defines the allowed types for questions (single or multiple choice)
type QuestionType string

const (
	SingleChoice QuestionType = "single"
	MultiChoice  QuestionType = "multi"
)

// Quiz represents a single quiz created by an admin
type Quiz struct {
	BaseModel   // Embed our custom base model
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
	BaseModel        // Embed our custom base model
	Text        string       `gorm:"not null" json:"text"`
	Type        QuestionType `gorm:"type:varchar(10);not null" json:"type"` // 'single' or 'multi'
	QuizID      uint         `gorm:"not null" json:"quizId"` // Foreign key to Quiz
	Choices     []Choice     `json:"choices,omitempty"`       // Has many relationship
}

// Choice represents a possible answer choice for a question
type Choice struct {
	BaseModel          // Embed our custom base model
	Text       string `gorm:"not null" json:"text"`
	IsCorrect  bool   `gorm:"not null;default:false" json:"isCorrect"`
	QuestionID uint   `gorm:"not null" json:"questionId"` // Foreign key to Question
}

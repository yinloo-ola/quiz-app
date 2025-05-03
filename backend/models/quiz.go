package models

import (
	"errors"
	"fmt"
	"strings"
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

// --- Helper Functions and Validation ---

// ParseQuestionType converts a string to a QuestionType enum.
func ParseQuestionType(typeStr string) (QuestionType, error) {
	switch strings.ToLower(typeStr) {
	case "single":
		return SingleChoice, nil
	case "multi":
		return MultiChoice, nil
	default:
		return "", fmt.Errorf("invalid question type string: '%s'", typeStr)
	}
}

// CorrectableChoice defines an interface for choice-like objects
// that can report whether they are marked as correct.
type CorrectableChoice interface {
	GetIsCorrect() bool
}

// GetIsCorrect implements the CorrectableChoice interface for the Choice model.
func (c Choice) GetIsCorrect() bool {
	return c.IsCorrect
}

// ValidateChoicesForQuestionType checks if the correctness pattern of choices
// is valid for the given question type.
func ValidateChoicesForQuestionType(qType QuestionType, choices []CorrectableChoice) error {
	if len(choices) < 2 {
		// Although the form prevents < 2 choices, API could technically receive it.
		// Depending on requirements, this might be allowed or disallowed.
		// For now, we assume the frontend validation enforces >= 2 choices.
		// return errors.New("a question must have at least two choices")
	}

	correctCount := 0
	for _, choice := range choices {
		if choice.GetIsCorrect() {
			correctCount++
		}
	}

	switch qType {
	case SingleChoice:
		if correctCount == 0 {
			return errors.New("a single-choice question must have exactly one correct answer, but none were marked correct")
		} else if correctCount > 1 {
			return errors.New("a single-choice question must have exactly one correct answer, but multiple were marked correct")
		}
	case MultiChoice:
		if correctCount == 0 {
			return errors.New("a multiple-choice question must have at least one correct answer, but none were marked correct")
		}
	default:
		// This case should ideally not be reached if ParseQuestionType is used correctly.
		return fmt.Errorf("unknown question type '%s' encountered during validation", qType)
	}

	return nil // Validation passed
}

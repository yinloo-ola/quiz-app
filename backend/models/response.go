package models

import (
	"time"

	"gorm.io/gorm"
)

// ResponderCredential represents a temporary credential for a user to take a specific quiz
type ResponderCredential struct {
	gorm.Model
	QuizID       uint       `gorm:"not null" json:"quizId"`
	Quiz         Quiz       `json:"quiz,omitempty"` // Belongs to Quiz
	Username     string     `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"not null;default:''" json:"-"` // Hashed password - Exclude!
	ExpiresAt    *time.Time `gorm:"not null" json:"expiresAt,omitempty"`
	Used         bool       `gorm:"default:false" json:"used"` // Has the credential been used?
	UsedAt       *time.Time `json:"usedAt,omitempty"`
	// CreatedAt is part of gorm.Model, handled automatically
}

// QuizResponse represents a single completed attempt of a quiz by a responder.
type QuizResponse struct {
	gorm.Model
	QuizID              uint                `gorm:"not null" json:"quizId"`
	Quiz                Quiz                `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	ResponderCredentialID *uint               `gorm:"uniqueIndex" json:"responderCredentialId,omitempty"`
	ResponderCredential *ResponderCredential `gorm:"foreignKey:ResponderCredentialID" json:"responderCredential,omitempty"`
	ResponderUsername   string              `gorm:"index" json:"responderUsername"`
	Score               *float64            `json:"score,omitempty"`
	SubmittedAt         time.Time           `gorm:"not null;default:CURRENT_TIMESTAMP" json:"submittedAt"`
	Answers             []Answer            `gorm:"foreignKey:QuizResponseID" json:"answers,omitempty"`
}

// Answer represents a single answer submitted by a responder for a specific question within a quiz response.
type Answer struct {
	gorm.Model
	QuizResponseID uint         `gorm:"not null;index" json:"quizResponseId"`
	QuestionID     uint         `gorm:"not null" json:"questionId"`
	Question       Question     `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	ChoiceID       *uint        `json:"choiceId,omitempty"`
	Choice         *Choice      `gorm:"foreignKey:ChoiceID" json:"choice,omitempty"`
	AnswerText     *string      `json:"answerText,omitempty"`
	IsCorrect      *bool        `json:"isCorrect,omitempty"`
}

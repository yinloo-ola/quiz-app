package models

import (
	"time"
)

// ResponderCredential represents a temporary credential for a user to take a specific quiz
type ResponderCredential struct {
	BaseModel               // Embed our custom base model
	QuizID       uint       `gorm:"not null;uniqueIndex:idx_quiz_username" json:"quizId"`
	Quiz         Quiz       `json:"quiz,omitempty"` // Belongs to Quiz
	Username     string     `gorm:"not null;uniqueIndex:idx_quiz_username" json:"username"` // Username unique per quiz
	PasswordHash string     `gorm:"not null;default:''" json:"-"` // Hashed password - Exclude!
	ExpiresAt    time.Time  `gorm:"not null" json:"expiresAt"`
	Used         bool       `gorm:"default:false" json:"used"` // Has the credential been used?
	UsedAt       *time.Time `json:"usedAt,omitempty"`
}

// QuizResponse represents a single completed attempt of a quiz by a responder.
type QuizResponse struct {
	BaseModel                                  // Embed our custom base model
	QuizID                uint                 `gorm:"not null" json:"quizId"`
	Quiz                  Quiz                 `gorm:"foreignKey:QuizID" json:"quiz,omitempty"`
	ResponderCredentialID *uint                `gorm:"uniqueIndex" json:"responderCredentialId,omitempty"`
	ResponderCredential   *ResponderCredential `gorm:"foreignKey:ResponderCredentialID" json:"responderCredential,omitempty"`
	ResponderUsername     string               `gorm:"index" json:"responderUsername"`
	Score                 *float64             `json:"score,omitempty"`
	StartedAt             *time.Time           `json:"startedAt,omitempty"`
	SubmittedAt           time.Time            `gorm:"not null;default:CURRENT_TIMESTAMP" json:"submittedAt"`
	TimeTakenSeconds      *int                 `json:"timeTakenSeconds,omitempty"`
	Answers               []Answer             `gorm:"foreignKey:QuizResponseID" json:"answers,omitempty"`
}

// Answer represents a single answer submitted by a responder for a specific question within a quiz response.
type Answer struct {
	BaseModel                   // Embed our custom base model
	QuizResponseID uint         `gorm:"column:response_id;not null;index" json:"quizResponseId"`
	// Add a field for the quiz_response_id column that maps to the same value as QuizResponseID
	DuplicateResponseID uint    `gorm:"column:quiz_response_id;not null;index" json:"-"`
	QuizResponse   QuizResponse `gorm:"foreignKey:QuizResponseID" json:"quizResponse,omitempty"`
	QuestionID     uint         `gorm:"not null" json:"questionId"`
	Question       Question     `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
	ChoiceID       *uint        `json:"choiceId,omitempty"`
	Choice         *Choice      `gorm:"foreignKey:ChoiceID" json:"choice,omitempty"`
	AnswerText     *string      `json:"answerText,omitempty"`
	IsCorrect      *bool        `json:"isCorrect,omitempty"`
	// Snapshot data to preserve the state at submission time
	QuestionTextSnapshot *string `json:"questionTextSnapshot,omitempty"`
	ChoicesSnapshot     *string `json:"choicesSnapshot,omitempty"` // JSON array of choice objects with ID, text, and isCorrect
}

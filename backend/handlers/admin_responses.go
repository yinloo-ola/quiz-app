package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/middleware"
	"github.com/yinloo-ola/quiz-app/backend/models"
	"gorm.io/gorm"
)

// QuizResponseView represents the data returned for a single quiz response in the list.
type QuizResponseView struct {
	ID                uint     `json:"id"`
	ResponderUsername string   `json:"responder_username"`
	Score             *float64 `json:"score"` // Pointer allows null
	SubmittedAt       time.Time `json:"submitted_at"`
}

// AnswerDetailView represents a single answer within the detailed response view.
type AnswerDetailView struct {
	QuestionText string  `json:"question_text"`
	SelectedChoiceText *string `json:"selected_choice_text,omitempty"` // Text of the selected choice
	AnswerText   *string `json:"answer_text,omitempty"`      // Text provided for open-ended
	IsCorrect    *bool   `json:"is_correct,omitempty"`       // If the answer was marked correct
}

// ResponseDetailView represents the full details of a single quiz response.
type ResponseDetailView struct {
	ID                uint               `json:"id"`
	QuizTitle         string             `json:"quiz_title"`
	ResponderUsername string             `json:"responder_username"`
	Score             *float64           `json:"score"`
	SubmittedAt       time.Time          `json:"submitted_at"`
	Answers           []AnswerDetailView `json:"answers"`
}

// ViewResponsesHandler handles fetching all responses for a specific quiz.
func ViewResponsesHandler(c *gin.Context) {
	// 1. Get Quiz ID
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 3. Verify Quiz exists and Admin owns it
	var quiz models.Quiz
	if err := tx.First(&quiz, quizID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + err.Error()})
		}
		return
	}
	if quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"}) // 404, not 403
		return
	}

	// 4. Fetch Responses for the Quiz
	var responses []models.QuizResponse
	query := tx.Where("quiz_id = ?", quizID).Order("submitted_at desc").Find(&responses)
	if query.Error != nil {
		tx.Rollback()
		log.Printf("Error fetching responses for quiz %d: %v", quizID, query.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responses"})
		return
	}

	// --- Commit Transaction (Read-only, but good practice) ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Prepare Response View
	responseViews := make([]QuizResponseView, len(responses))
	for i, resp := range responses {
		responseViews[i] = QuizResponseView{
			ID:                resp.ID,
			ResponderUsername: resp.ResponderUsername,
			Score:             resp.Score,
			SubmittedAt:       resp.SubmittedAt,
		}
	}

	c.JSON(http.StatusOK, responseViews)
}

// ViewResponseDetailsHandler handles fetching detailed information for a single response.
func ViewResponseDetailsHandler(c *gin.Context) {
	// 1. Get Response ID
	responseIDStr := c.Param("response_id")
	var responseID uint
	_, err := fmt.Sscan(responseIDStr, &responseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid response ID format"})
		return
	}

	// 2. Get Admin User ID
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 3. Fetch Response, Preloading related data for verification and details
	var response models.QuizResponse
	query := tx.Preload("Quiz"). // For ownership check
				 Preload("Answers"). // The list of answers
				 Preload("Answers.Question"). // Text of the question answered
				 Preload("Answers.Choice"). // Text of the choice selected
				 First(&response, responseID)

	if query.Error != nil {
		tx.Rollback()
		if query.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Response not found"})
		} else {
			log.Printf("Error finding response %d: %v", responseID, query.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding response"})
		}
		return
	}

	// 4. Verify Admin owns the Quiz associated with this Response
	if response.Quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Response not found"}) // 404, not 403
		return
	}

	// --- Commit Transaction (Read-only, but good practice) ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Prepare Detailed Response View
	answerViews := make([]AnswerDetailView, len(response.Answers))
	for i, ans := range response.Answers {
		var choiceText *string
		if ans.Choice != nil {
			choiceText = &ans.Choice.Text // Get text from preloaded Choice
		}
		answerViews[i] = AnswerDetailView{
			QuestionText: ans.Question.Text, // Get text from preloaded Question
			SelectedChoiceText: choiceText,
			AnswerText:   ans.AnswerText,
			IsCorrect:    ans.IsCorrect,
		}
	}

	detailedView := ResponseDetailView{
		ID:                response.ID,
		QuizTitle:         response.Quiz.Title,
		ResponderUsername: response.ResponderUsername,
		Score:             response.Score,
		SubmittedAt:       response.SubmittedAt,
		Answers:           answerViews,
	}

	c.JSON(http.StatusOK, detailedView)
}

package handlers

import (
	"encoding/json"
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
	ID                uint      `json:"id"`
	ResponderUsername string    `json:"responder_username"`
	Score             *float64  `json:"score"` // Pointer allows null
	SubmittedAt       time.Time `json:"submitted_at"`
	TimeTakenSeconds  *int      `json:"time_taken_seconds"`
}

// ChoiceView represents a single choice in a question
type ChoiceView struct {
	ID        uint   `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"isCorrect"`
}

// AnswerDetailView represents a single answer within the detailed response view.
type AnswerDetailView struct {
	QuestionText       string       `json:"question_text"`
	SelectedChoiceText *string      `json:"selected_choice_text,omitempty"` // Text of the selected choice
	AnswerText         *string      `json:"answer_text,omitempty"`          // Text provided for open-ended
	IsCorrect          *bool        `json:"isCorrect,omitempty"`           // If the answer was marked correct
	AllChoices         []ChoiceView `json:"all_choices,omitempty"`          // All available choices for the question
	CorrectChoiceIDs   []uint       `json:"correct_choice_ids,omitempty"`   // IDs of the correct choices
	SelectedChoiceIDs  []uint       `json:"selected_choice_ids,omitempty"`  // IDs of all selected choices
}

// ResponseDetailView represents the full details of a single quiz response.
type ResponseDetailView struct {
	ID                 uint               `json:"id"`
	QuizTitle          string             `json:"quiz_title"`
	ResponderUsername  string             `json:"responder_username"`
	Score              *float64           `json:"score"`
	StartedAt          *time.Time         `json:"started_at,omitempty"`
	SubmittedAt        time.Time          `json:"submitted_at"`
	TimeTakenSeconds   *int               `json:"time_taken_seconds,omitempty"`
	TimeTakenFormatted string             `json:"time_taken_formatted,omitempty"`
	Answers            []AnswerDetailView `json:"answers"`
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
			TimeTakenSeconds:  resp.TimeTakenSeconds,
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
					Preload("Answers").                  // The list of answers
					Preload("Answers.Question").         // Text of the question answered
					Preload("Answers.Question.Choices"). // All choices for the question
					Preload("Answers.Choice").           // Text of the choice selected
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
		var selectedChoiceIDs []uint
		var questionText string
		var allChoices []ChoiceView
		var correctChoiceIDs []uint

		// For single-choice questions
		if ans.Choice != nil {
			choiceText = &ans.Choice.Text // Get text from preloaded Choice
			selectedChoiceIDs = append(selectedChoiceIDs, ans.Choice.ID) // Add to the array
		}

		// For multi-choice questions, parse AnswerText if it exists and looks like JSON
		if ans.AnswerText != nil && len(*ans.AnswerText) > 0 && (*ans.AnswerText)[0] == '[' {
			// Try to parse it as a JSON array of choice IDs
			var choiceIDs []uint
			err := json.Unmarshal([]byte(*ans.AnswerText), &choiceIDs)
			if err == nil && len(choiceIDs) > 0 {
				selectedChoiceIDs = choiceIDs
			}
		}

		// Use snapshot data if available, otherwise use current question/choices
		if ans.QuestionTextSnapshot != nil {
			questionText = *ans.QuestionTextSnapshot
		} else {
			questionText = ans.Question.Text
		}

		// Use choices snapshot if available
		if ans.ChoicesSnapshot != nil {
			// Parse the choices snapshot
			var choiceSnapshots []map[string]interface{}
			err := json.Unmarshal([]byte(*ans.ChoicesSnapshot), &choiceSnapshots)
			if err == nil {
				// Convert snapshot to ChoiceView objects
				allChoices = make([]ChoiceView, len(choiceSnapshots))
				for j, choice := range choiceSnapshots {
					id, _ := choice["id"].(float64)      // JSON numbers are float64
					text, _ := choice["text"].(string)
					isCorrect, _ := choice["isCorrect"].(bool)

					allChoices[j] = ChoiceView{
						ID:        uint(id),
						Text:      text,
						IsCorrect: isCorrect,
					}

					// Collect correct choice IDs
					if isCorrect {
						correctChoiceIDs = append(correctChoiceIDs, uint(id))
					}
				}
			} else {
				log.Printf("Error parsing choices snapshot: %v", err)
				// Fall back to current choices if snapshot parsing fails
				allChoices = make([]ChoiceView, len(ans.Question.Choices))
				for j, choice := range ans.Question.Choices {
					allChoices[j] = ChoiceView{
						ID:        choice.ID,
						Text:      choice.Text,
						IsCorrect: choice.IsCorrect,
					}

					// Collect correct choice IDs
					if choice.IsCorrect {
						correctChoiceIDs = append(correctChoiceIDs, choice.ID)
					}
				}
			}
		} else {
			// Use current choices if no snapshot is available
			allChoices = make([]ChoiceView, len(ans.Question.Choices))
			for j, choice := range ans.Question.Choices {
				allChoices[j] = ChoiceView{
					ID:        choice.ID,
					Text:      choice.Text,
					IsCorrect: choice.IsCorrect,
				}

				// Collect correct choice IDs
				if choice.IsCorrect {
					correctChoiceIDs = append(correctChoiceIDs, choice.ID)
				}
			}
		}

		answerViews[i] = AnswerDetailView{
			QuestionText:       questionText,
			SelectedChoiceText: choiceText,
			SelectedChoiceIDs:  selectedChoiceIDs,  // Array of selected choice IDs
			AnswerText:         ans.AnswerText,
			IsCorrect:          ans.IsCorrect,
			AllChoices:         allChoices,
			CorrectChoiceIDs:   correctChoiceIDs,
		}
	}

	// Calculate time taken if StartedAt is available
	var timeTakenSeconds *int
	var timeTakenFormatted string
	if response.StartedAt != nil {
		// Calculate time taken in seconds
		duration := response.SubmittedAt.Sub(*response.StartedAt)
		seconds := int(duration.Seconds())
		timeTakenSeconds = &seconds

		// Format time taken as HH:MM:SS
		hours := seconds / 3600
		minutes := (seconds % 3600) / 60
		secs := seconds % 60
		timeTakenFormatted = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	} else if response.TimeTakenSeconds != nil {
		// If we have TimeTakenSeconds but no StartedAt, use that
		timeTakenSeconds = response.TimeTakenSeconds

		// Format time taken as HH:MM:SS
		hours := *timeTakenSeconds / 3600
		minutes := (*timeTakenSeconds % 3600) / 60
		secs := *timeTakenSeconds % 60
		timeTakenFormatted = fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
	}

	detailedView := ResponseDetailView{
		ID:                 response.ID,
		QuizTitle:          response.Quiz.Title,
		ResponderUsername:  response.ResponderUsername,
		Score:              response.Score,
		StartedAt:          response.StartedAt,
		SubmittedAt:        response.SubmittedAt,
		TimeTakenSeconds:   timeTakenSeconds,
		TimeTakenFormatted: timeTakenFormatted,
		Answers:            answerViews,
	}

	c.JSON(http.StatusOK, detailedView)
}

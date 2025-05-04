package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/middleware"
	"github.com/yinloo-ola/quiz-app/backend/models"
	"gorm.io/gorm"
)

// GetQuizForResponderHandler handles fetching quiz details for a responder.
// This endpoint returns the quiz without revealing the correct answers.
func GetQuizForResponderHandler(c *gin.Context) {
	// Get the quiz ID from the URL parameter
	quizIDStr := c.Param("quiz_id")
	quizID, err := strconv.ParseUint(quizIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// Get the responder's quiz ID from the context (set by middleware)
	responderQuizID, exists := c.Get(middleware.ResponderQuizIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Responder quiz ID not found in context"})
		return
	}

	// Verify the responder is authorized for this specific quiz
	if uint(quizID) != responderQuizID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to access this quiz"})
		return
	}

	// Fetch the quiz with questions and choices
	var quiz models.Quiz
	result := database.DB.Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC") // Order questions by their order field
	}).First(&quiz, quizID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			log.Printf("Database error fetching quiz %d: %v", quizID, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz details"})
		}
		return
	}

	// For each question, fetch its choices
	for i := range quiz.Questions {
		if err := database.DB.Model(&quiz.Questions[i]).
			Preload("Choices", func(db *gorm.DB) *gorm.DB {
				return db.Order("\"order\" ASC").Select("id, question_id, text, \"order\"") // Exclude is_correct flag
			}).
			First(&quiz.Questions[i], quiz.Questions[i].ID).Error; err != nil {
			log.Printf("Error fetching choices for question %d: %v", quiz.Questions[i].ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question choices"})
			return
		}
	}

	// Return the quiz data
	c.JSON(http.StatusOK, quiz)
}

// QuizSubmissionInput defines the structure for quiz submission request.
type QuizSubmissionInput struct {
	Answers []AnswerInput `json:"answers" binding:"required"`
	StartedAt *string `json:"started_at,omitempty"` // ISO8601 timestamp when the quiz was started
}

// AnswerInput defines the structure for an individual answer in the submission.
type AnswerInput struct {
	QuestionID  uint   `json:"question_id" binding:"required"`
	ChoiceIDs   []uint `json:"choice_ids" binding:"required"`
}

// SubmissionResult defines the structure for the quiz submission response.
type SubmissionResult struct {
	Score           float64                 `json:"score"`
	TotalQuestions  int                     `json:"total_questions"`
	CorrectAnswers  int                     `json:"correct_answers"`
	CorrectChoices  map[uint][]uint         `json:"correct_choices"`
}

// SubmitQuizHandler handles the submission of a quiz by a responder.
func SubmitQuizHandler(c *gin.Context) {
	// Get the quiz ID from the URL parameter
	quizIDStr := c.Param("quiz_id")
	quizID, err := strconv.ParseUint(quizIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// Get the responder's quiz ID and credential ID from the context (set by middleware)
	responderQuizID, exists := c.Get(middleware.ResponderQuizIDKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Responder quiz ID not found in context"})
		return
	}

	// Verify the responder is authorized for this specific quiz
	if uint(quizID) != responderQuizID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to submit answers for this quiz"})
		return
	}

	credentialID, exists := c.Get(middleware.ResponderAuthPayloadKey)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Responder credential ID not found in context"})
		return
	}

	// Bind the submission data
	var input QuizSubmissionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission data: " + err.Error()})
		return
	}

	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}

	// Fetch the quiz to verify it exists and get the responder credential
	var quiz models.Quiz
	if err := tx.First(&quiz, quizID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			log.Printf("Database error fetching quiz %d: %v", quizID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz"})
		}
		return
	}

	// Fetch the responder credential to mark it as used
	var credential models.ResponderCredential
	if err := tx.First(&credential, credentialID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch responder credential"})
		return
	}

	// Check if the credential has already been used
	if credential.Used {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "This quiz has already been submitted"})
		return
	}

	// Fetch all questions for the quiz to validate the submission
	var questions []models.Question
	if err := tx.Where("quiz_id = ?", quizID).Find(&questions).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error fetching questions for quiz %d: %v", quizID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quiz questions"})
		return
	}

	// Create a map of question IDs for quick lookup
	questionMap := make(map[uint]models.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// Validate that all submitted answers correspond to questions in this quiz
	for _, answer := range input.Answers {
		if _, exists := questionMap[answer.QuestionID]; !exists {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Answer submitted for question not in this quiz"})
			return
		}
	}

	// Fetch all correct choices for scoring
	var allChoices []models.Choice
	if err := tx.Where("question_id IN (?)", getQuestionIDs(questions)).Find(&allChoices).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error fetching choices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch question choices"})
		return
	}

	// Create maps for correct choices and all choices by question
	correctChoicesByQuestion := make(map[uint][]uint)
	allChoicesByQuestion := make(map[uint][]models.Choice)
	
	for _, choice := range allChoices {
		allChoicesByQuestion[choice.QuestionID] = append(allChoicesByQuestion[choice.QuestionID], choice)
		if choice.IsCorrect {
			correctChoicesByQuestion[choice.QuestionID] = append(correctChoicesByQuestion[choice.QuestionID], choice.ID)
		}
	}

	// Calculate score
	score, correctAnswers := calculateScore(input.Answers, correctChoicesByQuestion, questionMap)
	
	// Parse the start time if provided
	var startedAt *time.Time
	var timeTakenSeconds *int
	submittedAt := time.Now()
	
	if input.StartedAt != nil {
		parsedTime, err := time.Parse(time.RFC3339, *input.StartedAt)
		if err == nil {
			startedAt = &parsedTime
			
			// Calculate time taken in seconds
			duration := submittedAt.Sub(parsedTime)
			seconds := int(duration.Seconds())
			timeTakenSeconds = &seconds
			
			log.Printf("Quiz started at %v, submitted at %v, took %d seconds", parsedTime, submittedAt, seconds)
		} else {
			log.Printf("Error parsing start time: %v", err)
		}
	}
	
	// Create a response record
	response := models.QuizResponse{
		QuizID:            uint(quizID),
		ResponderUsername: credential.Username,
		StartedAt:         startedAt,
		SubmittedAt:       submittedAt,
		TimeTakenSeconds:  timeTakenSeconds,
		Score:             &score,
	}

	if err := tx.Create(&response).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error creating response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quiz response"})
		return
	}

	// Create answer records for each submitted answer
	for _, answer := range input.Answers {
		// Convert choice IDs to JSON
		choiceIDsJSON, err := json.Marshal(answer.ChoiceIDs)
		if err != nil {
			tx.Rollback()
			log.Printf("Error marshalling choice IDs: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process answer data"})
			return
		}

		// For now, we'll store the first choice ID if available, and store the JSON in AnswerText
		var choiceID *uint
		if len(answer.ChoiceIDs) > 0 {
			choiceID = &answer.ChoiceIDs[0]
		}
		
		choiceIDsJSONStr := string(choiceIDsJSON)
		
		// Get the question text for snapshot
		question, exists := questionMap[answer.QuestionID]
		if !exists {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Question not found for answer"})
			return
		}
		questionText := question.Text
		
		// Get the choices for this question and create a snapshot
		questionChoices := allChoicesByQuestion[answer.QuestionID]
		choiceSnapshots := make([]map[string]interface{}, len(questionChoices))
		for i, choice := range questionChoices {
			choiceSnapshots[i] = map[string]interface{}{
				"id":        choice.ID,
				"text":      choice.Text,
				"isCorrect": choice.IsCorrect,
			}
		}
		
		// Convert choices snapshot to JSON
		choicesSnapshotJSON, err := json.Marshal(choiceSnapshots)
		if err != nil {
			tx.Rollback()
			log.Printf("Error marshalling choices snapshot: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process choices data"})
			return
		}
		choicesSnapshotStr := string(choicesSnapshotJSON)
		
		answerRecord := models.Answer{
			QuizResponseID:      response.ID,
			DuplicateResponseID: response.ID, // Set both fields to the same value
			QuestionID:         answer.QuestionID,
			ChoiceID:           choiceID,
			AnswerText:         &choiceIDsJSONStr,  // Store all selected choices as JSON in AnswerText
			QuestionTextSnapshot: &questionText,
			ChoicesSnapshot:     &choicesSnapshotStr,
		}

		if err := tx.Create(&answerRecord).Error; err != nil {
			tx.Rollback()
			log.Printf("Database error creating answer record: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save answer"})
			return
		}
	}

	// Mark the credential as used
	credential.Used = true
	if err := tx.Save(&credential).Error; err != nil {
		tx.Rollback()
		log.Printf("Database error updating credential: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update credential status"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Database error committing transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save quiz submission"})
		return
	}

	// Return the result
	result := SubmissionResult{
		Score:           score,
		TotalQuestions:  len(questions),
		CorrectAnswers:  correctAnswers,
		CorrectChoices:  correctChoicesByQuestion,
	}

	c.JSON(http.StatusOK, result)
}

// Helper function to extract question IDs from a slice of questions
func getQuestionIDs(questions []models.Question) []uint {
	ids := make([]uint, len(questions))
	for i, q := range questions {
		ids[i] = q.ID
	}
	return ids
}

// Helper function to calculate the score based on submitted answers
func calculateScore(answers []AnswerInput, correctChoices map[uint][]uint, questions map[uint]models.Question) (float64, int) {
	totalQuestions := len(questions)
	if totalQuestions == 0 {
		return 0, 0
	}

	correctAnswers := 0

	// Create a map of submitted answers for quick lookup
	submittedAnswers := make(map[uint][]uint)
	for _, answer := range answers {
		submittedAnswers[answer.QuestionID] = answer.ChoiceIDs
	}

	// Check each question
	for questionID, correctChoiceIDs := range correctChoices {
		// Skip if no answer was submitted for this question
		submitted, exists := submittedAnswers[questionID]
		if !exists {
			continue
		}

		// For single-choice questions, check if the single submitted answer matches any correct choice
		if questions[questionID].Type == "single" {
			if len(submitted) == 1 && contains(correctChoiceIDs, submitted[0]) {
				correctAnswers++
			}
		} else if questions[questionID].Type == "multi" {
			// For multi-choice questions, check if all submitted choices are correct and all correct choices are submitted
			if setsEqual(submitted, correctChoiceIDs) {
				correctAnswers++
			}
		}
	}

	// Calculate score as percentage
	score := float64(correctAnswers) / float64(totalQuestions) * 100
	return score, correctAnswers
}

// Helper function to check if a slice contains a value
func contains(slice []uint, val uint) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// Helper function to check if two slices contain the same elements (regardless of order)
func setsEqual(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[uint]bool)
	for _, val := range a {
		aMap[val] = true
	}

	for _, val := range b {
		if !aMap[val] {
			return false
		}
	}

	return true
}

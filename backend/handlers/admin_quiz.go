package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/yinloo-ola/quiz-app/backend/database"
	"github.com/yinloo-ola/quiz-app/backend/middleware"
	"github.com/yinloo-ola/quiz-app/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// --- Structs for Create Quiz Request --- //

// CreateChoiceRequest defines the structure for a choice within a question creation request.
type CreateChoiceRequest struct {
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"is_correct"` // Defaults to false if omitted
}

// CreateQuestionRequest defines the structure for a question within a quiz creation request.
type CreateQuestionRequest struct {
	Text    string                `json:"text" binding:"required"`
	Type    string                `json:"type" binding:"required"`               // Add this line
	Choices []CreateChoiceRequest `json:"choices" binding:"required,min=1,dive"` // Must have at least one choice
	// TODO: Add Order fields if needed later
}

// CreateQuizRequest defines the structure for the entire create quiz request body.
type CreateQuizRequest struct {
	Title       string                  `json:"title" binding:"required"`
	Description string                  `json:"description"`
	TimeLimit   *uint                   `json:"time_limit"`                              // Optional time limit in minutes
	Questions   []CreateQuestionRequest `json:"questions" binding:"required,min=1,dive"` // Must have at least one question
}

// --- Structs for Update Quiz Request --- //

// ChoiceInput defines the structure for a choice in create/update requests.
type ChoiceInput struct {
	ID        *uint  `json:"id"` // Optional ID from frontend, not used in delete/recreate
	Text      string `json:"text" binding:"required"`
	IsCorrect bool   `json:"isCorrect"` // Matches frontend camelCase
}

// QuestionInput defines the structure for a question in create/update requests.
type QuestionInput struct {
	ID      *uint         `json:"id"` // Optional ID from frontend, not used in delete/recreate
	Text    string        `json:"text" binding:"required"`
	Type    string        `json:"type" binding:"required,oneof=single multi"` // Validate type
	Choices []ChoiceInput `json:"choices" binding:"required,dive"`            // Use dive for nested validation
}

// UpdateQuizRequest defines the expected JSON body for updating a quiz.
// Uses pointers for optional fields.
type UpdateQuizRequest struct {
	Title            *string          `json:"title"`              // Optional update
	Description      *string          `json:"description"`        // Optional update
	TimeLimitSeconds *uint            `json:"time_limit_seconds"` // Optional update, use pointer to distinguish 0/nil from not provided
	Questions        *[]QuestionInput `json:"questions"`          // Pointer to slice to detect if questions array was provided for update
}

// --- Structs for Update Question Request --- //

// UpdateQuestionRequest defines the structure for updating a question.
// For now, only allowing text update for simplicity.
type UpdateQuestionRequest struct {
	Text *string `json:"text"` // Use pointer for optional update
}

// --- Handler Function --- //

// CreateQuizHandler handles the creation of a new quiz.
func CreateQuizHandler(c *gin.Context) {
	var req CreateQuizRequest

	// 1. Bind JSON request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// 2. Get Admin User ID from context (set by middleware)
	adminUserIDAny, exists := c.Get(middleware.AuthPayloadKey)
	if !exists {
		// This should ideally not happen if middleware is applied correctly
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID not found in context"})
		return
	}
	adminUserID, ok := adminUserIDAny.(uint)
	if !ok {
		// This indicates a problem with how the ID was stored in the context
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin User ID has incorrect type in context"})
		return
	}

	// 3. Create Quiz, Questions, and Choices within a Transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// Create the Quiz record
	newQuiz := models.Quiz{
		Title:       req.Title,
		Description: req.Description,
		TimeLimit:   req.TimeLimit,
		AdminUserID: adminUserID,
	}
	if err := tx.Create(&newQuiz).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quiz: " + err.Error()})
		return
	}

	// Create Questions and Choices
	for _, questionReq := range req.Questions {
		newQuestion := models.Question{
			Text:   questionReq.Text,
			Type:   models.QuestionType(questionReq.Type), // Convert string to QuestionType
			QuizID: newQuiz.ID,
			// Order: // Set order if needed
			// Type:  // Set type if needed
		}
		if err := tx.Create(&newQuestion).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question: " + err.Error()})
			return
		}

		for _, choiceReq := range questionReq.Choices {
			newChoice := models.Choice{
				Text:       choiceReq.Text,
				IsCorrect:  choiceReq.IsCorrect,
				QuestionID: newQuestion.ID,
			}
			if err := tx.Create(&newChoice).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create choice: " + err.Error()})
				return
			}
		}
	}

	// 4. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		// Rollback already attempted implicitly by GORM on commit failure, but log anyway
		// tx.Rollback() // Not strictly needed here
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 5. Respond with created Quiz ID (or full details if preferred)
	// Fetching the full quiz with preloads might be inefficient here.
	// Returning just the ID is simpler for a create operation.
	c.JSON(http.StatusCreated, gin.H{
		"message": "Quiz created successfully",
		"quiz_id": newQuiz.ID,
	})
}

// GetQuizzesHandler handles listing quizzes for the logged-in admin.
func GetQuizzesHandler(c *gin.Context) {
	// 1. Get Admin User ID from context
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

	// 2. Find quizzes belonging to the admin
	var quizzes []models.Quiz
	// Preload Questions and their Choices to include them in the response
	result := database.DB.Preload("Questions.Choices").Where("admin_user_id = ?", adminUserID).Order("created_at desc").Find(&quizzes)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch quizzes: " + result.Error.Error()})
		return
	}

	// 3. Return the list of quizzes
	// If no quizzes found, return an empty list, not an error
	c.JSON(http.StatusOK, quizzes)
}

// GetQuizDetailsHandler handles retrieving details for a specific quiz.
func GetQuizDetailsHandler(c *gin.Context) {
	// 1. Get Quiz ID from path parameter
	quizIDStr := c.Param("quiz_id")
	// TODO: Add proper error handling for non-integer quiz_id
	var quizID uint
	// Cheap conversion for now, replace with strconv.ParseUint if more robust handling needed
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Find the quiz by ID, preload questions and choices, and verify ownership
	var quiz models.Quiz
	result := database.DB.Preload("Questions").Preload("Questions.Choices").First(&quiz, quizID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + result.Error.Error()})
		}
		return
	}

	// 4. Verify the quiz belongs to the logged-in admin
	if quiz.AdminUserID != adminUserID {
		// Return 404 rather than 403 to avoid revealing existence
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	// 5. Return the full quiz details
	c.JSON(http.StatusOK, quiz)
}

// UpdateQuizHandler handles updating an existing quiz.
func UpdateQuizHandler(c *gin.Context) {
	// 1. Get Quiz ID from path parameter
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Bind JSON request body
	var req UpdateQuizRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Provide more specific validation errors if possible
		validationErrs, ok := err.(validator.ValidationErrors)
		if ok {
			errorsMap := make(map[string]string)
			for _, fieldErr := range validationErrs {
				errorsMap[fieldErr.Field()] = fmt.Sprintf("Validation failed on '%s' tag", fieldErr.Tag())
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": errorsMap})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		}
		return
	}

	// 4. Start Transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Printf("Failed to start transaction: %v", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction"})
		return
	}
	// Defer Rollback in case of panic or error returns
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			// Optionally re-panic or log
			panic(r) // Re-panic if needed
		} else if tx.Error != nil {
			// If tx.Error is set during commit or other operations
			log.Printf("Rolling back transaction due to error: %v", tx.Error)
			tx.Rollback()
		}
	}()

	// 5. Find the existing quiz (within transaction)
	var existingQuiz models.Quiz
	if err := tx.First(&existingQuiz, quizID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback() // Rollback before sending response
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + err.Error()})
		}
		return
	}

	// 6. Verify ownership
	if existingQuiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found or not owned by user"}) // Keep consistent with GET detail
		return
	}

	// 7. Update top-level quiz fields if provided
	updates := make(map[string]interface{})
	needsUpdate := false
	if req.Title != nil {
		updates["title"] = *req.Title
		needsUpdate = true
	}
	if req.Description != nil {
		updates["description"] = *req.Description
		needsUpdate = true
	}
	// Use the correct JSON key `time_limit_seconds` and DB field `TimeLimitSeconds`
	if req.TimeLimitSeconds != nil {
		updates["time_limit_seconds"] = *req.TimeLimitSeconds // Use the correct field name from the DB model
		needsUpdate = true
	}

	if needsUpdate {
		if err := tx.Model(&existingQuiz).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quiz details: " + err.Error()})
			return
		}
	}

	// 8. Handle Question/Choice updates only if 'questions' key was present in the request
	if req.Questions != nil {
		log.Printf("Updating questions for quiz ID: %d", quizID)

		// 8a. Delete existing choices associated with the quiz
		if err := tx.Where("question_id IN (SELECT id FROM questions WHERE quiz_id = ?)", quizID).Delete(&models.Choice{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing choices: " + err.Error()})
			return
		}

		// 8b. Delete existing questions associated with the quiz
		if err := tx.Where("quiz_id = ?", quizID).Delete(&models.Question{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete existing questions: " + err.Error()})
			return
		}

		// 8c. Create new questions and choices based on the request
		for _, questionReq := range *req.Questions { // Dereference pointer
			questionType, err := models.ParseQuestionType(questionReq.Type)
			if err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid question type '%s' provided: %v", questionReq.Type, err)})
				return
			}

			newQuestion := models.Question{
				Text:   questionReq.Text,
				Type:   questionType,
				QuizID: quizID, // Link to the current quiz
			}
			if err := tx.Create(&newQuestion).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question: " + err.Error()})
				return
			}

			// Convert []ChoiceInput to []models.CorrectableChoice for validation
			choicesForValidation := make([]models.CorrectableChoice, len(questionReq.Choices))
			for i := range questionReq.Choices {
				choicesForValidation[i] = questionReq.Choices[i] // Assigning concrete type that satisfies interface
			}

			if err := models.ValidateChoicesForQuestionType(questionType, choicesForValidation); err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid choices for question '%s': %v", questionReq.Text, err)})
				return
			}

			for _, choiceReq := range questionReq.Choices {
				newChoice := models.Choice{
					Text:       choiceReq.Text,
					IsCorrect:  choiceReq.IsCorrect,
					QuestionID: newQuestion.ID, // Link to the newly created question
				}
				if err := tx.Create(&newChoice).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create choice: " + err.Error()})
					return
				}
			}
		}
	}

	// 9. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		// Rollback already deferred, just log and return error
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 10. Return the updated quiz (refetch outside transaction to ensure consistency)
	var updatedQuiz models.Quiz
	fetchResult := database.DB.Preload("Questions").Preload("Questions.Choices").First(&updatedQuiz, quizID)
	if fetchResult.Error != nil {
		// Log error, but might still return OK if update succeeded
		log.Printf("Error refetching quiz %d after update: %v", quizID, fetchResult.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refetch quiz after update"})
		return
	}

	c.JSON(http.StatusOK, updatedQuiz)
}

// DeleteQuizHandler handles soft-deleting a quiz.
func DeleteQuizHandler(c *gin.Context) {
	// 1. Get Quiz ID from path parameter
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Find the existing quiz to verify ownership before deleting
	// Important: Use Unscoped() if you need to potentially find/delete an already soft-deleted record,
	// but for a standard delete, finding only non-deleted records is correct.
	var quizToDelete models.Quiz
	result := database.DB.First(&quizToDelete, quizID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + result.Error.Error()})
		}
		return
	}

	// 4. Verify ownership
	if quizToDelete.AdminUserID != adminUserID {
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		return
	}

	// 5. Perform the soft delete
	// GORM's Delete() performs a soft delete if the model has a DeletedAt field
	deleteResult := database.DB.Delete(&quizToDelete)
	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz: " + deleteResult.Error.Error()})
		return
	}

	if deleteResult.RowsAffected == 0 {
		// This shouldn't normally happen if the record was found and ownership verified,
		// but could indicate a race condition or an issue with the delete operation itself.
		log.Printf("Delete requested for quiz %d, but no rows were affected.", quizID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete quiz, zero rows affected"})
		return
	}

	// 6. Respond with No Content
	c.Status(http.StatusNoContent)
}

// AddQuestionHandler handles adding a new question to an existing quiz.
func AddQuestionHandler(c *gin.Context) {
	// 1. Get Quiz ID from path parameter
	quizIDStr := c.Param("quiz_id")
	var quizID uint
	_, err := fmt.Sscan(quizIDStr, &quizID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quiz ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Bind JSON request body (reusing CreateQuestionRequest)
	var req CreateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 4. Find the target quiz and verify ownership (within transaction for locking)
	var targetQuiz models.Quiz
	// Use .Locking(clause.LockingOptions{Strength: "UPDATE"}) if high concurrency risk, but simple .First should suffice here.
	if err := tx.First(&targetQuiz, quizID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz: " + err.Error()})
		}
		return
	}

	if targetQuiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Quiz not found"}) // Return 404, not 403
		return
	}

	// 5. Create the Question record
	newQuestion := models.Question{
		Text:   req.Text,
		Type:   models.QuestionType(req.Type), // Convert string to QuestionType
		QuizID: targetQuiz.ID,
		// Order/Type can be added later if needed
	}
	if err := tx.Create(&newQuestion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create question: " + err.Error()})
		return
	}

	// 6. Create the Choice records
	for _, choiceReq := range req.Choices {
		newChoice := models.Choice{
			Text:       choiceReq.Text,
			IsCorrect:  choiceReq.IsCorrect,
			QuestionID: newQuestion.ID,
		}
		if err := tx.Create(&newChoice).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create choice: " + err.Error()})
			return
		}
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 7. Respond with the created question details (including choices)
	// Need to reload the question with choices to return them
	var createdQuestionWithChoices models.Question
	if err := database.DB.Preload("Choices").First(&createdQuestionWithChoices, newQuestion.ID).Error; err != nil {
		// Log the error, but maybe return a simpler success message if refetch fails?
		log.Printf("Error refetching question %d after creation: %v", newQuestion.ID, err)
		c.JSON(http.StatusCreated, gin.H{
			"message":     "Question created successfully, but failed to refetch details",
			"question_id": newQuestion.ID,
		})
		return
	}

	c.JSON(http.StatusCreated, createdQuestionWithChoices)
}

// UpdateQuestionHandler handles updating an existing question.
func UpdateQuestionHandler(c *gin.Context) {
	// 1. Get Question ID from path parameter
	questionIDStr := c.Param("question_id")
	var questionID uint
	_, err := fmt.Sscan(questionIDStr, &questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Bind JSON request body
	var req UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// --- Transaction Start ---
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start database transaction: " + tx.Error.Error()})
		return
	}

	// 4. Find the target question and verify ownership via its quiz
	var targetQuestion models.Question
	// Fetch the question first
	if err := tx.First(&targetQuestion, questionID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding question: " + err.Error()})
		}
		return
	}

	// Fetch the associated quiz to check ownership
	var quiz models.Quiz
	if err := tx.First(&quiz, targetQuestion.QuizID).Error; err != nil {
		tx.Rollback()
		// If the quiz isn't found, something is wrong (orphaned question?)
		// Treat as 'Question not found' for the user.
		if err == gorm.ErrRecordNotFound {
			log.Printf("Error: Question %d found, but associated Quiz %d not found.", targetQuestion.ID, targetQuestion.QuizID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found (or associated quiz missing)"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz for question: " + err.Error()})
		}
		return
	}

	if quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"}) // Return 404 for consistency
		return
	}

	// 5. Update the Question record if fields are provided
	updates := make(map[string]interface{})
	if req.Text != nil {
		updates["text"] = *req.Text
	}

	if len(updates) == 0 {
		tx.Rollback() // No changes needed, rollback transaction
		// Optionally return the existing question data or just 200 OK
		c.JSON(http.StatusOK, targetQuestion) // Return current data if no update needed
		return
	}

	if err := tx.Model(&targetQuestion).Updates(updates).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update question: " + err.Error()})
		return
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 6. Respond with the updated question details
	// Need to reload to get the final state, including any default values set by DB
	var updatedQuestion models.Question
	// Still preload choices for the response
	if err := database.DB.Preload("Choices").First(&updatedQuestion, targetQuestion.ID).Error; err != nil {
		log.Printf("Error refetching question %d after update: %v", targetQuestion.ID, err)
		c.JSON(http.StatusOK, gin.H{
			"message":     "Question updated successfully, but failed to refetch details",
			"question_id": targetQuestion.ID,
		})
		return
	}

	c.JSON(http.StatusOK, updatedQuestion)
}

// DeleteQuestionHandler handles soft-deleting an existing question.
func DeleteQuestionHandler(c *gin.Context) {
	// 1. Get Question ID from path parameter
	questionIDStr := c.Param("question_id")
	var questionID uint
	_, err := fmt.Sscan(questionIDStr, &questionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid question ID format"})
		return
	}

	// 2. Get Admin User ID from context
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

	// 3. Find the target question
	var targetQuestion models.Question
	if err := tx.First(&targetQuestion, questionID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding question: " + err.Error()})
		}
		return
	}

	// 4. Find the associated quiz to check ownership
	var quiz models.Quiz
	if err := tx.First(&quiz, targetQuestion.QuizID).Error; err != nil {
		tx.Rollback()
		log.Printf("Error finding quiz %d for question %d during delete: %v", targetQuestion.QuizID, targetQuestion.ID, err)
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Question not found (or associated quiz missing)"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error finding quiz for question: " + err.Error()})
		}
		return
	}

	// 5. Verify ownership
	if quiz.AdminUserID != adminUserID {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"}) // Return 404
		return
	}

	// 6. Soft-delete the question
	// GORM's default Delete performs a soft delete if the model has DeletedAt
	if err := tx.Delete(&targetQuestion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete question: " + err.Error()})
		return
	}

	// --- Commit Transaction ---
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction: " + err.Error()})
		return
	}

	// 7. Respond with No Content
	c.Status(http.StatusNoContent)
}

// Helper function to validate choices based on question type
// This should ideally live in the models package or a validation service
func (ci ChoiceInput) GetIsCorrect() bool { // Helper for interface matching if needed
	return ci.IsCorrect
}

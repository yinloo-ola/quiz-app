// Define shared types for the frontend application

// Type for question choices (single or multiple correct answers)
export type QuestionType = 'single' | 'multi';

// Interface for a choice input when creating/editing a quiz
export interface ChoiceInput {
  text: string;
  is_correct: boolean;
}

// Interface for a question input when creating/editing a quiz
export interface QuestionInput {
  text: string;
  type: QuestionType;
  choices: ChoiceInput[];
}

// Interface for the payload when creating a new quiz via the API
export interface QuizCreatePayload {
  title: string;
  description?: string;
  time_limit_seconds?: number;
  questions: QuestionInput[];
}

// You can add more types here as needed, for example:
// - Types for quiz data received from the backend
// - Types for API responses
// - Types for user authentication data

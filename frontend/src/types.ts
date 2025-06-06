// Define shared types for the frontend application

// Type for question choices (single or multiple correct answers)
export type QuestionType = 'single' | 'multi';

// --- Backend Data Structures ---

export interface Choice {
  id: number;
  text: string;
  isCorrect: boolean;
}

export interface Question {
  id: number;
  text: string;
  type: QuestionType;
  choices: Choice[];
}

export interface Quiz {
  id: number;
  title: string;
  description?: string;
  timeLimit?: number; // Time limit in seconds
  status?: string; // Added quiz status (e.g., 'draft', 'published')
  questions: Question[];
  created_at?: string; // Optional metadata
  updated_at?: string; // Optional metadata
}

// --- API Payloads ---

// Interface for a choice input when creating/editing a quiz
export interface ChoiceInput {
  id?: number; // Optional ID for updates
  text: string;
  isCorrect: boolean;
}

// Interface for a question input when creating/editing a quiz
export interface QuestionInput {
  id?: number; // Optional ID for updates
  text: string;
  type: QuestionType;
  choices: ChoiceInput[];
}

// Represents the full Quiz object structure received from the backend
export type QuizInput = {
  title: string;
  description?: string;
  timeLimit?: number; // Time limit in seconds
  questions: QuestionInput[];
};

// Interface for the payload when creating a new quiz via the API
export interface QuizCreatePayload {
  title: string;
  description?: string;
  timeLimit?: number; // Time limit in seconds
  questions: QuestionInput[];
}

// Interface for the payload when updating an existing quiz via the API
// Often similar to create, but might require ID or handle partial updates
// For now, let's assume it's the same as create for simplicity, but include the ID
export interface QuizUpdatePayload extends QuizCreatePayload {
  id: number;
}

// Represents a responder credential generated by an admin
export interface ResponderCredential {
  id: number;
  createdAt: string; // ISO Date string
  updatedAt: string; // ISO Date string
  deletedAt?: string | null; // ISO Date string or null
  quizId: number;
  username: string;
  expiresAt?: string | null; // ISO Date string or null
  used: boolean;
  usedAt?: string | null; // ISO 8601 string or null
}

// Interface for the response when generating new credentials
export interface GenerateCredentialsResponse {
  username: string;
  password: string;         // Plaintext password, only available on creation
  expiresAt?: string | null; // ISO 8601 string or null
  credential_id: number;    // Changed from uint to number for TS
}

// Represents a summary of a submitted quiz response for the admin list view
export interface QuizResponseSummary {
  id: number;
  quizId: number;
  // Support both camelCase and snake_case property names
  responderUsername?: string;
  responder_username: string; // The username of the responder who took the quiz
  score?: number | null; // Calculated score (might be null if not graded)
  submittedAt?: string;
  submitted_at: string; // ISO Date string when the response was submitted
  createdAt?: string;
  created_at?: string; // ISO Date string
  // Time tracking fields
  startedAt?: string;
  started_at?: string;
  timeTakenSeconds?: number;
  time_taken_seconds?: number;
  timeTakenFormatted?: string;
  time_taken_formatted?: string;
}

// Represents a choice in a question
export interface Choice {
  id: number;
  text: string;
  isCorrect: boolean;
}

// Represents an answer submitted by a responder for a specific question
export interface SubmittedAnswer {
  // Fields from the frontend submission (camelCase)
  questionId?: number;
  chosenChoiceIds?: number[]; // Array of choice IDs selected by the responder
  questionText?: string; // The text of the question
  selectedChoiceText?: string | null; // The text of the selected choice
  answerText?: string | null; // The raw answer text (for multiple choices, this might be JSON)

  // Fields from the backend response (snake_case)
  question_text?: string; // The text of the question
  selected_choice_text?: string | null; // The text of the selected choice
  selected_choice_ids?: number[]; // IDs of all selected choices
  answer_text?: string | null; // The raw answer text (for multiple choices, this might be JSON)
  isCorrect?: boolean; // Whether the answer is correct
  all_choices?: Choice[]; // All available choices for the question
  correct_choice_ids?: number[]; // IDs of the correct choices
}

// Represents the detailed view of a submitted quiz response, including answers
export interface QuizResponseDetail {
  // Support both camelCase (frontend) and snake_case (API) formats
  id: number;
  quizId?: number;
  quiz_id?: number;
  quizTitle?: string;
  quiz_title: string; // Title of the quiz taken
  responderUsername?: string;
  responder_username: string;
  startedAt?: string;
  started_at?: string;
  submittedAt?: string;
  submitted_at: string;
  createdAt?: string;
  created_at?: string;
  score?: number;
  timeTakenSeconds?: number;
  time_taken_seconds?: number;
  timeTakenFormatted?: string;
  time_taken_formatted?: string;
  answers: SubmittedAnswer[];
  // We might also need the original quiz questions/choices here to display them
  // alongside the answers, or fetch them separately.
  // Let's assume for now the backend includes enough info, or we enhance later.
  quiz?: Quiz; // Optional: Include full quiz structure if backend provides it
}

// You can add more types here as needed, for example:
// - Types for API responses
// - Types for user authentication data

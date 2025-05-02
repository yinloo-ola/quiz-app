// Define shared types for the frontend application

// Type for question choices (single or multiple correct answers)
export type QuestionType = 'single' | 'multi';

// --- Backend Data Structures ---

export interface Choice {
  id: number;
  text: string;
  is_correct: boolean;
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
  time_limit_seconds?: number;
  questions: Question[];
  created_at?: string; // Optional metadata
  updated_at?: string; // Optional metadata
}

// --- API Payloads ---

// Interface for a choice input when creating/editing a quiz
export interface ChoiceInput {
  id?: number; // Optional ID for updates
  text: string;
  is_correct: boolean;
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
  time_limit_seconds?: number; // Optional on creation, backend might default
  questions: QuestionInput[];
};

// Interface for the payload when creating a new quiz via the API
export interface QuizCreatePayload {
  title: string;
  description?: string;
  time_limit_seconds?: number;
  questions: QuestionInput[];
}

// Interface for the payload when updating an existing quiz via the API
// Often similar to create, but might require ID or handle partial updates
// For now, let's assume it's the same as create for simplicity, but include the ID
export interface QuizUpdatePayload extends QuizCreatePayload {
  id: number;
}

// You can add more types here as needed, for example:
// - Types for API responses
// - Types for user authentication data

import axios from 'axios';
import { useAuthStore } from '@/stores/auth'; // Use '@/stores' alias
import type {
  Quiz,
  QuizInput,
  ResponderCredential,
  QuizResponseSummary,
  QuizResponseDetail,
  GenerateCredentialsResponse // Import the new type
} from '@/types'; // Import Quiz and ResponderCredential types

const apiClient = axios.create({
  baseURL: 'http://localhost:8081', // Your Go backend URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// --- Axios Interceptor --- 
// Automatically add the JWT token to requests if it exists
apiClient.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore();
    const token = authStore.currentToken;

    // DEBUG: Log the request URL and the token being used
    console.log(`[API Interceptor] Request URL: ${config.url}, Token: ${token}`);

    // Add Authorization header only if the token exists and the request
    // isn't for the login endpoints themselves.
    if (token && config.url !== '/admin/login' && config.url !== '/responder/login') {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// --- API Functions --- 

// Define the expected login credentials structure
interface AdminLoginCredentials {
  username: string;
  password: string;
}

// Define the expected login response structure
interface AdminLoginResponse {
  token: string;
  // Include other fields if your backend login returns more data
}

// Define the responder login credentials structure
interface ResponderLoginCredentials {
  username: string;
  password: string;
}

// Define the responder login response structure
interface ResponderLoginResponse {
  token: string;
}

// Define the quiz submission input structure
interface QuizSubmissionInput {
  answers: {
    question_id: number;
    choice_ids: number[];
  }[];
  started_at?: string; // ISO string of when the quiz was started
}

// Define the quiz submission result structure
interface QuizSubmissionResult {
  score: number;
  total_questions: number;
  correct_answers: number;
  correct_choices: Record<number, number[]>; // Map of question_id to array of correct choice_ids
}

/**
 * Sends admin login credentials to the backend.
 * @param credentials - The admin username and password.
 * @returns Promise resolving with the login response (including token).
 */
export const adminLogin = async (credentials: AdminLoginCredentials): Promise<AdminLoginResponse> => {
  try {
    const response = await apiClient.post<AdminLoginResponse>('/admin/login', credentials);
    return response.data;
  } catch (error: any) {
    console.error('Admin login failed:', error);
    // Re-throw the error so the calling component/store can handle it
    // You might want to parse the error response for specific messages
    throw error.response?.data || error;
  }
};

/**
 * Sends responder login credentials to the backend.
 * @param credentials - The responder username and password.
 * @returns Promise resolving with the login response (including token).
 */
export const responderLogin = async (credentials: ResponderLoginCredentials): Promise<ResponderLoginResponse> => {
  try {
    const response = await apiClient.post<ResponderLoginResponse>('/responder/login', credentials);
    return response.data;
  } catch (error: any) {
    console.error('Responder login failed:', error);
    throw error.response?.data || error;
  }
};

/**
 * Fetches quiz details for a responder (without correct answer flags).
 * @param quizId - The ID of the quiz to fetch.
 * @returns Promise resolving with the quiz data.
 */
export const getResponderQuiz = async (quizId: number): Promise<Quiz> => {
  try {
    const response = await apiClient.get<Quiz>(`/quizzes/${quizId}`);
    return response.data;
  } catch (error: any) {
    console.error(`Failed to fetch responder quiz ${quizId}:`, error);
    throw error.response?.data || error;
  }
};

/**
 * Submits quiz answers for a responder.
 * @param quizId - The ID of the quiz being submitted.
 * @param answers - The answers to submit.
 * @returns Promise resolving with the submission result.
 */
export const submitQuiz = async (quizId: number, answers: QuizSubmissionInput): Promise<QuizSubmissionResult> => {
  try {
    const response = await apiClient.post<QuizSubmissionResult>(`/quizzes/${quizId}/submit`, answers);
    return response.data;
  } catch (error: any) {
    console.error(`Failed to submit quiz ${quizId}:`, error);
    throw error.response?.data || error;
  }
};

// Function to create a new admin quiz
export const createAdminQuiz = async (quizData: QuizInput) => {
  try {
    console.log('Sending create quiz payload:', quizData);
    const response = await apiClient.post('/admin/quizzes', quizData);
    return response.data; // Return the created quiz data from backend
  } catch (error) {
    console.error('Error creating admin quiz:', error);
    // Rethrow or handle error appropriately for the UI
    if (axios.isAxiosError(error) && error.response) {
      // Extract backend error message if available
      throw new Error(error.response.data.error || 'Failed to create quiz');
    } else {
      // Handle non-Axios errors or errors without a response
      throw new Error('An unexpected error occurred while creating the quiz.');
    }
  }
};

// Fetch quizzes created by the authenticated admin
export const fetchAdminQuizzes = async (): Promise<any[]> => {
  try {
    const response = await apiClient.get('/admin/quizzes');
    return response.data; // Assuming the backend returns an array of quizzes
  } catch (error: any) {
    console.error('Failed to fetch admin quizzes:', error);
    throw error.response?.data || error;
  }
};

/**
 * Fetch a specific quiz by its ID for the admin.
 * @param quizId The ID of the quiz to fetch.
 * @returns Promise resolving with the quiz data.
 */
export const getAdminQuiz = async (quizId: number): Promise<Quiz> => {
  try {
    const response = await apiClient.get<Quiz>(`/admin/quizzes/${quizId}`);
    console.log(`[API] Fetched quiz ${quizId}:`, response.data);
    return response.data;
  } catch (error: any) {
    console.error(`Failed to fetch admin quiz ${quizId}:`, error);
    throw error.response?.data || error;
  }
};

/**
 * Update an existing quiz.
 * @param quizId The ID of the quiz to update.
 * @param quizData The updated quiz data.
 * @returns Promise resolving with the updated quiz data from the backend.
 */
export const updateAdminQuiz = async (quizId: number, quizData: QuizInput): Promise<Quiz> => {
  try {
    const response = await apiClient.put<Quiz>(`/admin/quizzes/${quizId}`, quizData);
    return response.data;
  } catch (error) {
    console.error(`Error updating admin quiz ${quizId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to update quiz ${quizId}`);
    } else {
      throw new Error('An unexpected error occurred while updating the quiz.');
    }
  }
};

/**
 * Delete a specific quiz by its ID.
 * @param quizId The ID of the quiz to delete.
 * @returns Promise resolving when the quiz is successfully deleted.
 */
export const deleteAdminQuiz = async (quizId: number): Promise<void> => {
  try {
    await apiClient.delete(`/admin/quizzes/${quizId}`);
    console.log(`[API] Quiz ${quizId} deleted successfully.`);
  } catch (error) {
    console.error(`Error deleting admin quiz ${quizId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      // Extract backend error message if available
      throw new Error(error.response.data.error || `Failed to delete quiz ${quizId}`);
    } else {
      // Handle non-Axios errors or errors without a response
      throw new Error('An unexpected error occurred while deleting the quiz.');
    }
  }
};

// --- Credential Management API Functions ---

// Function to fetch credentials for a specific quiz
export const getAdminQuizCredentials = async (quizId: number): Promise<ResponderCredential[]> => {
  try {
    const response = await apiClient.get<ResponderCredential[]>(`/admin/quizzes/${quizId}/credentials`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching credentials for quiz ${quizId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to fetch credentials for quiz ${quizId}`);
    } else {
      throw new Error('An unexpected error occurred while fetching credentials.');
    }
  }
};

// Interface for the payload to generate credentials
interface GenerateCredentialsPayload {
  count: number;
  expiryHours?: number; // Optional expiry in hours
  username?: string; // Optional specific username
}

// Function to generate new credentials for a specific quiz
export const generateAdminQuizCredentials = async (
  quizId: number,
  payload: GenerateCredentialsPayload
): Promise<GenerateCredentialsResponse> => { // Update return type
  try {
    // Map frontend camelCase expiryHours to backend snake_case expiry_hours if present
    // Include username if provided
    const backendPayload = {
      count: payload.count,
      ...(payload.expiryHours !== undefined && { expiry_hours: payload.expiryHours }),
      ...(payload.username && { username: payload.username }), // Add username if present
    };
    console.log('[API] Generating credentials with payload:', backendPayload); // Log the payload being sent
    const response = await apiClient.post<GenerateCredentialsResponse>( // Expect single object response
      `/admin/quizzes/${quizId}/credentials`,
      backendPayload
    );
    return response.data; // Return the response data directly
  } catch (error) {
    console.error(`Error generating credentials for quiz ${quizId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to generate credentials for quiz ${quizId}`);
    } else {
      throw new Error('An unexpected error occurred while generating credentials.');
    }
  }
};

// --- Response Management API Functions ---

// Function to fetch response summaries for a specific quiz
export const getAdminQuizResponses = async (quizId: number): Promise<QuizResponseSummary[]> => {
  try {
    const response = await apiClient.get<QuizResponseSummary[]>(`/admin/quizzes/${quizId}/responses`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching responses for quiz ${quizId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to fetch responses for quiz ${quizId}`);
    } else {
      throw new Error('An unexpected error occurred while fetching responses.');
    }
  }
};

// Function to fetch detailed information for a specific response
export const getAdminResponseDetails = async (responseId: number): Promise<QuizResponseDetail> => {
  try {
    const response = await apiClient.get<QuizResponseDetail>(`/admin/responses/${responseId}`);
    return response.data;
  } catch (error) {
    console.error(`Error fetching details for response ${responseId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to fetch details for response ${responseId}`);
    } else {
      throw new Error('An unexpected error occurred while fetching response details.');
    }
  }
};

/**
 * Delete/revoke a specific credential by its ID.
 * @param credentialId The ID of the credential to delete.
 * @returns Promise resolving when the credential is successfully deleted.
 */
export const deleteCredential = async (credentialId: number): Promise<void> => {
  try {
    await apiClient.delete(`/admin/credentials/${credentialId}`);
    console.log(`[API] Credential ${credentialId} deleted successfully.`);
  } catch (error) {
    console.error(`Error deleting credential ${credentialId}:`, error);
    if (axios.isAxiosError(error) && error.response) {
      throw new Error(error.response.data.error || `Failed to delete credential ${credentialId}`);
    } else {
      throw new Error('An unexpected error occurred while deleting the credential.');
    }
  }
};

export default apiClient; // Export the configured instance if needed elsewhere

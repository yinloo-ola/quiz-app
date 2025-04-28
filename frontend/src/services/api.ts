import axios from 'axios';
import { useAuthStore } from '@/stores/auth'; // Use '@/stores' alias

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
    const token = authStore.token;

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

// Interface for the basic quiz data needed for creation
interface CreateQuizPayload {
  title: string;
  description?: string;
  // Add timeLimit if needed later
}

// Function to create a new admin quiz
export const createAdminQuiz = async (quizData: CreateQuizPayload) => {
  try {
    // Construct the full payload including a placeholder question
    // to satisfy backend validation (binding:"required,min=1,dive")
    const payload = {
      ...quizData,
      questions: [
        {
          text: 'Placeholder Question',
          choices: [{ text: 'Placeholder Choice', is_correct: true }],
        },
      ],
      // timeLimit: quizData.timeLimit // Add if implementing time limit
    };
    console.log('Sending create quiz payload:', payload);
    const response = await apiClient.post('/admin/quizzes', payload);
    return response.data; // Return the created quiz data from backend
  } catch (error) {
    console.error('Error creating admin quiz:', error);
    // Rethrow or handle error appropriately for the UI
    if (axios.isAxiosError(error) && error.response) {
      // Extract backend error message if available
      throw new Error(error.response.data.error || 'Failed to create quiz');
    } else {
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

// Add other API functions here (e.g., getQuizDetails, etc.)

export default apiClient; // Export the configured instance if needed elsewhere

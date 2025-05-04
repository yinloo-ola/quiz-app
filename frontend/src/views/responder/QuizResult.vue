<template>
  <div class="quiz-result-container">
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading results...</p>
    </div>
    
    <div v-else-if="error" class="error-container">
      <div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        {{ error }}
      </div>
      <button 
        @click="goHome" 
        class="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700"
      >
        Return to Home
      </button>
    </div>
    
    <div v-else class="result-content">
      <div class="result-header">
        <h1 class="text-3xl font-bold mb-2 text-gray-900">Quiz Results</h1>
        
        <div class="score-card">
          <div class="score-value">{{ Math.round(result.score) }}%</div>
          <div class="score-details">
            <p>{{ result.correct_answers }} correct out of {{ result.total_questions }} questions</p>
          </div>
        </div>
      </div>
      
      <div class="questions-review">
        <h2 class="text-2xl font-semibold mb-4 text-gray-900">Review Your Answers</h2>
        
        <div 
          v-for="(question, index) in quiz?.questions || []" 
          :key="question.id"
          class="question-card"
          :class="{ 
            'correct': isQuestionCorrect(question.id), 
            'incorrect': !isQuestionCorrect(question.id) 
          }"
        >
          <h3 class="text-xl font-medium mb-2 text-gray-900">
            Question {{ index + 1 }}: {{ question.text }}
          </h3>
          
          <div class="choices-container">
            <div 
              v-for="choice in question.choices" 
              :key="choice.id"
              class="choice-item"
              :class="{
                'selected': userSelectedChoice(question.id, choice.id),
                'correct-choice': isCorrectChoice(question.id, choice.id),
                'incorrect-choice': userSelectedChoice(question.id, choice.id) && !isCorrectChoice(question.id, choice.id)
              }"
            >
              <div class="choice-label">
                <span class="choice-indicator">
                  <svg v-if="userSelectedChoice(question.id, choice.id) && isCorrectChoice(question.id, choice.id)" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                  <svg v-else-if="userSelectedChoice(question.id, choice.id) && !isCorrectChoice(question.id, choice.id)" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                  <svg v-else-if="isCorrectChoice(question.id, choice.id)" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </span>
                <span class="choice-text text-gray-800">{{ choice.text }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="actions-container">
        <button 
          @click="logout" 
          class="logout-button"
        >
          Finish & Logout
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { getResponderQuiz } from '@/services/api';
import type { Quiz } from '@/types';

interface QuizSubmissionResult {
  score: number;
  total_questions: number;
  correct_answers: number;
  correct_choices: Record<number, number[]>; // Map of question_id to array of correct choice_ids
}

const props = defineProps<{
  id: string;
}>();

const router = useRouter();
const authStore = useAuthStore();

const quiz = ref<Quiz | null>(null);
const loading = ref(true);
const error = ref('');

// Get result from localStorage
const result = computed<QuizSubmissionResult>(() => {
  // Check localStorage
  const storedResult = localStorage.getItem('quiz_result');
  if (storedResult) {
    try {
      return JSON.parse(storedResult);
    } catch (e) {
      console.error('Error parsing stored result:', e);
    }
  }
  
  // Default empty result
  return {
    score: 0,
    total_questions: 0,
    correct_answers: 0,
    correct_choices: {}
  };
});

// User's selected answers (from localStorage)
const userAnswers = computed<Record<number, number[]>>(() => {
  const storedAnswers = localStorage.getItem('user_answers');
  if (storedAnswers) {
    try {
      return JSON.parse(storedAnswers);
    } catch (e) {
      console.error('Error parsing stored answers:', e);
    }
  }
  return {};
});

// Fetch quiz data for review
const fetchQuiz = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const quizId = parseInt(props.id);
    const quizData = await getResponderQuiz(quizId);
    quiz.value = quizData;
    
    // Result is already stored in localStorage by QuizTaker.vue
    
  } catch (err: any) {
    console.error('Error fetching quiz for review:', err);
    if (typeof err === 'object' && err.error) {
      error.value = err.error;
    } else {
      error.value = 'Failed to load quiz results. Please try again.';
    }
  } finally {
    loading.value = false;
  }
};

// Check if a question was answered correctly
const isQuestionCorrect = (questionId: number): boolean => {
  const correctChoices = result.value.correct_choices[questionId] || [];
  const userChoices = userAnswers.value[questionId] || [];
  
  // For single-choice questions
  if (quiz.value?.questions.find(q => q.id === questionId)?.type === 'single') {
    return userChoices.length === 1 && correctChoices.includes(userChoices[0]);
  }
  
  // For multi-choice questions, all correct choices must be selected and no incorrect ones
  if (userChoices.length !== correctChoices.length) return false;
  
  return userChoices.every(choice => correctChoices.includes(choice));
};

// Check if a choice was selected by the user
const userSelectedChoice = (questionId: number, choiceId: number): boolean => {
  return (userAnswers.value[questionId] || []).includes(choiceId);
};

// Check if a choice is a correct answer
const isCorrectChoice = (questionId: number, choiceId: number): boolean => {
  return (result.value.correct_choices[questionId] || []).includes(choiceId);
};

// Navigate to home
const goHome = () => {
  router.push({ name: 'home' });
};

// Logout and return to login
const logout = () => {
  // Clear result and answers from localStorage
  localStorage.removeItem('quiz_result');
  localStorage.removeItem('user_answers');
  
  // Logout from auth store
  authStore.logout();
  
  // Redirect to login
  router.push({ name: 'responder-login' });
};

// Fetch quiz data on component mount
onMounted(() => {
  fetchQuiz();
});
</script>

<style scoped>
.quiz-result-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 1rem;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border-left-color: #3b82f6;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.result-header {
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #e5e7eb;
}

.score-card {
  background-color: #f3f4f6;
  border-radius: 0.5rem;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  margin-top: 1rem;
}

.score-value {
  font-size: 2.5rem;
  font-weight: bold;
  color: #4f46e5;
  margin-right: 1.5rem;
}

.score-details {
  font-size: 1.125rem;
  color: #4b5563;
}

.questions-review {
  margin-bottom: 2rem;
}

.question-card {
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.question-card.correct {
  border-left: 4px solid #10b981;
}

.question-card.incorrect {
  border-left: 4px solid #ef4444;
}

.choices-container {
  margin-top: 1rem;
}

.choice-item {
  margin-bottom: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.25rem;
}

.choice-label {
  display: flex;
  align-items: center;
}

.choice-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.5rem;
  height: 1.5rem;
  margin-right: 0.75rem;
  color: #6b7280;
}

.choice-text {
  flex: 1;
}

.selected {
  background-color: #f3f4f6;
}

.correct-choice {
  background-color: #d1fae5;
}

.correct-choice .choice-indicator {
  color: #10b981;
}

.incorrect-choice {
  background-color: #fee2e2;
}

.incorrect-choice .choice-indicator {
  color: #ef4444;
}

.actions-container {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.logout-button {
  background-color: #4f46e5;
  color: white;
  font-weight: 500;
  padding: 0.75rem 2rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s;
}

.logout-button:hover {
  background-color: #4338ca;
}
</style>

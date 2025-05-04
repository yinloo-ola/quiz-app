<template>
  <div class="quiz-taker-container bg-gray-50 min-h-screen">
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p class="text-gray-700 font-medium mt-4">Loading quiz...</p>
    </div>

    <div v-else-if="error" class="error-container bg-white p-8 rounded-xl shadow-md">
      <div class="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 rounded-md mb-6">
        <p class="flex items-center">
          <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
          </svg>
          {{ error }}
        </p>
      </div>
      <button
        @click="fetchQuiz"
        class="bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition-colors duration-200 font-medium flex items-center justify-center w-full md:w-auto"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        Try Again
      </button>
    </div>

    <!-- Fixed timer that stays visible when scrolling (only shown when there's a time limit) -->
    <div v-if="quiz && quiz.timeLimit && quiz.timeLimit > 0" class="fixed-timer">
      <div class="fixed-timer-content">
        <div class="text-sm font-medium text-white mb-1 flex justify-between items-center">
          <span>Time Remaining:</span>
          <span class="text-xs opacity-80">Limit: {{ formatTimeLimit(quiz.timeLimit) }}</span>
        </div>
        <Timer
          :initial-time="quiz.timeLimit"
          @time-up="handleTimeUp"
          ref="timerRef"
        />
      </div>
    </div>

    <div class="quiz-content max-w-4xl mx-auto px-4 py-8">
      <div class="quiz-header bg-white p-6 rounded-xl shadow-md mb-8">
        <h1 class="text-3xl font-bold mb-3 text-indigo-900">{{ quiz?.title }}</h1>
        <p class="text-gray-700 text-lg">{{ quiz?.description }}</p>

        <div class="timers-container mt-6">
          <!-- Only keep the elapsed timer since the countdown timer is fixed at the top right -->
          <div class="elapsed-timer-container bg-blue-50 p-4 rounded-lg border border-blue-100">
            <div class="flex flex-col">
              <div class="text-sm font-medium text-blue-700 mb-2 flex items-center justify-between">
                <span>Your Time:</span>
                <span v-if="elapsedSeconds > 0" class="text-xs font-normal text-blue-600">
                  {{ getTimeMessage(elapsedSeconds) }}
                </span>
              </div>
              <ElapsedTimer
                :start-time="startTime"
                @update-elapsed="updateElapsedTime"
                ref="elapsedTimerRef"
              />
            </div>
          </div>
        </div>
      </div>

      <div class="questions-container space-y-8">
        <div
          v-for="(question, index) in quiz?.questions || []"
          :key="question.id"
          class="question-card bg-white p-6 rounded-xl shadow-md transition-all duration-200 hover:shadow-lg"
        >
          <h3 class="text-xl font-semibold mb-4 text-indigo-800 flex items-start">
            <span class="bg-indigo-100 text-indigo-800 rounded-full h-8 w-8 flex items-center justify-center mr-3 flex-shrink-0">{{ index + 1 }}</span>
            <span>{{ question.text }}</span>
          </h3>

          <div class="choices-container mt-5 space-y-3">
            <div
              v-for="choice in question.choices"
              :key="choice.id"
              class="choice-item"
            >
              <label
                :class="{
                  'choice-label': true,
                  'selected': isChoiceSelected(question.id, choice.id)
                }"
              >
                <input
                  v-if="question.type === 'single'"
                  type="radio"
                  :name="`question-${question.id}`"
                  :value="choice.id"
                  @change="selectSingleChoice(question.id, choice.id)"
                  :checked="isChoiceSelected(question.id, choice.id)"
                  class="h-5 w-5 text-indigo-600 mr-3"
                />
                <input
                  v-else
                  type="checkbox"
                  :name="`question-${question.id}`"
                  :value="choice.id"
                  @change="toggleMultiChoice(question.id, choice.id)"
                  class="h-5 w-5 text-indigo-600 mr-3"
                />
                <span class="choice-text text-gray-800 text-lg pl-4">{{ choice.text }}</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      <div class="submit-container mt-10 flex justify-center">
        <button
          @click="handleSubmit"
          :disabled="submitting"
          class="submit-button bg-indigo-600 hover:bg-indigo-700 text-white font-bold py-3 px-8 rounded-lg shadow-md hover:shadow-lg transition-all duration-200 text-lg flex items-center"
          :class="{'opacity-75 cursor-not-allowed': submitting}"
        >
          <span v-if="submitting" class="mr-2">
            <svg class="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </span>
          {{ submitting ? 'Submitting...' : 'Submit Quiz' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { getResponderQuiz, submitQuiz } from '@/services/api';
import type { Quiz } from '@/types';
import Timer from '@/components/Timer.vue';
import ElapsedTimer from '@/components/ElapsedTimer.vue';

// Define QuizSubmissionResult interface locally since it's not in types.ts
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
const authStore = useAuthStore(); // Used for authentication state
const timerRef = ref<InstanceType<typeof Timer> | null>(null);
const elapsedTimerRef = ref<InstanceType<typeof ElapsedTimer> | null>(null);

const quiz = ref<Quiz | null>(null);
const loading = ref(false);
const error = ref('');
const submitting = ref(false);
const selectedAnswers = ref<Record<number, number[]>>({});
const startTime = ref<string | undefined>(undefined); // Track when the quiz was started
const elapsedSeconds = ref(0); // Track elapsed time in seconds

// Fetch quiz data
const fetchQuiz = async () => {
  loading.value = true;
  error.value = '';
  try {
    // Fetch quiz data from API
    quiz.value = await getResponderQuiz(parseInt(props.id));
    
    // Record the start time
    startTime.value = new Date().toISOString();
    console.log('Quiz started at:', startTime.value);
    
    // Initialize selected answers
    if (quiz.value && quiz.value.questions) {
      quiz.value.questions.forEach(question => {
        // Initialize with empty array for each question
        selectedAnswers.value[question.id] = [];
      });
    }
  } catch (err: any) {
    console.error('Error fetching quiz:', err);
    if (typeof err === 'object' && err.error) {
      error.value = err.error;
    } else {
      error.value = 'Failed to load quiz. Please try again.';
    }
    quiz.value = null;
  } finally {
    loading.value = false;
  }
};

// Check if a choice is selected
const isChoiceSelected = (questionId: number, choiceId: number): boolean => {
  return selectedAnswers.value[questionId]?.includes(choiceId) || false;
};

// Handle single choice selection
const selectSingleChoice = (questionId: number, choiceId: number) => {
  selectedAnswers.value[questionId] = [choiceId];
};

// Handle multi-choice toggle
const toggleMultiChoice = (questionId: number, choiceId: number) => {
  if (!selectedAnswers.value[questionId]) {
    selectedAnswers.value[questionId] = [];
  }

  const index = selectedAnswers.value[questionId].indexOf(choiceId);
  if (index === -1) {
    // Add the choice
    selectedAnswers.value[questionId].push(choiceId);
  } else {
    // Remove the choice
    selectedAnswers.value[questionId].splice(index, 1);
  }
};

// Track elapsed time from the ElapsedTimer component
const updateElapsedTime = (seconds: number) => {
  elapsedSeconds.value = seconds;
};

// Get a contextual message based on elapsed time
const getTimeMessage = (seconds: number): string => {
  if (seconds < 60) {
    return "Just started";
  } else if (seconds < 300) { // 5 minutes
    return "Good pace";
  } else if (seconds < 600) { // 10 minutes
    return "Taking your time";
  } else if (seconds < 1200) { // 20 minutes
    return "Careful consideration";
  } else {
    return "Deep thinking";
  }
};

// Format time limit in a user-friendly way
const formatTimeLimit = (seconds: number): string => {
  if (!seconds) return "None";
  
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  
  if (hours > 0) {
    return `${hours}h ${minutes > 0 ? minutes + 'm' : ''}`;
  } else if (minutes > 0) {
    return `${minutes}m`;
  } else {
    return `${seconds}s`;
  }
};

// Handle quiz submission
const handleSubmit = async () => {
  if (!quiz.value) return;

  submitting.value = true;

  try {
    // Format answers for submission
    const answers = Object.entries(selectedAnswers.value).map(([questionId, choiceIds]) => ({
      question_id: parseInt(questionId),
      choice_ids: choiceIds
    }));

    // Create submission payload with start time
    const submissionPayload: {
      answers: { question_id: number; choice_ids: number[] }[];
      started_at?: string;
    } = {
      answers
    };
    
    // Only include started_at if it's not undefined
    if (startTime.value) {
      submissionPayload.started_at = startTime.value;
    }

    // Submit the quiz
    const result = await submitQuiz(parseInt(props.id), submissionPayload);

    // Stop both timers if they're running
    if (timerRef.value) {
      timerRef.value.stopTimer();
    }
    if (elapsedTimerRef.value) {
      elapsedTimerRef.value.stopTimer();
    }

    // Navigate to results page with the result data
    // Store the result in localStorage since router state isn't fully supported in Vue Router 4
    localStorage.setItem('quiz_result', JSON.stringify(result));
    localStorage.setItem('user_answers', JSON.stringify(selectedAnswers.value));

    router.push({
      name: 'quiz-result',
      params: { id: props.id }
    });

  } catch (err: any) {
    console.error('Error submitting quiz:', err);
    if (typeof err === 'object' && err.error) {
      error.value = err.error;
    } else {
      error.value = 'Failed to submit quiz. Please try again.';
    }
    submitting.value = false;
  }
};

// Handle timer expiry
const handleTimeUp = () => {
  // Auto-submit when time is up
  handleSubmit();
};

// Fetch quiz data on component mount
onMounted(() => {
  fetchQuiz();
});
</script>

<style scoped>
/* Fixed timer styles */
.fixed-timer {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  max-width: 200px;
}

.fixed-timer-content {
  background-color: rgba(79, 70, 229, 0.9);
  border-radius: 8px;
  padding: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(4px);
  border: 1px solid rgba(165, 180, 252, 0.5);
}

/* Override the timer component styles when in fixed position */
.fixed-timer :deep(.time-display) {
  background-color: transparent !important;
  color: white !important;
  font-size: 1.75rem;
  padding: 0;
  box-shadow: none;
}

/* Pulse animation for danger state in fixed timer */
.fixed-timer :deep(.time-display.danger) {
  background-color: transparent !important;
  color: #fecaca !important;
  animation: pulse 1s infinite;
}

/* Warning state in fixed timer */
.fixed-timer :deep(.time-display.warning) {
  background-color: transparent !important;
  color: #fef3c7 !important;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 2rem;
}

.spinner {
  border: 4px solid rgba(79, 70, 229, 0.1);
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border-left-color: #4f46e5;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Choice styling */
.choice-label {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid #e5e7eb;
  background-color: #f9fafb;
}

.choice-label:hover {
  background-color: #f3f4f6;
  border-color: #d1d5db;
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.choice-label.selected {
  background-color: #e0e7ff;
  border-color: #a5b4fc;
  box-shadow: 0 0 0 1px #c7d2fe;
}

/* Improve radio and checkbox appearance */
input[type="radio"], input[type="checkbox"] {
  appearance: none;
  background-color: #fff;
  margin: 0;
  font: inherit;
  color: currentColor;
  width: 1.25em;
  height: 1.25em;
  border: 1px solid #d1d5db;
  border-radius: 50%;
  transform: translateY(-0.075em);
  display: grid;
  place-content: center;
}

input[type="checkbox"] {
  border-radius: 0.25em;
}

input[type="radio"]:checked,
input[type="checkbox"]:checked {
  background-color: #4f46e5;
  border-color: #4f46e5;
}

input[type="radio"]:checked::before {
  content: "";
  width: 0.5em;
  height: 0.5em;
  border-radius: 50%;
  background-color: white;
  transform: scale(1);
  transition: 120ms transform ease-in-out;
  box-shadow: inset 1em 1em white;
}

input[type="checkbox"]:checked::before {
  content: "";
  width: 0.5em;
  height: 0.5em;
  transform: scale(1);
  transform-origin: center;
  clip-path: polygon(14% 44%, 0 65%, 50% 100%, 100% 16%, 80% 0%, 43% 62%);
  background-color: white;
}

/* Responsive design improvements */
@media (max-width: 640px) {
  .quiz-content {
    padding: 1rem;
  }

  .question-card {
    padding: 1.25rem;
  }

  .choice-label {
    padding: 0.5rem 0.75rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  h3 {
    font-size: 1.125rem;
  }
}
</style>

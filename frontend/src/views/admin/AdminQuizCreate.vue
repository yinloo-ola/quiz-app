<template>
  <div class="p-6">
    <h1 class="text-2xl font-semibold mb-6 text-gray-100">Create New Quiz</h1>

    <form @submit.prevent="createQuiz">
      <!-- Title Field -->
      <div class="mb-4">
        <label for="title" class="block text-sm font-medium text-gray-300 mb-1">Title</label>
        <input
          type="text"
          id="title"
          v-model="title"
          required
          class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400"
          placeholder="Enter quiz title"
        />
        <div v-if="titleError" class="text-red-500 text-sm mt-1">{{ titleError }}</div>
      </div>

      <!-- Description Field -->
      <div class="mb-6">
        <label for="description" class="block text-sm font-medium text-gray-300 mb-1">Description (Optional)</label>
        <textarea
          id="description"
          v-model="description"
          rows="3"
          class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400"
          placeholder="Enter a brief description for the quiz"
        ></textarea>
      </div>

      <!-- Time Limit Field -->
      <div class="mb-4">
        <label for="timeLimit" class="block text-sm font-medium text-gray-300 mb-1">Time Limit (seconds)</label>
        <input
          type="number"
          id="timeLimit"
          v-model.number="timeLimitSeconds"
          min="0"
          class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400"
          placeholder="0 for no limit"
        />
        <div v-if="timeLimitError" class="text-red-500 text-sm mt-1">{{ timeLimitError }}</div>
      </div>

      <!-- Questions Section -->
      <div class="mt-6 border-t border-slate-700 pt-6">
        <h2 class="text-xl font-semibold mb-4 text-gray-100">Questions</h2>
        <div v-if="questions.length === 0" class="text-gray-500 italic mb-4">
          No questions added yet.
        </div>
        <div v-if="questionsError" class="text-red-500 text-sm mb-4">{{ questionsError }}</div>

        <div v-for="(question, index) in questions" :key="index" class="mb-6 p-4 border border-slate-700 rounded-md bg-slate-800">
          <div class="flex justify-between items-center mb-3">
            <h3 class="font-medium text-lg text-gray-100">Question {{ index + 1 }}</h3>
            <button
              type="button"
              @click="removeQuestion(index)"
              class="text-red-400 hover:text-red-300 font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-red-500"
            >
              Remove Question
            </button>
          </div>

          <!-- Question Text -->
          <div class="mb-3">
            <label :for="`question-text-${index}`" class="block text-sm font-medium text-gray-300 mb-1">Question Text</label>
            <textarea
              :id="`question-text-${index}`"
              v-model="question.text"
              rows="2"
              required
              class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400"
              placeholder="Enter the question"
            ></textarea>
            <div v-if="questionErrors[index]" class="text-red-500 text-sm mt-1">{{ questionErrors[index] }}</div>
          </div>

          <!-- Question Type -->
          <div class="mb-3">
            <label class="block text-sm font-medium text-gray-300 mb-1">Question Type</label>
            <div class="flex items-center space-x-4">
              <label :for="`q-type-single-${index}`" class="flex items-center cursor-pointer">
                <input
                  type="radio"
                  :id="`q-type-single-${index}`"
                  :name="`question-type-${index}`"
                  value="single"
                  v-model="question.type"
                  class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-600 bg-slate-700"
                />
                <span class="ml-2 text-sm text-gray-300">Single Choice (Correct answer is one)</span>
              </label>
              <label :for="`q-type-multi-${index}`" class="flex items-center cursor-pointer">
                <input
                  type="radio"
                  :id="`q-type-multi-${index}`"
                  :name="`question-type-${index}`"
                  value="multi"
                  v-model="question.type"
                  class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-600 bg-slate-700"
                />
                <span class="ml-2 text-sm text-gray-300">Multiple Choice (Correct answer can be many)</span>
              </label>
            </div>
          </div>

          <!-- Choices Section -->
          <div class="mt-4 pl-4 border-l-2 border-slate-700">
            <h4 class="text-md font-semibold mb-2 text-gray-100">Choices</h4>
            <div v-if="question.choices.length === 0" class="text-sm text-gray-500 italic mb-2">
              Add at least one choice.
            </div>

            <div v-for="(choice, choiceIndex) in question.choices" :key="choiceIndex" class="flex items-center space-x-3 mb-2">
              <!-- Choice Text Input -->
              <input
                type="text"
                v-model="choice.text"
                required
                class="flex-grow mt-1 block w-full px-3 py-1.5 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400"
                placeholder="Enter choice text"
              />
              <!-- Correct Checkbox -->
              <label :for="`choice-correct-${index}-${choiceIndex}`" class="flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  :id="`choice-correct-${index}-${choiceIndex}`"
                  v-model="choice.is_correct"
                  class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-600 rounded bg-slate-700"
                />
                <span class="ml-2 text-sm text-gray-300">Correct</span>
              </label>
              <!-- Remove Choice Button -->
              <button
                type="button"
                @click="removeChoice(index, choiceIndex)"
                class="text-red-500 hover:text-red-400 p-1 rounded hover:bg-slate-700 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-offset-slate-700 focus:ring-red-500"
                title="Remove choice"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <button
              type="button"
              @click="addChoice(index)"
              class="mt-1 text-sm bg-slate-600 hover:bg-slate-700 text-gray-200 font-bold py-1 px-3 rounded focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-green-500"
            >
              Add Choice
            </button>
          </div>

        </div>
        <button
          type="button"
          @click="addQuestion"
          class="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-blue-500"
        >
          Add Question
        </button>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-end space-x-3 mt-8 border-t border-slate-700 pt-6">
        <button
          type="button"
          @click="cancel"
          class="bg-slate-600 hover:bg-slate-700 text-gray-200 font-bold py-2 px-4 rounded focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-slate-500"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="isLoading"
          class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded disabled:bg-teal-800 disabled:text-gray-400 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-teal-500"
        >
          {{ isLoading ? 'Saving...' : 'Save Quiz' }}
        </button>
      </div>
    </form>

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { createAdminQuiz } from '@/services/api';
import type { QuestionType, ChoiceInput, QuestionInput, QuizCreatePayload } from '@/types'; 

const router = useRouter();

// Define reactive variables for form fields
const title = ref('');
const description = ref('');
const timeLimitSeconds = ref<number | null>(null); // Time limit in seconds
const questions = ref<QuestionInput[]>([]); // Array to hold questions

const isLoading = ref(false);
const error = ref<string | null>(null);

// Validation errors
const titleError = ref<string | null>(null);
const timeLimitError = ref<string | null>(null);
const questionsError = ref<string | null>(null); // General error for questions structure
const questionErrors = ref<(string | null)[]>([]); // Errors for specific questions
const choiceErrors = ref<({ questionIndex: number; choiceIndex: number; message: string } | null)[]>([]); // Detailed choice errors

// --- Question Management ---
const addQuestion = () => {
  questions.value.push({
    text: '',
    type: 'single' as QuestionType, // Default type, ensure type safety
    choices: [],     // Start with empty choices
  });
};

const removeQuestion = (index: number) => {
  if (confirm(`Are you sure you want to remove Question ${index + 1}?`)) {
    questions.value.splice(index, 1);
  }
};

// --- Choice Management ---
const addChoice = (questionIndex: number) => {
  questions.value[questionIndex].choices.push({
    text: '',
    is_correct: false,
  } as ChoiceInput); // Add type assertion
};

const removeChoice = (questionIndex: number, choiceIndex: number) => {
  questions.value[questionIndex].choices.splice(choiceIndex, 1);
};

// Function to handle form submission
const createQuiz = async () => {
  isLoading.value = true;
  error.value = null;
  console.log('Form submitted');
  console.log('Title:', title.value);
  console.log('Description:', description.value);

  // Clear previous errors
  titleError.value = null;
  timeLimitError.value = null;
  questionsError.value = null;
  questionErrors.value = [];
  choiceErrors.value = [];

  // ** Perform Validation **
  if (!validateQuiz()) {
    isLoading.value = false;
    return; // Stop submission if validation fails
  }

  try {
    const newQuizData: QuizCreatePayload = {
      title: title.value,
      description: description.value || undefined,
      time_limit_seconds: timeLimitSeconds.value === null || timeLimitSeconds.value <= 0 ? undefined : timeLimitSeconds.value,
      questions: questions.value,
    };
    const createdQuiz = await createAdminQuiz(newQuizData);
    console.log('Quiz created successfully:', createdQuiz);
    alert('Quiz created successfully!');
    router.push({ name: 'admin-quiz-list' });

  } catch (err: any) {
    console.error('Failed to create quiz:', err);
    error.value = err.message || 'An unknown error occurred.';
    alert(`Error creating quiz: ${error.value}`);

  } finally {
    isLoading.value = false;
  }
};

// Function to cancel and go back
const cancel = () => {
  if (confirm('Are you sure you want to cancel? Any unsaved changes will be lost.')) {
    router.push({ name: 'admin-quiz-list' });
  }
};

// --- Validation ---
const validateQuiz = (): boolean => {
  let isValid = true;

  // Validate Title
  if (!title.value.trim()) {
    titleError.value = 'Quiz title is required.';
    isValid = false;
  }

  // Validate Time Limit
  if (timeLimitSeconds.value === null || timeLimitSeconds.value === undefined) {
    timeLimitError.value = 'Time limit is required.';
    isValid = false;
  } else if (typeof timeLimitSeconds.value !== 'number' || timeLimitSeconds.value <= 0 || !Number.isInteger(timeLimitSeconds.value)) {
    timeLimitError.value = 'Time limit must be a positive whole number (seconds).';
    isValid = false;
  }

  // Validate Questions structure
  if (questions.value.length === 0) {
    questionsError.value = 'A quiz must have at least one question.';
    isValid = false;
  }

  // Initialize/Reset question-specific errors
  questionErrors.value = Array(questions.value.length).fill(null);

  questions.value.forEach((question, qIndex) => {
    let correctChoiceCount = 0;

    // Validate Question Text
    if (!question.text.trim()) {
      questionErrors.value[qIndex] = 'Question text cannot be empty.';
      isValid = false;
    }

    // Validate Choices
    if (question.choices.length < 2) {
      questionErrors.value[qIndex] = (questionErrors.value[qIndex] ? questionErrors.value[qIndex] + ' ' : '') + 'Each question must have at least two choices.';
      isValid = false;
    }

    question.choices.forEach((choice, cIndex) => {
      // Validate Choice Text
      if (!choice.text.trim()) {
        // Add specific choice error - implementation needed later for display
        questionErrors.value[qIndex] = (questionErrors.value[qIndex] ? questionErrors.value[qIndex] + ' ' : '') + `Choice ${cIndex + 1} text cannot be empty.`;
        isValid = false;
      }
      if (choice.is_correct) {
        correctChoiceCount++;
      }
    });

    // Validate Correct Choice Count
    if (correctChoiceCount === 0) {
      questionErrors.value[qIndex] = (questionErrors.value[qIndex] ? questionErrors.value[qIndex] + ' ' : '') + 'At least one choice must be marked as correct.';
      isValid = false;
    } else if (question.type !== 'multi' && correctChoiceCount > 1) { 
      questionErrors.value[qIndex] = (questionErrors.value[qIndex] ? questionErrors.value[qIndex] + ' ' : '') + 'Only one choice can be marked as correct for single-choice questions.';
      isValid = false;
    }

  });

  return isValid;
};

</script>

<style scoped>
/* Add custom styles if needed */
input[type='radio']:focus,
input[type='checkbox']:focus {
  box-shadow: 0 0 0 2px theme('colors.teal.500'); /* Custom focus ring for radio/checkbox */
}

/* Removed scoped override for button focus - handled by global reset */
</style>

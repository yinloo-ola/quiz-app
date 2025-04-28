<template>
  <div class="p-6">
    <h1 class="text-2xl font-semibold mb-6">Create New Quiz</h1>

    <form @submit.prevent="createQuiz">
      <!-- Title Field -->
      <div class="mb-4">
        <label for="title" class="block text-sm font-medium text-gray-700 mb-1">Title</label>
        <input
          type="text"
          id="title"
          v-model="title"
          required
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          placeholder="Enter quiz title"
        />
      </div>

      <!-- Description Field -->
      <div class="mb-6">
        <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description (Optional)</label>
        <textarea
          id="description"
          v-model="description"
          rows="3"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          placeholder="Enter a brief description for the quiz"
        ></textarea>
      </div>

      <!-- Time Limit Field -->
      <div class="mb-4">
        <label for="timeLimit" class="block text-sm font-medium text-gray-700 mb-1">Time Limit (seconds)</label>
        <input
          type="number"
          id="timeLimit"
          v-model.number="timeLimitSeconds"
          min="0"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          placeholder="0 for no limit"
        />
      </div>

      <!-- Questions Section -->
      <div class="mt-6 border-t pt-6">
        <h2 class="text-xl font-semibold mb-4">Questions</h2>
        <div v-if="questions.length === 0" class="text-gray-500 italic mb-4">
          No questions added yet.
        </div>

        <div v-for="(question, index) in questions" :key="index" class="mb-6 p-4 border rounded-md bg-gray-50">
          <div class="flex justify-between items-center mb-3">
            <h3 class="font-medium text-lg">Question {{ index + 1 }}</h3>
            <button
              type="button"
              @click="removeQuestion(index)"
              class="text-red-500 hover:text-red-700 font-medium"
            >
              Remove Question
            </button>
          </div>

          <!-- Question Text -->
          <div class="mb-3">
            <label :for="`question-text-${index}`" class="block text-sm font-medium text-gray-700 mb-1">Question Text</label>
            <textarea
              :id="`question-text-${index}`"
              v-model="question.text"
              rows="2"
              required
              class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
              placeholder="Enter the question"
            ></textarea>
          </div>

          <!-- Question Type -->
          <div class="mb-3">
            <label class="block text-sm font-medium text-gray-700 mb-1">Question Type</label>
            <div class="flex items-center space-x-4">
              <label :for="`q-type-single-${index}`" class="flex items-center">
                <input
                  type="radio"
                  :id="`q-type-single-${index}`"
                  :name="`question-type-${index}`"
                  value="single"
                  v-model="question.type"
                  class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300"
                />
                <span class="ml-2 text-sm text-gray-700">Single Choice (Correct answer is one)</span>
              </label>
              <label :for="`q-type-multi-${index}`" class="flex items-center">
                <input
                  type="radio"
                  :id="`q-type-multi-${index}`"
                  :name="`question-type-${index}`"
                  value="multi"
                  v-model="question.type"
                  class="focus:ring-indigo-500 h-4 w-4 text-indigo-600 border-gray-300"
                />
                <span class="ml-2 text-sm text-gray-700">Multiple Choice (Correct answer can be many)</span>
              </label>
            </div>
          </div>

          <!-- Choices Section -->
          <div class="mt-4 pl-4 border-l-2 border-gray-200">
            <h4 class="text-md font-semibold mb-2">Choices</h4>
            <div v-if="question.choices.length === 0" class="text-sm text-gray-500 italic mb-2">
              Add at least one choice.
            </div>

            <div v-for="(choice, choiceIndex) in question.choices" :key="choiceIndex" class="flex items-center space-x-3 mb-2">
              <!-- Choice Text Input -->
              <input
                type="text"
                v-model="choice.text"
                required
                class="flex-grow px-2 py-1 border border-gray-300 rounded-md shadow-sm text-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                placeholder="Enter choice text"
              />
              <!-- Is Correct Checkbox -->
              <label class="flex items-center space-x-1.5 cursor-pointer">
                <input
                  type="checkbox"
                  v-model="choice.is_correct"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <span class="text-sm font-medium text-gray-700">Correct</span>
              </label>
              <!-- Remove Choice Button -->
              <button
                type="button"
                @click="removeChoice(index, choiceIndex)"
                class="text-red-500 hover:text-red-700 text-sm font-medium p-1 rounded-full hover:bg-red-100"
                title="Remove Choice"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <button
              type="button"
              @click="addChoice(index)"
              class="mt-1 text-sm bg-gray-200 hover:bg-gray-300 text-gray-700 font-bold py-1 px-3 rounded"
            >
              Add Choice
            </button>
          </div>

        </div>
        <button
          type="button"
          @click="addQuestion"
          class="mt-2 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Add Question
        </button>
      </div>

      <!-- Action Buttons -->
      <div class="flex justify-end space-x-3 mt-8 border-t pt-6">
        <button
          type="button"
          @click="cancel"
          class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="isLoading"
          class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isLoading ? 'Saving...' : 'Save Quiz (Basic)' }}
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

</script>

<style scoped>
/* Add component-specific styles here */
</style>

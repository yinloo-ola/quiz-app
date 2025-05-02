<template>
  <form @submit.prevent="handleSubmit">
    <!-- Title Field -->
    <div class="mb-4">
      <label for="title" class="block text-sm font-medium text-gray-300 mb-1">Title</label>
      <input
        type="text"
        id="title"
        v-model="title"
        required
        :disabled="isLoading"
        class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
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
        :disabled="isLoading"
        class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
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
        :disabled="isLoading"
        class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
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

      <div v-for="(question, index) in questions" :key="question.id || index" class="mb-6 p-4 border border-slate-700 rounded-md bg-slate-800">
        <div class="flex justify-between items-center mb-3">
          <h3 class="font-medium text-lg text-gray-100">Question {{ index + 1 }}</h3>
          <button
            type="button"
            @click="removeQuestion(index)"
            :disabled="isLoading"
            class="text-red-400 hover:text-red-300 font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
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
            :disabled="isLoading"
            class="mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
            placeholder="Enter the question"
          ></textarea>
          <div v-if="questionErrors[index]" class="text-red-500 text-sm mt-1">{{ questionErrors[index]?.text }}</div>
        </div>

        <!-- Question Type -->
        <div class="mb-3">
          <label class="block text-sm font-medium text-gray-300 mb-1">Question Type</label>
          <div class="flex items-center space-x-4">
            <label :for="`q-type-single-${index}`" class="flex items-center cursor-pointer">
              <input
                type="radio"
                :id="`q-type-single-${index}`"
                value="single"
                v-model="question.type"
                :name="`question-type-${index}`"
                :disabled="isLoading"
                @change="handleChoiceCorrectness(index)"
                class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
              />
              <span class="ml-2 text-sm text-gray-300">Single Choice</span>
            </label>
            <label :for="`q-type-multi-${index}`" class="flex items-center cursor-pointer">
              <input
                type="radio"
                :id="`q-type-multi-${index}`"
                value="multi"
                v-model="question.type"
                :name="`question-type-${index}`"
                :disabled="isLoading"
                class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
              />
              <span class="ml-2 text-sm text-gray-300">Multiple Choice</span>
            </label>
          </div>
        </div>

        <!-- Choices Section -->
        <div class="mt-4 pl-4 border-l-2 border-slate-700">
          <h4 class="text-md font-semibold mb-2 text-gray-100">Choices</h4>
          <div v-if="question.choices.length === 0" class="text-sm text-gray-500 italic mb-2">
            Add at least one choice.
          </div>
          <div v-if="questionErrors[index]?.choices" class="text-red-500 text-sm mb-2">{{ questionErrors[index]?.choices }}</div>

          <div v-for="(choice, choiceIndex) in question.choices" :key="choice.id || choiceIndex" class="flex items-center space-x-3 mb-2">
            <!-- Choice Text Input -->
            <input
              type="text"
              v-model="choice.text"
              required
              :disabled="isLoading"
              class="flex-grow mt-1 block w-full px-3 py-1.5 border border-slate-700 rounded-md bg-slate-600 text-gray-100 focus:outline-none focus:ring-teal-500 focus:border-teal-500 sm:text-sm placeholder-gray-400 disabled:opacity-50 disabled:cursor-not-allowed"
              placeholder="Enter choice text"
            />

            <!-- Correct Choice Checkbox/Radio -->
            <label :for="`choice-correct-${index}-${choiceIndex}`" class="flex items-center space-x-2 cursor-pointer">
              <input
                :type="question.type === 'single' ? 'radio' : 'checkbox'"
                :id="`choice-correct-${index}-${choiceIndex}`"
                :name="`correct-choice-${index}`" 
                :value="choiceIndex" 
                v-model="choice.is_correct"
                :checked="choice.is_correct"
                :disabled="isLoading"
                @change="handleChoiceCorrectness(index, choiceIndex)"
                class="focus:ring-teal-500 h-4 w-4 text-teal-600 border-gray-500 rounded disabled:opacity-50 disabled:cursor-not-allowed"
              />
              <span class="text-sm text-gray-300">Correct</span>
            </label>

            <!-- Remove Choice Button -->
            <button
              type="button"
              @click="removeChoice(index, choiceIndex)"
              :disabled="isLoading || question.choices.length <= 1"
              class="text-red-500 hover:text-red-400 p-1 rounded hover:bg-slate-700 focus:outline-none focus:ring-1 focus:ring-offset-1 focus:ring-offset-slate-700 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
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
            :disabled="isLoading"
            class="mt-1 text-sm bg-slate-600 hover:bg-slate-700 text-gray-200 font-bold py-1 px-3 rounded focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Add Choice
          </button>
        </div>

      </div>
      <button
        type="button"
        @click="addQuestion"
        :disabled="isLoading"
        class="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Add Question
      </button>
    </div>

    <!-- Global Error Display -->
    <div v-if="globalError" class="mt-4 p-3 bg-red-900 border border-red-700 text-red-200 rounded text-sm">
      {{ globalError }}
    </div>

    <!-- Action Buttons -->
    <div class="flex justify-end space-x-3 mt-8 border-t border-slate-700 pt-6">
      <button
        type="button"
        @click="handleCancel"
        :disabled="isLoading"
        class="px-4 py-2 border border-slate-600 rounded text-gray-300 hover:bg-slate-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-slate-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Cancel
      </button>
      <button
        type="submit"
        :disabled="isLoading"
        class="px-4 py-2 bg-teal-600 text-white rounded hover:bg-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-teal-500 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {{ isLoading ? 'Saving...' : submitButtonText }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import type { Quiz, Question, Choice, QuestionType, QuestionInput, ChoiceInput, QuizCreatePayload } from '@/types';

// --- Props --- 

interface Props {
  initialQuizData?: Quiz | null; 
  isLoading?: boolean; 
  submitButtonText?: string;
  globalError?: string | null; 
}

const props = withDefaults(defineProps<Props>(), {
  initialQuizData: null,
  isLoading: false,
  submitButtonText: 'Save Quiz',
  globalError: null,
});

// --- Emits --- 

const emit = defineEmits<{ 
  (e: 'submit', quizData: QuizCreatePayload): void
  (e: 'cancel'): void
}>();

// --- Form State --- 

const title = ref('');
const description = ref('');
const timeLimitSeconds = ref<number | null>(0);
const questions = ref<QuestionInput[]>([]); 

// --- Validation Errors --- 
const titleError = ref<string | null>(null);
const timeLimitError = ref<string | null>(null);
const questionsError = ref<string | null>(null); 
interface QuestionValidationError {
  text?: string | null;
  choices?: string | null; 
  choiceText?: (string | null)[]; 
}
const questionErrors = ref<QuestionValidationError[]>([]);

// --- Utility Functions (Add/Remove Questions/Choices) --- 

const addQuestion = () => {
  if (props.isLoading) return;
  questions.value.push({
    text: '',
    type: 'single', 
    choices: [{ text: '', is_correct: false }, { text: '', is_correct: false }], 
  });
  questionsError.value = null; 
};

const removeQuestion = (index: number) => {
  if (props.isLoading) return;
  if (confirm('Are you sure you want to remove this question and its choices?')) {
    questions.value.splice(index, 1);
    questionErrors.value.splice(index, 1); 
  }
};

const addChoice = (questionIndex: number) => {
  if (props.isLoading) return;
  questions.value[questionIndex].choices.push({ text: '', is_correct: false }); 
  if (questionErrors.value[questionIndex]?.choices?.includes('at least two choices')) {
    questionErrors.value[questionIndex]!.choices = null; 
  }
};

const removeChoice = (questionIndex: number, choiceIndex: number) => {
  if (props.isLoading) return;
  if (questions.value[questionIndex].choices.length <= 1) {
    alert('A question must have at least one choice.');
    return;
  }
  if (confirm('Are you sure you want to remove this choice?')) {
    questions.value[questionIndex].choices.splice(choiceIndex, 1);
    handleChoiceCorrectness(questionIndex); 
  }
};

const handleChoiceCorrectness = (questionIndex: number, changedChoiceIndex?: number) => {
  const question = questions.value[questionIndex];
  if (!question) return;

  if (question.type === 'single') {
    question.choices.forEach((choice: ChoiceInput) => { 
      if (changedChoiceIndex !== undefined) {
        choice.is_correct = (choice === question.choices[changedChoiceIndex]);
      } 
      else if (choice.is_correct && question.choices.findIndex((c: ChoiceInput) => c.is_correct) !== question.choices.indexOf(choice)) { 
         choice.is_correct = false;
      }
    });
  }

  if (questionErrors.value[questionIndex]?.choices?.includes('least one choice must be marked as correct')) {
    questionErrors.value[questionIndex]!.choices = null;
  } 
  if (questionErrors.value[questionIndex]?.choices?.includes('Only one choice can be marked as correct')) {
     questionErrors.value[questionIndex]!.choices = null;
  }
};

const handleSubmit = () => {
  if (props.isLoading) return;
  titleError.value = null;
  timeLimitError.value = null;
  questionsError.value = null;
  questionErrors.value = [];

  if (validateQuiz()) {
    const quizData: QuizCreatePayload = {
      title: title.value,
      description: description.value || undefined,
      time_limit_seconds: timeLimitSeconds.value === null || timeLimitSeconds.value <= 0 ? undefined : timeLimitSeconds.value,
      questions: questions.value.map((q: QuestionInput) => ({ 
        id: q.id, 
        text: q.text,
        type: q.type,
        choices: q.choices.map((c: ChoiceInput) => ({ 
          id: c.id, 
          text: c.text,
          is_correct: !!c.is_correct 
        }))
      })),
    };
    emit('submit', quizData);
  }
};

const handleCancel = () => {
  if (!props.isLoading) {
    emit('cancel');
  }
};

const validateQuiz = (): boolean => {
  let isValid = true;
  questionErrors.value = Array(questions.value.length).fill({}).map(() => ({ text: null, choices: null, choiceText: [] })); 

  if (!title.value.trim()) {
    titleError.value = 'Quiz title is required.';
    isValid = false;
  }

  const timeLimit = timeLimitSeconds.value;
  if (timeLimit === null || timeLimit === undefined || typeof timeLimit !== 'number' || timeLimit < 0 || !Number.isInteger(timeLimit)) {
      timeLimitError.value = 'Time limit must be a non-negative whole number (0 for no limit).';
      isValid = false;
  }

  if (questions.value.length === 0) {
    questionsError.value = 'A quiz must have at least one question.';
    isValid = false;
  }

  questions.value.forEach((question: QuestionInput, qIndex: number) => { 
    let correctChoiceCount = 0;
    const currentQuestionErrors: QuestionValidationError = { text: null, choices: null, choiceText: Array(question.choices.length).fill(null) };

    if (!question.text.trim()) {
      currentQuestionErrors.text = 'Question text cannot be empty.';
      isValid = false;
    }

    if (question.choices.length < 2) {
      currentQuestionErrors.choices = (currentQuestionErrors.choices ? currentQuestionErrors.choices + ' ' : '') + 'Each question must have at least two choices.';
      isValid = false;
    }

    let hasEmptyChoice = false;
    question.choices.forEach((choice: ChoiceInput) => { 
      if (!choice.text.trim()) {
        hasEmptyChoice = true;
        isValid = false;
      }
      if (choice.is_correct) {
        correctChoiceCount++;
      }
    });

    if (hasEmptyChoice) {
        currentQuestionErrors.choices = (currentQuestionErrors.choices ? currentQuestionErrors.choices + ' ' : '') + 'Choice text cannot be empty.';
    }

    if (correctChoiceCount === 0) {
      currentQuestionErrors.choices = (currentQuestionErrors.choices ? currentQuestionErrors.choices + ' ' : '') + 'At least one choice must be marked as correct.';
      isValid = false;
    } else if (question.type === 'single' && correctChoiceCount > 1) { 
      currentQuestionErrors.choices = (currentQuestionErrors.choices ? currentQuestionErrors.choices + ' ' : '') + 'Only one choice can be marked as correct for single-choice questions.';
      isValid = false;
    }

    questionErrors.value[qIndex] = currentQuestionErrors;

  });

  return isValid;
};

watch(() => props.initialQuizData, (newData: Quiz | null | undefined) => { 
  if (newData) {
    console.log("QuizForm received initial data:", JSON.parse(JSON.stringify(newData)));
    title.value = newData.title || '';
    description.value = newData.description || '';
    timeLimitSeconds.value = newData.time_limit_seconds === undefined || newData.time_limit_seconds === null ? 0 : newData.time_limit_seconds;
    questions.value = JSON.parse(JSON.stringify(newData.questions || []));
    questions.value.forEach((q: QuestionInput) => { 
      q.choices.forEach((c: ChoiceInput) => { c.is_correct = !!c.is_correct; }); 
    });
     titleError.value = null;
    timeLimitError.value = null;
    questionsError.value = null;
    questionErrors.value = [];
  } else {
    title.value = '';
    description.value = '';
    timeLimitSeconds.value = 0;
    questions.value = [];
    titleError.value = null;
    timeLimitError.value = null;
    questionsError.value = null;
    questionErrors.value = [];
  }
}, { immediate: true, deep: true });

</script>

<style scoped>
input[type='radio']:focus,
input[type='checkbox']:focus {
  box-shadow: 0 0 0 2px theme('colors.teal.500'); 
}

/* Removed scoped override for button focus - handled by global reset */
</style>

<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-4 text-gray-100">Response Details</h1>

    <div v-if="isLoading" class="text-center text-gray-400 py-8">Loading response details...</div>
    <div v-else-if="fetchError" class="alert alert-error shadow-lg">
      <div>
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current flex-shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>Error loading response details: {{ fetchError }}</span>
      </div>
    </div>
    <div v-else-if="responseDetails" class="space-y-6">
      <!-- Summary Section -->
      <div class="bg-gray-800 p-4 rounded-lg shadow">
        <h2 class="text-xl font-semibold text-gray-200 mb-3">Summary</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
          <div><strong>Quiz:</strong> <span class="text-gray-300">{{ responseDetails.quiz_title }}</span></div>
          <div><strong>Responder:</strong> <span class="text-gray-300 font-mono">{{ responseDetails.responder_username }}</span></div>
          <div><strong>Submitted:</strong> <span class="text-gray-300">{{ formatDateTime(responseDetails.submitted_at) }}</span></div>
          <div><strong>Started:</strong> <span class="text-gray-300">{{ responseDetails.started_at ? formatDateTime(responseDetails.started_at) : 'Not recorded' }}</span></div>
          <div><strong>Time Taken:</strong> <span class="text-gray-300">{{ responseDetails.time_taken_formatted || formatTimeTaken(responseDetails.time_taken_seconds) || 'Not recorded' }}</span></div>
          <div><strong>Response ID:</strong> <span class="text-gray-300 font-mono">{{ responseDetails.id }}</span></div>
          <div><strong>Score:</strong> <span class="text-gray-300">{{ responseDetails.score !== null && responseDetails.score !== undefined ? `${responseDetails.score}%` : 'Not scored' }}</span></div>
        </div>
      </div>

      <!-- Answers Section -->
      <div class="bg-gray-800 p-4 rounded-lg shadow">
        <h2 class="text-xl font-semibold text-gray-200 mb-4">Answers</h2>
        <div class="space-y-5">
          <div v-for="(answer, index) in responseDetails.answers" :key="index" class="border-b border-gray-700 pb-4 last:border-b-0 last:pb-0">
            <p class="font-semibold text-gray-200 mb-2">{{ index + 1 }}. {{ answer.question_text }}</p>

            <!-- All choices for this question -->
            <div class="space-y-2 pl-4 mt-3">
              <p class="text-sm text-gray-300 mb-1">All Choices:</p>
              <div v-for="(choice, choiceIndex) in answer.all_choices" :key="choiceIndex"
                   class="flex items-center text-sm p-1 rounded"
                   :class="{
                     'bg-green-900/30': choice.isCorrect,
                     'border-l-4 border-blue-500': isChoiceSelected(answer, choice.id),
                     'opacity-80': !choice.isCorrect && !isChoiceSelected(answer, choice.id)
                   }">
                <!-- Indicator for selected choice (check both single and multiple selections) -->
                <span v-if="isChoiceSelected(answer, choice.id)"
                      class="inline-block w-5 h-5 mr-2 border rounded flex items-center justify-center flex-shrink-0"
                      :class="{
                        'bg-green-600 border-green-500 text-white': choice.isCorrect,
                        'bg-red-600 border-red-500 text-white': !choice.isCorrect,
                      }">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path v-if="choice.isCorrect" stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                    <path v-else stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </span>

                <!-- Indicator for correct choice -->
                <span v-else-if="choice.isCorrect"
                      class="inline-block w-5 h-5 mr-2 border rounded flex items-center justify-center flex-shrink-0
                             bg-green-600/50 border-green-500/50 text-white">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                  </svg>
                </span>

                <!-- Empty indicator for other choices -->
                <span v-else class="inline-block w-5 h-5 mr-2 border border-gray-500 rounded flex-shrink-0"></span>

                <!-- Choice text with indicators -->
                <span :class="{
                  'text-green-400 font-medium': choice.isCorrect && isChoiceSelected(answer, choice.id),
                  'text-red-400': !choice.isCorrect && isChoiceSelected(answer, choice.id),
                  'text-green-300/80': choice.isCorrect && !isChoiceSelected(answer, choice.id),
                  'text-gray-400': !choice.isCorrect && !isChoiceSelected(answer, choice.id)
                }">
                  {{ choice.text }}
                  <span v-if="choice.isCorrect" class="text-xs text-green-500 ml-1">(Correct)</span>
                  <span v-if="isChoiceSelected(answer, choice.id)" class="text-xs text-blue-500 ml-1">(Selected)</span>
                </span>
              </div>

              <!-- If no choices are available or it's a text answer -->
              <div v-if="(!answer.all_choices || answer.all_choices.length === 0) && answer.answer_text" class="text-sm mt-2">
                <span class="text-gray-300">Text Answer: </span>
                <span class="text-gray-200 font-medium">{{ answer.answer_text }}</span>
              </div>
            </div>

            <!-- Overall result indicator -->
            <div v-if="answer.isCorrect !== undefined" class="mt-3 pl-4 text-sm font-semibold"
                 :class="answer.isCorrect ? 'text-green-400' : 'text-red-400'"
            >
              {{ answer.isCorrect ? '✓ Answer Correct' : '✗ Answer Incorrect' }}
            </div>
          </div>
        </div>
      </div>
    </div>
     <div v-else class="text-center text-gray-500 py-8">
       Response data is incomplete or missing associated quiz information.
     </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { getAdminResponseDetails } from '@/services/api';
import type { QuizResponseDetail } from '@/types';

const route = useRoute();
const responseId = ref<number | null>(null);
const responseDetails = ref<QuizResponseDetail | null>(null);
const isLoading = ref(false);
const fetchError = ref<string | null>(null);

const formatDateTime = (dateTimeString: string | null | undefined): string => {
  if (!dateTimeString) return '-';
  try {
    return new Date(dateTimeString).toLocaleString();
  } catch (e) {
    return 'Invalid Date';
  }
};

const formatTimeTaken = (seconds: number | null | undefined): string => {
  if (seconds === null || seconds === undefined) return '';

  // Format time taken as HH:MM:SS
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;

  return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
};

// Helper function to check if a choice is selected
const isChoiceSelected = (answer: any, choiceId: number): boolean => {
  // Check if the choice ID is in the selected_choice_ids array
  return answer.selected_choice_ids &&
         Array.isArray(answer.selected_choice_ids) &&
         answer.selected_choice_ids.includes(choiceId);
};

const fetchResponseDetails = async () => {
  if (!responseId.value) return;
  isLoading.value = true;
  fetchError.value = null;
  try {
    responseDetails.value = await getAdminResponseDetails(responseId.value);
    // Log to see if quiz data is included
    console.log('Fetched Response Details:', responseDetails.value);
    if (!responseDetails.value?.quiz) {
        console.warn('Response details fetched but missing associated quiz data.');
        // Optionally set an error or handle this state in the template
    }
  } catch (error: any) {
    console.error('Error fetching response details:', error);
    fetchError.value = error.message || 'Failed to fetch response details.';
  } finally {
    isLoading.value = false;
  }
};

// These helper functions are no longer needed with the new data format
// We're now directly using the answer objects from the backend

onMounted(() => {
  const idParam = route.params.responseId; // Make sure param name matches router config
  if (typeof idParam === 'string') {
    const parsedId = parseInt(idParam, 10);
    if (!isNaN(parsedId)) {
      responseId.value = parsedId;
      fetchResponseDetails();
    } else {
      fetchError.value = 'Invalid Response ID provided in URL.';
    }
  } else {
    fetchError.value = 'Response ID not found in URL parameters.';
  }
});
</script>

<style scoped>
/* Minimal scoped styles, rely on UnoCSS */
</style>

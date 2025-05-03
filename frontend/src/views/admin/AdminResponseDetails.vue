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
    <div v-else-if="responseDetails && responseDetails.quiz" class="space-y-6">
      <!-- Summary Section -->
      <div class="bg-gray-800 p-4 rounded-lg shadow">
        <h2 class="text-xl font-semibold text-gray-200 mb-3">Summary</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
          <div><strong>Quiz:</strong> <span class="text-gray-300">{{ responseDetails.quizTitle }}</span></div>
          <div><strong>Responder:</strong> <span class="text-gray-300 font-mono">{{ responseDetails.responderUsername }}</span></div>
          <div><strong>Submitted:</strong> <span class="text-gray-300">{{ formatDateTime(responseDetails.submittedAt) }}</span></div>
          <div><strong>Response ID:</strong> <span class="text-gray-300 font-mono">{{ responseDetails.id }}</span></div>
          <!-- Add Score here if available -->
        </div>
      </div>

      <!-- Answers Section -->
      <div class="bg-gray-800 p-4 rounded-lg shadow">
        <h2 class="text-xl font-semibold text-gray-200 mb-4">Answers</h2>
        <div class="space-y-5">
          <div v-for="(question, index) in responseDetails.quiz.questions" :key="question.id" class="border-b border-gray-700 pb-4 last:border-b-0 last:pb-0">
            <p class="font-semibold text-gray-200 mb-2">{{ index + 1 }}. {{ question.text }}</p>
            <div class="space-y-1 pl-4">
              <div v-for="choice in question.choices" :key="choice.id" class="flex items-center text-sm">
                 <span
                  class="inline-block w-5 h-5 mr-2 border rounded flex items-center justify-center flex-shrink-0"
                  :class="{
                    'bg-green-600 border-green-500 text-white': isSelected(question.id, choice.id) && choice.is_correct,
                    'bg-red-600 border-red-500 text-white': isSelected(question.id, choice.id) && !choice.is_correct,
                    'border-gray-500': !isSelected(question.id, choice.id)
                  }"
                 >
                    <svg v-if="isSelected(question.id, choice.id)" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                    </svg>
                 </span>
                 <span :class="{'text-gray-400': !choice.is_correct, 'text-green-400 font-medium': choice.is_correct }">
                   {{ choice.text }}
                   <span v-if="choice.is_correct" class="text-xs text-green-500 ml-1">(Correct)</span>
                 </span>
              </div>
            </div>
             <div v-if="submittedAnswerForQuestion(question.id)?.isCorrect !== undefined" class="mt-2 pl-4 text-xs font-semibold"
               :class="submittedAnswerForQuestion(question.id)?.isCorrect ? 'text-green-400' : 'text-red-400'"
             >
               {{ submittedAnswerForQuestion(question.id)?.isCorrect ? 'Answer Correct' : 'Answer Incorrect' }}
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
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { getAdminResponseDetails } from '@/services/api';
import type { QuizResponseDetail, SubmittedAnswer } from '@/types';

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

// Helper to find the submitted answer for a given question ID
const submittedAnswerForQuestion = (questionId: number): SubmittedAnswer | undefined => {
    return responseDetails.value?.answers.find(a => a.questionId === questionId);
};

// Helper to check if a specific choice was selected for a given question
const isSelected = (questionId: number, choiceId: number): boolean => {
    const answer = submittedAnswerForQuestion(questionId);
    return answer?.chosenChoiceIds.includes(choiceId) ?? false;
};

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

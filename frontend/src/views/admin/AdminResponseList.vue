<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-4 text-gray-100">Responses for Quiz #{{ quizId }}</h1>

    <div v-if="isLoading" class="text-center text-gray-400">Loading responses...</div>
    <div v-else-if="fetchError" class="text-red-400 text-center">Error loading responses: {{ fetchError }}</div>
    <div v-else-if="responses.length === 0" class="text-center text-gray-500">No responses submitted yet.</div>
    <div v-else class="overflow-x-auto">
      <table class="table-auto w-full bg-gray-800 rounded-lg shadow">
        <thead>
          <tr class="bg-gray-700">
            <th class="p-3 text-left text-gray-300 font-medium">Responder</th>
            <th class="p-3 text-center text-gray-300 font-medium">Score</th>
            <th class="p-3 text-left text-gray-300 font-medium">Time Taken</th>
            <th class="p-3 text-left text-gray-300 font-medium">Submitted At</th>
            <th class="p-3 text-center text-gray-300 font-medium">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="response in responses" :key="response.id" class="hover:bg-gray-700/50 border-t border-gray-700">
            <td class="p-3 text-sm text-gray-200 font-mono">{{ response.responder_username }}</td>
            <td class="p-3 text-sm text-gray-200 text-center">
              {{ response.score !== null && response.score !== undefined ? `${response.score}%` : 'N/A' }}
            </td>
            <td class="p-3 text-sm text-gray-200">
              <span :class="getTimeTakenClass(response.time_taken_seconds)">
                {{ formatTimeTaken(response.time_taken_seconds) || 'Not recorded' }}
              </span>
            </td>
            <td class="p-3 text-sm text-gray-200">{{ formatDateTime(response.submitted_at) }}</td>
            <td class="p-3 text-sm text-gray-200 text-center">
              <router-link
                :to="{ name: 'AdminResponseDetails', params: { responseId: response.id } }"
                class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white rounded-md inline-flex items-center transition-colors"
              >
                <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
                View Details
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { getAdminQuizResponses } from '@/services/api';
import type { QuizResponseSummary } from '@/types';

const route = useRoute();
const quizId = ref<number | null>(null);
const responses = ref<QuizResponseSummary[]>([]);
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

// Format time taken in seconds to a human-readable format (HH:MM:SS)
const formatTimeTaken = (seconds: number | null | undefined): string => {
  if (seconds === null || seconds === undefined) return '';
  
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  
  return [
    hours.toString().padStart(2, '0'),
    minutes.toString().padStart(2, '0'),
    secs.toString().padStart(2, '0')
  ].join(':');
};

// Return CSS classes based on time taken (color coding)
const getTimeTakenClass = (seconds: number | null | undefined): string => {
  if (seconds === null || seconds === undefined) return 'text-gray-400';
  
  // Color code based on time taken
  if (seconds < 60) return 'text-green-400'; // Under 1 minute
  if (seconds < 300) return 'text-blue-400'; // Under 5 minutes
  if (seconds < 600) return 'text-yellow-400'; // Under 10 minutes
  return 'text-orange-400'; // Over 10 minutes
};

const fetchResponses = async () => {
  if (!quizId.value) return;
  isLoading.value = true;
  fetchError.value = null;
  try {
    responses.value = await getAdminQuizResponses(quizId.value);
  } catch (error: any) {
    console.error('Error fetching responses:', error);
    fetchError.value = error.message || 'Failed to fetch responses.';
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  const idParam = route.params.id;
  if (typeof idParam === 'string') {
    const parsedId = parseInt(idParam, 10);
    if (!isNaN(parsedId)) {
      quizId.value = parsedId;
      fetchResponses(); // Fetch responses when component mounts
    } else {
      fetchError.value = 'Invalid Quiz ID provided in URL.';
    }
  } else {
    fetchError.value = 'Quiz ID not found in URL parameters.';
  }
});
</script>

<style scoped>
/* Using UnoCSS utility classes primarily, minimal scoped styles needed */
.table th {
    @apply font-medium;
}
</style>

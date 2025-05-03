<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-4 text-gray-100">Responses for Quiz #{{ quizId }}</h1>

    <div v-if="isLoading" class="text-center text-gray-400">Loading responses...</div>
    <div v-else-if="fetchError" class="text-red-400 text-center">Error loading responses: {{ fetchError }}</div>
    <div v-else-if="responses.length === 0" class="text-center text-gray-500">No responses submitted yet.</div>
    <div v-else class="overflow-x-auto">
      <table class="table table-zebra table-sm w-full bg-gray-800 rounded-lg shadow">
        <thead>
          <tr class="bg-gray-700">
            <th class="p-2 text-gray-300">Responder</th>
            <th class="p-2 text-gray-300">Score</th>
            <th class="p-2 text-gray-300">Submitted At</th>
            <th class="p-2 text-gray-300">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="response in responses" :key="response.id" class="hover:bg-gray-700/50">
            <td class="border-t border-gray-700 p-2 text-sm text-gray-200 align-middle font-mono">{{ response.responderUsername }}</td>
            <td class="border-t border-gray-700 p-2 text-sm text-gray-200 align-middle">
              {{ response.score !== null && response.score !== undefined ? response.score : 'N/A' }}
            </td>
            <td class="border-t border-gray-700 p-2 text-sm text-gray-200 align-middle">{{ formatDateTime(response.submittedAt) }}</td>
            <td class="border-t border-gray-700 p-2 text-sm text-gray-200 align-middle">
              <router-link
                :to="{ name: 'AdminResponseDetails', params: { responseId: response.id } }"
                class="btn btn-xs btn-outline btn-info"
              >
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

<template>
  <div class="p-6">
    <h1 v-if="quizData" class="text-2xl font-semibold mb-6 text-gray-100">Edit Quiz: {{ quizData.title }}</h1>
    <h1 v-else class="text-2xl font-semibold mb-6 text-gray-100">Loading Quiz...</h1>

    <QuizForm
      v-if="quizData || !isLoading" 
      :initial-quiz-data="quizData"
      :is-loading="isLoading"
      submit-button-text="Update Quiz"
      :global-error="globalError"
      @submit="handleUpdateQuiz"
      @cancel="handleCancel"
    />

    <div v-if="isLoading && !globalError" class="text-center mt-4 text-gray-400">
      Loading quiz data...
    </div>

     <div v-if="!isLoading && globalError && !quizData" class="mt-4 p-3 bg-red-900 border border-red-700 text-red-200 rounded text-sm">
      Error loading quiz: {{ globalError }}
      <button @click="fetchQuizData" class="ml-2 text-teal-400 hover:text-teal-300 underline">Retry</button>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import QuizForm from '@/components/admin/QuizForm.vue';
import { getAdminQuiz, updateAdminQuiz } from '@/services/api'; // Assuming these API functions exist
import type { Quiz, QuizInput } from '@/types'; // Use Quiz for fetched data, QuizInput for submit

const router = useRouter();
const route = useRoute();

const isLoading = ref(false);
const globalError = ref<string | null>(null);
const quizData = ref<Quiz | null>(null);
const quizId = ref<number | null>(null);

const fetchQuizData = async () => {
  if (!quizId.value) return;
  isLoading.value = true;
  globalError.value = null;
  quizData.value = null; // Clear previous data before fetching

  try {
    const response = await getAdminQuiz(quizId.value);
    quizData.value = response; // Assuming the API returns the full Quiz structure
  } catch (error: any) {
    console.error('Error fetching quiz:', error);
    globalError.value = error.response?.data?.error || error.message || 'An unexpected error occurred while fetching the quiz data.';
  } finally {
    isLoading.value = false;
  }
};

const handleUpdateQuiz = async (updatedQuizData: QuizInput) => {
  if (!quizId.value) return;
  isLoading.value = true;
  globalError.value = null;
  try {
    console.log(`Updating quiz ${quizId.value} with data:`, JSON.stringify(updatedQuizData, null, 2));
    // The update function might need the ID along with the payload
    await updateAdminQuiz(quizId.value, updatedQuizData);
    // Navigate back to the list on success
    // Consider adding a success notification
    router.push({ name: 'admin-quiz-list' });
  } catch (error: any) {
    console.error('Error updating quiz:', error);
    globalError.value = error.response?.data?.error || error.message || 'An unexpected error occurred while updating the quiz.';
  } finally {
    isLoading.value = false;
  }
};

const handleCancel = () => {
  router.push({ name: 'admin-quiz-list' });
};

onMounted(() => {
  const idParam = route.params.id;
  if (typeof idParam === 'string') {
    const parsedId = parseInt(idParam, 10);
    if (!isNaN(parsedId)) {
      quizId.value = parsedId;
      fetchQuizData();
    } else {
      globalError.value = 'Invalid Quiz ID provided in URL.';
    }
  } else {
     globalError.value = 'Quiz ID not found in URL parameters.';
  }
});

</script>

<style scoped>
/* View-specific styles if needed */
</style>

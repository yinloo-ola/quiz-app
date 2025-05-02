<template>
  <div class="p-6">
    <h1 class="text-2xl font-semibold mb-6 text-gray-100">Create New Quiz</h1>

    <QuizForm
      :is-loading="isLoading"
      submit-button-text="Create Quiz"
      :global-error="globalError"
      @submit="handleCreateQuiz"
      @cancel="handleCancel"
    />

  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import QuizForm from '@/components/admin/QuizForm.vue'; // Import the reusable form component
import { createAdminQuiz } from '@/services/api'; // Import the API service function
import type { QuizInput } from '@/types'; // Import the necessary payload type

const router = useRouter();

const isLoading = ref(false);
const globalError = ref<string | null>(null); // Define globalError

const handleCreateQuiz = async (quizData: QuizInput) => {
  isLoading.value = true;
  globalError.value = null;
  try {
    console.log('Creating quiz with data:', JSON.stringify(quizData, null, 2)); // Log payload
    await createAdminQuiz(quizData); // Call the API service function
    // Navigate to the quiz list or the newly created quiz page on success
    // Add success notification if desired
    router.push({ name: 'admin-quiz-list' });
  } catch (error: any) {
    console.error('Error creating quiz:', error);
    // Attempt to get a specific error message from the API response, otherwise use a generic one
    globalError.value = error.response?.data?.error || error.message || 'An unexpected error occurred while creating the quiz.';
  } finally {
    isLoading.value = false;
  }
};

const handleCancel = () => {
  // Optionally add a confirmation dialog if needed
  router.push({ name: 'admin-quiz-list' }); // Navigate back to the list
};

</script>

<style scoped>
/* Add any view-specific styles if needed */
input[type='radio']:focus,
input[type='checkbox']:focus {
  box-shadow: 0 0 0 2px theme('colors.teal.500'); /* Custom focus ring for radio/checkbox */
}

/* Removed scoped override for button focus - handled by global reset */
</style>

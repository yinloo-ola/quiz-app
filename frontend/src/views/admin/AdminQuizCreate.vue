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

      <!-- TODO: Add fields for Time Limit, Questions, Choices etc. -->
      <p class="text-gray-500 text-sm mb-6">More fields (time limit, questions) will be added later.</p>

      <!-- Action Buttons -->
      <div class="flex justify-end space-x-3">
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

const router = useRouter();

// Define reactive variables for form fields
const title = ref('');
const description = ref('');
const isLoading = ref(false);
const error = ref<string | null>(null);

// Function to handle form submission
const createQuiz = async () => {
  isLoading.value = true;
  error.value = null;
  console.log('Form submitted');
  console.log('Title:', title.value);
  console.log('Description:', description.value);

  try {
    const newQuizData = {
      title: title.value,
      description: description.value || undefined,
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

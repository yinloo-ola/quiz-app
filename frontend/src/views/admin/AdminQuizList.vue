<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-semibold">Manage Quizzes</h1>
      <button
        @click="goToCreateQuiz"
        class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
      >
        Create New Quiz
      </button>
    </div>

    <!-- Quiz List Table (placeholder) -->
    <div class="bg-white shadow-md rounded my-6">
      <table class="min-w-full table-auto">
        <thead>
          <tr class="bg-gray-200 text-gray-600 uppercase text-sm leading-normal">
            <th class="py-3 px-6 text-left">Title</th>
            <th class="py-3 px-6 text-left">Status</th>
            <th class="py-3 px-6 text-center">Actions</th>
          </tr>
        </thead>
        <tbody class="text-gray-600 text-sm font-light">
          <!-- Loading State -->
          <tr v-if="isLoading">
            <td colspan="3" class="py-3 px-6 text-center">Loading quizzes...</td>
          </tr>
          <!-- Error State -->
          <tr v-else-if="error">
             <td colspan="3" class="py-3 px-6 text-center text-red-600">Error loading quizzes: {{ error }}</td>
          </tr>
          <!-- Empty State -->
          <tr v-else-if="quizzes.length === 0">
             <td colspan="3" class="py-3 px-6 text-center">No quizzes found. Create one!</td>
          </tr>
          <!-- Quiz rows will be populated here -->
          <tr v-else v-for="quiz in quizzes" :key="quiz.id" class="border-b border-gray-200 hover:bg-gray-100">
            <td class="py-3 px-6 text-left whitespace-nowrap">
              <span>{{ quiz.title }}</span>
            </td>
            <td class="py-3 px-6 text-left">
              <span>{{ quiz.status || 'Draft' }}</span> <!-- Placeholder status -->
            </td>
            <td class="py-3 px-6 text-center">
               <!-- Action Buttons Placeholder -->
               <button class="bg-blue-500 hover:bg-blue-700 text-white text-xs py-1 px-2 rounded mr-1">View</button>
               <button class="bg-yellow-500 hover:bg-yellow-700 text-white text-xs py-1 px-2 rounded mr-1">Edit</button>
               <button class="bg-red-500 hover:bg-red-700 text-white text-xs py-1 px-2 rounded">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { fetchAdminQuizzes } from '@/services/api';

const router = useRouter();
const route = useRoute();

interface Quiz {
  id: number;
  title: string;
  description?: string;
  status: string;
  // Add other relevant quiz properties like CreatedAt, UpdatedAt if needed
}

const quizzes = ref<Quiz[]>([]);
const isLoading = ref(false);
const error = ref<string | null>(null);

// Refactored function to load quizzes
const loadQuizzes = async () => {
  console.log('Attempting to load quizzes...');
  isLoading.value = true;
  error.value = null;
  try {
    const data = await fetchAdminQuizzes();
    console.log('Quizzes loaded:', data);
    // Ensure backend data matches frontend interface (e.g., ID vs id)
    // Map data if necessary, or adjust interface
    quizzes.value = data;
  } catch (err: any) {
    console.error('Failed to load quizzes in component:', err);
    error.value = err.message || 'An unexpected error occurred while fetching quizzes.';
  } finally {
    isLoading.value = false;
  }
};

// Function to navigate to the create quiz page
const goToCreateQuiz = () => {
  router.push({ name: 'admin-quiz-create' });
};

// Call loadQuizzes when the component is first mounted
onMounted(loadQuizzes);

// Watch for route changes to ensure data is fresh when navigating back
watch(
  () => route.name,
  (newName) => {
    console.log('Route name changed to:', newName);
    if (newName === 'admin-quiz-list') {
      console.log('Route is admin-quiz-list, reloading quizzes...');
      loadQuizzes();
    }
  },
  { immediate: false }
);

</script>

<style scoped>
/* Component styles */
</style>

<template>
  <div class="p-6">
    <div class="flex justify-between items-center mb-4">
      <h1 class="text-2xl font-semibold text-gray-100">Manage Quizzes</h1>
      <button
        @click="goToCreateQuiz"
        class="bg-teal-600 hover:bg-teal-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-teal-500"
      >
        Create New Quiz
      </button>
    </div>

    <!-- Quiz List Table -->
    <div class="bg-slate-800 border border-slate-700 rounded my-6">
      <table class="min-w-full table-auto">
        <thead>
          <tr class="bg-slate-700 text-gray-300 uppercase text-sm leading-normal">
            <th class="py-3 px-6 text-left">Title</th>
            <th class="py-3 px-6 text-left">Status</th>
            <th class="py-3 px-6 text-center">Actions</th>
          </tr>
        </thead>
        <tbody class="text-gray-200 text-sm font-light">
          <!-- Loading State -->
          <tr v-if="isLoading">
            <td colspan="3" class="py-3 px-6 text-center text-gray-400">Loading quizzes...</td>
          </tr>
          <!-- Error State -->
          <tr v-else-if="error">
             <td colspan="3" class="py-3 px-6 text-center text-red-400">Error loading quizzes: {{ error }}</td>
          </tr>
          <!-- Empty State -->
          <tr v-else-if="quizzes.length === 0">
             <td colspan="3" class="py-3 px-6 text-center text-gray-400">No quizzes found. Create one!</td>
          </tr>
          <!-- Quiz rows -->
          <tr v-else v-for="quiz in quizzes" :key="quiz.id" class="border-b border-slate-700 hover:bg-slate-700">
            <td class="py-3 px-6 text-left whitespace-nowrap">
              <span>{{ quiz.title }}</span>
            </td>
            <td class="py-3 px-6 text-left">
              <span class="capitalize">{{ quiz.status || 'draft' }}</span> <!-- Ensure consistent casing -->
            </td>
            <td class="py-3 px-6 text-center">
               <!-- Action Buttons -->
               <button @click="goToEditQuiz(quiz.id)" class="bg-yellow-500 hover:bg-yellow-600 text-gray-900 text-xs py-1 px-2 rounded mr-1 focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-yellow-400">Edit</button>
               <button @click="viewResponses(quiz.id)" class="bg-sky-600 hover:bg-sky-700 text-white text-xs py-1 px-2 rounded mr-1 focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-sky-500">Responses</button>
               <button @click="manageCredentials(quiz.id)" class="bg-indigo-600 hover:bg-indigo-700 text-white text-xs py-1 px-2 rounded mr-1 focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-indigo-500">Credentials</button>
               <button @click="deleteQuiz(quiz.id)" class="bg-red-600 hover:bg-red-700 text-white text-xs py-1 px-2 rounded focus:outline-none focus:ring-2 focus:ring-offset-1 focus:ring-offset-slate-800 focus:ring-red-500">Delete</button>
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

// --- Action Handlers ---

const goToEditQuiz = (quizId: number) => {
  router.push({ name: 'admin-quiz-edit', params: { id: quizId } });
};

const viewResponses = (quizId: number) => {
  // TODO: Ensure 'admin-quiz-responses' route exists
  router.push({ name: 'admin-quiz-responses', params: { id: quizId } });
};

const manageCredentials = (quizId: number) => {
  // TODO: Ensure 'admin-quiz-credentials' route exists
  router.push({ name: 'admin-quiz-credentials', params: { id: quizId } });
};

const deleteQuiz = async (quizId: number) => {
  // TODO: Add confirmation dialog
  // TODO: Call API service function: await deleteAdminQuiz(quizId);
  console.warn(`Placeholder: Delete quiz ${quizId}. API call and confirmation needed.`);
  // TODO: Reload quizzes after deletion: loadQuizzes();
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
/* Removed scoped override for button focus - handled by global reset */
/* Component styles */
</style>

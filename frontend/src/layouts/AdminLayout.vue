<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-gray-800 text-white p-4">
      <h2 class="text-xl font-bold mb-6">Quiz Admin</h2>
      <nav>
        <ul>
          <li class="mb-2">
            <router-link 
              :to="{ name: 'admin-dashboard' }" 
              class="block px-4 py-2 rounded hover:bg-gray-700"
              active-class="bg-gray-900"
            >
              Dashboard
            </router-link>
          </li>
          <!-- Quiz Management Link will go here -->
          <li class="mb-2">
             <router-link 
              :to="{ name: 'admin-quiz-list' }" 
              class="block px-4 py-2 rounded hover:bg-gray-700"
              active-class="bg-gray-900"
            >
              Quizzes
            </router-link>
          </li>
        </ul>
      </nav>
      <div class="mt-auto">
        <!-- Logout Button Placeholder -->
        <button 
          @click="handleLogout" 
          class="w-full mt-6 px-4 py-2 rounded bg-red-600 hover:bg-red-700">
          Logout
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-6 overflow-auto">
      <router-view /> <!-- Nested routes will render here -->
    </main>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'; // Import store
import { useRouter } from 'vue-router'; // Import router

const authStore = useAuthStore();
const router = useRouter();

const handleLogout = () => {
  authStore.logout();
  router.push({ name: 'admin-login' }); // Redirect to login page
};
</script>

<style scoped>
/* Add specific styles for the layout if needed */
</style>

<template>
  <div class="flex h-screen"> 
    <!-- Sidebar -->
    <aside class="w-64 bg-slate-900 text-white p-4 flex flex-col"> 
      <h2 class="text-xl font-bold mb-6">Quiz Admin</h2>
      <nav class="flex-grow"> 
        <ul class="list-none pl-0">
          <li class="mb-2">
            <router-link 
              :to="{ name: 'admin-dashboard' }" 
              :class="getLinkClass('admin-dashboard')"
            >
              Dashboard
            </router-link>
          </li>
          <!-- Quiz Management Link will go here -->
          <li class="mb-2">
             <router-link 
              :to="{ name: 'admin-quiz-list' }" 
              :class="getLinkClass('admin-quiz-list')"
            >
              Quizzes
            </router-link>
          </li>
        </ul>
      </nav>
      <!-- Logout Button -->
      <div class="mt-auto">
        <button 
          @click="handleLogout"
          class="logout-button w-full mt-6 px-4 py-2 rounded border border-transparent bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-900 focus:ring-red-500"
        >
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
import { useRouter, useRoute } from 'vue-router'; // Import router

const authStore = useAuthStore();
const router = useRouter();
const route = useRoute(); // Get the current route object

// Function to determine link classes based on exact route name match
const getLinkClass = (targetRouteName: string) => {
  return {
    'bg-teal-600': route.name === targetRouteName, // Apply highlight only if names match exactly
    'block px-4 py-2 rounded hover:bg-slate-700': true // Always apply base classes
  };
};

const handleLogout = () => {
  authStore.logout();
  router.push({ name: 'admin-login' }); // Redirect to login page
};
</script>

<style scoped>
/* Removed scoped override for logout button focus - handled by global reset */
</style>

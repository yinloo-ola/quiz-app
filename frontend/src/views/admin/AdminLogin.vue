<template>
  <div class="p-4 max-w-md mx-auto">
    <h2 class="text-2xl font-bold mb-4">Admin Login</h2>
    <form @submit.prevent="onSubmit" class="space-y-4">
      <div>
        <label for="username" class="block text-sm font-medium">Username</label>
        <input
          id="username"
          v-model="username"
          type="text"
          required
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
      <div>
        <label for="password" class="block text-sm font-medium">Password</label>
        <input
          id="password"
          v-model="password"
          type="password"
          required
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
      
      <div v-if="authStore.error" class="text-red-600 text-sm">
        {{ authStore.error }}
      </div>

      <button
        type="submit"
        :disabled="isLoading"
        class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
      >
        {{ isLoading ? 'Logging in...' : 'Login' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth'; 
import { adminLogin } from '@/services/api'; 

const username = ref('');
const password = ref('');
const isLoading = ref(false);

const router = useRouter();
const authStore = useAuthStore();

const onSubmit = async () => {
  isLoading.value = true;
  authStore.setError(null); 

  try {
    const response = await adminLogin({ 
        username: username.value, 
        password: password.value 
    });
    
    authStore.setToken(response.token);
    // Optional: Fetch user details if needed
    // authStore.setUser(someUserData); 
    
    // Redirect to the admin dashboard
    router.push({ name: 'admin-dashboard' }); 

  } catch (error: any) {
    console.error('Login failed:', error); // Log the full error object for debugging
    
    let errorMessage = 'Login failed. Please check your credentials or server status.';
    if (error.response && error.response.data && error.response.data.error) {
      errorMessage = error.response.data.error; // Use backend error message
    } else if (error.message) {
      errorMessage = error.message; // Use generic error message
    }
    authStore.setError(errorMessage);

  } finally {
    isLoading.value = false;
  }
};
</script>

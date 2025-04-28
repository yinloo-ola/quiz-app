<template>
  <div class="min-h-screen flex items-center justify-center p-4"> 
    <div class="p-6 sm:p-8 max-w-md w-full bg-slate-800 rounded-lg shadow-md border border-slate-700"> 
      <h2 class="text-2xl font-bold mb-6 text-center text-gray-100">Admin Login</h2>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-300">Username</label>
          <input
            id="username"
            v-model="username"
            type="text"
            required
            class="box-border mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:border-teal-500 sm:text-sm placeholder-gray-400"
            placeholder="Enter username"
          />
        </div>
        <div>
          <label for="password" class="block text-sm font-medium text-gray-300">Password</label>
          <input
            id="password"
            v-model="password"
            type="password"
            required
            class="box-border mt-1 block w-full px-3 py-2 border border-slate-700 rounded-md bg-slate-700 text-gray-100 focus:outline-none focus:border-teal-500 sm:text-sm placeholder-gray-400"
            placeholder="Enter password"
          />
        </div>
        
        <div v-if="authStore.error" class="text-red-400 text-sm pt-2">
          {{ authStore.error }}
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-teal-600 hover:bg-teal-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-slate-800 focus:ring-teal-500 disabled:bg-teal-800 disabled:text-gray-400 disabled:cursor-not-allowed focus:ring-teal-500 focus:ring-offset-slate-800"
        >
          {{ isLoading ? 'Logging in...' : 'Login' }}
        </button>
      </form>
    </div>
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

<style scoped>
/* Removed scoped override for button focus - handled by global reset */
</style>

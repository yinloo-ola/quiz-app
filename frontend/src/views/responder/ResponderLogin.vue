<template>
  <div class="responder-login-container">
    <div class="login-card">
      <h1 class="text-3xl font-bold mb-6 text-gray-900">Quiz Login</h1>
      
      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
        {{ error }}
      </div>
      
      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-900">Username</label>
          <input 
            id="username" 
            v-model="username" 
            type="text" 
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        
        <div>
          <label for="password" class="block text-sm font-medium text-gray-900">Password</label>
          <input 
            id="password" 
            v-model="password" 
            type="password" 
            required
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
          />
        </div>
        
        <div>
          <button 
            type="submit" 
            :disabled="isLoading"
            class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
          >
            <span v-if="isLoading">Logging in...</span>
            <span v-else>Login</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { responderLogin } from '@/services/api';
import { jwtDecode } from 'jwt-decode';

const router = useRouter();
const authStore = useAuthStore();

const username = ref('');
const password = ref('');
const error = ref('');
const isLoading = ref(false);

interface JwtPayload {
  responder_credential_id: number;
  quiz_id: number;
  exp: number;
}

const handleLogin = async () => {
  error.value = '';
  isLoading.value = true;
  
  try {
    const response = await responderLogin({
      username: username.value,
      password: password.value
    });
    
    // Decode the JWT to get credential ID and quiz ID
    const decoded = jwtDecode<JwtPayload>(response.token);
    
    // Store token and credential info in auth store
    authStore.setResponderToken(
      response.token, 
      decoded.responder_credential_id, 
      decoded.quiz_id
    );
    
    // Redirect to the quiz page
    router.push({ 
      name: 'quiz-taker', 
      params: { id: decoded.quiz_id.toString() } 
    });
    
  } catch (err: any) {
    console.error('Login error:', err);
    if (typeof err === 'object' && err.error) {
      error.value = err.error;
    } else {
      error.value = 'Login failed. Please check your credentials and try again.';
    }
  } finally {
    isLoading.value = false;
  }
};
</script>

<style scoped>
.responder-login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f3f4f6;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style>

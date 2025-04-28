import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  // --- State --- S
  // Attempt to load token from localStorage on initial load
  const token = ref<string | null>(localStorage.getItem('admin_token'));
  const adminUser = ref<any | null>(null); // Store admin user details if needed later
  const error = ref<string | null>(null); // Store login errors

  // --- Getters --- (Computed properties)
  const isAuthenticated = computed(() => !!token.value);

  // --- Actions --- (Functions to modify state)
  function setToken(newToken: string) {
    token.value = newToken;
    localStorage.setItem('admin_token', newToken); // Persist token
    error.value = null; // Clear errors on successful login
  }

  function setUser(user: any) {
    adminUser.value = user;
  }

  function setError(errorMessage: string | null) { // Allow null to clear error
    error.value = errorMessage;
  }

  function logout() {
    token.value = null;
    adminUser.value = null;
    error.value = null;
    localStorage.removeItem('admin_token'); // Remove token from storage
    // Potentially redirect using router (can be done in the component calling logout)
  }

  // --- Return --- (Expose state, getters, actions)
  return {
    token,
    adminUser,
    error,
    isAuthenticated,
    setToken,
    setUser,
    setError,
    logout,
  };
});

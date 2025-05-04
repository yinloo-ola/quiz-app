import { defineStore } from 'pinia';
import { ref, computed } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  // --- State ---
  // Attempt to load tokens from localStorage on initial load
  const adminToken = ref<string | null>(localStorage.getItem('admin_token'));
  const responderToken = ref<string | null>(localStorage.getItem('responder_token'));
  const adminUser = ref<any | null>(null); // Store admin user details if needed later
  const responderCredentialId = ref<number | null>(null); // Store responder credential ID
  const responderQuizId = ref<number | null>(null); // Store quiz ID for responder
  const error = ref<string | null>(null); // Store login errors
  const authMode = ref<'admin' | 'responder' | null>(null); // Track current auth mode

  // Initialize auth mode based on stored tokens
  if (adminToken.value) {
    authMode.value = 'admin';
  } else if (responderToken.value) {
    authMode.value = 'responder';
  }

  // --- Getters --- (Computed properties)
  const isAdminAuthenticated = computed(() => !!adminToken.value);
  const isResponderAuthenticated = computed(() => !!responderToken.value);
  const isAuthenticated = computed(() => isAdminAuthenticated.value || isResponderAuthenticated.value);
  const currentToken = computed(() => authMode.value === 'admin' ? adminToken.value : responderToken.value);

  // --- Actions --- (Functions to modify state)
  function setAdminToken(newToken: string) {
    adminToken.value = newToken;
    authMode.value = 'admin';
    localStorage.setItem('admin_token', newToken); // Persist token
    // Clear responder token if it exists
    if (responderToken.value) {
      responderToken.value = null;
      localStorage.removeItem('responder_token');
    }
    error.value = null; // Clear errors on successful login
  }
  
  function setResponderToken(newToken: string, credentialId: number, quizId: number) {
    responderToken.value = newToken;
    responderCredentialId.value = credentialId;
    responderQuizId.value = quizId;
    authMode.value = 'responder';
    localStorage.setItem('responder_token', newToken); // Persist token
    localStorage.setItem('responder_credential_id', credentialId.toString());
    localStorage.setItem('responder_quiz_id', quizId.toString());
    // Clear admin token if it exists
    if (adminToken.value) {
      adminToken.value = null;
      localStorage.removeItem('admin_token');
    }
    error.value = null; // Clear errors on successful login
  }

  function setUser(user: any) {
    adminUser.value = user;
  }

  function setError(errorMessage: string | null) { // Allow null to clear error
    error.value = errorMessage;
  }

  function logout() {
    if (authMode.value === 'admin') {
      adminToken.value = null;
      adminUser.value = null;
      localStorage.removeItem('admin_token');
    } else if (authMode.value === 'responder') {
      responderToken.value = null;
      responderCredentialId.value = null;
      responderQuizId.value = null;
      localStorage.removeItem('responder_token');
      localStorage.removeItem('responder_credential_id');
      localStorage.removeItem('responder_quiz_id');
    }
    
    authMode.value = null;
    error.value = null;
    // Potentially redirect using router (can be done in the component calling logout)
  }

  // --- Return --- (Expose state, getters, actions)
  return {
    adminToken,
    responderToken,
    currentToken,
    adminUser,
    responderCredentialId,
    responderQuizId,
    error,
    authMode,
    isAuthenticated,
    isAdminAuthenticated,
    isResponderAuthenticated,
    setAdminToken,
    setResponderToken,
    setUser,
    setError,
    logout,
  };
});

<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-4 text-gray-100">Manage Credentials for Quiz #{{ quizId }}</h1>

    <!-- Generate Credentials Form -->
    <div class="bg-gray-800 p-4 rounded-lg shadow mb-6">
      <h2 class="text-lg font-medium text-gray-200 mb-3">Generate New Credentials</h2>
      <form @submit.prevent="handleGenerateCredentials" class="flex flex-wrap items-end gap-3">
        <div>
          <label for="username" class="block text-sm font-medium text-gray-400 mb-1">Username (optional):</label>
          <input type="text" id="username" v-model.trim="generateForm.username" class="input input-sm" placeholder="Leave blank for random"/>
        </div>
        <div>
          <label for="expiry" class="block text-sm font-medium text-gray-400 mb-1">Expiry (hours, optional):</label>
          <input type="number" id="expiry" v-model.number="generateForm.expiryHours" min="1" class="input input-sm" />
        </div>
        <button type="submit" :disabled="isGenerating" class="btn btn-primary btn-sm">
          <span v-if="isGenerating" class="loading loading-spinner loading-xs"></span>
          {{ isGenerating ? 'Generating...' : 'Generate' }}
        </button>
      </form>
      <p v-if="generationError" class="text-red-400 text-sm mt-2">{{ generationError }}</p>
      <div v-if="newlyGeneratedCredential" class="mt-4 p-3 bg-green-900/50 border border-green-700 rounded">
        <h3 class="text-md font-semibold text-green-300 mb-2">Generated Credentials (Password shown once):</h3>
        <div class="text-green-200 text-sm font-mono">
          <span>User: {{ newlyGeneratedCredential.username }}</span> | <span class="font-bold">Pass: {{ newlyGeneratedCredential.password }}</span>
        </div>
        <button @click="copyGeneratedCredentials" class="btn btn-xs btn-outline btn-accent mt-2">
          Copy Credentials
        </button>
      </div>
    </div>

    <!-- Existing Credentials List -->
    <h2 class="text-lg font-medium text-gray-200 mb-3">Existing Credentials</h2>
    <div v-if="isLoading" class="text-center text-gray-400">Loading credentials...</div>
    <div v-else-if="fetchError" class="text-red-400 text-center">Error loading credentials: {{ fetchError }}</div>
    <div v-else-if="credentials.length === 0" class="text-center text-gray-500">No credentials generated yet.</div>
    <div v-else class="overflow-x-auto">
      <table class="table table-zebra table-sm w-full">
        <thead>
          <tr>
            <th class="text-left">Username</th>
            <th class="text-left">Expires At</th>
            <th class="text-left w-16">Used</th>
            <th class="text-left">Used At</th>
            <th class="text-left">Created At</th>
            <th class="text-center">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cred in credentials" :key="cred.id" class="hover:bg-gray-700/50">
            <td class="font-mono">{{ cred.username }}</td>
            <td>{{ cred.expiresAt ? formatDateTime(cred.expiresAt) : 'Never' }}</td>
            <td class="text-center">
              <span :class="['badge badge-sm', cred.used ? 'badge-success' : 'badge-warning']">
                {{ cred.used ? 'Yes' : 'No' }}
              </span>
            </td>
            <td>{{ cred.usedAt ? formatDateTime(cred.usedAt) : '-' }}</td>
            <td>{{ formatDateTime(cred.createdAt) }}</td>
            <td class="text-center">
              <button 
                @click="removeCredential(cred.id)" 
                class="btn btn-xs btn-error" 
                :disabled="isDeleting === cred.id"
                title="Delete credential"
              >
                <span v-if="isDeleting === cred.id" class="loading loading-spinner loading-xs"></span>
                <span v-else>Delete</span>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { getAdminQuizCredentials, generateAdminQuizCredentials, deleteCredential } from '@/services/api';
import type { ResponderCredential, GenerateCredentialsResponse } from '@/types';

const route = useRoute();
const quizId = ref<number | null>(null);
const credentials = ref<ResponderCredential[]>([]);
const newlyGeneratedCredential = ref<GenerateCredentialsResponse | null>(null);
const isLoading = ref(false);
const isGenerating = ref(false);
const isDeleting = ref<number | null>(null); // Holds ID of credential being deleted
const fetchError = ref<string | null>(null);
const generationError = ref<string | null>(null);
const deleteError = ref<string | null>(null);

const generateForm = reactive({
  username: '',
  expiryHours: undefined as number | undefined,
});

const formatDateTime = (dateTimeString: string | null | undefined): string => {
  if (!dateTimeString) return '-';
  try {
    return new Date(dateTimeString).toLocaleString();
  } catch (e) {
    return 'Invalid Date';
  }
};

const fetchCredentials = async () => {
  if (!quizId.value) return;
  isLoading.value = true;
  fetchError.value = null;
  try {
    credentials.value = await getAdminQuizCredentials(quizId.value);
  } catch (error: any) {
    console.error('Error fetching credentials:', error);
    fetchError.value = error.message || 'Failed to fetch credentials.';
  } finally {
    isLoading.value = false;
  }
};

const handleGenerateCredentials = async () => {
  if (!quizId.value) return;
  isGenerating.value = true;
  generationError.value = null;
  newlyGeneratedCredential.value = null; // Clear previous result

  try {
    // Payload always includes count: 1 now
    const payload: { count: 1; expiryHours?: number; username?: string} = {
      count: 1,
      // Only include expiryHours if it's a positive number
      ...(generateForm.expiryHours && generateForm.expiryHours > 0
        ? { expiryHours: generateForm.expiryHours }
        : {}),
      // Include username if provided
      ...(generateForm.username ? { username: generateForm.username } : {}),
    };
    const generated = await generateAdminQuizCredentials(quizId.value, payload);
    newlyGeneratedCredential.value = generated; // Assign the single object
    // Refresh the list of existing credentials
    await fetchCredentials();

    // Reset form
    generateForm.username = ''; // Reset username
    generateForm.expiryHours = undefined;
  } catch (error: any) {
    console.error('Error generating credentials:', error);
    generationError.value = error.message || 'Failed to generate credentials.';
  } finally {
    isGenerating.value = false;
  }
};

const copyGeneratedCredentials = () => {
  if (!newlyGeneratedCredential.value) return; // Check if object exists

  const cred = newlyGeneratedCredential.value;
  const credentialsText = `Username: ${cred.username}\nPassword: ${cred.password || 'N/A'}`;

  navigator.clipboard.writeText(credentialsText)
    .then(() => {
      console.log('Credentials copied to clipboard');
      // Optional: show a temporary success message to the user
    })
    .catch(err => {
      console.error('Failed to copy credentials: ', err);
      // Optional: show an error message to the user
    });
};

// Function to remove a credential
const removeCredential = async (credentialId: number) => {
  if (!confirm('Are you sure you want to delete this credential? This action cannot be undone.')) {
    return; // User canceled the operation
  }
  
  isDeleting.value = credentialId;
  deleteError.value = null;
  
  try {
    await deleteCredential(credentialId);
    // Remove the credential from the local array to update UI
    credentials.value = credentials.value.filter(cred => cred.id !== credentialId);
  } catch (error: any) {
    console.error('Error removing credential:', error);
    deleteError.value = error.message || 'Failed to delete credential.';
  } finally {
    isDeleting.value = null;
  }
};

onMounted(() => {
  const idParam = route.params.id;
  if (typeof idParam === 'string') {
    const parsedId = parseInt(idParam, 10);
    if (!isNaN(parsedId)) {
      quizId.value = parsedId;
      fetchCredentials(); // Fetch credentials when component mounts
    } else {
      fetchError.value = 'Invalid Quiz ID provided in URL.';
    }
  } else {
    fetchError.value = 'Quiz ID not found in URL parameters.';
  }
});
</script>

<style scoped>
.input {
  @apply bg-gray-700 border border-gray-600 text-gray-200 rounded px-3 py-1.5 focus:ring-teal-500 focus:border-teal-500;
}
.btn {
    @apply font-semibold transition duration-150 ease-in-out;
}
.table {
    @apply bg-gray-800 rounded-lg shadow;
    th {
        @apply bg-gray-700 text-gray-300 font-medium p-2;
    }
    td {
        @apply border-t border-gray-700 p-2 text-sm text-gray-200 align-middle;
    }
}
/* Add other styles if needed */
</style>

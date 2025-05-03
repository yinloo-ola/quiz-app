<template>
  <div class="p-4 md:p-6">
    <h1 class="text-2xl font-semibold mb-4 text-gray-100">Manage Credentials for Quiz #{{ quizId }}</h1>

    <!-- Generate Credentials Form -->
    <div class="bg-gray-800 p-4 rounded-lg shadow mb-6">
      <h2 class="text-lg font-medium text-gray-200 mb-3">Generate New Credentials</h2>
      <form @submit.prevent="handleGenerateCredentials" class="flex flex-wrap items-end gap-3">
        <div>
          <label for="count" class="block text-sm font-medium text-gray-400 mb-1">Number to Generate:</label>
          <input type="number" id="count" v-model.number="generateForm.count" min="1" max="100" required class="input input-sm" />
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
      <div v-if="newlyGeneratedCredentials.length > 0" class="mt-4 p-3 bg-green-900/50 border border-green-700 rounded">
        <h3 class="text-md font-semibold text-green-300 mb-2">Generated Credentials:</h3>
        <ul class="list-disc list-inside text-green-200 text-sm">
          <li v-for="cred in newlyGeneratedCredentials" :key="cred.id">{{ cred.username }}</li>
        </ul>
         <button @click="copyGeneratedUsernames" class="btn btn-xs btn-outline btn-accent mt-2">
           Copy Usernames
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
            <th>Username</th>
            <th>Expires At</th>
            <th>Used</th>
            <th>Used At</th>
            <th>Created At</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="cred in credentials" :key="cred.id" class="hover:bg-gray-700/50">
            <td class="font-mono">{{ cred.username }}</td>
            <td>{{ cred.expiresAt ? formatDateTime(cred.expiresAt) : 'Never' }}</td>
            <td>
              <span :class="['badge badge-sm', cred.used ? 'badge-success' : 'badge-warning']">
                {{ cred.used ? 'Yes' : 'No' }}
              </span>
            </td>
            <td>{{ cred.usedAt ? formatDateTime(cred.usedAt) : '-' }}</td>
            <td>{{ formatDateTime(cred.createdAt) }}</td>
          </tr>
        </tbody>
      </table>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';
import { useRoute } from 'vue-router';
import { getAdminQuizCredentials, generateAdminQuizCredentials } from '@/services/api';
import type { ResponderCredential } from '@/types';

const route = useRoute();
const quizId = ref<number | null>(null);
const credentials = ref<ResponderCredential[]>([]);
const newlyGeneratedCredentials = ref<ResponderCredential[]>([]);
const isLoading = ref(false);
const isGenerating = ref(false);
const fetchError = ref<string | null>(null);
const generationError = ref<string | null>(null);

const generateForm = reactive({
  count: 1,
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
  if (!quizId.value || generateForm.count < 1) return;
  isGenerating.value = true;
  generationError.value = null;
  newlyGeneratedCredentials.value = []; // Clear previous results

  try {
    const payload = {
      count: generateForm.count,
      // Only include expiryHours if it's a positive number
      ...(generateForm.expiryHours && generateForm.expiryHours > 0
        ? { expiryHours: generateForm.expiryHours }
        : {}),
    };
    const generated = await generateAdminQuizCredentials(quizId.value, payload);
    newlyGeneratedCredentials.value = generated;
    // Optionally refresh the main list immediately or show a success message
    // For simplicity, we'll just show the newly generated ones for now.
    // await fetchCredentials(); // Uncomment to refresh the main list
    generateForm.count = 1; // Reset form
    generateForm.expiryHours = undefined;
  } catch (error: any) {
    console.error('Error generating credentials:', error);
    generationError.value = error.message || 'Failed to generate credentials.';
  } finally {
    isGenerating.value = false;
  }
};

const copyGeneratedUsernames = () => {
  const usernames = newlyGeneratedCredentials.value.map(c => c.username).join('\n');
  navigator.clipboard.writeText(usernames)
    .then(() => {
      // Optional: show a success message
      console.log('Usernames copied to clipboard');
    })
    .catch(err => {
      console.error('Failed to copy usernames: ', err);
      // Optional: show an error message
    });
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

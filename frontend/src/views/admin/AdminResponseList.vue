<template>
  <div class="p-6">
    <h1 class="text-2xl font-semibold mb-6 text-gray-100">Quiz Responses</h1>
    <p class="text-gray-400">Placeholder for displaying quiz responses (Quiz ID: {{ quizId }}).</p>
    <!-- Response listing logic will go here -->
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const quizId = ref<number | null>(null);

onMounted(() => {
  const idParam = route.params.id;
  if (typeof idParam === 'string') {
    const parsedId = parseInt(idParam, 10);
    if (!isNaN(parsedId)) {
      quizId.value = parsedId;
      // Fetch responses based on quizId here
    }
  }
});
</script>

<style scoped>
/* Styles for response list */
</style>

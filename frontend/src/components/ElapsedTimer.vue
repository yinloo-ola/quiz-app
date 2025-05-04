<template>
  <div class="elapsed-timer">
    <div class="timer-header">
      <div class="timer-label">Time Elapsed:</div>
      <div class="timer-progress-container">
        <div class="timer-progress" :style="{ width: progressWidth }"></div>
      </div>
    </div>
    <div class="time-display" :class="{ 'milestone-1': milestone1, 'milestone-2': milestone2, 'milestone-3': milestone3 }">
      <span class="time-icon" aria-hidden="true">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd" />
        </svg>
      </span>
      {{ formattedTime }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps<{
  startTime?: string; // Optional ISO timestamp when the quiz was started
}>();

const emit = defineEmits<{
  (e: 'update-elapsed', seconds: number): void;
}>();

const secondsElapsed = ref(0);
const timerInterval = ref<number | null>(null);

// If startTime was provided, calculate initial elapsed time
onMounted(() => {
  if (props.startTime) {
    const startTimeMs = new Date(props.startTime).getTime();
    const currentTimeMs = new Date().getTime();
    secondsElapsed.value = Math.floor((currentTimeMs - startTimeMs) / 1000);
  }
  startTimer();
});

// Computed properties for formatting and styling
const formattedTime = computed(() => {
  const hours = Math.floor(secondsElapsed.value / 3600);
  const minutes = Math.floor((secondsElapsed.value % 3600) / 60);
  const seconds = secondsElapsed.value % 60;
  
  if (hours > 0) {
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  } else {
    return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
  }
});

// Visual indicators for time milestones
const milestone1 = computed(() => secondsElapsed.value >= 60 && secondsElapsed.value < 300);  // 1-5 minutes
const milestone2 = computed(() => secondsElapsed.value >= 300 && secondsElapsed.value < 600); // 5-10 minutes
const milestone3 = computed(() => secondsElapsed.value >= 600); // 10+ minutes

// Progress bar width calculation (max at 20 minutes = 1200 seconds)
const progressWidth = computed(() => {
  const maxSeconds = 1200; // 20 minutes
  const percentage = Math.min((secondsElapsed.value / maxSeconds) * 100, 100);
  return `${percentage}%`;
});

// Start the timer
const startTimer = () => {
  if (timerInterval.value !== null) return; // Don't start if already running
  
  timerInterval.value = window.setInterval(() => {
    secondsElapsed.value++;
    emit('update-elapsed', secondsElapsed.value);
  }, 1000);
};

// Stop the timer
const stopTimer = () => {
  if (timerInterval.value !== null) {
    clearInterval(timerInterval.value);
    timerInterval.value = null;
  }
};

// Reset the timer
const resetTimer = () => {
  stopTimer();
  secondsElapsed.value = 0;
};

// Get current elapsed seconds
const getElapsedSeconds = () => {
  return secondsElapsed.value;
};

// Expose methods to parent component
defineExpose({
  startTimer,
  stopTimer,
  resetTimer,
  getElapsedSeconds
});

// Clean up on component unmount
onBeforeUnmount(() => {
  stopTimer();
});
</script>

<style scoped>
.elapsed-timer {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.timer-header {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.timer-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #4b5563;
}

.timer-progress-container {
  width: 100%;
  height: 4px;
  background-color: #e5e7eb;
  border-radius: 2px;
  overflow: hidden;
}

.timer-progress {
  height: 100%;
  background-color: #4f46e5;
  transition: width 1s linear;
}

.time-display {
  font-family: monospace;
  font-size: 1.5rem;
  font-weight: bold;
  padding: 0.5rem 0.75rem;
  border-radius: 0.375rem;
  background-color: #dbeafe;
  color: #1e40af;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.time-icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.time-display.milestone-1 {
  background-color: #e0e7ff;
  color: #4338ca;
}

.time-display.milestone-2 {
  background-color: #c7d2fe;
  color: #4f46e5;
}

.time-display.milestone-3 {
  background-color: #a5b4fc;
  color: #4f46e5;
  animation: subtle-pulse 2s infinite;
}

@keyframes subtle-pulse {
  0% { opacity: 1; }
  50% { opacity: 0.9; }
  100% { opacity: 1; }
}
</style>

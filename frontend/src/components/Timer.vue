<template>
  <div class="timer">
    <div class="time-display" :class="{ 'warning': isWarning, 'danger': isDanger }">
      {{ formattedTime }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps<{
  initialTime: number; // Time in seconds
}>();

const emit = defineEmits<{
  (e: 'time-up'): void;
}>();

const secondsRemaining = ref(props.initialTime);
const timerInterval = ref<number | null>(null);

// Computed properties for formatting and styling
const formattedTime = computed(() => {
  const minutes = Math.floor(secondsRemaining.value / 60);
  const seconds = secondsRemaining.value % 60;
  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
});

const isWarning = computed(() => {
  return secondsRemaining.value <= 300 && secondsRemaining.value > 60; // Warning when 5 minutes or less
});

const isDanger = computed(() => {
  return secondsRemaining.value <= 60; // Danger when 1 minute or less
});

// Start the timer
const startTimer = () => {
  if (timerInterval.value !== null) return; // Don't start if already running
  
  timerInterval.value = window.setInterval(() => {
    if (secondsRemaining.value <= 0) {
      stopTimer();
      emit('time-up');
    } else {
      secondsRemaining.value--;
    }
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
  secondsRemaining.value = props.initialTime;
};

// Expose methods to parent component
defineExpose({
  startTimer,
  stopTimer,
  resetTimer
});

// Start timer on component mount
onMounted(() => {
  startTimer();
});

// Clean up on component unmount
onBeforeUnmount(() => {
  stopTimer();
});
</script>

<style scoped>
.timer {
  display: flex;
  justify-content: center;
  align-items: center;
}

.time-display {
  font-family: monospace;
  font-size: 1.5rem;
  font-weight: bold;
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  background-color: #e5e7eb;
  color: #1f2937;
}

.time-display.warning {
  background-color: #fef3c7;
  color: #92400e;
}

.time-display.danger {
  background-color: #fee2e2;
  color: #b91c1c;
  animation: pulse 1s infinite;
}

@keyframes pulse {
  0% { opacity: 1; }
  50% { opacity: 0.8; }
  100% { opacity: 1; }
}
</style>

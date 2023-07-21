<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {Remote} from "udap-ui/remote";
import Countdown from "@/components/Countdown.vue";

const props = defineProps<{
  remote: Remote
  secondary?: boolean
}>()


</script>

<template>
  <div class="lh-1">
    <Countdown :more="remote.client.attempts"
               :percent=" remote.client.nextAttempt/2000 "></Countdown>
  </div>
</template>

<style scoped>
.pulse {
  animation: pulse-in 1s ease-in-out infinite;
}

@keyframes pulse-in {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 0.6;
    filter: blur(0.5px);
  }
  100% {
    opacity: 1;
  }
}

.surface {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 1rem 0.25rem;
}
</style>
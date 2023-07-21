<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {reactive, watchEffect} from "vue";

const props = defineProps<{
  until?: number,
  since?: number,
}>()

const state = reactive({
  display: "0",
  value: 0,
  target: 0,
})

watchEffect(() => {
  if (props.until) {
    state.target = props.until;
    updateDisplay()
  }
})

watchEffect(() => {
  if (props.since) {
    state.target = props.since;
    updateDisplay()
  }
})

function updateDisplay() {
  state.display = convertNanosecondsDurationToString(state.value)
}

function convertNanosecondsDurationToString(nanoseconds: number): string {
  if (nanoseconds < 0) {
    nanoseconds = Math.abs(nanoseconds)
    // throw new Error("Input must be a non-negative number.");
  }

  if (nanoseconds === 0) {
    return "0ns";
  }

  const units = ["ns", "µs", "ms", "s", "ks"];
  let value = nanoseconds;
  let unitIndex = 0;

  while (value >= 1000 && unitIndex < units.length - 1) {
    if (unitIndex >= 3) {
      let seconds = value % 60
      let minutes = (value / 60) % 60
      let hours = (value / 60 / 60) % 60
      let days = Math.round((value / 60 / 60 / 24) % 24)
      let years = Math.round((value / 60 / 60 / 24 / 365))
      return `${years > 0 ? (years + 'y') : ''} ${(days > 0) ? (Math.round(days) + 'd') : ''} ${Math.round(hours)}h  ${Math.round(minutes)}m ${Math.round(seconds)}s`;

    }
    value /= 1000;
    unitIndex++;
  }

  return `${value.toFixed(2)} ${units[unitIndex]}`;
}
</script>

<template>
  <div>
    {{ state.display }}
  </div>
</template>

<style lang="scss" scoped>

</style>
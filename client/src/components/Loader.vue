<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
interface Loader {
  size?: string
}

let props = defineProps<Loader>()
</script>

<template>
  <span class="loader"></span>
</template>

<style lang="scss" scoped>
.loader {
  width: 0.6rem;
  height: 0.6rem;
  border: 3px solid rgba(128, 128, 128, 0.8);
  border-bottom-color: transparent;
  border-radius: 50%;
  display: inline-block;
  box-sizing: border-box;
  animation: rotation 1s linear infinite;
}

@keyframes rotation {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}
</style>
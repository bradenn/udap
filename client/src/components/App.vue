<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

interface App {
  name: string
  icon?: string
}

let props = defineProps<App>()

</script>

<template>
  <div class="app-container">
    <div class="app-icon element">
      <i v-if="props.icon" :class="props.icon" class="fa-solid fa-fw"></i>
      <i v-else class="fa-solid fa-circle fa-fw"></i>
    </div>
    <div class="app-name">{{ props.name }}</div>
  </div>
</template>

<style lang="scss" scoped>


.app-container {
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
}

.app-icon {

  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-self: center;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1/1 !important;
  border-radius: 0.6rem;
  font-size: 0.9rem;
  color: rgba(200, 200, 200, 0.6);
}

.app-container:active {


}

.app-name {
  font-size: 0.6rem;
  font-family: "SF Pro Display", sans-serif;
  font-weight: 400;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.55)
}
</style>
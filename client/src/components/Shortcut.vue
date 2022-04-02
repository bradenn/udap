<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

interface App {
  name: string
  icon?: string
}

let props = defineProps<App>()

</script>

<template>
  <div class="shortcut-container">
    <div class="shortcut-icon element">

      <div class="shortcut-name">
        <i v-if="props.icon" :class="props.icon" class="fa-solid fa-fw"></i>
        <i v-else class="fa-solid fa-circle fa-fw"></i>
        {{ props.name }}
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>


.shortcut-container {
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;
}

.shortcut-icon {
  height: 2.25rem;
  flex-grow: 1;
  width: 100%;

  display: flex;
  align-self: center;
  align-items: center;
  justify-content: center;

  border-radius: 0.6rem;
  font-size: 0.8rem;
  color: rgba(200, 200, 200, 0.6);
}

.shortcut-container:active {


}

.shortcut-name {
  font-size: 0.7rem;
  line-height: 0.5rem;
  font-family: "SF Pro Display", sans-serif;
  font-weight: 400;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.55)
}
</style>
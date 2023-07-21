<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

const props = defineProps<{
  name: string,
  icon: string,
  to: string,
}>()

</script>

<template>
  <router-link :to="props.to" active-class="element" class="surface px-3" style="border-radius: 8px">
    <div class="sf-icon label-o4 lh-1">{{ props.icon }}</div>
    <div class="app-link-header lh-1">{{ props.name }}</div>
  </router-link>
</template>

<style scoped>

.surface {
  border: none;
  background: rgba(22, 22, 22, 0.6);
  backdrop-filter: blur(40px);
  box-shadow: inset 0 0 1px 1.5px rgba(37, 37, 37, 0.6), 0 0 3px 1px rgba(22, 22, 22, 0.6);
  /* Note: backdrop-filter has minimal browser support */
  aspect-ratio: 2.71828183/0.8;
  border-radius: 11.5px;
  -webkit-backdrop-filter: blur(40px) !important;
  display: flex;
  justify-content: center;


  transition: background-color 100ms;

}

.router-link-active {
  outline: 1px solid red;
}


a {
  text-decoration: none;
}

.app-link-header {
  font-size: 0.95rem;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.65);

}
</style>
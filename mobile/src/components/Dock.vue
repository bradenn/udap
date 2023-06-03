<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
const props = defineProps<{
  aspect?: number
  secondary?: boolean
}>()
</script>

<template>
  <div class="element dock">
    <router-link active-class="dock-on" class="subplot" to="/home/dashboard">
      <div class="sf-icon">􀎟</div>
      <div class="label-c4">Home</div>
    </router-link>
    <router-link active-class="dock-on" class="subplot" to="/home/macros">
      <div class="sf-icon">􀞿</div>
      <div class="label-c4">Macros</div>
    </router-link>
    <router-link active-class="dock-on" class="subplot" to="/home/monitor">
      <div class="sf-icon">􁀘</div>
      <div class="label-c4">Monitor</div>
    </router-link>
  </div>
</template>

<style scoped>
.subplot {
  height: 100%;
  flex-grow: 1;
  border-radius: 0.3rem !important;
  font-weight: 800;
  background-color: transparent !important;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
//box-shadow: 0 0 8px 2px rgba(0, 0, 0, 0.05); //border: 1px solid rgba(255, 255, 255, 0.05); color: rgba(255, 255, 255, 0.55); font-weight: 500; font-size: 1rem; background-color: rgba(255, 255, 255, 0.025);
}

a {
  text-decoration: none !important;
}

.subplot.dock-on {
  background-color: rgba(46, 46, 48, 0.4) !important;
  color: rgba(255, 255, 255, 0.9) !important;

  .sf-icon {
    color: rgba(255, 255, 255, 0.9) !important;
    filter: drop-shadow(0 0 4px rgba(255, 255, 255, 0.2)) !important;
  }

}


.dock {
  gap: 0.5rem;
  height: 4rem;
  padding: 0.5rem;
  display: flex;
  justify-content: space-between;
}
</style>
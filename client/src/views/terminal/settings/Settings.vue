<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup></script>

<template>

  <div class="d-flex flex-row align-content-start align-items-start gap h-100">
    <!--    <Plot :rows="8" :cols="1" class="gap-0" style="width: 10rem;">-->
    <!--      <Subplot v-for="route in $route.matched[1].children" :icon="route.icon || 'bars-progress'" :name="route.name"-->
    <!--               :to="route.path"></Subplot>-->
    <!--    </Plot>-->

    <div class="flex-grow-1">
      <router-view/>
    </div>
  </div>
</template>

<style scoped>

.sidebar-container {
  width: 14rem;
}
</style>
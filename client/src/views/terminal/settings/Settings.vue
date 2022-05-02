<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup></script>

<template>

  <div class="d-flex flex-row align-content-start align-items-start gap h-100">
    <!--    <Plot :rows="8" :cols="1" class="gap-0" style="width: 10rem;">-->
    <!--      <Subplot v-for="route in $route.matched[1].children" :icon="route.icon || 'bars-progress'" :name="route.name"-->
    <!--               :to="route.path"></Subplot>-->
    <!--    </Plot>-->

    <div class="flex-grow-1">
      <router-view v-slot="{ Component }">
        <div class="animateIn">
          <component :is="Component"/>
        </div>
      </router-view>
    </div>
  </div>
</template>

<style scoped>


/*
  Enter and leave animations can use different
  durations and timing functions.
*/
.slide-fade-enter-active {
  transition: all 600ms ease-out;
}

.slide-fade-leave-active {
  transition: all 300ms cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  position: absolute;
  transform: scale(0.98);
}

.sidebar-container {
  width: 14rem;
}
</style>
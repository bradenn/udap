<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Plot from "@/components/plot/Plot.vue"</script>
<template>
  <div class="h-75">
    <Plot :cols="3" :rows="1">
      <router-link class="subplot justify-content-center" draggable="false" to="/terminal/timing/stopwatch">
        <i class="fa-solid fa-stopwatch fa-fw"></i>&nbsp;Stopwatch
      </router-link>
      <router-link class="subplot justify-content-center" draggable="false" to="/terminal/timing/timer">
        <i class="fa-solid fa-hourglass fa-fw"></i>&nbsp;Timer
      </router-link>
      <router-link class="subplot justify-content-center" draggable="false" to="/terminal/timing/conversion">
        <i class="fa-solid fa-0 fa-fw"></i>&nbsp;Conversion
      </router-link>
    </Plot>

    <router-view></router-view>
  </div>
</template>

<style scoped>

</style>
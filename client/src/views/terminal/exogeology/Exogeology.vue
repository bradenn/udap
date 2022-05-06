<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup></script>

<template>

  <div class="d-flex flex-row align-content-start align-items-start gap h-100">
    <!--    <div class="d-flex flex-column">-->
    <!--      <Plot :rows="1" :cols="3" style="width: 10rem;">-->
    <!--        <Subplot v-for="route in $route.matched[1].children" :icon="route.icon || 'earth-americas'" :name="route.name"-->
    <!--                 :to="route.path"></Subplot>-->
    <!--      </Plot>-->
    <!--    </div>-->
    <div class="w-100 h-100">
      <router-view/>
    </div>
  </div>
</template>

<style scoped>

</style>
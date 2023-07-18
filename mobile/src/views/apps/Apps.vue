<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";


const apps = [
  {
    name: "Monitor",
    link: "/home/monitor",
    icon: '􁂥'
  },
  {
    name: "Beam Sim",
    link: "/home/beampath",
    icon: '􁀘'
  },
  {
    name: "Lasers",
    link: "/home/sentry",
    icon: '􀑃'
  },
  {
    name: "Lights",
    link: "/home/lights",
    icon: '􁓼'
  }
]

</script>

<template>

  <div class="app-list">
    <div v-for="app in apps" :key="app.link" class="d-flex flex-column align-items-center justify-content-center">
      <Element :to="app.link" class="d-flex align-items-center justify-content-center"
               mutable
               style="aspect-ratio: 1/1; width: 100%; box-shadow: inset 0 0 1px 1px rgba(255,255,255,0.08) !important;">
        <div class="sf label-c0 label-o4" style="font-size: 1.8rem">{{ app.icon }}</div>
      </Element>
      <div class="label-c5 label-o6 label-w600">{{ app.name }}</div>
    </div>
  </div>
</template>

<style lang="scss">

.app-list {
  display: grid;
  grid-template-rows: repeat(8, minmax(1rem, 1fr));
  grid-template-columns: repeat(4, minmax(2rem, 6rem));
  grid-gap: 1.5rem;
}
</style>
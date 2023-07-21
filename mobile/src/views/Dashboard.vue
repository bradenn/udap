<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";
import type {Entity as EntityType, Zone} from "udap-ui/types";
// @ts-ignore
import {registerSW} from "virtual:pwa-register";

const remote = core.remote()


const state = reactive({
  entities: [] as EntityType[],
  zones: [] as Zone[],
})

onMounted(() => {
  state.entities = remote.entities
  state.zones = remote.zones.filter(z => z.pinned)

})

watchEffect(() => {
  state.zones = remote.zones.filter(z => z.pinned)
  return remote.entities
})
if ("serviceWorker" in navigator) {
  // && !/localhost/.test(window.location)) {
  registerSW();
}

</script>

<template>

  <div>
    <router-view></router-view>
  </div>
</template>

<style scoped>

.home-grid {
  display: grid;
  gap: 0.25rem;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  grid-template-rows: repeat(1, minmax(0, 1fr));
}
</style>
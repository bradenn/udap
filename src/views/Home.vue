<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";
import type {Entity as EntityType, Zone} from "udap-ui/types";
import Entity from "@/components/Entity.vue";
import ThermostatV2 from "@/components/ThermostatV2.vue";

const remote = core.remote()

const state = reactive({
  entities: [] as EntityType[],
  zones: [] as Zone[],
  thermostat: {} as EntityType,
})

onMounted(() => {
  state.entities = remote.entities
  state.zones = remote.zones.filter(z => z.pinned)
})

watchEffect(() => {
  state.zones = remote.zones.filter(z => z.pinned)
  let e = remote.entities.find(e => e.name === "thermostat")
  if (e) state.thermostat = e
  return remote.entities
})


</script>

<template>
  <div class="d-flex flex-column gap-3 justify-content-between h-100">
    <div>
      <div class="label-c5 label-w700 label-o5 px-2">Climate</div>
      <div class="w-100">
        <ThermostatV2 :entity="state.thermostat"></ThermostatV2>
      </div>
    </div>
    <div v-for="zone in state.zones">
      <div>
        <div class="label-c5 label-w700 label-o5 px-2">{{ zone.name }}</div>
        <div class="home-grid">
          <Entity v-for="e in zone.entities" :key="e.id" :entity="e"></Entity>
        </div>
      </div>
    </div>


  </div>
</template>

<style scoped>

.home-grid {
  display: grid;
  gap: 0.25rem;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  grid-template-rows: repeat(1, minmax(0, 1fr));
}
</style>
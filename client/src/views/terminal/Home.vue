<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Light from "@/components/Light.vue";
import {inject, onMounted, reactive, watch} from "vue";
import Weather from "@/components/Weather.vue";

// Define the local reactive data for this view
let state = reactive({
  lights: []
})

// Compare the names of the entities to sort them accordingly
function compareName(a: any, b: any): number {
  if (a.name < b.name)
    return -1;
  if (a.name > b.name)
    return 1;
  return 0;
}

// Inject the remote udap context
const remote: any = inject('remote')

// Force the light state to be read on load
onMounted(() => {
  updateLights(remote.entities)
})

// Update the Lights based on the remote injection changes
watch(remote.entities, (newEntities, oldEntities) => {
  updateLights(newEntities)
})

// Update the current set of lights based on the entities provided
function updateLights(entities: any) {
  // Find all applicable entities
  let candidates = entities.filter((f: any) => f.type === 'spectrum' || f.type === 'switch' || f.type === 'dimmer');
  // Sort and assign them to the reactive object
  state.lights = candidates.sort(compareName)
}

</script>

<template>
  <div>
    <div class="cluster gap">
      <div v-for="light in state.lights">
        <Light :entity="light"></Light>
      </div>
    </div>
    <div class="cluster gap">
      <Weather></Weather>
    </div>
  </div>
</template>

<style scoped>


</style>

<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted} from "vue";
import PaneMenu from "@/components/pane/PaneMenu.vue";
import PaneMenuItem from "@/components/pane/PaneMenuItem.vue";
import type {Remote} from "@/types";
import Slider from "@/components/Slider.vue";
import Switch from "@/components/Switch.vue";

onMounted(() => {

})

onMounted(() => {

})

const remote = inject("remote") as Remote

</script>

<template>
  <div class="remote-grid">
    <PaneMenu title="Devices">
      <PaneMenuItem v-for="entity in remote.entities" :active="false" :subtext="entity.module" :title="entity.name"
                    to=""></PaneMenuItem>
    </PaneMenu>
    <div class="d-flex flex-column gap-1 remote-controls ">
      <Switch :active="false" name="Power"></Switch>
      <Slider :max="100" :min="0" :step="5" name="Brightness" unit="%"></Slider>

    </div>
  </div>
</template>

<style scoped>

.remote-controls {
  grid-column: span 2 / 4;
}

.remote-grid {
  display: grid;
  grid-column-gap: 0.5rem;
  grid-row-gap: 0.5rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}
</style>
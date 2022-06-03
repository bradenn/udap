<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote, Zone} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import CreateZone from "@/views/terminal/settings/zone/CreateZone.vue";
import ZonePreview from "@/components/zone/ZonePreview.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  zones: {} as Zone[],
  deleted: {} as Zone[],
  showDeleted: false,
  loading: true,
  mode: "list",
  model: "",
})


onMounted(() => {
  state.loading = true
  handleUpdates(remote.zones)
})
watchEffect(() => handleUpdates(remote.zones))

function handleUpdates(zones: Zone[]) {
  state.zones = zones.filter(z => state.showDeleted ? true : !z.deleted)
  state.loading = false
}

function setMode(mode: string) {
  state.mode = mode
}

function toggleShowDeleted() {
  state.showDeleted = !state.showDeleted
}

</script>

<template>
  <div>
    <div class="d-flex justify-content-start py-2 px-0">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-map fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Zones</div>
      <div class="flex-fill"></div>

      <Plot :cols="2" :rows="1" small style="width: 13rem;">


        <Radio :active="false" :fn="() => toggleShowDeleted()"
               :title="state.showDeleted?'Hide Deleted':'Show Deleted'"></Radio>
        <Radio :active="false" :fn="() => setMode('create')"
               :title="state.mode === 'create'?'Cancel':'New Zone'"></Radio>
      </Plot>
    </div>
    <div v-if="state.loading">
      <div class="element p-2">
        <div class="label-c1 label-o4 d-flex align-content-center gap-1">
          <div>
            <Loader size="sm"></Loader>
          </div>
          <div class="">Loading...</div>
        </div>
      </div>
    </div>
    <div v-else-if="state.mode === 'list'">
      <div class="zones">
        <div v-for="zone in state.zones">
          <ZonePreview :key="zone.id" :zone="zone"></ZonePreview>
        </div>
      </div>

    </div>
    <div v-else-if="state.mode === 'create'" class="d-flex justify-content-center btn-outline-primary">
      <CreateZone :done="() => state.mode = 'list'">
      </CreateZone>
    </div>
    <div v-else>


    </div>
  </div>
</template>

<style lang="scss" scoped>
.zones {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-rows: repeat(4, 1fr);
  grid-template-columns: repeat(4, 1fr);
}

</style>
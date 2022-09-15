<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Device, Endpoint, Remote} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import Loader from "@/components/Loader.vue";
import PaneList from "@/components/pane/PaneList.vue";
import ListInput from "@/components/pane/ListInput.vue";
import PaneMenuItem from "@/components/pane/PaneMenuItem.vue";
import endpointService from "@/services/endpointService";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  endpoints: {} as Endpoint[],
  devices: {} as Device[],
  loading: true,
  toCreate: {
    name: ""
  } as Endpoint,
  mode: "list"
})


onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.endpoints = remote.endpoints
  state.devices = remote.devices
  state.loading = false
  return remote
}

function setMode(mode: string) {
  state.mode = mode
}

function createEndpoint() {
  endpointService.createEndpoint(state.toCreate).then(res => {
    state.toCreate.name = `${res}`
  }).catch(err => {
    state.toCreate.name = `${err}`
  })
}

</script>

<template>
  <div v-if="!state.loading">
    <Plot :cols="1" :rows="1" class="mb-1" small style="width: 6rem;">
      <Radio :active="false" :fn="() => setMode(state.mode === 'create'?'list':'create')"
             :title="state.mode !== 'list'?'Cancel':'New Endpoint'"></Radio>
    </Plot>
    <div v-if="state.mode === 'list'">

      <div class="endpoint-container w-100">
        <div v-for="endpoint in state.endpoints"
             :key="endpoint.id" class="">
          <Plot :alt="endpoint.key" :cols="3" :rows="1" :title="endpoint.name">
            <div></div>
          </Plot>

        </div>

      </div>
    </div>
    <div v-else-if="state.mode === 'create'">

      <PaneList style="width: 15rem;" title="Configuration">
        <ListInput :change="(s: string) => state.toCreate.name = s"
                   :value="state.toCreate.name"
                   description="The name of the endpoint"
                   name="Endpoint Name" type="text"></ListInput>
        <PaneMenuItem :active="false" :fn="createEndpoint" subtext="" title="Create"></PaneMenuItem>
      </PaneList>

    </div>
  </div>

  <div v-else>
    <div class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-expand fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Endpoints</div>
      <div class="flex-fill"></div>

    </div>
    <div class="element p-2">
      <div class="label-c1 label-o4 d-flex align-content-center gap-1">
        <div>
          <Loader size="sm"></Loader>
        </div>
        <div class="">Loading...</div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.endpoint-container {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(2, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

</style>
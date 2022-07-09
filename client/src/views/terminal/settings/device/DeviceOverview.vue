<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Device, Remote} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import Loader from "@/components/Loader.vue";
import axios from "axios";
import Subplot from "@/components/plot/Subplot.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  devices: {} as Device[],
  histories: new Map<string, number[]>(),
  maxes: new Map<string, number>(),
  loading: true,
  mode: "list"
})


onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function sortAlpha(a: Device, b: Device) {
  if (a.name > b.name) {
    return -1
  } else {
    return 1
  }
}

function handleUpdates(remote: Remote) {
  state.devices = remote.devices
  state.devices = state.devices.sort(sortAlpha)
  state.devices.forEach(t => {
    let local = state.histories.get(t.mac) || []
    let cand = t.latency / 1000
    if (local.includes(cand)) return
    local.push(cand)
    if (local.length > 20) {
      local = local.slice(1)
    }
    state.histories.set(t.mac, local)
  })
  state.loading = false
  return remote
}

function setMode(mode: string) {
  state.mode = mode
}

function isOnline(ns: string): boolean {
  return (new Date().valueOf() - new Date(ns).valueOf()) < 30 * 1000
}

function nsToMs(ns: number): number {
  return Math.round(ns / 1000 / 1000 * 10) / 10
}

function rename(name: string, device: Device) {
  device.name = name
  axios.post("http://10.0.1.18:3020/devices/update", JSON.stringify(device))
}

</script>

<template>
  <div v-if="!state.loading">
    <div class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-expand fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Devices</div>
      <div class="flex-fill"></div>

      <Plot :cols="1" :rows="1" small style="width: 6rem;">
        <Radio :active="false" :fn="() => setMode(state.mode === 'create'?'list':'create')"
               :title="state.mode === 'create'?'Cancel':'New Endpoint'"></Radio>
      </Plot>
    </div>
    <div v-if="state.mode === 'list'">

      <div class="device-container">
        <Plot v-for="device in state.devices" :cols="2" :rows="1">
          <div class="subplot subplot-inline justify-content-between px-0 w-100">
            <div class="d-flex align-items-center">
              <div :style="`background-color: rgba(${device.state==='ONLINE'?'25, 135, 84':'135, 100, 2'}, 0.53);`"
                   class="status-marker"></div>
              <div>
                <div class="label-c1 label-o4 label-r lh-1 w-100">
                  <div>{{ device.name || device.ipv4 }}</div>
                </div>
                <div class="label-c4  label-o3 label-r py-0 w-100 d-flex justify-content-between"
                     style="line-height: 0.55rem">
                  {{ device.mac }}

                </div>


              </div>
            </div>

            <div class="label-c3 label-o4 d-flex flex-column justify-content-end align-items-end">
              <div :class="`${ device.state === 'ONLINE'?'text-success':''}`" class="label-o3">
                &nbsp;{{ device.state === "ONLINE" ? `${nsToMs(device.latency)} ms` : device.state }}
              </div>
              <div class="d-flex gap-1">

                <div v-if="state.histories" class="label-c3 label-o3 d-flex flex-row align-items-end time-marker-line">

                  <div v-for="marker in state.histories.get(device.mac)?.map(d => d / 1000)"
                       :style="`height:${Math.log(marker)}px;`"
                       class="time-marker"></div>
                </div>
              </div>

            </div>

          </div>
          <div class="d-flex gap-1 text-success justify-content-center">
            <Subplot :active="true" :fn="() => $router.push(`/terminal/settings/devices/${device.id}`)"
                     name="Configure" sf="cog" style="width: 100%;"></Subplot>

          </div>
        </Plot>


      </div>
    </div>
    <div v-else-if="state.mode === 'create'">
      Create mode!
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

.time-marker {
  width: 2px;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  height: 1px;
}

.time-marker-line {
  display: flex;
  flex-direction: row;
  justify-content: center;
  height: 20px;
  width: 75px;
  align-items: center;
  gap: 1px;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}

.overflow-ellipse {
  display: block;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis !important;
  text-wrap: none !important;
  max-width: 34rem;
}

.device-container {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-auto-flow: column;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.status-marker {
  width: 4px !important;
  border-radius: 4px;
  height: 28px;

  margin-right: 14px;
  margin-left: 8px;


  background-color: rgba(25, 135, 84, 0.53);
}

.device-unit {
  padding: 0.25rem 0.5rem;
  display: flex;
  align-content: center;
  align-items: center;
}


</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Device, Remote, Utilization} from "@/types";
import {useRoute} from 'vue-router'
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";

interface NewZone {
  name: string
  user: string
  entities: string[]
}

const route = useRoute()
let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  loading: true,
  util: {} as Utilization,
  model: "",
  mode: "toggles",
  device: {} as Device,
  name: "",
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.device = remote.devices.find(device => device.id === route.params.device) as Device
  if (!state.device) return
  state.util = state.device.utilization
  state.loading = false
  return remote.devices
}

interface CreateZoneProps {
  done: () => {}
}

let props = defineProps<CreateZoneProps>();

function enterChar(char: string) {
  // If the incoming char is a backspace, decrement the cursor and clear the value.
  if (char === "{bksp}") {
    // Only decrement the cursor if it is bigger than zero
    if (state.name.length > 0) state.name = state.name.slice(0, state.name.length - 1)
    // Exit the function
    return
  }

  char = char.replace("{space}", " ")
  // Return if the security code block is full
  if (state.name.length >= 64) return
  // Add the provided char to the cursor position
  state.name += char
}

function nextStep() {
  state.mode = 'select'
}

// https://stackoverflow.com/questions/10420352/converting-file-size-in-bytes-to-human-readable-string
function bytesToString(bytes: number, si: boolean = false, dp: number = 2): string {
  const thresh = si ? 1000 : 1024;

  if (Math.abs(bytes) < thresh) {
    return bytes + ' B';
  }

  const units = si
      ? ['kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
      : ['KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB'];
  let u = -1;
  const r = 10 ** dp;

  do {
    bytes /= thresh;
    ++u;
  } while (Math.round(Math.abs(bytes) * r) / r >= thresh && u < units.length - 1);


  return bytes.toFixed(dp) + ' ' + units[u];
}


</script>

<template>
  <div v-if="!state.loading">
    <div class="d-flex justify-content-start py-2 gap-2">
      <Plot :cols="1" :rows="1">
        <Subplot :active="state.mode === 'toggles'" :fn="() => $router.push(`/terminal/settings/devices`)"
                 icon="chevron-left"
                 name="Back"></Subplot>
      </Plot>
      <div class="label-w500 label-o4 label-xxl lh-1"><i :class="`fa-solid fa-share-nodes fa-fw`"></i></div>

      <div>
        <div class="label-lg label-o4 label-r lh-1 w-100">
          <div>{{ state.device.name || "Unnamed" }}</div>
        </div>
        <div class="label-xs  label-o3 label-r py-0 w-100 d-flex justify-content-between"
             style="line-height: 1rem">
          {{ state.device.mac }} -
          {{ state.device.ipv4 }}

        </div>


      </div>
      <div class="flex-fill"></div>
      <Plot :cols="3" :rows="1" style="width: 13rem;">
        <Subplot :active="state.mode === 'toggles'" :fn="() => state.mode = 'toggles'"
                 name="Toggles"></Subplot>
        <Subplot :active="state.mode === 'name'" :fn="() => state.mode = 'name'"
                 name="Rename"></Subplot>
      </Plot>
    </div>
    <div class="device-container">

      <Plot :alt="`${state.util.cpu.cores} Cores`" :cols="2" :rows="4" title="CPU">
        <div v-for="(value, core) in state.util.cpu.usage" class="subplot">
          <div class="label-c2 label-o3 label-w600" style="width: 0.5rem">{{ core + 1 }}</div>&nbsp;
          <div class="tick-bar flex-grow-1">


            <div class="tick-overlay">
              <div class="d-flex justify-content-end">{{ Math.round(value * 10) / 10 }}%</div>
            </div>
            <div :style="`width:${value}%;`" class="tick"></div>
          </div>
        </div>
      </Plot>
      <Plot :cols="1" :rows="1" alt="" title="Memory">
        {{ bytesToString((state.util.memory.used / 100.0) * state.util.memory.total) }} /
        {{ bytesToString(state.util.memory.total) }}
      </Plot>

    </div>
  </div>

</template>

<style lang="scss" scoped>
.tick-overlay {
  position: relative !important;
  height: 0;
  z-index: 10 !important;

  mix-blend-mode: color-dodge;
  padding-right: 0.25rem;
}

.tick-bar {
  justify-content: start !important;
  width: 100%;
}

.tick {

  background-color: rgba(22, 94, 176, 1);
  border-radius: 3px;
  transition: width 200ms ease-out;
  height: 1rem;
}

.usage-meter {
  background-color: #6f42c1;
  height: 100%;
}

.device-container {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-auto-flow: column;
  grid-template-rows: repeat(8, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.cursor-blink {
  height: 1rem;
  width: 2px;
  border-radius: 2px;
  background-color: rgba(255, 255, 255, 0.125);
  margin-left: 2px;
  animation: cursor-flash 1s ease infinite;
}

@keyframes cursor-flash {
  0% {
    opacity: 1;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.125;
  }
}

</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Compute, Device, Remote, Utilization} from "@/types";
import {useRoute} from 'vue-router'
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import Meter from "@/components/Meter.vue";

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
  util: {compute: [] as Compute[]} as Utilization,
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

      <div class="label-w500 label-o4 label-xxl lh-1 pt-1"><i :class="`fa-solid fa-share-nodes fa-fw`"></i></div>

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
    </div>
    <div class="device-container">
      <div class="d-flex flex-column gap-1">
        <Plot :alt="`${state.util.cpu.cores} Cores`" :cols="10" :rows="1" title="CPU Utilization">
          <Meter v-for="(value, core) in state.util.cpu.usage" :danger="95" :percent="value" :thresh="80"
                 :ticks="20" :vertical="true" value="">
          </Meter>
          <Meter v-for="b in Array(10-state.util.cpu.cores)" :blank="true" :percent="0" :ticks="20"
                 :vertical="true"
                 value="">
          </Meter>

          <div v-for="b in Array(10-state.util.cpu.cores)" class="utility-stack-blank-v subplot">
          </div>
        </Plot>

        <Plot
            :alt="`${bytesToString(state.util.memory.total*state.util.memory.used/100)} / ${bytesToString(state.util.memory.total)}`"
            :cols="1" :rows="1"
            title="Memory Usage">
          <Meter :danger="95" :percent="state.util.memory.used" :thresh="80" :ticks="73" :vertical="false"
                 value=""></Meter>
        </Plot>
        <Plot
            :alt="`${bytesToString(state.util.disk.total*state.util.disk.used/100)} / ${bytesToString(state.util.disk.total)}`"
            :cols="1" :rows="1"
            title="Disk Usage">
          <Meter :danger="90" :percent="state.util.disk.used" :thresh="75" :ticks="73" :vertical="false"
                 value=""></Meter>
        </Plot>
      </div>
      <div class="d-flex flex-column gap-1">
        <Plot v-for="gpu in state.util.compute" v-if="state.util.compute"
              :alt="`${gpu.processes.length} Process${gpu.processes.length > 1?'es':''}`"
              :cols="1"
              :rows="3"
              title="GPU Utilization">

          <div class="d-flex w-100 gap-1 h-100 w-100">
            <Meter :alt="`${gpu.utilization.gpu} %`" :percent="gpu.utilization.gpu" :ticks="35"
                   :vertical="false" title="GPU" value="">
            </Meter>

            <Meter :alt="`${gpu.utilization.memory} %`" :percent="gpu.utilization.memory" :ticks="35"
                   :vertical="false" title="Memory" value="">
            </Meter>
          </div>

          <div class="d-flex w-100 gap-1 h-100">
            <Meter :alt="`${gpu.temperature.current} C / ${gpu.temperature.max} C`"
                   :danger="(gpu.temperature.target/gpu.temperature.max)*100"
                   :percent="(gpu.temperature.current/gpu.temperature.max)*100"
                   :thresh="(60/gpu.temperature.max)*100" :ticks="35" :vertical="false"
                   title="Temperature"
                   value="">

            </Meter>
            <Meter :alt="`${gpu.fanSpeed} %`" :danger="90" :percent="gpu.fanSpeed"
                   :thresh="75" :ticks="35" :vertical="false"
                   title="Fan Speed"
                   value="">
            </Meter>
          </div>

          <div class="d-flex w-100 gap-1 h-100">
            <Meter :alt="`${gpu.power.draw} W / ${gpu.power.max} W`" :danger="(200/gpu.power.max)*100"
                   :percent="(gpu.power.draw/gpu.power.max)*100"
                   :thresh="(160/gpu.power.max)*100" :ticks="35" :vertical="false"
                   title="Power"
                   value="">
            </Meter>
            <Meter :alt="`${gpu.memory.used + gpu.memory.reserved} MiB / ${gpu.memory.total} MiB`" :danger="90"
                   :percent="((gpu.memory.used + gpu.memory.reserved)/gpu.memory.total)*100"
                   :thresh="75" :ticks="35"
                   :vertical="false"
                   title="Memory"
                   value="">
            </Meter>
          </div>

        </Plot>

        <Plot v-for="gpu in state.util.compute" v-if="state.util.compute" :cols="1" :rows="2" title="GPU Clocks">
          <div class="d-flex w-100 gap-1 h-100">
            <Meter :alt="`${gpu.clocks.graphics.current} MHz / ${gpu.clocks.graphics.max} MHz`"
                   :percent="(gpu.clocks.graphics.current/gpu.clocks.graphics.max)*100" :ticks="35"
                   :vertical="false" title="Graphics"
                   value="">
            </Meter>

            <Meter :alt="`${gpu.clocks.video.current} MHz / ${gpu.clocks.video.max} MHz`"
                   :percent="(gpu.clocks.video.current/gpu.clocks.video.max)*100" :ticks="35"
                   :vertical="false" title="Video" value="">
            </Meter>
          </div>
          <div class="d-flex w-100 gap-1 h-100">
            <Meter :alt="`${gpu.clocks.memory.current} MHz / ${gpu.clocks.memory.max} MHz`"
                   :percent="(gpu.clocks.memory.current/gpu.clocks.memory.max)*100" :ticks="35"
                   :vertical="false" title="Memory" value="">
            </Meter>

            <Meter :alt="`${gpu.clocks.streaming.current} MHz / ${gpu.clocks.streaming.max} MHz`"
                   :percent="(gpu.clocks.streaming.current/gpu.clocks.streaming.max)*100" :ticks="35"
                   :vertical="false" title="Streaming"
                   value="">
            </Meter>
          </div>
        </Plot>
      </div>

    </div>

  </div>

</template>

<style lang="scss" scoped>


.tick-overlay > :not(.label-o3) {
  text-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
}

.tick-overlay {
  color: rgba(255, 255, 255, 0.5);
  position: relative !important;
  height: 0;
  z-index: 10 !important;
  //mix-blend-mode: luminosity;
  font-size: 0.6rem;
  display: flex;
  justify-content: space-between;
  padding-inline: 0.2rem;
  line-height: 1rem;
}

.tick-bar {
  justify-content: start !important;
  width: 100%;
}

.tick {
  margin-left: 2.2rem;
  background-color: rgba(22, 94, 176, 0.9);
  opacity: 1;
  border-radius: 5px;
  transition: width 200ms ease-out;
  height: 1rem;
}

.tick-inline {
  position: relative;
  background-color: rgba(22, 94, 176, 1);
  border-radius: 6px;
  transition: width 200ms ease-out;
  height: 1.25rem;
}

.tick-inline :nth-child(2) {
  position: absolute !important;
  top: 0;
  width: 100%;
}

.usage-meter {
  background-color: #6f42c1;
  height: 100%;
}

.device-container {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-auto-flow: row;
  grid-template-rows: repeat(3, 1fr);
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
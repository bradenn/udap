<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Device, Remote} from "@/types";
import SimpleKeyboard from "@/components/Keyboard.vue";
import axios from "axios";
import router from "@/router";
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
  state.loading = false
  return remote.entities
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

function updateName(name: string, device: Device) {
  device.name = name
  let payload = JSON.stringify(device)
  axios.post("http://10.0.1.18:3020/devices/update", payload).then(res => {
    router.push("/terminal/settings/devices")
  }).catch(err => {
    console.log(err)
  })
}

function toggleQueryable(queryable: boolean, device: Device) {
  device.isQueryable = queryable
  let payload = JSON.stringify(device)
  axios.post("http://10.0.1.18:3020/devices/update", payload).then(res => {

  }).catch(err => {
    console.log(err)
  })
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
    <div v-if="state.mode === 'name'">

      <div class="d-flex justify-content-center">
        <div class="element p-2" style="width: 20rem;">
          <div class="text-input w-100"
               style="font-size: 0.9rem;">
            <div v-html="state.name"></div>
            <div class="cursor-blink"></div>
            <div class="flex-fill"></div>
            <div class="label-o3 label-c1" @mousedown="state.name = ''">
              Clear
            </div>
          </div>
        </div>
        <div
            class="element label-o3 label-c1 label-r d-flex justify-content-center align-items-center px-3 mx-1"
            @click="() => updateName(state.name, state.device)">Rename</div>
      </div>
      <SimpleKeyboard :input="enterChar" class="position-absolute" keySet="d"
                      keyboardClass="simple-keyboard"
      ></SimpleKeyboard>
    </div>
    <div v-else-if="state.mode === 'toggles'">
      <Plot :cols="2" :rows="1" alt="Hosts UDAP Monitoring Daemon on :5050?" title="Queryable">
        <Subplot :active="state.device.isQueryable" :fn="() => toggleQueryable(!state.device.isQueryable, state.device)"
                 name="Queryable"></Subplot>
        <Subplot :active="!state.device.isQueryable"
                 :fn="() => toggleQueryable(!state.device.isQueryable, state.device)"
                 name="Not Queryable"></Subplot>
      </Plot>
    </div>
  </div>

</template>

<style lang="scss" scoped>
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
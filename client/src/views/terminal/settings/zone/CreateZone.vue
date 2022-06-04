<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Entity, Remote, Zone} from "@/types";
import SimpleKeyboard from "@/components/Keyboard.vue";
import Plot from "@/components/plot/Plot.vue";
import axios from "axios";


interface NewZone {
  name: string
  user: string
  entities: string[]
}

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  loading: true,
  model: "",
  mode: "name",
  zone: {
    name: "",
    entities: [] as string[],
    user: ""
  } as NewZone,
  entities: [] as Entity[]
})

onMounted(() => {
  state.loading = true
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.entities = remote.entities
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
    if (state.zone.name.length > 0) state.zone.name = state.zone.name.slice(0, state.zone.name.length - 1)
    // Exit the function
    return
  }

  char = char.replace("{space}", " ")
  // Return if the security code block is full
  if (state.zone.name.length >= 64) return
  // Add the provided char to the cursor position
  state.zone.name += char
}

function nextStep() {
  state.mode = 'select'
}

function createZone() {
  let protoZone: Zone = {
    name: state.zone.name,
    user: state.zone.user,
    entities: remote.entities.filter((e) => state.zone.entities.includes(e.id))
  } as Zone
  let payload = JSON.stringify(protoZone)
  axios.post("http://localhost:3020/zones/create", payload).then(res => {
    props.done();
  }).catch(err => {
    console.log(err)
  })

}


function toggleEntity(id: string) {
  if (state.zone.entities.includes(id)) {
    state.zone.entities = state.zone.entities.filter((eid: string) => eid !== id)
  } else {
    state.zone.entities.push(id)
  }
}


</script>

<template>
  <div v-if="state.mode === 'name'">
    <div class="d-flex justify-content-center">
      <div class="element p-2" style="width: 20rem;">
        <div class="text-input w-100"
             style="font-size: 0.9rem;">
          <div v-html="state.zone.name"></div>
          <div class="cursor-blink"></div>
          <div class="flex-fill"></div>
          <div class="label-o3 label-c1" @mousedown="state.zone.name = ''">
            Clear
          </div>
        </div>
      </div>
      <div
          class="element label-o3 label-c1 label-r d-flex justify-content-center align-items-center px-3 mx-1"
          @click="() => nextStep()">Next</div>
    </div>
    <SimpleKeyboard :input="enterChar" class="position-absolute" keySet="d"
                    keyboardClass="simple-keyboard"
    ></SimpleKeyboard>

  </div>
  <div v-else-if="state.mode === 'select'" class="h-100">
    <div class="d-flex flex-row gap">
      <Plot :cols="1" :rows="5" style="width: 12rem; height: 100%">
        <div class="label-lg label-r label-o5 px-1">{{ state.zone.name }}</div>
        <div class="subplot">
          <div>Entities</div>
          <div class="flex-fill"></div>
          <div class="px-1">{{ state.zone.entities.length }}</div>
          <div class="label-o2">􀆓</div>
        </div>

        <div class="subplot">
          <div>Manager</div>
          <div class="flex-fill"></div>
          <div>Braden Nicholson</div>
        </div>
      </Plot>
      <Plot :cols="1" :rows="10" style="width: 12rem; height: 100%">
        <div class="label-lg label-r label-o5 px-1">Entities</div>
        <div v-for="entity in state.entities" class="subplot" @click="() => toggleEntity(entity.id)">

          <div>{{ entity.name }}</div>
          <div class="flex-fill"></div>
          <div>{{ state.zone.entities.includes(entity.id) ? '􀃳' : '􀂒' }}</div>

        </div>
      </Plot>
      <Plot :cols="1" :rows="1">
        <div class="subplot" @click="() => createZone()">
          <div>Create</div>
        </div>
      </Plot>

    </div>
    <div>

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
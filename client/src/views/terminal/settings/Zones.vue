<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint, Remote} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import SimpleKeyboard from "@/components/Keyboard.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  endpoints: {} as Endpoint[],
  loading: true,
  mode: "list",
  model: "",
})


onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.endpoints = remote.endpoints
  state.loading = false
  return remote.endpoints
}

function setMode(mode: string) {
  state.mode = mode
}

function enterChar(char: string) {
  // If the incoming char is a backspace, decrement the cursor and clear the value.
  if (char === "{bksp}") {
    // Only decrement the cursor if it is bigger than zero
    if (state.model.length > 0) state.model = state.model.slice(0, state.model.length - 1)
    // Exit the function
    return
  }

  char = char.replace("{space}", " ")
  // Return if the security code block is full
  if (state.model.length >= 64) return
  // Add the provided char to the cursor position
  state.model += char
}

</script>

<template>
  <div class="d-flex justify-content-start py-2 px-0">
    <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-map fa-fw`"></i></div>
    <div class="label-w500 opacity-100 label-xxl px-2">Zones</div>
    <div class="flex-fill"></div>

    <Plot :cols="2" :rows="1" small style="width: 13rem;">
      <div></div>
      <Radio :active="false" :fn="() => setMode(state.mode === 'create'?'list':'create')"
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
    Listing
  </div>
  <div v-else-if="state.mode === 'create'" class="d-flex justify-content-center">
    <div class="element p-2" style="width: 20rem">
      <div>
        <input v-model="state.model" class="subplot p-1 w-100 text-center label-lg" placeholder="Zone Identifier"
               style="font-size: 1.2rem; " type="text">

      </div>
    </div>
    <SimpleKeyboard :input="enterChar" keySet="d" keyboardClass="simple-keyboard"
    ></SimpleKeyboard>
  </div>
  <div v-else>


  </div>
</template>

<style scoped>


</style>
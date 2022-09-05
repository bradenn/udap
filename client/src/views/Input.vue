<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import PanePopup from "@/components/pane/PanePopup.vue";
import Keyboard from "@/components/Keyboard.vue";
import {reactive} from "vue";

interface InputProps {
  name: string
  value: string
  description: string
  icon?: string
  type?: string
  close: () => void
  apply: (a: string) => void
}

let props = defineProps<InputProps>()

let state = reactive({
  edit: props.value
})

function enterKey(key: string) {
  if (key === '{bksp}') {
    let def = state.edit
    state.edit = def.substring(0, def.length - 1);
  } else {
    state.edit += key
  }

}

function applyChanges() {
  props.apply(state.edit)
}

</script>

<template>
  <div class="context-rect">
    <div class="context context-full d-flex justify-content-center">
      <PanePopup :apply="() => applyChanges()" :close="() => close()" :title="`Config`"
                 style="width: 20.5rem; height: 9rem">
        <div>
          <div class="label-xxs px-1 label-w500 label-o5 lh-1 mb-1">{{ props.name }}</div>

          <div class="d-flex mb-1">
            <div class="subplot p-2 mt-0 label-c2 w-100">
              <div class="text-input w-100">
                <div class="label-xs" v-text="state.edit"></div>
                <div class="cursor-blink"></div>
                <div class="flex-fill"></div>
              </div>
            </div>
          </div>

          <div class="label-c2 label-o4 px-1">{{ props.description }}</div>
        </div>
      </PanePopup>

      <Keyboard :input="enterKey" class="keyboard-location" keySet="d"
                keyboardClass="simple-keyboard"
      ></Keyboard>

    </div>
  </div>
</template>

<style lang="scss">

</style>
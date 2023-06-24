<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {Macro} from "udap-ui/types";
import macroService from "@/services/macroService";
import {reactive} from "vue";

const props = defineProps<{
  macro?: Macro
}>()

function runMacro() {
  if (!props.macro?.id) return
  macroService.runMacro(props.macro?.id)
}

const state = reactive({
  pressed: false,
  holding: false,
  timeout: 0,
  down: 0,
})


function mouseDown(e: TouchEvent) {
  e.preventDefault();
  state.down = Date.now();
  state.holding = false;
  state.pressed = true
  state.timeout = window.setTimeout(() => {
    state.holding = true;
    // Long Click
  }, 250); // Adjust the timeout as per your requirement
}

function mouseUp(e: TouchEvent) {
  e.preventDefault();
  clearTimeout(state.timeout);
  state.pressed = false
  const elapsedTime = Date.now() - state.down;

  if (elapsedTime < 250 && !state.holding) {
    runMacro()
  }
}

</script>

<template>
  <div v-if="props.macro" :class="`${state.pressed?'pressed':''}`" class="element px-3" @touchend="mouseUp"
       @touchstart="mouseDown">

    <div class="label-w500 label-c4 label-o7 d-flex gap-2">
      <div class="sf-icon">􀋦</div>
      {{ props.macro.name }}
    </div>

  </div>
</template>

<style scoped>

.element.pressed {
  transform: scale(0.99); /* Scale down the button when pressed */
//box-shadow: inset 0 0 5px rgba(0, 0, 0, 0.1);
}

</style>
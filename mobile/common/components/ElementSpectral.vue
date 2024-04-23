<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import List from "./List.vue";
import ElementHeader from "./ElementHeader.vue";
import ElementPair from "./ElementPair.vue";
import {onMounted, reactive, watchEffect} from "vue";
import Slider from "./Slider.vue";

interface Spectral {
  on: boolean
  dim: number
  cct: number
  hue: number
  mode: string
}

const props = defineProps<{
  value?: string
  change?: (val: string) => void
}>()

const state = reactive({
  lastCCT: 5000,
  spectral: {
    on: false,
    dim: 1,
    mode: "cct",
    hue: 50,
    cct: 5000
  } as Spectral,
  loaded: false
})

function fromCommand(command: string): Spectral {
  if (!command) {
    return {
      on: false,
      dim: 1,
      mode: "cct",
      hue: 50,
      cct: 5000
    } as Spectral
  }
  let chunks = command.split(";")
  let on = chunks[0] === "on";
  let modeChunk = chunks[1].split(":")
  let mode = modeChunk[0]
  let modeValue = modeChunk[1]
  return {
    on: on,
    mode: mode,
    cct: mode === "cct" ? parseInt(modeValue) : 0,
    hue: mode === "dim" ? parseInt(modeValue) : 5000,
    dim: parseInt(chunks[2]),
  } as Spectral
}

function asCommand(spectral: Spectral): string {

  let mode = ""
  let value = 0
  if (spectral.mode == "cct") {
    mode = "cct";
    value = spectral.cct
  } else if (spectral.mode == "hue") {
    mode = "hue"
    value = spectral.hue
  }

  value = Math.round(value)

  let brightness = 1

  if (spectral.dim > 0 && spectral.dim <= 100) {
    brightness = spectral.dim
  }


  return `${spectral.on ? "on" : "off"};${mode}:${value};${brightness}`
}

watchEffect(() => {
  state.spectral = fromCommand(props?.value || "")
  if (state.spectral.cct == -1) {
    state.lastCCT = 5000
  } else {
    state.lastCCT = state.spectral.cct
  }


  return props.value
})

onMounted(() => {
  state.spectral = fromCommand(props?.value || "")
  state.loaded = true
  if (!props.change) return
  props.change(asCommand(state.spectral))
})

function togglePower() {
  state.spectral.on = !state.spectral.on
  update()
}

function toggleMode() {
  if (state.spectral.mode == "cct") {
    state.spectral.mode = "hue"
    update()
    return
  } else if (state.spectral.mode == "hue") {
    state.spectral.mode = "cct"
    update()
    return
  } else {
    state.spectral.mode = "cct"
    return
  }

}

function toggleLuminosityEdit() {
  state.spectral.on = !state.spectral.on
}

function update() {
  if (props.change && state.loaded) {
    props.change(asCommand(state.spectral))
  }
}

function toggleAutoCCT() {
  if (state.spectral.cct == -1) {
    state.spectral.cct = state.lastCCT
  } else {
    state.lastCCT = state.spectral.cct
    state.spectral.cct = -1;
  }
  state.spectral.mode = 'cct';
  update()
}

function selectCCTIcon(state: number): string {
  switch (state) {
    case 4100:
      return "􀆺 "
    case 5000:
      return "􀇕 "
    case 5000:
      return "􀆮 "
    default:
      return ""
  }
}

</script>

<template>
  <List>
    <ElementPair :cb="togglePower" :value="state.spectral.on?'ON':'OFF'" icon="􀆨" title="Power"></ElementPair>
    <ElementHeader title="Brightness"></ElementHeader>
    <ElementPair :value="`${state.spectral.dim}%`" icon="􀇯" title="Luminosity"></ElementPair>
    <Slider :change="(d) => {state.spectral.dim = d; update()}" :max="100" :min="1" :step="1"
            :value="state.spectral.dim"
            bg="dim"></Slider>
    <ElementHeader title="Color Mode"></ElementHeader>
    <ElementPair :cb="toggleMode"
                 :value="`${selectCCTIcon(state.spectral.cct)}${state.spectral.mode == 'cct'?(state.spectral.cct == -1?'Auto':`${state.spectral.cct}&deg;`):state.spectral.hue} - ${state.spectral.mode == 'cct'?'White':'Color'}`"
                 icon="􁒏"
                 title="Mode">
    </ElementPair>
    <ElementHeader title="Color Temperature"></ElementHeader>
    <ElementPair :cb="toggleAutoCCT"
                 :value="`${state.spectral.mode == 'cct' && state.spectral.cct == -1?`Circadian`:`Manual`}`"
                 icon="􁿌"
                 title="Color Temperature Mode"></ElementPair>

    <Slider
        :change="(d) => {state.spectral.mode = 'cct';state.spectral.cct = d;  update(); console.log(state.spectral)}"
        :max="6500" :min="2700" :step="100"
        :value="state.spectral.cct"
        bg="cct"
        continuous></Slider>


    <ElementHeader title="Color"></ElementHeader>
    <Slider :change="(d) => {state.spectral.hue = d; state.spectral.mode = 'hue'; update()}" :max="360" :min="0"
            :step="1" :value="state.spectral.hue"
            bg="hue"></Slider>

  </List>
</template>

<style lang="scss" scoped>

</style>
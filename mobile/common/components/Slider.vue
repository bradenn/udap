<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {onBeforeMount, reactive} from "vue";
import Element from "./Element.vue";
import core from "../core";


const props = defineProps<{
  min: number
  max: number
  step: number
  value: number
  continuous?: boolean
  disabled?: boolean
  change: (value: number) => void
  bg: string
}>()


const state = reactive({
  delta: 0,
  value: 0,
  dx: 0,
  data: {},
  width: 0,
  lastEdit: Date.now(),
  reset: 0 as number
})

onBeforeMount(() => {
  state.value = props.value
  document.body.style.setProperty("--thumb-color-preference", prefs.accent)
})

const prefs = core.preferences();


function setup() {
  state.dx = (state.width - 8) / (props.max - props.min)
}

function locateTouch(value: number): number {

  let step = Math.floor(value / (state.dx))
  step = Math.max(step, props.min)
  step = Math.min(step, props.max - 1)

  return step
}

function touchEnd(e: TouchEvent) {
  console.log(state.value)
  return

}

function dragUpdate(value: Event) {
  value = value as InputEvent
  let input = value.target as HTMLInputElement
  state.value = parseInt(input.value)
  state.lastEdit = Date.now()

  if (props.change) {
    props.change(Math.round(state.value))
  }
}

function handleUpdate(value: Event) {
  value = value as InputEvent
  let input = value.target as HTMLInputElement
  state.value = parseInt(input.value)
  state.lastEdit = Date.now()
  if (props.change) {
    props.change(Math.round(state.value))
  }
  // clearTimeout(state.reset)
  // //@ts-ignore
  // state.reset = setTimeout(() => {
  //   if (Date.now() - state.lastEdit > 250) {
  //     if (props.change) {
  //       props.change(Math.round(state.value))
  //     }
  //   }
  // }, 250)

  // let ie = value.target as InputEvent
  // console.log(ie.target)
  // state.value = value
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

</script>


<template>
  <Element class="p-0 w-100" foreground style="padding: 0 !important; height: 3rem">
    <input :class="`element-slider-${props.bg}`" :disabled="props.disabled" :max="props.max"
           :min="props.min"
           :step="props.step"
           :value="state.value"
           class="element-slider w-100"

           type="range" @change="dragUpdate" @input="dragUpdate"/>
    <!--    <div class="slider-preview"-->
    <!--         :style="`width:calc(${map_range(state.value, props.min, props.max, 0, 95)}% + 0.5rem);`"></div>-->

  </Element>
</template>

<style scoped>

input:disabled {
  filter: blur(1px) brightness(80%) saturate(80%);
}

.slider-preview {
  position: relative;
  z-index: -1;
  top: calc(-2.65rem);
  opacity: 0.3;
  left: 5px;
  width: 10%;
  height: calc(3rem - 1.5rem);
  background-color: var(--thumb-color-preference);
  border-radius: calc(0.375rem - 4px);
}

.slider-preview-inverse {
  position: relative;
  z-index: -1;
  top: calc(-3rem + 2px);
  opacity: 0.3;

  width: 10%;
  height: 2rem;
  background-color: rgba(255, 255, 255, 0.1);
}

.element-slider {

  -webkit-appearance: none;
  width: 100%;
  height: calc(3rem);

  margin: 0;
  background-color: transparent;
  border-radius: calc(0.375rem);
  padding: 0.375rem !important;
  transition: all 100ms ease-in-out;
}

.element-slider-dim {
  background: linear-gradient(90deg, rgba(0, 0, 0, 0.4) 1%, rgba(213, 230, 255, 0.5) 100%);
}

.element-slider-cct {
  background: linear-gradient(90deg, rgba(255, 160, 75, 0.5) 0%, rgba(255, 249, 230, 0.5) 50%, rgba(213, 230, 255, 0.5) 100%) !important;
}


.element-slider-hue {
  background: linear-gradient(to right,
  hsla(0deg, 100%, 40%, 0.5) 0%,
  hsla(60deg, 100%, 40%, 0.5) 16.67%,
  hsla(120deg, 100%, 40%, 0.5) 33.33%,
  hsla(180deg, 100%, 40%, 0.5) 50%,
  hsla(240deg, 100%, 40%, 0.5) 66.67%,
  hsla(320deg, 100%, 40%, 0.5) 83.33%,
  hsla(360deg, 100%, 40%, 0.5) 100%
  ) !important;

}

.element-slider:hover {
  opacity: 1;
}

.element-slider::-webkit-slider-container {
  opacity: 1;
}

.element-slider::-webkit-slider-thumb {
  background: rgba(255, 255, 255, 0.5) !important;
  box-shadow: inset 0 0 2px 1px rgba(0, 0, 0, 0.3) !important;
  backdrop-filter: contrast(50%) invert(50%) blur(25px);
  z-index: -1;
  height: calc(3rem - 0.75rem);
  width: 1.5rem;
  border-radius: calc(0.375rem - 4px);
  -webkit-appearance: none;
  appearance: none;
  transition: all 500ms ease;
}
</style>
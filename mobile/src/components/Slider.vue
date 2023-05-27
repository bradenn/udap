<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {onBeforeMount, reactive} from "vue";

const props = defineProps<{
  min: number
  max: number
  step: number
  value: number
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
})

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

function handleUpdate(value: Event) {
  value = value as InputEvent
  let input = value.target as HTMLInputElement
  state.value = parseInt(input.value)
  state.lastEdit = Date.now()
  clearTimeout(state.reset)
  //@ts-ignore
  state.reset = setTimeout(() => {
    if (Date.now() - state.lastEdit > 250) {
      if (props.change) {
        props.change(Math.round(state.value))
      }
    }
  }, 250)

  // let ie = value.target as InputEvent
  // console.log(ie.target)
  // state.value = value
}

</script>


<template>
  <input :class="`slider-${props.bg}`" :max="props.max" :min="props.min" :step="props.step" :value="state.value"
         class="slider element"
         type="range" @input="handleUpdate"
  >
</template>

<style scoped>

.slider-range {
  width: 100%;
}

.grip {
  position: relative;
  height: 4rem;
  /*border: 1px solid rgba(255, 255, 255, 0.1);*/
  border-radius: 8.24px;
  width: 48px;
}

.slider {
  padding: 4px 4px;
}
</style>
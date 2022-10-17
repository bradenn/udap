<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import {v4 as uuidv4} from "uuid";
import Ticks from "@/components/Ticks.vue";


const props = defineProps<{
  min: number,
  max: number,
  step: number,
  unit?: string,
  name: string,
  value: number,
  tags?: string[],
  change?: (a: number) => void
  live?: boolean
  confirm?: boolean
}>()

let state = reactive({
  value: 0,
  track: 0,
  dragging: false,
  confirm: false,
  since: 0,
  target: 0,
  change: 0,
  thumbWidth: 0,
  transform: 0,
  ready: false,
  timer: 0,
  lastSend: 0,
  uuid: uuidv4().toString(),
  stops: ((props.max - props.min) / props.step) + 1,
  ws: {} as WebSocket,
})

const haptics = inject("haptic") as (a: number, b: number, c: number) => void
onMounted(() => {

  let r = document.getElementById(`track-${state.uuid}`) as HTMLElement

  let w = r.getBoundingClientRect().width - 16 - 51
  // let np = ev.clientX - r.offsetLeft
  state.thumbWidth = (w / state.stops)
  state.transform = (state.stops) * (state.thumbWidth) + 16
  state.value = Math.floor(props.value / props.step)
  state.track = Math.floor(state.value * state.thumbWidth) + (state.value !== 0 ? 10 : 0)
  state.target = state.value
  state.lastSend = Date.now()
})

watchEffect(() => {

  // let ivid = setInterval(() => {
  //
  //   state.track = (state.change + state.track) / 2
  //   console.log("RESOURCES")
  //   if (Math.abs(state.track - state.change) <= 0.1) {
  //     clearInterval(ivid)
  //   }
  // }, 16)

  state.track = state.change

  return state.change
})

watchEffect(() => {

  if (state.target == state.value) return
  if (state.value % 5 != 0 && Date.now() - state.lastSend <= 40) return
  if (state.value == 0 || state.value == state.stops - 1) {

    // haptics(0, 1, 25)
    haptics(2, 1, 100)
    haptics(0, 1, 100)
    haptics(1, 1, 100)
    // haptics(1, 2, 100)
  } else {
    // haptics(1, 1, 10)
    haptics(1, 2, 25)

  }
  state.target = state.value

  state.lastSend = Date.now()
  // tick(1, 2)
  return state.value
})

function handleDrag(ev: MouseEvent) {
  state.since = Date.now()
  if (!state.dragging) return
  let r = document.getElementById(`track-${state.uuid}`) as HTMLElement

  let w = r.getBoundingClientRect().width - 16 - 51
  let np = ev.clientX - r.offsetLeft - 25
  let dx = (w / state.stops)
  state.thumbWidth = dx
  if (np / dx <= 0) {
    // state.value = 0
    // state.track = 0
  } else if (np >= w) {
    // state.value = Math.floor(np/dx)
    // state.track = Math.floor(np/dx) * (props.max-1)
    // state.track = dx * props.max
    // tick(0, 2, 4095)
  } else {
    state.value = Math.floor(np / dx)
    state.change = Math.floor(np / dx) * dx
  }
  if (props.live) {
    props.change(state.value * props.step)
  }
}

function mouseDown() {
  state.dragging = true

}

function stopDrag() {
  if (state.dragging) {
    if (!props.confirm) {
      props.change(state.value * props.step)
    } else {
      state.confirm = true
    }
    state.dragging = false
  }
}

</script>

<template>
  <div class="d-flex justify-content-center w-100">
    <div :id="`track-${state.uuid}`" class="element p-0" style="width: 100%"
         @mousedown="mouseDown" @mousemove="handleDrag"
         @mouseleave="stopDrag" @mouseup="stopDrag">
      <div class="d-flex justify-content-between lh-1 " style="padding-top:4px">
        <div class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">{{ props.name }}</div>
        <div v-if="props?.tags" class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">
          {{ props.tags[state.value] }}
        </div>
        <div v-else class="label-c1 label-o3 label-w500 p-1 px-2 pb-0">
          {{ props.min + state.value * props.step }}{{ props.unit }}
        </div>
      </div>
      <Ticks v-if="state.transform > 0" :active="state.value" :min="props.min" :series="5" :step="props.step"
             :style="`width: ${state.transform}px; top: 0px;left: ${10}px; position:relative; height:1.25rem;`"
             :tags="tags" :ticks="state.stops"></Ticks>
      <div class="shuttle-path subplot"></div>
      <div :style="`left: ${state.track}px;  position: relative; width: ${state.thumbWidth + 50}px;`"
           class="shuttle-cock subplot m-1 mt-1">
        <div class="shuttle"></div>
        <div class="shuttle-center">􀌇</div>


      </div>
      <!--    <input type="range" class="subplot slider slider-small" min="0" max="100" v-model="state.value"/>-->


    </div>
  </div>

</template>

<style lang="scss" scoped>
.shuttle-path {
  position: absolute;
  width: calc(100% - 0.5rem - 2px);
  margin: 0.25rem;
  //background-color: #0a58ca;
  height: 2rem;
}

.shuttle-text {
  position: absolute;
  left: 8px;
}

.shuttle-cock {

  height: 2rem;
  width: 100px;
  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);

}


.arrow-down {

}

.shuttle {
  width: 3px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.25);
  border-radius: 3px;
  position: relative;

  top: -28px
}

.shuttle-center {
  display: flex;
  justify-content: center;
  width: 100%;
  left: 0;
  transform: rotateZ(90deg);
  font-size: 1rem;
  //outline: 1px solid white;
  color: rgba(255, 255, 255, 0.1);
  font-weight: 300;
  font-family: "SF Pro", sans-serf, serif;
  //text-shadow: 0 0 4px rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  position: absolute;


}

.slider {
}
</style>
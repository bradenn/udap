<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, reactive} from "vue";

const props = defineProps<{
  horizontal?: boolean
}>()

let state = reactive({
  dragging: false,
  target: {} as HTMLElement,
  xDecelerate: 0,
  yDecelerate: 0,

  xStart: 0,
  yStart: 0,

  xActual: 0,
  yActual: 0,

  xVelocity: 0,
  yVelocity: 0,

  xLastPoll: 0,
  yLastPoll: 0,

  xLastPos: 0,
  yLastPos: 0,

  ySilence: false,

  lastTap: 0,
})

function dragStart(e: MouseEvent) {
  state.xVelocity = 0
  state.yVelocity = 0
  clearInterval(state.xDecelerate)
  clearInterval(state.yDecelerate)
  state.xLastPos = 0
  state.yLastPos = 0
  state.xLastPoll = 0
  state.yLastPoll = 0
  state.dragging = false
  state.yStart = e.clientY
  state.xStart = e.clientX
}

const haptics = inject("haptic") as (a: number, b: number, c: number) => void

function drag(e: MouseEvent) {

  if (!state.dragging && Math.abs(state.yStart - e.clientY) > 5 || (props.horizontal && Math.abs(state.yStart - e.clientY) > 5)) {
    state.dragging = true
  }

  if (!state.dragging) return
  let element = e.currentTarget as HTMLElement
  if (!element) return;
  state.target = element

  // haptics(1, 1, 50)


  if (props.horizontal) {
    let poll = new Date().valueOf()
    let dt = state.xLastPoll - poll
    let dx = state.xLastPos - e.clientX

    element.scrollLeft += state.xStart - e.clientX
    state.xActual = dx
    state.xStart = e.clientX
    if (state.xStart - state.xLastPos) {
      if (Date.now().valueOf() - state.lastTap >= 50)
        mouseDown()
      // element.style.position = "relative"
    }
    // if (Date.now().valueOf() - state.lastTap >= 50) {
    //   mouseDown()
    // }
    state.xVelocity = dx / dt
    state.xLastPos = e.clientX
    state.xLastPoll = poll
    //state.xStart - e.clientX

    // if()

  } else {
    let poll = new Date().valueOf()
    let dt = state.yLastPoll - poll
    let dx = state.yLastPos - e.clientY

    element.scrollTop += state.yStart - e.clientY
    state.yActual = dx
    state.yStart = e.clientY


    state.yVelocity = dx / dt
    state.yLastPos = e.clientY
    state.yLastPoll = poll
    if (element.scrollTop <= 0 && !state.ySilence) {
      haptics(1, 1, 100)
      state.ySilence = true
    } else if (element.scrollTop > 0) {
      state.ySilence = false
    }
  }

}

function mouseDown() {
  state.lastTap = Date.now().valueOf()
  haptics(1, 1, 50)
}

function dragStop(e: MouseEvent) {
  state.yStart = 0
  state.xStart = 0
  click(e)
  let element = e.currentTarget as HTMLElement
  if (!element) return;

  state.xDecelerate = setInterval(() => {
    state.xVelocity = state.xVelocity * 0.80
    element.scrollLeft -= state.xVelocity * 20
    state.xStart = state.xStart * 0.8

    if (Math.abs(state.xVelocity) <= 0.01 && state.xStart <= 0.1) {
      state.xVelocity = 0
      clearInterval(state.xDecelerate)
    }
  }, 16)

  state.yDecelerate = setInterval(() => {
    state.yVelocity = state.yVelocity * 0.95
    element.scrollTop -= state.yVelocity * 20
    state.yStart = state.yStart * 0.8
    if (element.scrollTop <= 0 && !state.ySilence) {
      haptics(1, 1, 100)
      state.ySilence = true
    }
    if (Math.abs(state.yVelocity) <= 0.01 && state.yStart <= 0.1) {
      state.yVelocity = 0
      clearInterval(state.yDecelerate)
    }
  }, 16)


  setTimeout(() => {
    state.dragging = false
  }, 250)

}

function click(e: MouseEvent) {
  if (state.dragging) {
    e.preventDefault()
    e.stopPropagation()
    e.stopImmediatePropagation()
  }
}

</script>
<template>

  <div v-on:click="click" v-on:mousedown="dragStart" v-on:mousemove="drag" v-on:mouseup="dragStop"
  >
    <!--    <div class="position-absolute top">{{ state.xStart }}</div>-->
    <slot></slot>
  </div>
</template>


<style scoped>

</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {reactive} from "vue";

const props = defineProps<{
  horizontal?: boolean
}>()

let state = reactive({
  dragging: false,
  target: {} as HTMLElement,
  xDecelerate: 0,
  yStart: 0,
  xStart: 0,
  xVelocity: 0,
  xLastPoll: 0,
  xLastPos: 0,
})

function dragStart(e: MouseEvent) {
  state.xVelocity = 0
  clearInterval(state.xDecelerate)
  state.xLastPos = 0
  state.xLastPoll = 0
  state.dragging = false
  state.yStart = e.clientY
  state.xStart = e.clientX
}

function drag(e: MouseEvent) {

  if (!state.dragging && Math.abs(state.yStart - e.clientY) > 5 || (props.horizontal && Math.abs(state.yStart - e.clientY) > 5)) {
    state.dragging = true
  }

  if (!state.dragging) return
  let element = e.currentTarget as HTMLElement
  if (!element) return;
  state.target = element
  if (props.horizontal) {
    let poll = new Date().valueOf()
    let dt = state.xLastPoll - poll
    let dx = state.xLastPos - e.clientX

    element.scrollLeft += state.xStart - e.clientX

    state.xStart = e.clientX
    state.xVelocity = dx / dt
    state.xLastPos = e.clientX
    state.xLastPoll = poll
    //state.xStart - e.clientX

  }
  element.scrollTop += state.yStart - e.clientY
  state.yStart = e.clientY
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
    if (Math.abs(state.xVelocity) <= 0.01) {
      state.xVelocity = 0
      clearInterval(state.xDecelerate)
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
  <div v-on:click="click" v-on:mousedown="dragStart" v-on:mousemove="drag" v-on:mouseup="dragStop">
    <slot></slot>
  </div>
</template>


<style scoped>

</style>
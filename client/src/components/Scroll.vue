<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {reactive} from "vue";

let state = reactive({
  dragging: false,
  yStart: 0
})

function dragStart(e: MouseEvent) {
  state.dragging = false
  state.yStart = e.clientY
}

function drag(e: MouseEvent) {
  if (!state.dragging && Math.abs(state.yStart - e.clientY) > 5) {
    state.dragging = true
  }
  if (!state.dragging) return
  let element = e.currentTarget as HTMLElement
  if (!element) return;
  element.scrollTop += state.yStart - e.clientY
  state.yStart = e.clientY
}

function dragStop(e: MouseEvent) {
  state.yStart = 0
  click(e)
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
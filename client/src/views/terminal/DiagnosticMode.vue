<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, onUnmounted, reactive} from "vue";

const state = reactive({
  ctx: {} as CanvasRenderingContext2D,
  width: 0,
  height: 0,
  runner: 0,
  touchPosition: {
    x: 0,
    y: 0,
    down: false
  },
  tick: 0
})

onMounted(() => {
  setupCanvas()
})

onUnmounted(() => {

})


function drawOverlay() {
  if (!state.touchPosition.down) return
  state.ctx.beginPath()
  state.ctx.strokeStyle = "rgba(255,255,255,1)"
  state.ctx.moveTo(0, state.touchPosition.y)
  state.ctx.lineTo(state.width, state.touchPosition.y)
  state.ctx.moveTo(state.touchPosition.x, 0)
  state.ctx.lineTo(state.touchPosition.x, state.height)
  state.ctx.closePath()
  state.ctx.stroke()
  const touchBoxRadius = 32
  state.ctx.clearRect(state.touchPosition.x - touchBoxRadius, state.touchPosition.y - touchBoxRadius, touchBoxRadius * 2, touchBoxRadius * 2)

}

function draw() {
  state.ctx.clearRect(0, 0, state.width, state.height)
  drawOverlay()

}

function redraw() {
  draw()
  requestAnimationFrame(redraw)
}

function setupCanvas() {
  const chart = document.getElementById(`root-diagnostic`) as HTMLCanvasElement
  if (!chart) return;

  const ctx = chart.getContext('2d')
  if (!ctx) return
  state.ctx = ctx
  let scale = 1
  ctx.scale(scale, scale)


  chart.width = chart.clientWidth * scale
  chart.height = chart.clientHeight * scale

  state.width = chart.width
  state.height = chart.height
  // ctx.translate(0, 0)
  // ctx.translate(0.5, 0.5);
  ctx.clearRect(0, 0, state.width, state.height)

  redraw()


}

function mouseUp(e: MouseEvent) {
  state.touchPosition.down = false
}

function mouseDown(e: MouseEvent) {
  state.touchPosition.down = true
  draw()
}

function mouseMove(e: MouseEvent) {
  state.touchPosition.x = e.x
  state.touchPosition.y = e.y
  draw()
}


</script>

<template>
  <canvas id="root-diagnostic" class="diagnostic-overlay" @mousedown="mouseDown" @mousemove="mouseMove"
          @mouseup="mouseUp">
  </canvas>
</template>


<style lang="scss">

.diagnostic-overlay {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 100;
  width: 100%;
  height: 100%;
}

</style>
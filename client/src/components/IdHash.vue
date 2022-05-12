<script lang="ts" setup>


import {onMounted, reactive} from "vue";


let state = reactive<{
  canvas: HTMLCanvasElement
  tick: number
}>({
  canvas: {} as HTMLCanvasElement,
  tick: 0
})

onMounted(() => {
  loadCanvas()
  drawCanvas()
  // setInterval(drawCanvas, 100)
})

function loadCanvas() {
  state.canvas = document.getElementById("canvas-id") as HTMLCanvasElement
  const ctx = state.canvas.getContext('2d')
  if (!ctx) return
  ctx.scale(0.5, 0.5)
}

function drawCanvas() {
  const ctx = state.canvas.getContext('2d')
  if (!ctx) return
  ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)

  let n = 24;

  let diam = 84;
  let div = (Math.PI * 2) / n;


  for (let i = 0; i < n; i++) {
    let ax = diam * Math.cos(i * div)
    let ay = diam * Math.sin(i * div)
    let bx = diam / 1.3 * Math.cos(i * div)
    let by = diam / 1.3 * Math.sin(i * div)
    ctx.strokeStyle = `rgba(255,255,255,${n * div % 1})`
    ctx.beginPath()
    ctx.moveTo(ax * 0.5 + diam / 2, ay * 0.5 + diam / 2)
    ctx.lineTo(bx * 0.5 + diam / 2, by * 0.5 + diam / 2)
    ctx.closePath()
    ctx.stroke()

  }
  state.tick += 1;


}

</script>

<template>
  <div class="canvas-container top">
    <canvas id="canvas-id" height="42" width="42"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.canvas-container {
  height: 42px;
  width: 42px;

  aspect-ratio: 1/1 !important;
}

</style>

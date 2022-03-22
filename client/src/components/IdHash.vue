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
  setInterval(drawCanvas, 200)
})

function loadCanvas() {
  state.canvas = document.getElementById("canvas-id") as HTMLCanvasElement

}

function drawCanvas() {
  const ctx = state.canvas.getContext('2d')
  if (!ctx) return
  let clk = 0;
  let w = 44
  let h = 44
  let ox = 1
  let oy = 1
  ctx.save()

  ctx.strokeStyle = 'rgba(255,255,255,0.5)';
  ctx.ellipse(w / 2, h / 2, w / 2, h / 2, 0, 0, Math.PI * 2, true)
  ctx.strokeStyle = 'rgba(255,255,255,0.3)';
  ctx.ellipse(w / 2, h / 2, w / 5, h / 5, 0, 0, Math.PI * 2, true)
  let div = 96
  let divisor = (Math.PI * 2) / div
  ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)

  let tick = state.tick
  for (let x = 0; x < div; x++) {
    if (x % 3 === 0) {

      ctx.strokeStyle = 'rgba(255,255,255,0.5)';
    } else {
      ctx.strokeStyle = 'rgba(224,224,224,0.4)';
    }

    let ax = w / 5 * Math.cos(x * divisor) + w / 2
    let ay = h / 5 * Math.sin(x * divisor) + h / 2

    let bx = w / 2 * Math.cos(x * divisor + tick) + w / 2
    let by = h / 2 * Math.sin(x * divisor + tick) + h / 2
    ctx.beginPath();
    ctx.moveTo(ax, ay);
    ctx.lineTo(bx, by);
    ctx.closePath();
    ctx.stroke();

  }
  state.tick += (state.tick + Math.PI / 8)
  ctx.restore()

}

</script>

<template>
  <div class="canvas-container">
    <canvas id="canvas-id"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.canvas-container {
  height: 42px;
  width: 42px;

  aspect-ratio: 1/1 !important;
}

</style>

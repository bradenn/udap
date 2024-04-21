<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, onUnmounted, reactive} from "vue";
import {v4 as uuidv4} from "uuid";
import {PreferencesRemote} from "udap-ui/persistent";

import Element from "udap-ui/components/Element.vue";

const props = defineProps<{
  x: number[]
  y: number[]
  z: number[]
  time: number[]
}>()

const preferences = inject("preferences") as PreferencesRemote

let historyLength = 5000;
const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  uuid: uuidv4(),
})

onMounted(() => {
  configureCanvas()
  animate()

})

onUnmounted(() => {

})


function animate() {
  requestAnimationFrame(animate)
  draw()
}

function configureCanvas() {
  const _canvas = document.getElementById(`pantilt-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d", {}) as CanvasRenderingContext2D

  let scale = 2
  state.ctx.scale(scale, scale)

  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale

  draw()
}

function averageArray(values: number[]): number {
  let sum = 0;
  values.forEach(v => sum += v)
  sum /= values.length
  return sum
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function lerp(a: number, b: number, t: number): number {
  return (1 - t) * a + t * b;
}

function draw() {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 1

  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
  // drawLegend()
  let w = ctx.canvas.width;
  let h = ctx.canvas.height;
  ctx.fillStyle = preferences.accent
  // ctx.fillRect(0, 0, w / 2, h / 2)

  let pos = mapPanTilt(averageArray(props.x), averageArray(props.y), averageArray(props.z))

  ctx.strokeStyle = "rgba(255,255,255,0.125)";
  ctx.beginPath()
  ctx.moveTo(0, h / 2)
  ctx.lineTo(w, h / 2)
  ctx.stroke()
  ctx.closePath()

  ctx.strokeStyle = "rgba(255,255,255,0.125)";
  ctx.beginPath()
  ctx.moveTo(w / 2, 0)
  ctx.lineTo(w / 2, h)
  ctx.stroke()
  ctx.closePath()

  let rings = 6
  let dx = (w * 0.5) / rings

  for (let i = 0; i < rings; i++) {
    let r = dx * i
    ctx.beginPath()
    ctx.ellipse(w / 2, h / 2, r, r, 0, 0, Math.PI * 2, false)
    ctx.stroke();
    ctx.closePath()
  }
  let angle = 5
  let offsetX = map_range(pos.p, -angle, angle, -w / 2, w / 2)
  let offsetY = map_range(pos.t, -angle, angle, -w / 2, w / 2)
  ctx.strokeStyle = preferences.accent

  for (let i = 1; i < rings; i++) {
    ctx.lineWidth = map_range(i, 0, rings, 1, 3)
    let multiplierX = offsetX / i
    let multiplierY = offsetY / i
    ctx.beginPath()
    ctx.ellipse(w / 2 + multiplierX, h / 2 + multiplierY, dx * i, dx * i, 0, 0, Math.PI * 2, false)
    ctx.stroke();
    ctx.closePath()
  }

  let centerPoint = 25

  ctx.beginPath()
  ctx.ellipse(w / 2 + offsetX, h / 2 + offsetY, centerPoint, centerPoint, 0, 0, Math.PI * 2, false)
  ctx.fill();
  ctx.closePath()


  ctx.fillStyle = "rgba(255,255,255,0.125)"
  ctx.beginPath()
  ctx.ellipse(w / 2 + offsetX * 2, h / 2 + offsetY * 2, centerPoint / 2, centerPoint / 2, 0, 0, Math.PI * 2, false)
  ctx.fill();
  ctx.closePath()

  ctx.font = "24px JetBrainsMono-Regular"
  let offsetP = `${pos.p > 0 ? '+' : ''}${(Math.round(pos.p * 100) / 100).toFixed(2)}° `;
  let offsetT = `${pos.t > 0 ? '+' : ''}${(Math.round(pos.t * 100) / 100).toFixed(2)}° `;
  let mxp = ctx.measureText(offsetP)
  let mxt = ctx.measureText(offsetT)
  ctx.fillText(offsetP, 10, 10 + mxp.actualBoundingBoxAscent);
  ctx.fillText(offsetT, 10, 40 + mxt.actualBoundingBoxAscent);


}

function mapPanTilt(x: number, y: number, z: number): { p: number, t: number } {

  let ax_g = x / 1000.0; // Convert to g (assuming ±2g scale)
  let ay_g = y / 1000.0; // Convert to g (assuming ±2g scale)
  let az_g = z / 1000.0; // Convert to g (assuming ±2g scale)

  let roll = Math.atan2(-ay_g, -az_g) * 180.0 / Math.PI;
  let pitch = Math.atan2(ax_g, Math.sqrt(ay_g * ay_g + az_g * az_g)) * 180.0 / Math.PI;

  return {p: pitch, t: roll}
}

</script>

<template>
  <Element style="aspect-ratio: 1/1">
    <div class="canvas-container w-100">
      <canvas :id="`pantilt-${state.uuid}`" class="inner-canvas w-100 h-100" style="aspect-ratio: 1/1"></canvas>
    </div>
  </Element>
</template>

<style lang="scss" scoped>

</style>
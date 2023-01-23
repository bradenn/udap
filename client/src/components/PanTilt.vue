<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, onUnmounted, reactive, watchEffect} from "vue";

const state = reactive({
  uuid: uuidv4(),
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  ready: false
})

let props = defineProps<{
  a: number,
  b: number
}>();

onMounted(() => {
  configureCanvas()
})

onUnmounted(() => {
  state.ctx.canvas.remove()
})

function configureCanvas() {
  const _canvas = document.getElementById(`pan-rig-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let scale = 1
  state.ctx.scale(scale, scale)

  state.ready = true
  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale

  draw()
}


function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}

watchEffect(() => {
  draw()
  return props
})


function draw() {
  let ctx = state.ctx;
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

  ctx.beginPath();
  ctx.strokeStyle = "rgba(255,255,255,0.25)";
  ctx.fillStyle = "rgba(255,255,255,0.25)";

  let w = ctx.canvas.width;
  let h = ctx.canvas.height;

  ctx.save()
  ctx.strokeStyle = "rgba(0,0,0,0.5)";
  ctx.translate(w / 4, h / 2);
  ctx.rotate(props.a * (Math.PI / 180))
  ctx.ellipse(0, 0, h / 3, h / 3, 0, 0, Math.PI * 2, false)
  ctx.fill()
  ctx.beginPath()
  ctx.moveTo(0, 0)
  ctx.lineTo(0, h / 3)
  ctx.closePath()
  ctx.stroke();
  ctx.restore()

  ctx.beginPath()
  ctx.moveTo(ctx.canvas.width / 2, 0)
  ctx.lineTo(ctx.canvas.width / 2, ctx.canvas.height)
  ctx.closePath()
  ctx.stroke()

  ctx.save()
  ctx.strokeStyle = "rgba(0,0,0,0.5)";
  ctx.translate(w / 2 + w / 4, h / 2);
  ctx.rotate(props.b * (Math.PI / 180))
  ctx.ellipse(0, 0, h / 3, h / 3, 0, 0, Math.PI * 2, false)
  ctx.fill()
  ctx.beginPath()
  ctx.moveTo(0, 0)
  ctx.lineTo(0, h / 3)
  ctx.closePath()
  ctx.stroke();
  ctx.restore()

}


</script>

<template>
  <div class="canvas-container">
    <canvas :id="`pan-rig-${state.uuid}`" class="inner-canvas"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.inner-canvas {
  width: 100%;
  height: 100%;

}

.canvas-container {
  display: flex;
  flex-direction: row;
  justify-content: center;
  height: 100%;
  width: 100%;
  align-items: center;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
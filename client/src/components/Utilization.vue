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
  percent: number,
}>();

onMounted(() => {
  configureCanvas()
  draw()
})

onUnmounted(() => {
  state.ctx.canvas.remove()

})

function configureCanvas() {
  const _canvas = document.getElementById(`perf-canvas-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let scale = 2
  // state.ctx.scale(scale, scale)

  state.ready = true
  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale

  draw()
}

watchEffect(() => {
  draw()
  return props.percent
})

const barWidth = 2;
const barSpacing = 1;

function lerp(a: number, b: number, t: number): number {
  return (1 - t) * a + t * b;
}


function draw() {

  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

  ctx.beginPath();
  ctx.strokeStyle = "rgba(255,255,255,0.25)";
  let w = ctx.canvas.width
  let h = ctx.canvas.height

  ctx.lineWidth = h;
  ctx.lineCap = "round"
  ctx.lineJoin = "round"
  let scaled = Math.log(props.percent * 100) / Math.log(100)
  if (props.percent == 0) {
    scaled = 0
  } else if (props.percent * 100 == 100) {
    scaled = 1
  }
  // ctx.lineJoin = "s"
  let paddingX = h / 2
  ctx.beginPath()
  ctx.moveTo(paddingX, h / 2)
  ctx.lineTo(paddingX + lerp(0, w - paddingX * 2, scaled), h / 2)
  ctx.closePath()

  ctx.stroke();
}

// function draw() {
//   if (!state.ready) return
//   // let max = -1
//   // let min = 9999
//   //
//   // for (let i = 0; i < props.values.length; i++) {
//   //   if(props.values[i] > max) max = props.values[i]
//   //   if(props.values[i] < min) min = props.values[i]
//   // }
//
//   let ctx = state.ctx
//   ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
//   ctx.fillStyle = "rgba(255,255,255,0.25)"
//   let last = 0;
//   for (let i = 0; i < props.values.length; i++) {
//     let height = Math.max(Math.log(props.values[i]), 1)
//     let run = lerp(last, height, 1);
//     ctx.fillRect(
//         i * (barWidth + barSpacing),
//         ctx.canvas.height - height,
//         barWidth,
//         height
//     );
//   }
//
//   ctx.fill()
//
// }


</script>

<template>
  <div class="canvas-container">
    <canvas :id="`perf-canvas-${state.uuid}`" class="inner-canvas"></canvas>
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
  height: 20px;
  width: 100%;
  align-items: center;
  border-radius: 20px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 4px
}
</style>
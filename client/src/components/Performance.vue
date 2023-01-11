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
  bars: number,
  values: number[]
}>();

onMounted(() => {
  configureCanvas()
})

onUnmounted(() => {
  state.ctx.canvas.remove()
})

function configureCanvas() {
  const _canvas = document.getElementById(`perf-canvas-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let scale = 1
  state.ctx.scale(scale, scale)

  state.ready = true
  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale

  draw()
}

watchEffect(() => {
  draw()
  return props.values
})

const barWidth = 2;
const barSpacing = 1;

function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}


function draw() {
  let ctx = state.ctx;
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

  ctx.beginPath();
  ctx.strokeStyle = "rgba(255,255,255,0.25)";

  let minY = Math.min(...props.values);
  let maxY = Math.max(...props.values);

  let lastX = 0;
  let lastY = ctx.canvas.height - (props.values[0] - minY) / (maxY - minY) * ctx.canvas.height;

  for (let i = 1; i < props.values.length; i++) {
    let x = i * (barWidth + barSpacing);
    let y = ctx.canvas.height - (props.values[i] - minY) / (maxY - minY) * ctx.canvas.height;

    for (let t = 0; t < 1; t += 0.01) {
      let interpolatedX = lerp(lastX, x, t);
      let interpolatedY = lerp(lastY, y, t);
      // ctx.strokeStyle = getColor(lerp(props.values[i - 1], props.values[i], t), minY, maxY);
      ctx.lineTo(interpolatedX, interpolatedY);
      ctx.stroke();
      ctx.beginPath();
      ctx.moveTo(interpolatedX, interpolatedY);
    }

    lastX = x;
    lastY = y;
  }

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
  width: 75px;
  align-items: center;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
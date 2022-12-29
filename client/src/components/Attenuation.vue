<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, reactive, watchEffect} from "vue";

const state = reactive({
  uuid: uuidv4(),
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  ready: false
})

let props = defineProps<{
  percent: number
  frequency: number,
  duty: number,
}>();

onMounted(() => {
  configureCanvas()
})

function configureCanvas() {
  const _canvas = document.getElementById(`attenuation-${state.uuid}`)
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
})

const barWidth = 2;
const barSpacing = 1;

function squareWave(f: number, a: number, d: number) {
  // Calculate the period of the wave (in seconds)
  const T = 1 / f;
  // Calculate the length of one half-period of the wave (in seconds)
  const t_half = T / 2;
  return function (t: number) {
    // Calculate the current phase of the wave (0-1)
    const phase = (t / T) % 1;
    // Return the appropriate value for the current phase
    if (phase < d) {
      return a;
    } else {
      return -a;
    }
  }
}

function drawSquareWave(amplitude: number, frequency: number, dutyCycle: number, offset: number) {

  let ctx = state.ctx
  let canvas = ctx.canvas

  // Set the stroke style
  ctx.strokeStyle = "rgba(255,255,255,0.25)"

  // Begin a new path
  ctx.beginPath();

  // Set the starting position
  ctx.moveTo(0, canvas.height / 2);

  let wave = squareWave(frequency, 1, dutyCycle)

  // Set the initial x-position
  let x = 0;

  // Set the step size
  const step = 0.01;

  // Set the y-position
  let y = 0;

  // Set the loop limit
  const limit = canvas.width;

  // Loop through the canvas width
  for (x = 0; x < limit; x += step) {
    // Calculate the y-position based on the duty cycle
    y = wave(x / limit);
    // if ((x / (1 / f)) % ((1 / f) / dutyCycle) < ((1 / f) / dutyCycle)) {
    // } else {
    //   y = -a;
    // }

    // Draw a line to the new position
    ctx.lineTo(x, canvas.height / 2 - y * (canvas.height / 4));
  }
  // Stroke the path
  ctx.stroke();
  ctx.closePath();
}

function draw() {
  if (!state.ready) return
  // let max = -1
  // let min = 9999
  //
  // for (let i = 0; i < props.values.length; i++) {
  //   if(props.values[i] > max) max = props.values[i]
  //   if(props.values[i] < min) min = props.values[i]
  // }

  let ctx = state.ctx
  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
  ctx.strokeStyle = "rgba(255,255,255,0.25)"
  drawSquareWave(0, 60, 0.126, 0);


}


</script>

<template>
  <div class="canvas-container element">
    <canvas :id="`attenuation-${state.uuid}`" class="inner-canvas"></canvas>
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
  height: 4rem;
  width: 100%;
  align-items: center;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
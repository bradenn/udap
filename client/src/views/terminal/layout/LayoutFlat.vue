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


interface Point {
  x: number
  y: number
}

function P(x: number, y: number): Point {
  return {x: x, y: y}
}

const points: Point[] = [
  P(1.4283, -1.4283),
  P(0, 0),
  P(0, 2.88),
  P(-0.2404, 3.1204),
  P(0.5374, 3.8982),
  P(3.6275, 0.8081),
  P(1.4283, -1.4283),
]
//
// new THREE.Vector2(1.4283, -1.4283),
//     new THREE.Vector2(0.0001, 0), // Zero Point (Between room and living room/patio)
//     new THREE.Vector2(0, 2.88),
//     new THREE.Vector2(-0.2404, 3.1204),
//     new THREE.Vector2(0.5374, 3.8982),
//     new THREE.Vector2(3.6275, 0.8081),

onMounted(() => {
  configureCanvas()
})

function configureCanvas() {
  const _canvas = document.getElementById(`layout-canvas-${state.uuid}`)
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
const scale = 10;

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
  ctx.strokeStyle = "rgba(255,255,255,1)"

  ctx.translate(scale, scale)


  ctx.beginPath()
  for (let i = 0; i < points.length; i++) {
    let p = points[i]
    ctx.lineTo(p.x * scale, p.y * scale)
    // ctx.moveTo(p.x * scale, p.y * scale)
    ctx.stroke();
  }
  ctx.closePath()
  ctx.translate(-scale, -scale)


}


</script>

<template>
  <div class="canvas-container">
    <canvas :id="`layout-canvas-${state.uuid}`" class="inner-canvas"></canvas>
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
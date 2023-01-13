<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, onUnmounted, reactive, watchEffect} from "vue";
import core from "@/core";
import type {Attribute} from "@/types";
import attributeService from "@/services/attributeService";

const CELLS = 28
const state = reactive({
  uuid: uuidv4(),
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  query: {request: ""} as Attribute,
  response: {value: ""} as Attribute,
  cells: [] as number[],
  copy: [] as number[],
  animation: 0,
  timeout: 0 as number,
  interval: 0 as number,
  ready: false
})

let props = defineProps<{}>();
const remote = core.remote();


onMounted(() => {
  configureCanvas()
  update()
})

onUnmounted(() => {
  state.ctx.canvas.remove()

})

function configureCanvas() {
  const _canvas = document.getElementById(`draw-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let scale = 2

  state.ctx.scale(scale, scale)


  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale
  for (let i = 0; i < CELLS * CELLS; i++) {
    state.cells[i] = 0
  }

  state.ctx.imageSmoothingEnabled = true
  state.ctx.imageSmoothingQuality = "high"
  state.ctx.stroke();
  redrawGrid()
  state.ready = true

}


const barWidth = 2;
const barSpacing = 1;

function lerp(a: number, b: number, t: number): number {
  return a + (b - a) * t;
}


function drawPixel(e: MouseEvent) {
  let div = Math.floor((state.ctx.canvas.width) / CELLS)
  let x = Math.floor((e.offsetX * 2) / div)
  let y = Math.floor((e.offsetY * 2) / div)

  let sl = 5
  for (let i = 0; i < sl; i++) {
    let dv = (Math.PI * 2) / sl
    let dx = Math.round(x + Math.cos(dv * i))
    let dy = Math.round(y + Math.sin(dv * i))
    if (dx >= 0 && dx < CELLS && dy >= 0 && dy < CELLS) {
      state.cells[dy * CELLS + dx] = Math.min((state.cells[dy * CELLS + dx] + 0.5 / 2), 1);
    }
  }

  state.cells[y * CELLS + x] = Math.min((state.cells[y * CELLS + x] + 0.75), 1)
  redrawGrid()

}

function mouseDown(e: MouseEvent) {
  if (!state.ready) return;
  cancelAnimationFrame(state.animation)
  clearTimeout(state.interval)
  if (state.interval != 0) return
  for (let i = 0; i < CELLS * CELLS; i++) {
    state.cells[i] = 0
  }

  redrawGrid()
  state.timeout = 250 / 16.66
}

function resetGrid() {
  for (let i = 0; i < CELLS * CELLS; i++) {
    state.cells[i] = 0
    state.copy[i] = 0
  }
  redrawGrid()
}

function dimGrid() {
  for (let i = 0; i < state.cells.length; i++) {
    state.cells[i] *= 0.85
    if (state.cells[i] < 0.05) state.cells[i] = 0
  }
}

let cd = 0

interface stream {
  stream: string
}

function update() {
  state.query = remote.attributes.find(a => a.id === "df94e5ce-d0ab-4bf0-b865-0fc053f8008d") as Attribute
  state.response = remote.attributes.find(a => a.id === "c6862fcd-d7ce-476c-bacc-752646e596b9") as Attribute
}

watchEffect(() => {
  update()
  return remote.attributes
})

function request() {
  let s = {} as stream
  if (!state.ready) return;
  s.stream = btoa(String.fromCharCode(...new Uint8Array(state.copy.map(c => Math.round(c * 255)))))
  if (!state.query) return
  const lq = state.query as Attribute
  lq.request = JSON.stringify(s)
  attributeService.request(lq).then(res => {
    console.log(res)
  }).catch(err => {
    console.log(err)
  })
}

function countDown() {
  dimGrid()
  redrawGrid()

  if (state.timeout <= 0) {
    state.timeout = 0
    cancelAnimationFrame(state.animation)
    request()
    resetGrid()
  } else {
    state.animation = requestAnimationFrame(countDown)
  }
  state.timeout -= 1
}

let to = 0

function mouseUp(e: MouseEvent) {
  if (!state.ready) return;
  state.copy = state.cells.map(c => c)
  clearTimeout(state.interval)
  state.interval = setTimeout(() => {
    countDown()
    clearTimeout(state.interval)
    state.interval = 0
  }, 250)

}

function redrawGrid() {
  state.ctx.clearRect(0, 0, state.ctx.canvas.width || 0, state.ctx.canvas.height || 0);
  const pad = 4


  let div = Math.floor((state.ctx.canvas.width) / CELLS)
  for (let i = 0; i < CELLS; i++) {
    state.ctx.beginPath()
    state.ctx.moveTo(0, i * div)
    state.ctx.lineTo(state.ctx.canvas.width, i * div)
    state.ctx.closePath()
    state.ctx.stroke()
    state.ctx.beginPath()
    state.ctx.moveTo(i * div, 0)
    state.ctx.lineTo(i * div, state.ctx.canvas.height)

    state.ctx.closePath()
    state.ctx.stroke()
    for (let j = 0; j < CELLS; j++) {
      let val = state.cells[i * CELLS + j]

      state.ctx.fillStyle = `rgba(255,128,1,${val})`;
      state.ctx.fillRect(j * div + pad / 2, i * div + pad / 2, div - pad, div - pad)

    }
  }


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
    <canvas :id="`draw-${state.uuid}`" class="inner-canvas"
            @mousedown="mouseDown" @mousemove="drawPixel"
            @mouseout="mouseUp" @mouseup="mouseUp"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.guess {
  font-size: 8rem;
  animation: show-guess 1s normal;
}

@keyframes show-guess {
  0% {
    scale: 0.8;
  }
  10% {
    scale: 1;
  }
  50% {
    scale: 0.8;
  }
  100% {
    scale: 0;
  }
}

.guess-box {

  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;

}

.inner-canvas {
  width: calc(12px * 28);
  height: calc(12px * 28);

}

$bg-color: rgba(34, 38, 45, 0);
$dot-color: rgba(255, 128, 1, 0.25);

// Dimensions
$dot-size: 28px;
$dot-space: 16px;

.whiteboard-canvas {

  background-color: $bg-color;
  opacity: 1;
  border-radius: 0.24rem;
  border: 1px solid $dot-color;
  background-image: radial-gradient($dot-color 0.980000000000000px, $bg-color 0.980000000000000px), radial-gradient($dot-color 0.98px, $bg-color 0.980000000000000px);
  background-size: $dot-size $dot-size;
  background-position: 0 0, calc($dot-size / 2) calc($dot-size / 2);
}

.canvas-container {
  display: flex;
  flex-direction: row;
  justify-content: center;
  height: 100%;
  width: 75px;
  align-items: center;
  border-radius: 6px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
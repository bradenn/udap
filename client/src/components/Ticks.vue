<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import {v4 as uuidv4} from "uuid";
import {Preference} from "@/preferences";

const props = defineProps<{
  active?: number,
  min: number,
  ticks: number,
  series: number,
  step?: number,
  tags?: string[]
}>()

const state = reactive({
  ctx: {} as CanvasRenderingContext2D,
  width: 0,
  height: 0,
  uuid: uuidv4().toString(),
  loaded: false,
})

onMounted(() => {
  setupCanvas()
  drawTicks()
})

watchEffect(() => {
  drawTicks()
  redraw()
})
const preferences = inject("preferences") as Preference

function drawTicks() {
  let ctx = state.ctx
  let series = props.series
  ctx.lineWidth = 2
  let offset = 22;
  let dx = ctx.canvas.width / (props.ticks);
  let step = props.step ? props.step : 1
  for (let i = 0; i < props.ticks; i++) {
    ctx.fillStyle = "rgba(255,255,255, 0.25)"
    ctx.strokeStyle = "rgba(255,255,255, 0.25)"
    let height = ctx.canvas.height;
    let text = `${props.min + i * step}`
    if (props.tags) {
      text = props.tags[i]
    }
    if (props.active == i) {
      ctx.fillStyle = `rgba(${preferences.ui.accent}, 0.8)`
      ctx.strokeStyle = `rgba(${preferences.ui.accent}, 0.4)`
      ctx.font = "400 40px SF Pro"
      let ms = ctx.measureText(text)
      let infill = 0;
      if (i == 0) {
        if (ms.width / 2 > offset) {
          infill = ms.width / 2 - offset
        }
      }


      ctx.fillText(`${text}`, offset + i * (dx) - ms.width / 2 + infill, height / 2 + (40 / 3))
    } else {
      ctx.font = "400 24px SF Pro"
      let ms = ctx.measureText(`${text}`)
      let infill = 0;
      if (i == 0) {
        if (ms.width / 2 > offset) {
          infill = ms.width / 2 - offset
        }
      }
      ctx.fillText(`${text}`, offset + i * (dx) - ms.width / 2 + infill, height / 2 + (24 / 3))
    }


    ctx.beginPath()
    ctx.moveTo(offset + i * (dx), height)
    ctx.lineTo(offset + i * (dx), height - 10)
    ctx.closePath()
    ctx.stroke()
  }
}

function redraw() {
  state.ctx.clearRect(0, 0, state.width, state.height)
  drawTicks()
}

function setupCanvas() {
  const chart = document.getElementById(`tick-canvas-${state.uuid}`) as HTMLCanvasElement
  if (!chart) return;

  const ctx = chart.getContext('2d')
  if (!ctx) return
  state.ctx = ctx
  let scale = 1.75
  ctx.scale(scale, scale)


  chart.width = chart.clientWidth * scale
  chart.height = chart.clientHeight * scale

  state.width = ctx.canvas.width
  state.height = ctx.canvas.height
  ctx.translate(0, 0)
  ctx.clearRect(0, 0, state.width, state.height)

  redraw()
}

</script>

<template>
  <div>

    <canvas :id="`tick-canvas-${state.uuid}`" style="width: 100%; height: 100%;"></canvas>
  </div>
</template>

<style lang="scss" scoped>
.shuttle-cock {

  height: 2rem;
  width: 100px;
}

@keyframes shuttle-cock-grow {
  0% {
    transform: scale(2) !important;
  }
  50% {
    transform: scale(2);
  }
  100% {
    transform: scale(2);
  }
}

.shuttle-cock:active {
  animation: shuttle-cock-grow 200ms forwards linear;
}

.slider {
}
</style>
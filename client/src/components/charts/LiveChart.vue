<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {onMounted, onUpdated, reactive, watchEffect} from "vue";

interface ChartValues {

}

interface ChartProps {
  name: string
  unit?: string
  values: number[]
  valuesPerSection: number
  color: string
  marker: number
  scale: number
}

onUpdated(() => {
  redraw()
})

watchEffect(() => {
  redraw()
  return props.values
})

let props = defineProps<ChartProps>()

let state = reactive({
  ctx: {} as CanvasRenderingContext2D,
  width: 0 as number,
  height: 0 as number,
  min: 0 as number,
  max: -10E9 as number,
})

onMounted(() => {

  loadDom()
})

function calcMinMax() {
  for (let i = 0; i < props.valuesPerSection; i++) {
    if (props.values[i] > state.max) state.max = props.values[i]
    else if (props.values[i] < state.min) state.min = props.values[i]
  }
}

function lerp(v0: number, v1: number, t: number): number {
  return (1 - t) * v0 + t * v1;
}

function loadDom() {
  const chart = document.getElementById(`chart-${props.name}`) as HTMLCanvasElement
  if (!chart) return;

  const ctx = chart.getContext('2d')
  if (!ctx) return
  state.ctx = ctx
  let scale = 1.75
  ctx.scale(scale, scale)


  chart.width = chart.clientWidth * scale
  chart.height = chart.clientHeight * scale

  state.width = chart.width
  state.height = chart.height
  ctx.translate(0, 0)
  ctx.clearRect(0, 0, state.width, state.height)

  redraw()

}


function redraw() {
  state.ctx.clearRect(0, 0, state.width, state.height)
  calcMinMax()
  drawHorizontalLines()
  drawSections()
  drawLerpLine()
  drawMarkers()
}

function drawLerpLine() {
  let ctx = state.ctx
  ctx.setLineDash([])
  ctx.strokeStyle = props.color
  ctx.beginPath()
  ctx.lineWidth = 2 * props.scale
  let middle = state.height / 2;
  // ctx.moveTo(0, state.height / 2)
  let chunk = Math.round(state.width / props.valuesPerSection);
  for (let i = 0; i < state.width; i++) {

    let idx = Math.floor(i / chunk)
    let dx = lerp(props.values[idx], props.values[idx + 1], (i % chunk) / chunk)
    ctx.lineTo(i, (state.height - 6) - ((dx - state.min) / (state.max - state.min)) * (state.height / 1.5))

  }
  ctx.moveTo(state.width, state.height / 2)
  ctx.closePath()
  ctx.stroke()
}

function drawHorizontalLines() {
  let ctx = state.ctx
  ctx.lineWidth = 2
  ctx.strokeStyle = `rgba(255,255,255,${0.25})`
  ctx.setLineDash([])

  let divisions = 1;
  let div = state.height / divisions;

  ctx.beginPath()
  for (let i = 0; i <= divisions; i++) {
    ctx.moveTo(0, div * i)
    ctx.lineTo(state.width, div * i)
  }
  ctx.closePath()
  ctx.stroke()


}

function drawMarkers() {
  let ctx = state.ctx
  ctx.lineWidth = 2
  ctx.setLineDash([])
  ctx.lineDashOffset = 2
  ctx.fillStyle = props.color
  ctx.strokeStyle = `rgba(255,255,255,${0.1})`
  ctx.font = "600 30px SF Pro Rounded"

  let chunk = (state.width / props.valuesPerSection);
  let dx = ((props.values[props.marker] - state.min) / (state.max - state.min))

  ctx.beginPath()
  ctx.moveTo(props.marker * chunk, 0)
  ctx.lineTo(props.marker * chunk, state.height)
  ctx.closePath()
  let metrics = ctx.measureText(`${props.values[props.marker]}${props.unit ? props.unit : ''}`)
  ctx.stroke()

  ctx.fillText(`${props.values[props.marker]}${props.unit ? props.unit : ''}`, props.marker * chunk - 4 - metrics.width, (state.height - 6) - dx * (state.height / 1.5) - 12)
  ctx.ellipse(props.marker * chunk, (state.height - 6) - dx * (state.height / 1.5), 6, 6, 0, 0, 2 * Math.PI)
  ctx.fill()
}

function drawSections() {
  let ctx = state.ctx
  ctx.lineWidth = 2
  ctx.setLineDash([4, 8])

  ctx.lineDashOffset = 2
  ctx.strokeStyle = `rgba(255,255,255,${0.2})`
  ctx.fillStyle = `rgba(255,255,255,${0.5})`
  let div = state.width
  ctx.font = "30px SF Pro Display"
  for (let i = 0; i <= 1; i++) {
    ctx.fillText(props.name, i * div + 10, 40)
    ctx.beginPath()
    ctx.moveTo(i * div, 0)
    ctx.lineTo(i * div, state.height)
    ctx.closePath()
    ctx.stroke()
  }

}


</script>

<template>
  <canvas :id="`chart-${props.name}`" class="horizontal-canvas"></canvas>
</template>


<style lang="scss">

.horizontal-canvas {
  padding: 0;
  margin: 0;
  height: 4rem;
  width: 100%;
}

</style>
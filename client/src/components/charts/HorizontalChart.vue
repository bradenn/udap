<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {onMounted, reactive} from "vue";

interface ChartValues {

}

interface ChartProps {
  name: string
  values: number[]
  valuesPerSection: number
  sectionsNames: string[]
  color: string
  scale: number
}

let props = defineProps<ChartProps>()

let state = reactive({
  ctx: {} as CanvasRenderingContext2D,
  width: 0 as number,
  height: 0 as number
})

onMounted(() => {
  loadDom()
})

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

  drawHorizontalLines()
  drawSections()
  drawLerpLine()
}

function drawLerpLine() {
  let ctx = state.ctx
  ctx.setLineDash([])
  ctx.strokeStyle = props.color
  ctx.beginPath()
  ctx.lineWidth = 6
  let middle = state.height / 2;
  // ctx.moveTo(0, state.height / 2)
  let chunk = (state.width / props.values.length);
  for (let i = 0; i < state.width; i++) {

    let idx = Math.floor(i / chunk)
    let dx = lerp(props.values[idx], props.values[idx + 1], (i % chunk) / chunk)
    ctx.lineTo(i, (state.height - 6) - dx * (state.height / 1.5))

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

function drawSections() {
  let ctx = state.ctx
  ctx.lineWidth = 2
  ctx.setLineDash([4, 8])
  ctx.lineDashOffset = 2
  ctx.strokeStyle = `rgba(255,255,255,${0.2})`
  ctx.fillStyle = `rgba(255,255,255,${0.5})`
  let div = state.width / props.sectionsNames.length;
  ctx.font = "30px SF Pro Display"
  for (let i = 0; i <= props.sectionsNames.length; i++) {
    ctx.fillText(props.sectionsNames[i], i * div + 10, 40)
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
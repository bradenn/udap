<!-- Copyright (c) 2023 Braden Nicholson -->
<script lang="ts" setup>

import {v4 as uuidv4} from "uuid";
import {onMounted, onUpdated, reactive} from "vue";

const props = defineProps<{
  min: number
  max: number
  weatherMin: number
  weatherMax: number
  inside: number
  outside: number
  outsideNext: number
  set: {
    heat: number,
    cool: number
  }
  ecoMin: number
  ecoMax: number
}>()

const state = reactive({
  uuid: uuidv4(),
  ctx: {} as CanvasRenderingContext2D
})

onMounted(() => {
  setupCanvas()
})

onUpdated(() => {
  render()
})


function lerp(v0: number, v1: number, t: number): number {
  return (1 - t) * v0 + t * v1;
}

function setupCanvas() {
  const canvas = document.getElementById(state.uuid) as HTMLCanvasElement
  if (!canvas) return
  const ctx = canvas.getContext("2d") as CanvasRenderingContext2D
  if (!ctx) return
  ctx.canvas.width = canvas.clientWidth * 4
  ctx.canvas.height = canvas.clientHeight * 4
  state.ctx = ctx
  render()
}

function render() {
  const ctx = state.ctx
  const w = ctx.canvas.width
  const h = ctx.canvas.height
  ctx.clearRect(0, 0, w, h)
  ctx.fillStyle = 'rgba(255,128,1,1)'
  let delta = props.max - props.min
  let dx = (w) / delta
  let dy = h / 4

  let inside = Math.round(props.inside - props.min)
  let outside = Math.round(props.outside - props.min)
  let next = Math.round(props.outsideNext - props.min)


  let setHeat = Math.round(props.set.heat - props.min)
  let setCool = Math.round(props.set.cool - props.min)

  let arrowPoint = Math.max(Math.abs(next - outside), dx / 2)
  let arrowDir = next - outside

  ctx.strokeStyle = 'rgba(44,44,46,1)'
  ctx.lineWidth = 4
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'

  for (let i = 0; i < delta; i += 0.5) {
    ctx.beginPath()
    ctx.moveTo(dx * i, h / 2 - dy + 5)
    ctx.lineTo(dx * i, h / 2 + dy - 5)
    ctx.closePath()
    ctx.stroke()
  }


// Set the fill style and draw a rectangle


  ctx.save()
  ctx.lineWidth = 32
  let emax = Math.round(props.weatherMax - props.min)
  let emin = Math.round(props.weatherMin - props.min)
  let gradient = ctx.createLinearGradient(dx * emin, 0, dx * emax, 0);

  gradient.addColorStop(0, "rgba(90,200,245,1)");
  gradient.addColorStop(1, "rgba(255,69,55,1)");
  ctx.strokeStyle = gradient;
  ctx.lineJoin = 'round'

  ctx.beginPath()
  ctx.moveTo(dx * emin, h / 2)
  ctx.lineTo(dx * emax, h / 2)
  ctx.closePath()
  ctx.stroke()
  ctx.restore()

  // ctx.save()
  // ctx.lineWidth = h / 16
  // let ecoMax = Math.round(props.ecoMax - props.min)
  // let ecoMin = Math.round(props.ecoMin - props.min)
  //
  // gradient = ctx.createLinearGradient(0, 0, w, 0);
  //
  // gradient.addColorStop(0, "rgba(48,209,88,0.8)");
  // gradient.addColorStop((dx * ecoMin) / w, "rgba(255,214,10,0.8)");
  // gradient.addColorStop(((dx * (ecoMax - 1)) / w), "rgba(255,112,10, 0.8)");
  // gradient.addColorStop(1, "rgba(99,219,91,0.8)");
  // ctx.strokeStyle = gradient;
  // ctx.lineJoin = 'round'
  //
  // ctx.beginPath()
  // ctx.moveTo(0, h / 2 + dy - 5 - h / 16 / 2)
  // ctx.lineTo(w, h / 2 + dy - 5 - h / 16 / 2)
  // ctx.closePath()
  // ctx.stroke()
  // ctx.restore()

  ctx.save()
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.lineWidth = 3
  ctx.fillStyle = arrowDir > 0 ? 'rgba(255,255,255,0.9)' : 'rgba(255,255,255,0.9)'
  ctx.strokeStyle = arrowDir > 0 ? 'rgba(255,255,255,1)' : 'rgba(255,255,255,1)'
  ctx.beginPath()
  ctx.moveTo(dx * outside, h / 2)
  ctx.lineTo(dx * outside, h / 2 + dy - 5)
  ctx.closePath()
  ctx.fill()
  ctx.stroke()
  ctx.beginPath()
  ctx.moveTo(dx * outside, h / 2 - h / 12)
  ctx.lineTo(dx * outside - (arrowDir < 0 ? arrowPoint : -arrowPoint), h / 2)
  ctx.lineTo(dx * outside, h / 2 + h / 12)
  ctx.closePath()
  ctx.fill()
  ctx.stroke()
  ctx.restore()


  // ctx.save()
  // ctx.strokeStyle = 'rgba(99,99,102,1)'
  // ctx.lineWidth = dy / 4
  // ctx.lineCap = 'round'
  // ctx.lineJoin = 'round'
  // ctx.beginPath()
  // ctx.moveTo(dx * inside, h / 2)
  // ctx.lineTo(dx * set, h / 2)
  // ctx.closePath()
  // ctx.stroke()
  // ctx.restore()

  ctx.lineWidth = 3
  ctx.strokeStyle = 'rgba(255,255,255,1)'
  ctx.beginPath()
  ctx.moveTo(dx * inside, h / 2 - dy + 5)
  ctx.lineTo(dx * inside, h / 2)
  ctx.closePath()
  ctx.stroke()
  //

  ctx.save()
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.lineWidth = 3
  ctx.fillStyle = "rgba(90,200,245,0.9)"
  ctx.strokeStyle = "rgba(90,200,245,0.9)"
  // ctx.beginPath()
  // ctx.moveTo(dx * setCool, h / 2)
  // ctx.lineTo(dx * setCool, h / 2 - dy + 5)
  // ctx.closePath()
  // ctx.fill()
  // ctx.stroke()
  ctx.beginPath()
  ctx.moveTo(dx * setCool, h / 2 - h / 12)
  ctx.lineTo(dx * setCool - dx / 2, h / 2)
  ctx.lineTo(dx * setCool, h / 2 + h / 12)
  ctx.closePath()
  ctx.fill()
  ctx.stroke()
  ctx.restore()

  ctx.save()
  ctx.lineCap = 'round'
  ctx.lineJoin = 'round'
  ctx.lineWidth = 3
  ctx.fillStyle = "rgba(255,69,55,1)"
  ctx.strokeStyle = "rgba(255,69,55,1)"
  // ctx.beginPath()
  // ctx.moveTo(dx * setHeat, h / 2)
  // ctx.lineTo(dx * setHeat, h / 2 - dy + 5)
  // ctx.closePath()
  // ctx.fill()
  // ctx.stroke()
  ctx.beginPath()
  ctx.moveTo(dx * setHeat, h / 2 - h / 12)
  ctx.lineTo(dx * setHeat + dx / 2, h / 2)
  ctx.lineTo(dx * setHeat, h / 2 + h / 12)
  ctx.closePath()
  ctx.fill()
  ctx.stroke()
  ctx.restore()

  // ctx.translate(0, -h / 4)
  ctx.fillStyle = 'rgba(255,255,255,0.7)'
  ctx.font = "600 32px SF Pro Display"
  let currentTemp = `${props.inside}° F`

  let mx = ctx.measureText(currentTemp)

  mx = ctx.measureText("􀎟 " + currentTemp)
  ctx.fillText('􀎟 ' + currentTemp, dx * inside - mx.width / 2, h / 4 - 5)
  currentTemp = `􀆮 ${props.outside}° F`
  mx = ctx.measureText(currentTemp)
  ctx.fillText(currentTemp, (dx * outside) - mx.width / 2, h / 2 + h / 4 + mx.actualBoundingBoxAscent + 5)

  ctx.font = "600 32px SF Pro Display"
  ctx.fillStyle = 'rgba(255,255,255,0.8)'
  mx = ctx.measureText("􀀁")
  ctx.fillText('􀆼', dx * emin - mx.width / 2, h / 2 + mx.fontBoundingBoxAscent / 2 - 4)
  ctx.fillStyle = 'rgba(255,255,255,0.8)'
  ctx.fillText('􀷏', dx * emax - mx.width / 2, h / 2 + mx.fontBoundingBoxAscent / 2 - 4)
}


</script>

<template>
  <div class="d-flex justify-content-between label-c6 label-w600 label-o3">
    <div class="d-flex justify-content-between label-c6 flex-grow-1">
      <canvas :id="`${state.uuid}`" class="hash-canvas"></canvas>
    </div>
  </div>

</template>

<style scoped>
.hash-canvas {
  background-color: rgba(255, 255, 255, 0);

  border-radius: 0.25rem;
  width: 100%;
  height: 2.5rem;
}
</style>
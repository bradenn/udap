<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, onUnmounted, reactive, watchEffect} from "vue";
import {v4 as uuidv4} from "uuid";
import {PreferencesRemote} from "udap-ui/persistent";

const props = defineProps<{
  time: number[],
  data: number[],
}>()


const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  animation: 0,
  uuid: uuidv4(),
})
watchEffect(() => {
  draw();
  return props.data
})
onMounted(() => {
  configureCanvas()
  draw()
})


onUnmounted(() => {
  // cancelAnimationFrame(state.animation)
})

const preferences = inject("preferences") as PreferencesRemote

function animate() {
  state.animation = requestAnimationFrame(animate)
  draw()
}

function configureCanvas() {
  const _canvas = document.getElementById(`linechart-${state.uuid}`)
  state.canvas = _canvas as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d", {}) as CanvasRenderingContext2D

  let scale = 2
  state.ctx.scale(scale, scale)

  state.canvas.width = state.canvas.clientWidth * scale
  state.canvas.height = state.canvas.clientHeight * scale

  draw()
}

function averageArray(values: number[]): number {
  let sum = 0;
  values.forEach(v => sum += v)
  sum /= values.length
  return sum
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function lerp(a: number, b: number, t: number): number {
  return (1 - t) * a + t * b;
}

function draw() {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 1

  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
  // drawLegend()
  let w = ctx.canvas.width;
  let h = ctx.canvas.height;
  let vh = h - 40
  ctx.fillStyle = preferences.accent
  // ctx.fillRect(0, 0, w / 2, h / 2)

  ctx.strokeStyle = "rgba(255,255,255,0.125)";
  ctx.beginPath()
  ctx.moveTo(0, h / 2)
  ctx.lineTo(w, h / 2)
  ctx.stroke()
  ctx.closePath()

  if (!props.data || !props.time) return;

  let min = Math.min(...props.data)
  let max = Math.max(...props.data)

  let minT = Math.min(...props.time)
  let maxT = Math.max(...props.time)

  let dT = maxT - minT;
  ctx.fillStyle = "rgba(255,255,255,0.4)"

  let dtM = Math.round(dT / (1000 * 60 * 30))
  ctx.font = "24px JetBrains Mono"

  ctx.lineWidth = 2;
  let lx = 0;
  let ly = h / 2;

  let pairs = []

  for (let i = 0; i < props.time.length; i++) {
    pairs.push({
      time: props.time[i],
      data: props.data[i],
    })
  }
  ctx.save()
  ctx.restore()

  let sortedTime = pairs.sort((a, b) => a.time - b.time)
  for (let i = 0; i < pairs.length; i++) {
    let t = sortedTime[i].time || 0
    let d = sortedTime[i].data || 0
    let x = map_range(t, minT, maxT, 0, w)
    let y = map_range(d, min, max, vh - 20, 20)
    if (t % (1000 * 60 * 60) == 0) {
      ctx.strokeStyle = "rgba(255,255,255,0.2)"
      ctx.lineWidth = 2;
      ctx.beginPath()
      ctx.moveTo(x, 0)
      ctx.lineTo(x, vh)
      ctx.stroke()
      ctx.closePath()
      let d = new Date(t)
      let time = `${d.getHours() % 12} ${d.getHours() >= 12 ? 'PM' : 'AM'}`
      let mxt = ctx.measureText(time)
      if (i == 0) {
        x += mxt.width / 2
      }
      ctx.font = "24px JetBrains Mono"
      ctx.fillText(time, x - mxt.width / 2, h - mxt.actualBoundingBoxAscent / 2);
    } else if (t % (1000 * 60 * 30) == 0) {
      ctx.strokeStyle = "rgba(255,255,255,0.1)"
      ctx.lineWidth = 1;
      ctx.beginPath()
      ctx.moveTo(x, 0)
      ctx.lineTo(x, vh)
      ctx.stroke()
      ctx.closePath()
    }

    ctx.lineWidth = 2;
    ctx.strokeStyle = "rgba(0,128,255,1)"

    ctx.beginPath()
    if (i == 0) {
      ctx.moveTo(x, y)
    } else if (x - lx > w / 10) {
      ctx.moveTo(x, y)
    } else {
      ctx.moveTo(lx, ly)

    }

    ctx.lineTo(x, y)
    lx = x
    ly = y
    ctx.stroke()
    ctx.closePath()
  }


  ctx.strokeStyle = preferences.accent


  // let offsetP = `${pos.p > 0 ? '+' : ''}${(Math.round(pos.p * 100) / 100).toFixed(2)}° `;
  // let offsetT = `${pos.t > 0 ? '+' : ''}${(Math.round(pos.t * 100) / 100).toFixed(2)}° `;
  // let mxp = ctx.measureText(offsetP)
  // let mxt = ctx.measureText(offsetT)
  // ctx.fillText(offsetP, 10, 10 + mxp.actualBoundingBoxAscent);
  // ctx.fillText(offsetT, 10, 40 + mxt.actualBoundingBoxAscent);


}

</script>

<template>
  <canvas :id="`linechart-${state.uuid}`"
          style="font-family: 'JetBrains Mono',serif; aspect-ratio: 2/1; width: 100%"></canvas>
</template>

<style lang="scss" scoped>
#loadFont {

}
</style>
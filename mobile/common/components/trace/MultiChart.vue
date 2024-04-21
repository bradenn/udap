<!-- Copyright (c) 2024 Braden Nicholson -->

<script setup lang="ts">

import {inject, onMounted, onUnmounted, reactive, watchEffect} from "vue";
import {v4 as uuidv4} from "uuid";
import {PreferencesRemote} from "udap-ui/persistent";
import {Trace} from "../services/traceService";

interface DataObject {
  name: string,
  metadata: any,
  x: number[],
  y: number[]
}

const props = defineProps<{
  traces: Trace[],
  transforms?: ((src: number) => number)[]
}>()


const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  animation: 0,
  uuid: uuidv4(),
})

watchEffect(() => {
  draw();
  return props.traces
})
onMounted(() => {
  configureCanvas()
  // animate()

})


onUnmounted(() => {
  cancelAnimationFrame(state.animation)
})

const preferences = inject("preferences") as PreferencesRemote

function animate() {
  // state.animation = requestAnimationFrame(animate)
  // draw()
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

interface TraceMetadata {
  valueMin: number,
  valueMax: number,
  timeMin: number,
  timeMax: number,
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

function renderLegend(md: TraceMetadata, trace: Trace, color: string, adjacent: number, index: number) {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 2
  let w0 = 100
  let w = ctx.canvas.width - 20;
  let h0 = 20
  let h = ctx.canvas.height - 70;
  ctx.fillStyle = `rgba(${color}, 1)`
  ctx.strokeStyle = `rgba(${color}, 0.5)`

  let di = (w - w0) / (3)
  let extent = index % 3
  let row = index > extent ? 1 : 0
  let x = di * extent
  let y = h + row * 26
  let text = `${trace.labels["serial"]}`
  let mx = ctx.measureText(text)
  ctx.fillText(text, w0 + x, y + mx.actualBoundingBoxAscent + 8)


}

function renderTrace(md: TraceMetadata, trace: Trace, color: string) {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 1
  let w0 = 100
  let w = ctx.canvas.width - 20;
  let h0 = 20
  let h = ctx.canvas.height - 100;
  ctx.fillStyle = preferences.accent
  ctx.fillStyle = `rgba(${color}, 1)`
  ctx.strokeStyle = `rgba(${color}, 0.5)`

  let lx = 0;
  let ly = 0;
  for (let i = 0; i < trace.time.length; i++) {
    let t = trace.time[i] || 0
    let d = trace.data[i] || 0
    let x = map_range(t, md.timeMin, md.timeMax, w0, w)
    let y = map_range(d, md.valueMin, md.valueMax, h - 20, h0)
    let s = (10 / trace.time.length) * i;
    if (i == 0) {
      lx = x;
      ly = y;
    }
    ctx.beginPath()
    ctx.moveTo(lx, ly)
    ctx.lineTo(x, y)
    ctx.stroke();
    ctx.closePath()
    lx = x;
    ly = y
    ctx.beginPath()
    ctx.ellipse(x, y, 2, 2, 0, 0, Math.PI * 2)
    ctx.fill();
    ctx.stroke();


  }

}

function alignTrace(trace: Trace): Trace {
  if (!trace.data || !trace.time) return;

  const indices = Array.from(trace.time.keys())
  indices.sort((a, b) => trace.time[a] - trace.time[b])
  return {
    time: indices.map(i => trace.time[i]),
    data: indices.map(i => trace.data[i]),
    name: trace.name,
    labels: trace.labels,
  }
}

function alignTraces(traces: Trace[]): Trace[] {
  return traces.map(t => alignTrace(t)).sort((a, b) => a.name.localeCompare(b.name))
}

function transformTrace(trace: Trace): Trace {
  if (!props.transforms) {
    return trace
  }

  for (let i = 0; i < props.transforms.length; i++) {
    let fn: ((src: number) => number) = props.transforms[i]

    // trace.time = trace.time.map(t => t % (1000 * 60 * 60 * 24))
    trace.data = trace.data.map(d => fn(d))
  }

  return trace
}

function transformTraces(traces: Trace[]): Trace[] {
  if (!props.transforms) {
    return traces
  }

  return traces.map(t => transformTrace(t))
}

function renderScales(md: TraceMetadata) {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 1

  let duration = md.timeMax - md.timeMin
  let w0 = 100;
  let h0 = 20;
  let w = ctx.canvas.width - 20;
  let h = ctx.canvas.height - 100;
  ctx.strokeStyle = "rgba(255,255,255,0.5)";
  // ctx.strokeRect(w0, h0, w - w0, h - h0)
  let min = 60 * 1000
  let hr = 60 * min
  let day = 24 * hr
  let dt = hr
  if (duration > day) {
    dt = duration / 15
  } else if (duration >= 12 * hr) {
    dt = hr * 2
  } else if (duration >= 6 * hr) {
    dt = hr
  } else if (duration >= 2 * hr) {
    dt = hr / 2
  } else if (duration >= hr) {
    dt = 10 * min
  }


  let dx = Math.round(duration / dt)
  console.log(md)
  for (let i = 0; i <= dx; i++) {

    let x = map_range(md.timeMin + i * dt, md.timeMin, md.timeMax, w0, w)
    let y = map_range(md.timeMin + i * dt, md.timeMin, md.timeMax, h0, h)
    let v = map_range(md.timeMin + i * dt, md.timeMin, md.timeMax, md.valueMax, md.valueMin)

    ctx.beginPath()
    ctx.moveTo(x, h0)
    ctx.lineTo(x, h)
    ctx.stroke()
    ctx.closePath()
    ctx.fillStyle = "rgba(255,255,255,0.5)";
    ctx.font = "24px JetBrains Mono"
    let value = Math.round(v * 100) / 100
    let text = `${value}`
    let mx = ctx.measureText(text)
    ctx.fillText(text, w0 - mx.width - 8, y + mx.actualBoundingBoxAscent / 2)

    ctx.beginPath()
    ctx.moveTo(w0, y)
    ctx.lineTo(w, y)
    ctx.stroke()
    ctx.closePath()

    let date = new Date(md.timeMin + i * dt)
    text = `${date.getHours()}:${date.getMinutes()}`
    mx = ctx.measureText(text)
    ctx.fillText(text, x - mx.width / 2, h + mx.actualBoundingBoxAscent + 8)
  }

}


function evaluateTraces(traces: Trace[]): TraceMetadata {
  if (traces.length < 1) {
    return {
      valueMin: 0,
      valueMax: 0,
      timeMin: 0,
      timeMax: 0,
    } as TraceMetadata
  }
  let valueMins = []
  let valueMaxes = []

  let timeMin = []
  let timeMaxes = []

  for (let i = 0; i < traces.length; i++) {

    valueMins.push(Math.min(...traces[i].data))
    valueMaxes.push(Math.max(...traces[i].data))

    timeMin.push(Math.min(...traces[i].time))
    timeMaxes.push(Math.max(...traces[i].time))
  }

  return {
    valueMin: Math.min(...valueMins),
    valueMax: Math.max(...valueMaxes),
    timeMin: Math.min(...timeMin),
    timeMax: Math.max(...timeMaxes),
  } as TraceMetadata

}

function draw() {
  let ctx = state.ctx;
  if (!ctx.canvas) return
  ctx.lineWidth = 1

  ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
  let index = 0;
  let hold = 0;
  const traceColors: string[] = ["10, 132, 255", "48, 209, 88", "255, 214, 10", "255,69,58", "191,90,242"]
  let traces = alignTraces(props.traces)
  let transformed = transformTraces(traces)
  let metadata = evaluateTraces(transformed)

  renderScales(metadata)

  for (let trace of transformed) {
    renderTrace(metadata, trace, traceColors[index])
    renderLegend(metadata, trace, traceColors[index], transformed.length, index)
    index = (index + 1) % traceColors.length;
  }

  // hold = (hold + 1) % 20;
  // if (hold == 0) {
  //   index = (index + 1) % traceColors.length;
  //
  // }


  // // drawLegend()
  // let w = ctx.canvas.width;
  // let h = ctx.canvas.height;
  // let vh = h - 40
  // ctx.fillStyle = preferences.accent
  // // ctx.fillRect(0, 0, w / 2, h / 2)
  //
  // ctx.strokeStyle = "rgba(255,255,255,0.125)";
  // ctx.beginPath()
  // ctx.moveTo(0, h / 2)
  // ctx.lineTo(w, h / 2)
  // ctx.stroke()
  // ctx.closePath()
  //
  // if (!props.data || !props.time) return;
  //
  // let min = Math.min(...props.data)
  // let max = Math.max(...props.data)
  //
  // let minT = Math.min(...props.time)
  // let maxT = Math.max(...props.time)
  //
  // let dT = maxT - minT;
  // ctx.fillStyle = "rgba(255,255,255,0.4)"
  //
  // let dtM = Math.round(dT / (1000 * 60 * 30))
  // ctx.font = "24px JetBrains Mono"
  //
  // ctx.lineWidth = 2;
  // let lx = 0;
  // let ly = h / 2;
  //
  // let pairs = []
  //
  // for (let i = 0; i < props.time.length; i++) {
  //   pairs.push({
  //     time: props.time[i],
  //     data: props.data[i],
  //   })
  // }
  // ctx.save()
  // ctx.restore()
  //
  // let sortedTime = pairs.sort((a, b) => a.time - b.time)
  // for (let i = 0; i < pairs.length; i++) {
  //   let t = sortedTime[i].time || 0
  //   let d = sortedTime[i].data || 0
  //   let x = map_range(t, minT, maxT, 0, w)
  //   let y = map_range(d, min, max, vh - 20, 20)
  //   if (t % (1000 * 60 * 60) == 0) {
  //     ctx.strokeStyle = "rgba(255,255,255,0.2)"
  //     ctx.lineWidth = 2;
  //     ctx.beginPath()
  //     ctx.moveTo(x, 0)
  //     ctx.lineTo(x, vh)
  //     ctx.stroke()
  //     ctx.closePath()
  //     let d = new Date(t)
  //     let time = `${d.getHours() % 12} ${d.getHours() >= 12 ? 'PM' : 'AM'}`
  //     let mxt = ctx.measureText(time)
  //     if (i == 0) {
  //       x += mxt.width / 2
  //     }
  //     ctx.font = "24px JetBrains Mono"
  //     ctx.fillText(time, x - mxt.width / 2, h - mxt.actualBoundingBoxAscent / 2);
  //   } else if (t % (1000 * 60 * 30) == 0) {
  //     ctx.strokeStyle = "rgba(255,255,255,0.1)"
  //     ctx.lineWidth = 1;
  //     ctx.beginPath()
  //     ctx.moveTo(x, 0)
  //     ctx.lineTo(x, vh)
  //     ctx.stroke()
  //     ctx.closePath()
  //   }
  //
  //   ctx.lineWidth = 2;
  //   ctx.strokeStyle = "rgba(0,128,255,1)"
  //
  //   ctx.beginPath()
  //   if (i == 0) {
  //     ctx.moveTo(x, y)
  //   } else if (x - lx > w / 10) {
  //     ctx.moveTo(x, y)
  //   } else {
  //     ctx.moveTo(lx, ly)
  //
  //   }
  //
  //   ctx.lineTo(x, y)
  //   lx = x
  //   ly = y
  //   ctx.stroke()
  //   ctx.closePath()
  // }
  //
  //
  // ctx.strokeStyle = preferences.accent
  //
  //
  // // let offsetP = `${pos.p > 0 ? '+' : ''}${(Math.round(pos.p * 100) / 100).toFixed(2)}° `;
  // // let offsetT = `${pos.t > 0 ? '+' : ''}${(Math.round(pos.t * 100) / 100).toFixed(2)}° `;
  // // let mxp = ctx.measureText(offsetP)
  // // let mxt = ctx.measureText(offsetT)
  // // ctx.fillText(offsetP, 10, 10 + mxp.actualBoundingBoxAscent);
  // // ctx.fillText(offsetT, 10, 40 + mxt.actualBoundingBoxAscent);
}

</script>

<template>
  <div class="h-100">
    <canvas :id="`linechart-${state.uuid}`"
            style="font-family: 'JetBrains Mono',serif; aspect-ratio: 2/1; width: 100%"></canvas>

  </div>
</template>

<style scoped lang="scss">
#loadFont {

}
</style>
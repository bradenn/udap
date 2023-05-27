<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, onUnmounted, reactive, watchEffect} from "vue";


const state = reactive({
    uuid: uuidv4(),
    canvas: {} as HTMLCanvasElement,
    ctx: {} as CanvasRenderingContext2D,
    cache: [] as number[][],
    depth: 0,
    stats: {
        min: 0,
        max: 0,
        avg: 0,
        samples: 0,
        duration: 0,
    },
    ready: false
})

function lerp(a: number, b: number, t: number): number {
    return (1 - t) * a + t * b;
}

interface Resolve {
    begin: number
    values: number[]
}

let props = defineProps<{
    values: number[],
    timecodes: number[],
    chunkSize: number
}>();

onMounted(() => {
    configureCanvas()
})

onUnmounted(() => {
    state.ctx.canvas.remove()
})


const numBack = 2

watchEffect(() => {
    // state.cache[state.depth] = props.values
    // state.depth = (state.depth + 1) % numBack
    draw()
    return props.values
})


function configureCanvas() {
    const _canvas = document.getElementById(`signal-${state.uuid}`)
    state.canvas = _canvas as HTMLCanvasElement
    state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
    let scale = 1
    state.ctx.scale(scale, scale)

    state.canvas.width = state.canvas.clientWidth * scale
    state.canvas.height = state.canvas.clientHeight * scale

    draw()
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
    return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

interface POS {
    x: number
    y: number
}

function drawPattern(ctx: CanvasRenderingContext2D, values: number[], depth: number) {


    let minY = 0;
    let maxY = 4095;
    ctx.fillStyle = "rgba(255,255,255,0.25)";
    ctx.lineWidth = 1
    ctx.strokeStyle = `rgba(255, 128, 1, 0.5)`;
    let w = ctx.canvas.width;
    let h = ctx.canvas.height;


    // ctx.moveTo(0, h / 1.25 - (values[0] - minY) / (maxY - minY) * (ctx.canvas.height) / 1.5)
    let lastX = 0;
    let lastY = h / 1.25 - (values[0] - minY) / (maxY - minY) * (ctx.canvas.height) / 1.5;
    let mass = w / (values.length)
    for (let i = 0; i < values.length; i++) {
        let x = i * mass;
        let y = h / 1.25 - (values[(values.length - 1) - i] - minY) / (maxY - minY) * (ctx.canvas.height) / 1.5
        ctx.moveTo(lastX, lastY)
        ctx.lineTo(x, y)
        lastX = x;
        lastY = y
    }


}

// function drawLegend() {
//   let ctx = state.ctx;
//   if (!ctx.canvas) return
//   let w = state.ctx.canvas.width;
//   let h = state.ctx.canvas.height;
//   let offset = w / props.timecodes.length
//
//   let min = Math.min(...props.timecodes);
//   let max = Math.max(...props.timecodes);
//   ctx.beginPath()
//   ctx.strokeStyle = "rgba(1,128,255,0.5)";
//   let last = props.timecodes[0];
//   ctx.moveTo(map_range(props.timecodes[0], min, max, 0, w), h / 1.25)
//   for (let i = 0; i < props.timecodes.length; i++) {
//     let x = map_range(props.timecodes[i], min, max, 0, w);
//     let y = h / 1.25;
//
//     ctx.moveTo(x, h / 1.25 + 5)
//     ctx.lineTo(x, h / 1.25 - 5)
//     ctx.fillText(`${(props.timecodes[i] - last) / 1000}ms`, x, y + 16)
//     last = props.timecodes[i]
//   }
//   ctx.stroke()
//   ctx.closePath();
// }

function convertHz(hz: number): string {
    if (hz < 1000) {
        return `${hz} Hz`;
    } else if (hz < 1000000) {
        return `${(hz / 1000).toFixed(2)} kHz`;
    } else if (hz < 1000000000) {
        return `${(hz / 1000000).toFixed(2)} MHz`;
    } else if (hz < 1000000000000) {
        return `${(hz / 1000000000).toFixed(2)} GHz`;
    } else {
        return `${(hz / 1000000000000).toFixed(2)} THz`;
    }
}

function draw() {
    let ctx = state.ctx;
    if (!ctx.canvas) return
    ctx.lineWidth = 2
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);


    if (props.values.length > 2) {

        state.stats.min = Math.min(...props.values);
        state.stats.max = Math.max(...props.values);
        state.stats.avg = props.values.reduce((a, b) => a + b) / props.values.length;
        state.stats.samples = props.values.length

    }
    // drawLegend()
    let w = ctx.canvas.width;
    let h = ctx.canvas.height;
    ctx.strokeStyle = "rgba(255,255,255,0.125)";
    ctx.beginPath()
    ctx.moveTo(0, h / 1.25)
    ctx.lineTo(w, h / 1.25)
    ctx.stroke()
    ctx.closePath()

    let scale = 25;

    for (let i = 0; i < w / scale; i++) {
        ctx.beginPath()
        ctx.moveTo(i * scale, h / 1.25 - 2)
        ctx.lineTo(i * scale, h / 1.25 + 2)
        ctx.stroke()
        ctx.closePath()
    }

    state.stats.duration = 0

    ctx.beginPath()
    drawPattern(ctx, props.values, 0)
    ctx.closePath()
    ctx.stroke()

}


</script>

<template>
    <div class="element d-flex">
        <div style="width: 8rem">
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">vMin</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{
                    state.stats.max <= 1 ? '<' : ''
                    }}{{ Math.round(state.stats.min * 100) / 100 || 0 }} mV
                </div>
            </div>
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">vMax</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{ state.stats.max <= 0 ? '<' : '' }}{{
                    Math.round(state.stats.max * 100) / 100 || 0
                    }} mV
                </div>
            </div>
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">vAvg</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{
                    Math.round(state.stats.avg * 100) / 100 || 0
                    }} mV
                </div>
            </div>
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">dS</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{
                    state.stats.samples
                    }} s
                </div>
            </div>
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">dT</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{
                    Math.round(state.stats.duration / 1000.0 * 10) / 10
                    }} ms
                </div>
            </div>
            <div class="d-flex justify-content-between px-2 py-1 align-items-center">
                <div class="label-c2 label-w500 label-o3">F</div>
                <div class="label-c2 label-w600 label-o4" style="font-family: JetBrains Mono,serif">
                    {{
                    convertHz((state.stats.samples / (state.stats.duration / 1000 / 1000)))
                    }}
                </div>
            </div>
        </div>
        <div class="canvas-container ">

            <canvas :id="`signal-${state.uuid}`" class="inner-canvas"></canvas>
        </div>
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
  width: 100%;
  height: 12rem;
  align-items: center;
  border-radius: 8px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
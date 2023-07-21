<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {v4 as uuidv4} from "uuid";

import {onMounted, onUnmounted, reactive, watchEffect} from "vue";

const state = reactive({
    uuid: uuidv4(),
    canvas: {} as HTMLCanvasElement,
    ctx: {} as CanvasRenderingContext2D,
    ready: false
})

let props = defineProps<{
    frequency?: number,
    frequencies?: number[],
    phase: number
}>();

onMounted(() => {
    configureCanvas()
})

onUnmounted(() => {
    state.ctx.canvas.remove()
})


function lerp(a: number, b: number, t: number): number {
    return a + (b - a) * t;
}

watchEffect(() => {
    draw()
    return props
})

function configureCanvas() {
    const _canvas = document.getElementById(`signal-${state.uuid}`)
    state.canvas = _canvas as HTMLCanvasElement
    state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
    let scale = 1
    state.ctx.scale(scale, scale)

    state.ready = true
    state.canvas.width = state.canvas.clientWidth * scale
    state.canvas.height = state.canvas.clientHeight * scale

    draw()
}

function draw() {
    let ctx = state.ctx;
    if (!ctx.canvas) return
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    ctx.lineWidth = 2


    let w = ctx.canvas.width;
    let h = ctx.canvas.height;
    ctx.strokeStyle = "rgba(255,255,255,0.125)";
    ctx.beginPath()
    ctx.moveTo(0, h / 2)
    ctx.lineTo(w, h / 2)
    ctx.stroke()
    ctx.closePath()

    let scale = 10;
    for (let i = 0; i < w / scale; i++) {
        ctx.beginPath()
        ctx.moveTo(i * scale, h / 2 - 5)
        ctx.lineTo(i * scale, h / 2 + 5)
        ctx.stroke()
        ctx.closePath()
    }

    ctx.strokeStyle = "rgba(255,128,1,0.5)";
    ctx.fillStyle = "rgba(255,255,255,0.25)";

    if (props.frequency) {
        ctx.beginPath()
        let dx = 1 / props.frequency
        ctx.moveTo(0, h / 2)
        for (let i = 0; i < w; i++) {
            let dy = Math.sin((Math.PI * 2) * (dx * i)) * h / 4
            ctx.lineTo(i, h / 2 + dy)
        }
        ctx.stroke()
        ctx.closePath()
    } else if (props.frequencies) {
        ctx.beginPath()

        ctx.moveTo(0, h / 2)
        for (let i = 0; i < w; i++) {
            let sumY = 0;
            for (let frequency of props.frequencies) {
                let dy = Math.sin((Math.PI * 2) * ((1 / frequency) * i)) * h / 4
                sumY += dy
            }

            ctx.lineTo(i, h / 2 + sumY)
        }
        ctx.stroke()
        ctx.closePath()
    }


}


</script>

<template>
    <div class="element">
        <div class="d-flex justify-content-between" style="height: 1rem">
            <div v-if="props.frequency" class="label-c2 label-o3 label-w600 px-1">
                Frequency: {{
                props.frequency
                }}Hz
            </div>
            <div v-else-if="props.frequencies"
                 class="label-c2 label-o3 label-w600 px-1">Frequencies: {{
                props.frequencies.join("Hz, ")
                }}Hz
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
  //height: 100%;

}

.canvas-container {
  display: flex;
  flex-direction: row;
  justify-content: center;
  //height: 100%;
  width: 100%;
  align-items: center;
  border-radius: 8px;
  background-color: hsla(214, 9%, 28%, 0.2);
  padding: 6px
}
</style>
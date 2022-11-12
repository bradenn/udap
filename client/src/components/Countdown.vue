<!-- Copyright (c) 2022 Braden Nicholson -->


<script lang="ts" setup>
import {onMounted, reactive} from "vue";

interface Loader {
    size?: string
    percent: number
}

let props = defineProps<Loader>()

let state = reactive({
    ctx: {} as CanvasRenderingContext2D,
    width: 0,
    height: 0,
})

onMounted(() => {
    mountCanvas()
})

function drawSemiArc() {

    let ctx = state.ctx
    let thickness = 6
    ctx.strokeStyle = "rgba(179,115,26,1)"
    ctx.fillStyle = "rgba(179,115,26,1)"
    ctx.lineWidth = thickness

    let radius = state.width / 2 - (thickness * 1.5)
    let current = (Math.PI / 180) * (270 + 360 * props.percent) % 360
    let start = (Math.PI / 180) * 270
    ctx.strokeStyle = "rgba(179,115,26,0.125)"
    ctx.beginPath()
    ctx.arc(state.width / 2, state.width / 2, radius, 0, (Math.PI / 180) * 359, false)
    ctx.stroke()
    ctx.closePath()
    ctx.strokeStyle = "rgba(179,115,26,1)"
    ctx.beginPath()
    ctx.arc(state.width / 2, state.width / 2, radius, start, current, false)
    ctx.stroke()
    ctx.closePath()

    let capRadius = thickness / 2

    ctx.beginPath()
    let cx = state.width / 2 + radius * Math.cos(current)
    let cy = state.width / 2 + radius * Math.sin(current)

    ctx.ellipse(cx, cy, capRadius, capRadius, 0, 0, (Math.PI / 180) * 360, false)

    cx = state.width / 2 + radius * Math.cos(start)
    cy = state.width / 2 + radius * Math.sin(start)

    ctx.ellipse(cx, cy, capRadius, capRadius, 0, 0, (Math.PI / 180) * 360, false)
    ctx.fill()

    ctx.closePath()
}


function render() {
    drawSemiArc()
}

function animate() {
    state.ctx.clearRect(0, 0, state.width, state.height)
    render()
    requestAnimationFrame(animate)
}

function mountCanvas() {

    const canvas = document.getElementById("countdown") as HTMLCanvasElement
    if (!canvas) return;

    const ctx = canvas.getContext('2d') as CanvasRenderingContext2D
    if (!ctx) return
    state.ctx = ctx
    let scale = 2
    ctx.scale(scale, scale)


    canvas.width = canvas.clientWidth * scale
    canvas.height = canvas.clientHeight * scale

    state.width = ctx.canvas.width
    state.height = ctx.canvas.height
    ctx.translate(0, 0)
    ctx.clearRect(0, 0, state.width, state.height)

    animate()

}


</script>

<template>
    <div>
        <canvas id="countdown" class="countdown-canvas"></canvas>
    </div>
</template>

<style lang="scss" scoped>
.countdown-canvas {
  width: 0.75rem;
  height: 0.75rem;
}
</style>
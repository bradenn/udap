<!-- Copyright (c) 2023 Braden Nicholson -->
<script lang="ts" setup>

import {v4 as uuidv4} from "uuid";
import {onMounted, reactive} from "vue";

const props = defineProps<{
    data: string
}>()

const state = reactive({
    uuid: uuidv4(),
    ctx: {} as CanvasRenderingContext2D
})

onMounted(() => {
    setupCanvas()
})

function setupCanvas() {
    const canvas = document.getElementById(state.uuid) as HTMLCanvasElement
    if (!canvas) return
    const ctx = canvas.getContext("2d") as CanvasRenderingContext2D
    if (!ctx) return

    state.ctx = ctx

    render()
}

function render() {
    const ctx = state.ctx
    const w = ctx.canvas.width
    const h = ctx.canvas.height


    const cw = 12, ch = 12

    const dx = w / cw
    const dy = h / ch

    const points = 128
    const dz = (Math.PI * 2) / points

    const sx = 12, sy = 12

    let last = Math.random()
    for (let i = 0; i < points; i++) {
        const t = i * dz
        const r = Math.cos(43 * t)
        const x = Math.cos(t) * r * w / 2
        const y = Math.sin(t) * r * h / 2
        if (Math.random() > 0.25) {
            ctx.fillStyle = "rgba(255,128, 1, 0.5)"
        } else {
            ctx.fillStyle = "rgba(1,128, 255, 0.5)"
        }
        ctx.fillRect(w / 2 + x - sx / 2, h / 2 + y - sy / 2, sx, sy)
    }
    ctx.fill()

}


</script>

<template>
    <canvas :id="`${state.uuid}`" class="hash-canvas"></canvas>
</template>

<style scoped>
.hash-canvas {
    width: 1rem;
    height: 1rem;
}
</style>
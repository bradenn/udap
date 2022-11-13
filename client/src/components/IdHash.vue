<script lang="ts" setup>


import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Preferences, Status} from "@/types";
import type {Remote} from "@/remote";

let state = reactive({
    canvas: {} as HTMLCanvasElement,
    tick: 0,
    menu: false,
    reloading: true,
    connected: false,
    zoneEntity: {} as Entity,
    zoneAttribute: {} as Attribute,
    status: {} as Status
})

let preferences: Preferences = inject("preferences") as Preferences
let remote: Remote = inject('remote') as Remote
let system: any = inject('system')

onMounted(() => {
    update()
    state.reloading = false
})

watchEffect(() => {
    state.connected = remote.connected
    update()
    return state.zoneAttribute
})

function update() {
    let entity = remote.entities.find(e => e.name === 'faces')
    if (!entity) return
    state.zoneEntity = entity

    let attr = remote.attributes.find(e => e.key === 'deskFace')
    if (!attr) return
    state.zoneAttribute = attr

    let stat = JSON.parse(attr.value) as Status
    if (!stat) return

    state.status = stat

}


onMounted(() => {
    loadCanvas()
    drawCanvas()
    setInterval(drawCanvas, 100)
})

function loadCanvas() {
    state.canvas = document.getElementById("canvas-id") as HTMLCanvasElement
    const ctx = state.canvas.getContext('2d')
    if (!ctx) return

    ctx.scale(1, 1)

}

function drawCanvas() {

    const ctx = state.canvas.getContext('2d')
    if (!ctx) return
    ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)

    let n = Math.round(300 * (1 - state.status.predictions[0].distance));

    let diam = 600;
    let div = (Math.PI * 2) / n;


    for (let i = 0; i < n; i++) {
        let ax = diam * Math.cos(i * div)
        let ay = diam * Math.sin(i * div)
        let bx = diam / 1.3 * Math.cos(i * div)
        let by = diam / 1.3 * Math.sin(i * div)
        ctx.strokeStyle = `rgba(255,255,255,${n * div % 1})`
        ctx.beginPath()
        ctx.moveTo(ax * 0.5 + diam / 2, ay * 0.5 + diam / 2)
        ctx.lineTo(bx * 0.5 + diam / 2, by * 0.5 + diam / 2)
        ctx.closePath()
        ctx.stroke()

    }

    state.tick += 1;


}

</script>

<template>
    <div class="canvas-container top">
        {{ state.status.predictions[0].name }}
        <canvas id="canvas-id" height="600" width="600"></canvas>
    </div>
</template>

<style lang="scss" scoped>
.canvas-container {
  height: 600px;
  width: 600px;

  aspect-ratio: 1/1 !important;
}

</style>

<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute} from "@/types";
import type {Remote} from "@/remote";

let state = reactive({
    raw: {} as Attribute,
    data: [] as number[][],
    context: {} as CanvasRenderingContext2D,
    loading: false,
    rotation: {
        x: 0,
        y: 3 * Math.PI / 4,
    }
})

let remote = inject("remote") as Remote

let map = new Map<string, number>();

onMounted(() => {
    state.loading = true
    map.set("year", 0)
    map.set("month", 1)
    map.set("day", 2)
    map.set("hour", 3)
    map.set("minute", 4)
    map.set("second", 5)
    map.set("millisecond", 6)
    map.set("bx", 7)
    map.set("by", 8)
    map.set("bz", 9)
    map.set("vx", 10)
    map.set("vy", 11)
    map.set("vz", 12)
    map.set("density", 13)
    map.set("temperature", 14)
    map.set("idk", 15)

    updateData()

    loadDom()

})

interface Point {
    x: number
    y: number
}

function OrthoProject(x1: number, y1: number, z1: number): Point {
    let top = -state.context.canvas.height / 2, bottom = state.context.canvas.height / 2;
    let left = -state.context.canvas.width / 2, right = state.context.canvas.width / 2;
    let S = 1 / (Math.tan((90 / 2) * (Math.PI / 180)))

    let r = state.context.canvas.width / 2 / state.context.canvas.height / 2
    let a = Math.PI
    let f = state.context.canvas.height / 2
    let n = -state.context.canvas.height / 2

    let x = Math.cos(state.rotation.y) * x1 + Math.sin((state.rotation.y)) * z1
    let y = y1;
    let z = -Math.sin(state.rotation.y) * x1 + Math.cos(state.rotation.y) * z1

    let px = (1 / r) * Math.atan(a / 2) * x
    let py = Math.atan(a / 2) * y
    let pz = z * (f / (f - n)) + 1
    let pw = z * (f * n / (f - n))

    return {
        x: px,
        y: py,
    }
}

watchEffect(() => {
    return updateData()
})


function updateData() {
    let raw = remote.attributes.find(attr => attr.key === "gauss")
    if (!raw) return
    state.raw = raw
    let parsed = JSON.parse(state.raw.value) as number[][]
    if (!parsed) return
    state.data = parsed
    state.loading = false
    drawData()
    return parsed

}


function loadDom() {
    let chart = document.getElementById(`magneto`) as HTMLCanvasElement
    if (!chart) return;
    let ctx = chart.getContext('2d') as CanvasRenderingContext2D
    if (!ctx) return
    let scale = 1.75
    ctx.scale(scale, scale)

    chart.width = chart.clientWidth * scale
    chart.height = chart.clientHeight * scale
    ctx.translate(0, 0)


    state.context = ctx

    setInterval(animate, 32)
}

function animate() {
    state.rotation.y = Math.PI
    drawData()
}

// Format for SWPC -> https://github.com/MSTEM-QUDA/SWMF/blob/stable/Param/SWPC/IMF.dat

function drawData() {
    if (state.loading) return
    const ctx = state.context
    let halfW = ctx.canvas.height / 2;
    ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height)
    ctx.translate(halfW * 2, halfW);
    let center = OrthoProject(0, 0, 0)
    let radius = ctx.canvas.height / 15
    let Re = ctx.canvas.height
    // year mo dy hr min sec msec bx by bz vx vy vz dens temp
    ctx.font = "500 32px SF Pro Rounded"
    ctx.fillStyle = "rgba(255,255,255,1)"
    let sunWard = OrthoProject(center.x + halfW / 2, center.y, halfW / 2)
    ctx.fillText("SUN", sunWard.x, sunWard.y)
    ctx.beginPath()
    ctx.strokeStyle = `rgba(255, 0, 0, 1)`
    ctx.moveTo(center.x, center.y)
    ctx.lineTo(sunWard.x, sunWard.y)
    ctx.stroke()
    ctx.closePath()

    ctx.fillText("TERMINATOR", center.x + 20, center.y - halfW + 30)
    ctx.beginPath()
    ctx.strokeStyle = `rgba(255, 0, 0, 1)`
    ctx.moveTo(center.x, center.y)
    ctx.lineTo(center.x, center.y - halfW)
    ctx.stroke()
    ctx.closePath()

    if (!state.data[7]) return

    let minDens = 10000;
    let maxDens = -10000;
    for (let i = 0; i < state.data.length; i++) {
        let dat = state.data[i]
        if (!dat) continue
        let d = dat.length > 13 ? dat[13] : 0;
        if (d < minDens) {
            minDens = d
        } else if (d > maxDens) {
            maxDens = d;
        }
    }
    ctx.lineWidth = 2

    ctx.beginPath()
    ctx.strokeStyle = `rgba(128, 255, 128, 0.8)`
    ctx.ellipse(center.x, center.y, radius, radius, 0, 0, Math.PI * 2)
    ctx.stroke()
    ctx.closePath()
    let width = ctx.canvas.width
    let scale = radius
    let dirLine = 10;
    for (let datum of state.data) {
        if (!datum) continue
        let x = (datum[7] * scale)
        let y = (datum[9] * scale)
        let z = (datum[8] * scale)
        let x2 = (datum[10] / 100)
        let y2 = (datum[12] / 100)
        let z2 = (datum[11] / 100)
        let pt = OrthoProject(x, y, z)
        let pt2 = OrthoProject(x + x2, y + y2, z + z2)
        // let y = halfW + (datum[12] / 50) * halfW / 2

        ctx.beginPath()
        ctx.strokeStyle = `rgba(${255 - ((datum[13] - minDens) / (maxDens - minDens)) * 255}, ${((datum[13] - minDens) / (maxDens - minDens)) * 255}, 64, 1)`
        ctx.moveTo(pt.x - 1, pt.y - 1)
        // + ((pt.x - center.x) / 20)  + ((pt.y - center.y) / 20)
        ctx.lineTo(pt.x + 1, pt.y + 1)

        ctx.closePath()
        ctx.stroke();

    }
    ctx.translate(-halfW * 2, -halfW);

}

</script>

<template>
    <canvas id="magneto" style="width: 100%; height: 100%"></canvas>

</template>

<style scoped>

</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onBeforeUnmount, onMounted, reactive, watchEffect} from "vue";
import type {Preferences} from "@/types";
import canvas from "@/composables/canvas";
import core from "@/core";

const state = reactive({
    thumb: {
        w: 0,
        h: 0,
        x: 0,
        y: 0
    },
    slider: {
        w: 0,
        h: 0,
        x: 0,
        y: 0
    },
    pos: {
        x: 0,
        y: 0
    },
    lastHaptic: 0,
    value: 0,
    result: 0,
    hold: 0,
    notches: 50
})
const haptics = core.haptics()

function requestHaptic() {
    if (Date.now().valueOf() - state.lastHaptic >= 25) {
        state.lastHaptic = Date.now().valueOf()
        // haptics.tap(2, 1, 12)

    }
}

function draw(c: CanvasRenderingContext2D) {
    c.strokeStyle = "rgba(0,0,0,0.125)"
    c.fillStyle = "rgba(255,255,255,0.25)"
    c.lineWidth = 2


    // for (let i = 1; i < state.notches - state.value; i++) {
    //     c.beginPath()
    //
    //     let width = 16
    //     let color = (i / numNotches) / 2
    //     c.fillStyle = `rgba(255,255,255,${color})`
    //     c.roundRect(xPadding, state.slider.h - i * dx + state.thumb.h / 2, state.thumb.w, state.thumb.h / 2, 8);
    //
    //     c.fill()
    //     c.stroke()
    //     c.closePath()
    // }
    // c.fillStyle = `rgba(255, 255, 255, 0.025)`
    // c.beginPath()
    // c.roundRect(xPadding, state.pos.y + state.thumb.h / 4, state.thumb.w, state.thumb.h * 2, 12);
    // c.stroke()
    // c.fill()
    // c.closePath()
    //
    // for (let i = state.notches - state.value + 1; i < state.notches; i++) {
    //     c.beginPath()
    //
    //     let width = 16
    //     let color = (i / numNotches) / 2
    //     c.fillStyle = `rgba(255, 255, 255, 0.01)`
    //     c.roundRect(xPadding, state.slider.h - i * dx + state.thumb.h / 2, state.thumb.w, state.thumb.h / 2, 8);
    //
    //     c.fill()
    //     // c.stroke()
    //     c.closePath()
    // }
    
    c.fillStyle = `rgba(255, 255, 255, 0.025)`
    c.beginPath()
    c.roundRect(0, state.pos.y, state.thumb.w, state.slider.h - state.pos.y, 12);
    c.stroke()
    c.fill()
    c.closePath()
}

function mouseDrag() {
    const element = document.getElementById(`${canvas.uuid}`) as HTMLCanvasElement
    if (!element) return;
    let rect = element.getBoundingClientRect()

    state.slider.x = rect.x
    state.slider.y = rect.y
    state.slider.h = rect.height * 2
    state.slider.w = rect.width * 2

    state.thumb.w = state.slider.w
    state.thumb.h = state.slider.h / state.notches
    state.thumb.x = 0
    state.thumb.y = 0


    element.addEventListener("mousemove", moveSlider)
}

function moveSlider(ev: MouseEvent) {


    let adjustedPosition = ev.offsetY * 2

    let dx = state.slider.h / state.notches
    let np = adjustedPosition - state.thumb.h / 2

    if (np / dx < 0) {
        state.value = 0

        state.pos.y = state.value * dx
        requestHaptic()
    } else if (np > state.slider.h - state.thumb.h * 4) {
        state.value = state.notches

        state.pos.y = state.value * dx
        requestHaptic()
    } else {
        state.value = Math.floor(np / dx)

        state.pos.y = state.value * dx
    }


}


onMounted(() => {
    mouseDrag()
    canvas.setupCanvas(draw)
})

onBeforeUnmount(() => {
    canvas.dispose()
})

watchEffect(() => {
    if (state.hold == state.value) return
    state.hold = state.value
    requestHaptic()
    state.result = state.notches - state.value
    return state.value
})

const preferences = inject("preferences") as Preferences


</script>

<template>
    <div class="element v-slider p-0">

        <div class="h-100" style="outline: 1px solid white">
            <canvas :id="`${canvas.uuid}`" style="height: 100%; width: 96px !important;"></canvas>
        </div>
    </div>
    {{ Math.round((state.result / state.notches) * 100) }}
</template>

<style lang="scss" scoped>
.element.v-slider {

  height: 100%;

}
</style>
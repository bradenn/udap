<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onBeforeUnmount, onMounted, reactive, watchEffect} from "vue";
import type {Preferences} from "@/types";
import {useCanvas} from "@/composables/canvas";
import {v4 as uuidv4} from "uuid";
import core from "@/core";

interface Props {
    change: (a: number) => void
}

defineProps<Props>()

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
    notches: 50,
    uuid: uuidv4()
})

const haptics = core.haptics()


function requestHaptic() {
    if (Date.now().valueOf() - state.lastHaptic >= 25) {
        state.lastHaptic = Date.now().valueOf()
        // haptics.tap(2, 1, 12)

    }
}

function draw(c: CanvasRenderingContext2D): void {
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

    c.fillStyle = `rgba(255, 149, 0, 0.3)`
    c.beginPath()
    c.roundRect(0, state.pos.y, state.thumb.w, state.slider.h - state.pos.y, 0);

    c.stroke()
    c.fill()
    c.closePath()
    c.beginPath()
    c.roundRect(state.thumb.w / 2 - state.thumb.w / 4 / 2, state.pos.y + 16, state.thumb.w / 4, 10, 12);
    c.stroke()
    c.fill()
    c.closePath()
}

function mouseDrag() {
    const element = document.getElementById(`${state.uuid}`) as HTMLCanvasElement
    if (!element) return;
    let rect = element.getBoundingClientRect()

    state.slider.x = rect.x
    state.slider.y = rect.y
    state.slider.h = rect.height * 2
    state.slider.w = rect.width * 2

    state.thumb.w = state.slider.w
    state.thumb.h = state.slider.h
    state.thumb.x = 0
    state.thumb.y = 0


    element.addEventListener("mousemove", moveSlider)
}

function moveSlider(ev: MouseEvent) {
    state.pos.y = ev.offsetY * 2
}


onMounted(() => {
    mouseDrag()
    useCanvas(state.uuid, draw)
})

onBeforeUnmount(() => {
})

watchEffect(() => {
    if (state.hold == state.value) return

    if (state.hold - state.value) {
        state.hold = state.value
    }

    requestHaptic()
    state.result = state.notches - state.value
    return state.value
})

const preferences = inject("preferences") as Preferences


</script>

<template>
    {{ state.result }}
    <div class="element v-slider p-0">
        <div class="h-100" style="outline: 1px solid white">
            <canvas :id="`${state.uuid}`" style="height: 100%; width: 128px !important;"></canvas>
        </div>
    </div>

</template>

<style lang="scss" scoped>
.element.v-slider {

  height: 100%;

}
</style>
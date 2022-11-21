<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {reactive} from "vue";
import core from "@/core";

const props = defineProps<{
    horizontal?: boolean
}>()

let state = reactive({
    dragging: false,
    target: {} as HTMLElement,
    xDecelerate: 0,
    yDecelerate: 0,

    xStart: 0,
    yStart: 0,

    xActual: 0,
    yActual: 0,

    xVelocity: 0,
    yVelocity: 0,

    xLastPoll: 0,
    yLastPoll: 0,

    xLastPos: 0,
    yLastPos: 0,

    xSilence: false,
    ySilence: false,

    lastTap: 0,
})

function dragStart(e: MouseEvent) {
    state.xVelocity = 0
    state.yVelocity = 0
    clearInterval(state.xDecelerate)
    clearInterval(state.yDecelerate)
    state.xLastPos = 0
    state.yLastPos = 0
    state.xLastPoll = 0
    state.yLastPoll = 0
    state.dragging = false
    state.yStart = e.clientY
    state.xStart = e.clientX
}

const haptics = core.haptics()

function drag(e: MouseEvent) {

    if (!state.dragging && Math.abs(state.yStart - e.clientY) > 5 || (props.horizontal && Math.abs(state.xStart - e.clientX) > 5)) {
        state.dragging = true
    }

    if (!state.dragging) return
    let element = e.currentTarget as HTMLElement
    if (!element) return;
    state.target = element

    // haptics(1, 1, 50)


    if (props.horizontal) {
        let poll = new Date().valueOf()
        let dt = state.xLastPoll - poll
        let dx = state.xLastPos - e.clientX

        element.scrollLeft += (state.xStart - e.clientX)
        state.xActual = element.scrollLeft
        state.xStart = e.clientX
        if (state.xActual <= 0 || state.xActual + element.clientWidth >= element.scrollWidth) {
            if (!state.xSilence) {
                state.xSilence = true
                haptics.tap(2, 1, 25)
            }
        } else {
            state.xSilence = false
        }

        state.xVelocity = dx / dt
        state.xLastPos = e.clientX
        state.xLastPoll = poll
    } else {
        clearInterval(state.yDecelerate)
        let poll = new Date().valueOf()
        let dt = state.yLastPoll - poll
        let dx = state.yLastPos - e.clientY

        element.scrollTop += state.yStart - e.clientY
        state.yActual = dx
        state.yStart = e.clientY


        state.yVelocity = dx / dt
        state.yLastPos = e.clientY
        state.yLastPoll = poll
        if (element.scrollTop <= 0 && !state.ySilence) {
            // haptics(1, 1, 100)
            state.ySilence = true
        } else if (element.scrollTop > 0) {
            state.ySilence = false
        }
    }

}

function mouseDown() {
    state.lastTap = Date.now().valueOf()
    haptics.tap(1, 1, 50)
}

function dragStop(e: MouseEvent) {
    state.yStart = 0
    state.xStart = 0
    let element = e.currentTarget as HTMLElement
    if (!element) return;

    state.xDecelerate = setInterval(() => {
        state.xVelocity = state.xVelocity * 0.9
        element.scrollLeft -= state.xVelocity * 50
        state.xStart = state.xStart * 0.9

        if (Math.abs(state.xVelocity) <= 0.01 && state.xStart <= 0.1) {
            state.xVelocity = 0
            clearInterval(state.xDecelerate)
        }
    }, 16)

    state.yDecelerate = setInterval(() => {

        state.yVelocity = state.yVelocity * 0.965
        element.scrollTop -= state.yVelocity * 12

        if (Math.abs(state.yVelocity) <= 0.0001) {
            state.yVelocity = 0;
            clearInterval(state.yDecelerate)
        }
    }, 16)


    setTimeout(() => {
        state.dragging = false
    }, 250)

}

function click(e: MouseEvent) {
    if (state.dragging) {
        e.preventDefault()
        e.stopPropagation()
    }
}

</script>
<template>

    <div v-on:mousedown="dragStart" v-on:mousemove="drag" v-on:mouseup="dragStop"
         v-on:click.capture="click">
        <slot></slot>

    </div>
</template>


<style scoped>

</style>
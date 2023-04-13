<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {onBeforeMount, reactive} from "vue";

const props = defineProps<{
    min: number
    max: number
    step: number
    value: number
    change: (value: number) => void
    bg: string
}>()


const state = reactive({
    delta: 0,
    value: 0,
    dx: 0,
    data: {},
    width: 0,
})

onBeforeMount(() => {
    state.value = props.value
})

function setup() {
    state.dx = (state.width - 8) / (props.max - props.min)
}

function locateTouch(value: number): number {

    let step = Math.floor(value / (state.dx))
    step = Math.max(step, props.min)
    step = Math.min(step, props.max - 1)

    return step
}

function touchEnd(e: TouchEvent) {
    props.change(state.value)
}

function handleMove(e: TouchEvent) {
    let c = e.target as HTMLDivElement
    let parent = c.parentElement
    state.width = parent?.clientWidth || 0
    setup()
    let w = parent?.clientWidth || 0

    let touches = e.targetTouches.item(0) || {} as Touch
    c.style.width = `${(state.dx)}px`
    let tx = locateTouch(touches.clientX - state.dx / 2)
    let size = 48
    c.style.marginLeft = `${(tx * state.dx)}px`
    state.value = tx

    state.data = tx
}

</script>


<template>
    <input :class="`slider-${props.bg}`" :value="state.value" class="slider element" type="range" @touchend="touchEnd">
</template>

<style scoped>

.slider-range {
    width: 100%;
}

.grip {
    position: relative;
    height: 4rem;
    /*border: 1px solid rgba(255, 255, 255, 0.1);*/
    border-radius: 8.24px;
    width: 48px;
}

.slider {
    padding: 4px 4px;
}
</style>
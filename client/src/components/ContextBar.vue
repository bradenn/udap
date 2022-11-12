<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>


import {reactive} from "vue";

let state = reactive({
    dragging: false,
    inert: false,
    selected: {} as HTMLElement,
    inner: ""
})

function mouseMove(event: MouseEvent) {
    let r = event.target as HTMLElement
    if (!r) return

    let bar = document.getElementById("the-content-bar") as HTMLElement
    if (!bar) return

    // if (r.parentElement !== bar) {
    //     return
    // }

    let box = bar.getBoundingClientRect()
    let percent = (event.clientX / box.width) * 24
    state.inner = "" + Math.floor(percent)
    state.selected = r
    r.style.gridColumn = state.inner
    r.style.gridRow = "1"

}

function mouseUp(event: MouseEvent) {
    let r = event.target as HTMLElement
    if (!r) return
    state.inert = false
    state.dragging = false
    state.selected = {} as HTMLElement
}

function mouseDown(event: MouseEvent) {
    let r = event.target as HTMLElement
    if (!r) return
    state.inert = true
    state.dragging = true
    state.selected = r
}

</script>

<template>
    <div id="the-content-bar" class="content-bar">
        <slot>

        </slot>
    </div>

</template>

<style lang="scss">
.content-bar > div {
  //outline: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  height: 2rem;
}

.content-bar {
  display: grid;
  grid-gap: 0.125rem;
  width: 100%;
  grid-template-columns: repeat(24, 1fr);
  grid-template-rows: repeat(1, 1fr);
  //outline: 1px solid white;
  //backdrop-filter: blur(24px);
  //box-shadow: inset 0 -2px 2px 0.25px rgba(255, 255, 255, 0.1);
}
</style>
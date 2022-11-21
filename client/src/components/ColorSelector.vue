<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import FixedScroll from "@/components/scroll/FixedScroll.vue";
import {inject, reactive} from "vue";
import type {Preferences} from "@/types";


const prefs = inject("preferences") as Preferences


interface Color {
    name: string
    hex: string
}

let state = reactive({
    selected: prefs.appdata.colors as number[]
})

function select(to: number) {
    if (state.selected.includes(to)) {
        state.selected = state.selected.filter(s => s != to)
        prefs.appdata.colors = state.selected
    } else {
        state.selected.push(to)
        prefs.appdata.colors = state.selected
    }
}


</script>

<template>
    <div class="element color-selector">
        <FixedScroll :horizontal="true" class="color-grid">
            <div v-for="c in Array(360).keys()" :key="c"
                 :class="`${state.selected.includes(c)?'selected':''}`"
                 :style="`background-color: hsl(${c}, 100%, 50%);`"
                 class="color" @click="() => select(c)">

            </div>
        </FixedScroll>
    </div>
</template>

<style lang="scss" scoped>
div.element.color-selector {

  div.color {
    opacity: 1;
    //border-radius: 0.125rem;
    padding: 0.25rem;

    aspect-ratio: 1.61803398874/1 !important;
    //width: 100%;

    height: 1.1rem;


  }


  div.color.selected {
    z-index: 10;
    opacity: 1;

    box-shadow: 0 0 32px 8px rgba(0, 0, 0, 0.5), inset 0 0 1px 1px rgba(255, 255, 255, 0.3);
    border-radius: 0.25rem;

    scale: 1.1;
    padding: 0.25rem;

  }

  .color-grid {
    flex-direction: column;
    flex-wrap: wrap;
    padding: 0.25rem;
    display: flex;
    //grid-auto-flow: column;
    grid-gap: 0.125rem;
    height: 26rem;
    width: 100%;

    overflow-x: scroll;
  }

  .color-text {

    //text-shadow: 0 0 2px rgba(0, 0, 0, 0.25);
    filter: invert(1) contrast(200%);
  }
}
</style>
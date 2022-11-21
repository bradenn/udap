<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import FixedScroll from "@/components/scroll/FixedScroll.vue";

import colors from '@/assets/colors.json'
import {inject, reactive} from "vue";
import type {Preferences} from "@/types";


const prefs = inject("preferences") as Preferences


interface Color {
    name: string
    hex: string
}

let state = reactive({
    selected: prefs.appdata.colors as string[]
})

function select(to: string) {
    if (state.selected.includes(to)) {
        state.selected = state.selected.filter(s => s != to)
        prefs.appdata.colors = state.selected
    } else {
        state.selected.push(to)
        prefs.appdata.colors = state.selected
    }

}


function hexToRgb(hex: string) {
    let result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : {r: 0, g: 0, b: 0};
}

function RGBToHSL(r: number, g: number, b: number): any {
    // Make r, g, and b fractions of 1
    r /= 255;
    g /= 255;
    b /= 255;

    // Find greatest and smallest channel values
    let cmin = Math.min(r, g, b),
        cmax = Math.max(r, g, b),
        delta = cmax - cmin,
        h = 0,
        s = 0,
        l = 0;
    // Calculate hue
    // No difference
    if (delta == 0)
        h = 0;
    // Red is max
    else if (cmax == r)
        h = ((g - b) / delta) % 6;
    // Green is max
    else if (cmax == g)
        h = (b - r) / delta + 2;
    // Blue is max
    else
        h = (r - g) / delta + 4;

    h = Math.round(h * 60);

    // Make negative hues positive behind 360°
    if (h < 0)
        h += 360;
    // Calculate lightness
    l = (cmax + cmin) / 2;

    // Calculate saturation
    s = delta == 0 ? 0 : delta / (1 - Math.abs(2 * l - 1));

    // Multiply l and s by 100
    s = +(s * 100).toFixed(1);
    l = +(l * 100).toFixed(1);


    return {
        h: h,
        s: s,
        l: l,
    }
}

function badColors(a: string): boolean {
    let j = hexToRgb(a)
    let l = RGBToHSL(j.r, j.g, j.b)

    if (l.s < 50) return false;
    if (l.l > 65 || l.l < 25) return false;

    return true

}

function sort(a: string, b: string) {
    let j = hexToRgb(a)
    let l = RGBToHSL(j.r, j.g, j.b)
    let k = hexToRgb(b)
    let p = RGBToHSL(k.r, k.g, k.b)

    if (l.h <= p.h) {
        return 1
    } else if (l.h > p.h) {
        return -1
    }
    return 0
}

</script>

<template>
    <div class="element color-selector">
        <FixedScroll :horizontal="true" class="color-grid">
            <div v-for="c in colors.map(a => a.hex).filter(badColors).sort(sort)" :key="c"
                 :class="`${state.selected.includes(c)?'selected':''}`"
                 :style="`background-color: ${c};`"
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

    height: 2.2rem;


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
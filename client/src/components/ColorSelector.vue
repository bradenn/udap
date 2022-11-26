<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import FixedScroll from "@/components/scroll/FixedScroll.vue";
import {inject, reactive} from "vue";
import type {Preferences} from "@/types";
import Context from "@/views/terminal/Context.vue";
import Header from "@/components/Header.vue";


const prefs = inject("preferences") as Preferences

let props = defineProps<{
    done: () => void
}>()

let state = reactive({
    selected: prefs.appdata.colors as number[],
    brightness: 45,
    maxHue: 360,
    frame: 0
})

function animate() {
    if (state.maxHue >= 360) {
        cancelAnimationFrame(state.frame)
        return
    }
    state.frame = requestAnimationFrame(animate)
    state.maxHue += 1
}


function select(to: number) {
    if (state.selected.includes(to)) {
        state.selected = state.selected.filter(s => s != to)
    } else {
        state.selected.push(to)
    }
}

function save() {
    prefs.appdata.colors = state.selected

    done()
}

function done() {
    prefs.appdata.colors = state.selected
    if (props.done) {
        props.done()
    }
}

</script>

<template>
    <Context>
        <div class="color-container d-flex flex-column justify-content-start">
            <Header :done="done" :save="save" name="Primary Colors"></Header>
            <div class="element color-selector" @click.stop>
                <FixedScroll :horizontal="true" class="color-grid">
                    <div v-for="c in Array(state.maxHue).keys()" :key="c"
                         :class="`${state.selected.includes(c)?'selected':''}`"
                         :style="`background-color: hsl(${c}, 100%, ${15 + state.brightness * 0.5}%);`"
                         class="color" @click="() => select(c)">

                    </div>
                </FixedScroll>
            </div>
        </div>
    </Context>
</template>

<style lang="scss" scoped>


div.color-container {
  z-index: 100;
  position: absolute;
  top: 3rem;
  margin: 1rem;
  width: calc(100vw - 2rem);
  height: calc(100vh);
  display: block;
}

div.color {
  opacity: 1;
  //border-radius: 0.125rem;
  padding: 0.25rem;

  aspect-ratio: 0.5/1 !important;
  //width: 100%;

  height: 1.6rem;

}

div.color.selected {
  z-index: 8;
  transform: translate3d(0, 0, 0);

  box-shadow: inset 0 0 1px 2px white;
  border-radius: 0.125rem;
  //transform: scale3d(105%, 105%, 105%);
  padding: 0.25rem;

}

.color-grid {
  flex-direction: column;
  flex-wrap: wrap;
  padding: 0.25rem;

  display: flex;
  //grid-auto-flow: column;
  grid-gap: 0.0625rem;
  height: 26rem;
  width: 100%;

  overflow-x: scroll;
}

.color-text {

}

div.element.color-selector {


}
</style>
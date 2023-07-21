<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, reactive} from "vue";

import Color from "@/components/Color.vue";
import type {Preferences} from "@/types";
import FixedScroll from "@/components/scroll/FixedScroll.vue";
import Loader from "@/components/Loader.vue";
import ColorSelector from "@/components/ColorSelector.vue";

const props = defineProps<{
  change?: (a: number) => void
  selected?: number,
  loader?: boolean
}>()

const state = reactive({
  showColorSelector: false,
  prefs: {} as Preferences
})


function setColor(color: number) {
  if (props.change) {
    props.change(color)
  }
}

onMounted(() => {
  let prefs = inject("preferences") as Preferences
  if (!prefs.appdata) {
    prefs.appdata = {
      colors: []
    }
  }
  state.prefs = prefs as Preferences
})


</script>

<template>
  <ColorSelector v-if="state.showColorSelector" :done="() => state.showColorSelector = false"></ColorSelector>
  <div :class="`${props.selected?'selected':''}`" class="element color-container">
    <div class="d-flex justify-content-between px-1 pb-1">
      <div class="label-c1 label-o3 label-w500">Color</div>
      <div class="label-c1 label-o3 label-w500">
        <Loader v-if="props.loader"></Loader>
        {{ props.selected }}
      </div>

    </div>
    <FixedScroll :horizontal="true" style="overflow-x: scroll">
      <div v-if="state.prefs.appdata" class="d-flex gap-1">
        <Color v-for="color in state.prefs.appdata.colors.sort()" :color="color"
               :selected="props.selected === color"
               @click="() => setColor(color)">
        </Color>
        <div :class="`${props.selected?'selected':''}`" class="color subplot surface">
          <div :style="`background-color: rgba(255,255,255,0.024);`"
               class="swatch" @click="() => {state.showColorSelector = !state.showColorSelector}">
            􀅼
          </div>
        </div>
      </div>
    </FixedScroll>
  </div>
</template>

<style lang="scss" scoped>
div.color.subplot.surface {

  border-radius: 0.25rem;
  padding: 0.125rem;

  aspect-ratio: 2/2 !important;
  //width: 100%;
  background-color: rgba(255, 255, 255, 0.025);
  height: 2.25rem;


}

div.swatch {
  border-radius: 0.1875rem;
  width: 100%;
  height: 100%;
  opacity: 0.6;
  display: flex;
  justify-content: center;
  align-items: center;

}

div.element.color-container {


  //width: 100%;


}
</style>
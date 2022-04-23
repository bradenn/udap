<script lang="ts" setup>

import {reactive} from "vue";

import type {Attribute} from "@/types";

const props = defineProps<{
  attribute: Attribute,
  commit: (value: any) => void,
  small: boolean
}>()

// Local state for the slider
const state = reactive({
  latest: new Date(),
  waiting: false,
  sliding: false,
  loading: false,
  local: props.attribute,
})

// The predefined styles for each attribute key
const styles: any = {
  dim: {
    name: 'Intensity',
    min: 0,
    max: 100,
    step: 5,
    unit: '%',
    icon: '􀇯',
    class: 'slider-dim',
  },
  cct: {
    name: 'Warmth',
    min: 2000,
    max: 8000,
    step: 100,
    unit: 'K',
    icon: '􀍽',
    class: 'slider-cct',
  },
  hue: {
    name: 'Color',
    min: 1,
    max: 360,
    step: 1,
    unit: '°',
    icon: '􀟗',
    class: 'slider-hue',
  },
  main: {
    min: 1,
    max: 100,
    step: 5,
    unit: '%',
    icon: '􀌆',
    class: 'slider-dim',
  }
}

// Lock the slider, kinda like a state mutex
function slideStart(_: MouseEvent) {
  state.sliding = true
}

// Send the changes when the user lifts their finger
function commitChanges(_: MouseEvent) {
  props.commit(state.local)
  state.waiting = true
  state.latest = new Date()
  // Prevent updates until after we send the state
  state.sliding = false
}

</script>

<template>
  <div v-if="!small" class="py-1" v-on:click.stop>
    <div class="d-flex justify-content-between align-content-center align-items-center">
      <div class="h-bar justify-content-start align-items-center align-content-center">
        <div class="label-xxs label-o2 label-w600">{{ styles[attribute.key].icon }}</div>
        <div class="label-xxs label-o4 label-w500 fixed-width-name">&nbsp;&nbsp; {{ styles[attribute.key].name }}</div>
        <div class="fill"></div>
        <div class="h-bar gap label-xxs label-o3 label-w400 px-2">
          <div class="label-xxs label-o3">{{ state.local.request }} {{ styles[attribute.key].unit }}</div>

        </div>
      </div>

      <input v-model="state.local.request"
             :class="`slider-${attribute.key}`"
             :max=styles[attribute.key].max
             :min=styles[attribute.key].min
             :step=styles[attribute.key].step
             class="range-slider slider element"
             type="range"
             v-on:mousedown="slideStart"
             v-on:mouseup="commitChanges">

    </div>
  </div>
  <div v-else @mousemove.stop>
    <div class="d-flex flex-column justify-content-between align-content-start align-items-start ">

      <input v-model="state.local.request"
             :class="`slider-${attribute.key}`"
             :max=styles[attribute.key].max
             :min=styles[attribute.key].min
             :step=styles[attribute.key].step
             class=" slider element slider-small"
             type="range"
             v-on:mousedown="slideStart"
             v-on:mouseup="commitChanges">
    </div>
  </div>
</template>

<style scoped>
.range-slider {
  width: 15rem;
}

.fixed-width-name {
  width: 3.5rem;
}
</style>

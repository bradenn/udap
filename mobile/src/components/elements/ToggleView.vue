<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";

let props = defineProps<{
  title?: string
  subtitle?: string
  elements?: number
  description?: number
}>()

let state = reactive({
  toggle: 0,
  elements: 2,
})

onMounted(() => {
  if (props.elements) state.elements = props.elements
})

function onTouch(e: TouchEvent) {
  state.toggle = (state.toggle + 1) % (state.elements)
}


</script>

<template>
  <div class="" @touchstart="onTouch">
    <div class="d-flex justify-content-between align-items-center px-2">
      <div class="d-flex align-items-center gap-1">
        <div class="label-c5 label-o6 label-w600">{{ props.title }}</div>
        <div v-if="props.subtitle" class="label-c5 label-o4 label-w600">- {{ props.subtitle }}</div>
      </div>
      <div class="d-flex align-items-center gap-2">
        <slot name="description"></slot>
        <div v-if="props.description">{{ props.description }}</div>
        <div
            style="height: 0.8rem; width: 1.5px; border-radius:1px; background-color: rgba(255,255,255,0.2)"></div>
        <div class="d-flex gap-1">
          <div v-for="arr in Array(state.elements).keys()">
            <div :class="`circle ${state.toggle == arr?'circle-fill':'circle-empty'}`"></div>
          </div>
        </div>
      </div>
    </div>
    <div v-if="state.toggle == 0" class="animated">
      <slot name="primary"></slot>
    </div>
    <div v-else class="animated">
      <slot name="secondary"></slot>
    </div>

  </div>
</template>

<style lang="scss" scoped>

.animated {
  animation: fade-slow 200ms ease-in-out forwards;
}

@keyframes fade-slow {
  0% {
    opacity: 0.8;
    scale: 0.98;
    filter: blur(2px);
  }
  25% {
    opacity: 0.95;
    scale: 0.98;

  }
  100% {
    opacity: 1;
    scale: 1;
    filter: blur(0px);
  }
}

.circle {
  width: 6px;
  height: 6px;
  border-radius: 6px;
}

$circleColor: rgba(255, 255, 255, 0.6);

.circle-fill {
  background-color: $circleColor;
}

.circle-empty {
  border: 1px solid $circleColor;
}
</style>
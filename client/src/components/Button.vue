<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, reactive} from "vue";


const props = defineProps<{
  active: boolean,
  text: string,
}>()

let state = reactive({})

const haptics = inject("haptic") as (a: number, b: number, c: number) => void

function push(_: Event) {
  haptics(1, 1, 75)
}

function release(_: Event) {
  // haptics(1, 1, 75)
  haptics(1, 1, 50)
}


</script>

<template>
  <div :class="`${props.active?'':'subplot-inline'}`" class="switch-space subplot switch-space-active" @mousedown="push"
       @mouseup="release">
    {{ props.text }}
  </div>
</template>

<style lang="scss" scoped>

.switch-space-active {
  color: rgba(255, 255, 255, 0.5) !important;

  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);
}

.switch-right {
  animation: switch-move-right 150ms forwards ease-out;
}

.switch-space:active {
  transform: scale3d(0.96, 0.96, 1);
  background-color: rgba(255, 255, 255, 0.025) !important;
}

.switch-space:active {
  //animation: switch-active 125ms forwards linear;
}

@keyframes switch-active {
  0% {

  }

  100% {
    transform: scale3d(1, 1, 1);

  }
}

.switch-left {
  animation: switch-move-left 150ms forwards ease;
}

@keyframes switch-move-left {
  0% {
    transform: translateX(16px);
  }
  100% {
    transform: translateX(0px);
  }
}

.switch-space {
  display: flex;
  justify-content: center;
  font-weight: 600;
  align-items: center;
  font-size: 0.5rem;
  height: 1.8rem;
  color: rgba(255, 255, 255, 0.4);
}

.switch-spots {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(2, 1fr);
}

.switch-path {
  position: absolute;
  width: calc(100% - 0.5rem);
  left: 8px;

  height: 1.8rem;
  //background-color: #0a58ca;
}

.shuttle-text {
  position: absolute;
  left: 8px;
}

.shuttle-cock {

  height: 2rem;
  //width: 100px;
  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);

}


.arrow-down {

}

.shuttle {
  width: 3px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.25);
  border-radius: 3px;
  position: relative;

  top: -28px
}

.shuttle-center {
  width: 3px;
  height: 8px;
  background-color: rgba(255, 255, 255, 0.3);
  border-radius: 3px;
  position: absolute;


}

.slider {
}
</style>
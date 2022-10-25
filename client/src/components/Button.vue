<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, reactive} from "vue";
import {useRouter} from "vue-router";


const props = defineProps<{
  active: boolean,
  text: string,
  to?: string,
}>()

let state = reactive({})
let router = useRouter()

const haptics = inject("haptic") as (a: number, b: number, c: number) => void

function push(_: Event) {
  haptics(2, 1, 75)
  if (props.to) {
    router.push(props.to)
  }
}

function release(_: Event) {
  // haptics(1, 1, 75)
  // haptics(1, 1, 5)
}


</script>

<template>
  <div :class="`${props.active?'':'subplot-inline'}`" class="switch-space subplot switch-space-active"
       @mousedown="push"
       @mouseup="release">
    <span :class="`${props.active?'':''}`">{{ props.text }}</span>
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
}

.switch-space {
  display: flex;
  justify-content: center;
  font-weight: 600;
  align-items: center;
  font-size: 0.65rem;
  line-height: 0.7rem;
  height: 1.8rem;
  color: rgba(255, 255, 255, 0.4);
}
</style>
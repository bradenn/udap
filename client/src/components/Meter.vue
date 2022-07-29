<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive, watch} from "vue";

interface Meter {
  ticks: number
  vertical: boolean
  percent: number
  title?: string
  alt?: string
  blank?: boolean
  thresh?: number
  danger?: number
}

let props = defineProps<Meter>()

let state = reactive({
  activeTicks: [] as string[],
  inactiveTicks: [] as string[],
})

onMounted(() => {
  calculateTicks()
})

watch(props, (o, n) => {
  calculateTicks()
})

function calculateTicks() {
  state.activeTicks = Array(Math.ceil((props.percent / 100) * props.ticks)) as string[]
  state.inactiveTicks = Array(Math.floor(props.ticks - (props.percent / 100) * props.ticks)) as string[]
  for (let i = 0; i < state.activeTicks.length; i++) {
    let p = (i / props.ticks) * 100
    state.activeTicks[i] = meterColor(p, props.thresh || 100, props.danger || 100)
  }
  for (let i = 0; i < state.inactiveTicks.length; i++) {
    let pI = (((!props.vertical ? state.activeTicks.length + i : props.ticks - i)) / props.ticks) * 100
    state.inactiveTicks[i] = meterColorUnmet(pI, props.thresh || 100, props.danger || 100)
  }
}

function meterColorUnmet(percent: number, l1: number, l2: number): string {
  if (percent >= l1 && percent < l2) {
    return "rgba(255,202,128,0.05)"
  } else if (percent >= l2) {
    return "rgba(255,158,153,0.05)"
  }
  return "rgba(133,188,255,0.05)"
}


function meterColor(percent: number, l1: number, l2: number): string {
  if (percent >= l1 && percent < l2) {
    return "rgba(255, 149, 0, 1)"
  } else if (percent >= l2) {
    return "rgba(255, 59, 48, 1)"
  }
  return "rgba(10, 122, 255, 1)"
}


</script>

<template>

  <div v-if="vertical">
    <div v-if="blank" class="subplot utility-stack-blank-v h-100">
    </div>
    <div v-else class="subplot utility-stack-v">
      <div v-for="b in state.inactiveTicks"
           :style="props.thresh?`background-color: ${b} !important;`:''"
           class="utility-bar-v"></div>
      <div v-for="b in state.activeTicks"
           :style="props.thresh?`background-color: ${b} !important;`:''"
           class="utility-bar-v utility-bar-fill"></div>
    </div>
  </div>
  <div v-else class="d-flex flex-column">
    <div v-if="props.title" class="d-flex justify-content-between w-100 px-1" style="height: 0.7rem">
      <div class="label-c3 label-o3 label-w600">{{ props.title }}</div>
      <div class="label-c3 label-o3 label-w600">{{ props.alt }}</div>
    </div>
    <div v-if="blank" class="subplot utility-stack-blank-v h-100">
    </div>
    <div class="subplot utility-stack-h h-100">
      <div v-for="b in state.activeTicks"
           :style="props.thresh?`background-color: ${b} !important;`:''"
           class="utility-bar-fill utility-bar-h">
      </div>
      <div v-for="b in state.inactiveTicks" :style="props.thresh?`background-color: ${b} !important;`:''"
           class="utility-bar-h"></div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.utility-bar-fill {
  background-color: #0075cb !important;
  box-shadow: inset 0 0.5px rgba(255, 255, 255, 0.04), 0 0.5px rgba(0, 0, 0, 0.3), 0 0 1px 1px rgba(1, 159, 255, 0.1) !important;
}

.utility-stack-v {
  width: 1.5rem;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.1);
  display: grid;
  grid-row-gap: 2px;
  grid-auto-flow: column;
  grid-template-rows: repeat(v-bind(ticks), 1fr);
  grid-template-columns: repeat(1, 1fr);
  padding: 0.25rem 0.25rem;
}

.utility-stack-blank-v {
  width: 1.5rem;
  background-color: rgba(0, 0, 0, 0.1);
  opacity: 0.8;
  display: grid;
  grid-row-gap: 2px;
  grid-auto-flow: column;
  grid-template-rows: repeat(v-bind(ticks), 1fr);
  grid-template-columns: repeat(1, 1fr);
  padding: 0.25rem 0.25rem;
}

.utility-bar-v {
  content: '';
  border-radius: 3px;
  width: 100%;
  height: 5px;
  box-shadow: inset 0 1px rgba(255, 255, 255, 0.025), 0 0.5px rgba(0, 0, 0, 0.3);
  background-color: rgba(255, 255, 255, 0.05);
}


.utility-stack-h {
  width: 100%;
  height: 100%;
  align-items: center;
  background-color: rgba(0, 0, 0, 0.1);
  display: grid;
  grid-column-gap: 2px;
  grid-auto-flow: column;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(v-bind(ticks), 1fr);
  padding: 0.25rem 0.25rem;
}

.utility-stack-blank-h {
  width: 100%;
  background-color: rgba(0, 0, 0, 0.1);
  opacity: 0.8;
  display: grid;
  grid-row-gap: 2px;
  grid-auto-flow: column;
  grid-template-rows: repeat(v-bind(ticks), 1fr);
  grid-template-columns: repeat(1, 1fr);
  padding: 0.25rem 0.25rem;
}

.utility-bar-h {
  content: '';
  border-radius: 3px;
  width: 6px;
  height: 100%;
  box-shadow: inset 0 1px rgba(255, 255, 255, 0.025), 0 0.5px rgba(0, 0, 0, 0.3);
  background-color: rgba(255, 255, 255, 0.05);
}
</style>
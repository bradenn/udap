<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {onMounted, onUnmounted, reactive, watchEffect} from "vue";
import {PreferencesRemote} from "../persistent";
import core from "../core";
import {TimeUnit} from "../time";

const props = defineProps<{
  until?: number,
  since?: number,
  live?: boolean,
  nano?: boolean
  speed?: number,
  precise?: boolean
}>()


const state = reactive({
  display: [] as TimeUnit[],
  value: 0,
  last: 0,
  target: 0,
  timeout: 0,
  toggle: false,
  animation: true,
  frame: 0,
  tick: 0,
  tickLimit: 10
})


onMounted(() => {
  if (props.speed) {
    state.tickLimit = 60;
  }
  if (props.live) {
    setTimeout(() => {
      //@ts-ignore
      state.frame = setInterval(animate, 1000 / state.tickLimit), 1000 - new Date().getMilliseconds()
    })
  }
})

const preferences = core.preferences() as PreferencesRemote

onUnmounted(() => {
  clearInterval(state.frame)
  state.animation = false
})

function animate() {

  if (state.animation) {
    // state.frame = setTimeout(animate, 1000 - new Date().getMilliseconds())
    state.tick = (state.tick + 1) % state.tickLimit
    if (state.tick == 0) {

    }
    if (props.since) {
      state.target = (Date.now().valueOf() - props.since)
    }
    state.last = state.value
    //map_range(state.tick, 0, state.tickLimit, state.last, state.target)
    state.value = state.target
    // console.log(state.value)
    state.display = convertNanosecondsDurationToString(state.value)
  }

}

watchEffect(() => {
  if (props.until) {
    state.target = props.until;
    updateDisplay()
  }
  return props.until
})

watchEffect(() => {
  if (props.since && !props.live) {
    state.target = props.since;
    updateDisplay()
  }
  return props.since
})


function updateDisplay() {
  let delay = 1000 - new Date().getMilliseconds()
  state.value = state.target
  if (state.timeout == 0) {
    //@ts-ignore
    state.timeout = setTimeout(() => {
      state.display = convertNanosecondsDurationToString(state.value)
      // state.toggle = !state.toggle
      state.timeout = 0
    }, delay) as number
  }

}


function convertNanosecondsDurationToString(nanoseconds: number): TimeUnit[] {
  if (nanoseconds < 0) {
    nanoseconds = Math.abs(nanoseconds)
    // throw new Error("Input must be a non-negative number.");
  }
  if (!props.nano) {
    let value = nanoseconds;
    let miliseconds = value % 1000
    let seconds = value / 1000 % 60
    let minutes = (value / 1000 / 60) % 60
    let hours = (value / 1000 / 60 / 60) % 60
    let days = Math.round((value / 1000 / 60 / 60 / 24) % 24)
    let years = Math.round((value / 1000 / 60 / 60 / 24 / 365))
    let timeUnits = [] as TimeUnit[]
    if (years > 0) {
      timeUnits.push({
        value: Math.round(years),
        units: "y"
      } as TimeUnit)
    }
    if (years > 0) {
      timeUnits.push({
        value: Math.round(days),
        units: "d"
      } as TimeUnit)
    }
    if (years > 0) {
      timeUnits.push({
        value: Math.round(hours),
        units: "h"
      } as TimeUnit)
    }
    if (years > 0) {
      timeUnits.push({
        value: Math.round(minutes),
        units: "m"
      } as TimeUnit)
    }
    timeUnits.push({
      value: Math.round(seconds),
      units: "s"
    } as TimeUnit)
    if (props.precise) {
      timeUnits.push({
        value: Math.round(miliseconds),
        units: "ms"
      } as TimeUnit)
    }

    return timeUnits;
  }

  if (nanoseconds === 0) {
    return [{value: 0, units: "ns"} as TimeUnit];
  }

  const units = ["ns", "µs", "ms", "s", "ks"];
  let value = nanoseconds;
  let unitIndex = 0;

  while (value >= 1000 && unitIndex < units.length - 1) {
    if (unitIndex >= 3) {
      let seconds = value % 60
      let minutes = (value / 60) % 60
      let hours = (value / 60 / 60) % 60
      let days = Math.round((value / 60 / 60 / 24) % 24)
      let years = Math.round((value / 60 / 60 / 24 / 365))
      let timeUnits = [] as TimeUnit[]
      if (years > 0) {
        timeUnits.push({
          value: Math.round(years),
          units: "y"
        } as TimeUnit)
      }
      if (years > 0) {
        timeUnits.push({
          value: Math.round(days),
          units: "d"
        } as TimeUnit)
      }
      if (years > 0) {
        timeUnits.push({
          value: Math.round(hours),
          units: "h"
        } as TimeUnit)
      }
      timeUnits.push({
        value: Math.round(minutes),
        units: "m"
      } as TimeUnit)
      timeUnits.push({
        value: Math.round(seconds),
        units: "s"
      } as TimeUnit)

      return timeUnits;

    }
    value /= 1000;
    unitIndex++;
  }
  return [{value: value, units: units[unitIndex]} as TimeUnit]

}

function pad(num: number, size: number): string {
  let s = num + "";
  while (s.length < size) s = "<span class='label-o4'>0</span>" + s;
  return s;
}
</script>

<template>
  <div :class="`${state.toggle?'time-frame':'time-frame-i'}`" class="mono">
    <div class="d-flex gap-2">
      <div v-for="t in state.display" class="d-flex " style="gap: 0.125rem">

        <div class="label-o5" v-html="pad(t.value, t.units == 'ms'?3:2)"></div>
        <div class="label-o3">{{ t.units }}</div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.time-frame {
  animation: animateIn 0.25s forwards;
}

.time-frame-i {
  animation: animateIn 0.25s forwards;
}

.gradient-text {
  //background-color: #f3ec78;
  background-image: linear-gradient(0deg, rgba(255, 255, 255, 0.55), rgba(255, 255, 255, 0.66), rgba(255, 255, 255, 0.65), rgba(255, 255, 255, 0.6));
  background-size: 100%;
  -webkit-background-clip: text;
  -moz-background-clip: text;
  -webkit-text-fill-color: transparent;
  -moz-text-fill-color: transparent;
}

@keyframes animateIn {
  0% {
    transform: translateY(-1px);

  }

  100% {

    transform: translateY(-0px);
  }
}
</style>
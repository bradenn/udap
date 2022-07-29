<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import type {Weather} from "@/weather";
import {onMounted, reactive, watchEffect} from "vue";


interface DaytimeProps {
  latest: Weather
}

let props = defineProps<DaytimeProps>()

let state = reactive({
  maxTemp: 0,
  minTemp: 0,
  deltas: {
    begin: Array(16) as number[],
    now: Array(16) as number[],
    end: Array(16) as number[],
    prefix: Array(16) as number[],
    suffix: Array(16) as number[],

  }
})

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

onMounted(() => {
  getDifference()
})

watchEffect(() => {
  return getDifference()
})

function getDifference() {
  let d = new Date();
  for (let i = 0; i < props.latest.hourly.temperature_2m.length; i++) {
    if (props.latest.hourly.temperature_2m[i] > state.maxTemp) {
      state.maxTemp = props.latest.hourly.temperature_2m[i]
    } else if (props.latest.hourly.temperature_2m[i] < state.minTemp) {
      state.minTemp = props.latest.hourly.temperature_2m[i]

    }
  }
  return props.latest
}


</script>
<template>
  <div v-if="props.latest.daily">

    <div class="day">
      <div v-for="a in props.latest.hourly.temperature_2m"
           :style="`height: ${map_range(a, state.minTemp, state.maxTemp, 0, 100)}%;`" class="hour"></div>
    </div>

    <br>


  </div>
</template>

<style scoped>


.day {
  display: flex;
  flex-direction: row;
  align-content: end;
  align-items: end;
  height: 3rem;
  width: 100%;
  gap: 2px;
}

.hour {
  content: ' ';
  width: 5px;
  border-radius: 5px;
  background-color: white;
}

.timeline-stick {
  height: 4rem;
  width: 8px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.hourly {

}

.timestick {
  content: ' ';
  width: 8px;
  height: 1px;
  border-radius: 2px;
  background-color: white;
}
</style>
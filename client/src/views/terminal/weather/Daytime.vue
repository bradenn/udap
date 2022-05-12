<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import type {Weather} from "@/weather";
import {onMounted, reactive, watchEffect} from "vue";


interface DaytimeProps {
  latest: Weather
}

let props = defineProps<DaytimeProps>()

let state = reactive({
  deltas: {
    begin: Array(16) as number[],
    now: Array(16) as number[],
    end: Array(16) as number[],
    prefix: Array(16) as number[],
    suffix: Array(16) as number[],

  }
})

onMounted(() => {
  getDifference()
})

watchEffect(() => {
  return getDifference()
})

function getDifference() {
  let d = new Date();
  let now = d.valueOf() - d.getHours() * 60 * 60 * 1000 - d.getMinutes() * 60 * 1000 - d.getSeconds() * 1000 - d.getMilliseconds()
  for (let i = 0; i < props.latest.daily.sunrise.length; i++) {
    let sunrise = (props.latest.daily.sunrise[i] * 1000)
    let sunset = (props.latest.daily.sunset[i] * 1000)
    state.deltas.begin[i] = ((props.latest.daily.sunrise[i] * 1000) - (now + 1000 * 60 * 60 * 24 * i)) / (1000 * 60 * 60 * 24)
    state.deltas.end[i] = (now + 1000 * 60 * 60 * 24 * (i + 1) - (props.latest.daily.sunset[i] * 1000)) / (1000 * 60 * 60 * 24)

    // = -props.latest.daily.sunrise[i]
  }
  return state.deltas
}

</script>
<template>
  <div v-if="props.latest.daily">

    <div class="d-flex gap-1 flex-row">

      <div v-for="a in Array(props.latest.daily.sunrise.length).keys()">
        state.deltas.begin[a]
        state.deltas.end[a]
      </div>

    </div>

    <br>


  </div>
</template>

<style scoped>
.timeline-stick {
  height: 4rem;
  width: 8px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.timestick {
  content: ' ';
  width: 8px;
  height: 1px;
  border-radius: 2px;
  background-color: white;
}
</style>
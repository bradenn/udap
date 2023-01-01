<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";
import {reactive} from "vue";
import moment from "moment";

let state = reactive({
  timer: {
    duration: 0,
    current: 0,
    started: 0
  },
  ctx: 0
})

function reset() {
  state.timer.duration = 0;
  state.timer.current = 0;
  state.timer.started = 0;
  clearInterval(state.ctx)
  state.ctx = 0
}

function setDuration(duration: number) {
  reset()
  state.timer.duration = duration
  start()
}

function start() {
  if (state.ctx != 0) return
  state.timer.started = new Date().valueOf()
  state.ctx = setInterval(updateTime, 100)
}


function updateTime() {
  state.timer.current = state.timer.duration - (new Date().valueOf() - state.timer.started)
}

</script>

<template>
  <div class=" h-100">
    <div class="element">
      <div class="d-flex">
        <div>{{ moment(state.timer.current).format('m') }}</div>
        <div>:</div>
        <div>{{ moment(state.timer.current).format('s') }}</div>
        <div>:</div>
        <div>{{ moment(state.timer.current).format('S') }}</div>
      </div>
    </div>
    <div class="d-flex gap mt-1">
      <Plot :cols="2" :rows="2" style="width: 12rem">
        <Radio :active="false" :fn="() => setDuration(1000 * 60)"
               title="1 min"></Radio>
        <Radio :active="false" :fn="() => setDuration(5 * 1000 * 60)"
               title="5 min"></Radio>
        <Radio :active="false" :fn="() => setDuration(15 * 1000 * 60)"
               title="15 min"></Radio>
        <Radio :active="false" :fn="() => setDuration(30 * 1000 * 60)"
               title="30 min"></Radio>
      </Plot>
      <Plot :cols="2" :rows="2" name="Drinks" style="width: 12rem">
        <Radio :active="false" :fn="() => setDuration(4*1000 * 60)"
               title="Earl Grey, Hot"></Radio>
      </Plot>
    </div>

  </div>
</template>

<style scoped>

</style>
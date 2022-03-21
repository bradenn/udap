<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive} from "vue";


let state = reactive({
  running: false,
  start: 0,
  stop: 0,
  delta: {
    milliseconds: "",
    seconds: "",
    minutes: "",
    hours: "",
    days: "",
    years: "",
    decades: "",
  },
  stops: [],
  time: '',
  interval: 0,
})

onMounted(() => {
  reset()
})


function reset() {
  stop()
  state.start = 0
  state.running = false
  tick()
}


function start() {

  if (state.start !== 0) {
    state.running = true
    state.start = new Date().valueOf() - (state.stop - state.start)
    state.interval = setInterval(tick, 3)
  } else {
    state.running = true
    state.start = new Date().valueOf()
    state.interval = setInterval(tick, 3)
  }
}


function stop() {
  state.running = false
  state.stop = new Date().valueOf()
  clearInterval(state.interval)
}

function tick() {

  let now = state.start > 0 ? new Date().valueOf() - state.start : 0


  state.delta.milliseconds = (Math.floor(now) % 1000).toString().padStart(3, "0")
  state.delta.seconds = (Math.floor(now / 1000) % 60).toString().padStart(2, "0")
  state.delta.minutes = (Math.floor(now / 1000 / 60) % 60).toString().padStart(2, "0")
  state.delta.hours = (Math.floor(now / 1000 / 60 / 60) % 24).toString().padStart(2, "0")
  state.delta.days = (Math.floor(now / 1000 / 60 / 60 / 24) % 365).toString().padStart(3, "0")
  state.delta.years = (Math.floor(now / 1000 / 60 / 60 / 24) % 10).toString().padStart(4, "0")

}

</script>

<template>
  <div class="d-flex flex-column justify-content-start align-items-center h-100">
    <div class="d-inline-block element px-5 py-3 my-3">
      <div class="d-flex flex-row align-items-center justify-content-center gap-4">
        <div class="d-flex flex-column align-items-end">
          <div class="timer">{{ state.delta.days }}</div>
          <div class="label-sm label-w400 label-o2">days</div>
        </div>

        <div class="spacer"></div>

        <div class="d-flex flex-column align-items-end">
          <div class="timer">{{ state.delta.hours }}</div>
          <div class="label-sm label-w400 label-o2">hours</div>
        </div>

        <div class="spacer"></div>

        <div class="d-flex flex-column align-items-end">
          <div class="timer">{{ state.delta.minutes }}</div>
          <div class="label-sm label-w400 label-o2">minutes</div>
        </div>

        <div class="spacer"></div>

        <div class="d-flex flex-column align-items-end">
          <div class="timer">{{ state.delta.seconds }}</div>
          <div class="label-sm label-w400 label-o2">seconds</div>
        </div>


        <div class="spacer"></div>

        <div class="d-flex flex-column align-items-end">
          <div class="timer">{{ state.delta.milliseconds }}</div>
          <div class="label-sm label-w400 label-o2">milliseconds</div>
        </div>

      </div>

    </div>

    <div class="element surface d-flex flex-row gap-2 p-2">
      <div v-if="!state.running" class="element px-3 button-xl element-green" @click="start">
        <div class="fa-solid fa-play fa-fw fa-2x element-green"></div>
        <div>Start</div>
      </div>
      <div v-else class="element px-3 button-xl element-red" @click="stop">
        <div class="fa-solid fa-stop fa-fw fa-2x element-green"></div>
        <div>Stop</div>
      </div>
      <div class="element px-3 button-lg" @click="reset">
        <div class="fa-solid fa-backward-step fa-fw fa-2x"></div>
        <div>Reset</div>
      </div>

    </div>
  </div>
</template>

<style scoped>
.button-lg {
  width: 10rem;
  aspect-ratio: 3/3;
  align-items: center;
  font-weight: 500;
  border-radius: 0.4rem;
  display: flex;
  color: rgba(255, 255, 255, 0.8);
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
  flex-direction: column;
  justify-content: center;
  gap: 1.5rem;
  text-transform: uppercase;
  box-shadow: 0 0 16px 1px rgba(0, 0, 0, 0.025);
}

.button-xl {
  width: 20rem;
  aspect-ratio: 6/3;
  align-items: center;
  border-radius: 0.4rem;
  font-weight: 500;
  display: flex;
  color: rgba(255, 255, 255, 0.8);
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.2);
  flex-direction: column;
  justify-content: center;
  gap: 1.5rem;
  text-transform: uppercase;
  box-shadow: 0 0 16px 1px rgba(0, 0, 0, 0.025);
}

.timer {
  font-family: "Roboto", sans-serif;
  letter-spacing: 1px !important;
  font-variant-numeric: tabular-nums lining-nums;
  font-weight: 350 !important;
  font-size: 4.5rem;
  line-height: 5rem;
  transition: all 1000ms ease-in;
  text-shadow: 0 0 8px rgba(0, 0, 0, 0.05);
}

.spacer {
  width: 0.125rem;
  height: 3rem;
  border-radius: 1rem;
  background-color: rgba(255, 255, 255, 0.125);
  box-shadow: 0 0 8px 2px rgba(0, 0, 0, 0.05);
  margin-bottom: 1rem;
  line-height: 5rem;
}
</style>
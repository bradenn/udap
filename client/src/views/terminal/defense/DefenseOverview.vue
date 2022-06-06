<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Keypad from "@/components/Keypad.vue";
import Plot from "@/components/plot/Plot.vue";
import {onMounted, reactive} from "vue";
import axios from "axios";
import Subplot from "@/components/plot/Subplot.vue";
import Confirm from "@/components/plot/Confirm.vue";

let state = reactive({
  laser: false,
  pan: 90,
  tilt: 180,
  runner: 0,
})

onMounted(() => {
  query()
})

function query() {
  axios.get(`http://10.0.1.60/query`).then(res => {
    state.laser = (res.data.laser === 1)
    state.pan = res.data.pan
    state.tilt = res.data.tilt
  }).catch(res => {
    console.log(res)
  })
}

function laserToggle() {
  axios.get(`http://10.0.1.60/laser/${!state.laser ? 1 : 0}`).then(res => {
    state.laser = res.data.value === 1
  }).catch(res => {
    console.log(res)
  })
}

function laserTilt(value: number) {
  axios.get(`http://10.0.1.60/tilt/${value}`).then(res => {
    state.tilt = res.data.value
  }).catch(res => {
    console.log(res)
  })
}

function laserPan(value: number) {
  axios.get(`http://10.0.1.60/pan/${value}`).then(res => {
    state.pan = res.data.value
  }).catch(res => {
    console.log(res)
  })
}

function laserHome() {
  laserPan(90)

  laserTilt(90)
}

function laserDoor() {
  laserPan(90)
  laserTilt(90)
}

function laserWall() {
  laserPan(105)
  laserTilt(140)
}

function map_range(value, low1, high1, low2, high2) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function laserRun() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  state.runner = setInterval(() => {
    let a1 = map_range(Math.floor(tick), 0, 100, 90, 105)
    let b1 = map_range(Math.floor(tick), 0, 100, 90, 135)

    tick += dir ? -0.125 : 0.125;

    if (tick >= 100) {
      dir = true;
    } else if (tick <= 0) {
      dir = false;
    }

    laserPan(Math.cos(Math.floor(tick)) * map_range(Math.floor((2 * Math.PI / 100) * tick), 0, 100, 15, 1) + a1)
    laserTilt(Math.sin(Math.floor(tick)) + b1)

  }, 50)

}

function laserCircle() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  state.runner = setInterval(() => {
    tick += dir ? -1 : 1;

    if (tick >= 100) {
      dir = true;
    } else if (tick <= 0) {
      dir = false;
    }

    laserPan(90 + Math.cos((2 * Math.PI / 100) * tick) * 90)
    laserTilt(90 + Math.sin((2 * Math.PI / 100) * tick) * 45)

  }, 50)

}

function laserStop() {
  clearInterval(state.runner)
  state.runner = 0
}

</script>

<template>
  <div class="h-100 d-flex gap">
    <Keypad></Keypad>
    <div class="d-flex gap justify-content-center">
      <div>LASER {{ state.laser }}</div>
      <div>PAN {{ state.pan }}</div>
      <div>TILT {{ state.tilt }}</div>
    </div>

    <Plot :cols="1" :rows="2" style="height: 8rem; width: 13rem">
      <Confirm v-if="!state.laser" :active="state.laser" :fn="laserToggle"
               :title="`LASER ${state.laser?'OFF':'ON'}`"></Confirm>
      <Subplot v-else :active="true" :fn="laserToggle" :name="`LASER ${state.laser?'OFF':'ON'}`"></Subplot>
    </Plot>

    <Plot :cols="2" :rows="4" style="height: 8rem; width: 13rem">
      <Subplot :active="true" :fn="laserHome" name="HOME"></Subplot>
      <div></div>
      <Subplot :active="true" :fn="() => laserPan(state.pan-5)" name="LEFT"></Subplot>
      <Subplot :active="true" :fn="() => laserPan(state.pan+5)" name="RIGHT"></Subplot>
      <Subplot :active="true" :fn="() => laserTilt(state.tilt-5)" name="DOWN"></Subplot>
      <Subplot :active="true" :fn="() => laserTilt(state.tilt+5)" name="UP"></Subplot>
    </Plot>

    <Plot :cols="2" :rows="3" style="width: 13rem; height: 8rem;">
      <Subplot :active="true" :fn="() => laserDoor()" name="DOOR"></Subplot>
      <Subplot :active="true" :fn="() => laserWall()" name="BED"></Subplot>
      <Subplot :active="true" :fn="() => laserRun()" name="RUN"></Subplot>
      <Subplot :active="true" :fn="() => laserCircle()" name="CIRCLE"></Subplot>
      <Subplot :active="true" :fn="() => laserStop()" name="STOP"></Subplot>
    </Plot>

  </div>
</template>

<style lang="scss" scoped>
</style>
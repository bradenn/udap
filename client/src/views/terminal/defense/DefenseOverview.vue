<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>


import Plot from "@/components/plot/Plot.vue";
import {inject, onMounted, onUnmounted, reactive, watch} from "vue";
import axios from "axios";
import Subplot from "@/components/plot/Subplot.vue";
import Confirm from "@/components/plot/Confirm.vue";
import DefenseAuth from "@/views/terminal/defense/DefenseAuth.vue";
import type {Session} from "@/types";

let state = reactive({
  laser: false,
  pan: 90,
  tilt: 180,
  runner: 0,
  speed: 1,
  auth: false,
})

let session = inject("session") as Session

onMounted(() => {
  query()
  verifyAuth(session)
})

onUnmounted(() => {
  laserStop()
})

watch(session, (current: Session, previous: Session) => {
  verifyAuth(current)
})

function verifyAuth(current: Session) {
  state.auth = (!!current.user.id)
}

function query() {
  axios.get(`http://10.0.1.60/query`).then(res => {
    state.laser = (res.data.laser === 1)
    state.pan = map_range(res.data.pan, 0, 1800, 0, 180)
    state.tilt = map_range(res.data.tilt, 0, 1800, 0, 180)
  }).catch(res => {
    console.log(res)
  })
}

function laserPower(on: boolean) {
  axios.get(`http://10.0.1.60/laser/${on ? 1 : 0}`).then(res => {
    state.laser = res.data.value === 1
  }).catch(res => {
    console.log(res)
  })
}

function laserToggle() {
  laserPower(!state.laser)
}

function laserTilt(value: number) {
  let tiltA = map_range(value, 0, 180, 0, 1800)
  axios.get(`http://10.0.1.60/tilt/${tiltA}`).then(res => {
    state.tilt = map_range(res.data.value, 0, 1800, 0, 180)
  }).catch(res => {
    console.log(res)
  })
}

function laserPan(value: number) {
  let panA = map_range(value, 0, 180, 0, 1800)
  axios.get(`http://10.0.1.60/pan/${panA}`).then(res => {
    state.pan = map_range(res.data.value, 0, 1800, 0, 180)
  }).catch(res => {
    console.log(res)
  })
}

function laserHome() {
  laserPanTilt(90, 90)
}

function laserWall() {
  laserPanTilt(105, 154)
}

function laserPanTilt(pan: number, tilt: number) {
  let panA = map_range(pan, 0, 180, 0, 1800)
  let tiltA = map_range(tilt, 0, 180, 0, 1800)
  axios.get(`http://10.0.1.60/pan/${panA}/tilt/${tiltA}`).then(res => {
    state.pan = map_range(res.data.pan, 0, 1800, 0, 180)
    state.tilt = map_range(res.data.tilt, 0, 1800, 0, 180)
  }).catch(res => {
    console.log(res)
  })
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function laserSpeed(speed: number) {
  state.speed = speed
}

function laserStopAll() {
  clearInterval(state.runner)
  state.runner = 0
  state.speed = 1
  laserHome()
  laserPower(false)
}

function getSpeed(): number {
  return state.speed
}

function laserRun() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  state.runner = setInterval(() => {
    let b1 = map_range(tick, 0, 1000, 90, 145)
    tick += dir ? -state.speed : state.speed;
    if (Math.floor(tick) % 250 == 0) {
      state.speed = 1 + Math.random() * 10
    }
    if (Math.floor(tick) >= 1000) {
      dir = true;
    } else if (Math.floor(tick) <= 0) {
      dir = false;
    }
    //
    // laserPan(Math.cos(Math.floor(tick)) * map_range(Math.floor((2 * Math.PI / 100) * tick), 0, 100, 15, 1) + a1)
    // laserTilt(Math.sin(Math.floor(tick)) + b1)
    if (state.speed > 5) {
      laserPanTilt(95 + Math.sin((2 * Math.PI / (250 - state.speed * 4) * tick)) * map_range(tick, 0, 1000, 20, 5), b1)
    } else {
      laserPanTilt(95 + Math.cos((2 * Math.PI / (250 - state.speed * 4) * tick)) * map_range(tick, 0, 1000, 20, 5), b1)
    }


  }, 65)

}

function laserCircle() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  state.runner = setInterval(() => {
    tick += 1;

    if (tick >= 1000) {
      tick = 0;
    } else if (tick <= 0) {
      tick = 1000;
    }

    // 0 - 180
    let panTo = map_range(Math.sin((2 * Math.PI / 1000) * tick), -1, 1, 0, 180)
    // 90 - 180
    let tiltTo = map_range(Math.cos((2 * Math.PI / 1000) * tick), -1, 1, 95, 110)

    laserPanTilt(panTo, tiltTo)

  }, 80)

}

function laserStop() {
  clearInterval(state.runner)
  state.runner = 0
}

</script>

<template>
  <div v-if="state.auth" class="d-flex gap flex-wrap mt-4">

    <Plot :cols="2" :rows="2" style="width: 13rem" title="Sentry">
      <Confirm :active="state.laser" :disabled="state.laser"
               :fn="laserToggle" :title="`${state.laser?'DISABLE':'ENABLE'} LASER`"></Confirm>

      <Subplot :active="true" :fn="laserStopAll" name="STOP ALL" theme="danger"></Subplot>
    </Plot>

    <Plot :cols="2" :rows="4" style="width: 13rem" title="Control">
      <Subplot :active="true" :fn="laserHome" name="HOME"></Subplot>
      <div></div>
      <Subplot :active="true" :fn="() => laserPan(state.pan-1)" name="LEFT"></Subplot>
      <Subplot :active="true" :fn="() => laserPan(state.pan+1)" name="RIGHT"></Subplot>
      <Subplot :active="true" :fn="() => laserTilt(state.tilt-1)" name="DOWN"></Subplot>
      <Subplot :active="true" :fn="() => laserTilt(state.tilt+1)" name="UP"></Subplot>
    </Plot>

    <Plot :cols="2" :rows="3" style="width: 13rem;" title="Programmed">
      <Subplot :active="true" :fn="() => laserWall()" name="Bed"></Subplot>
      <Subplot :active="true" :fn="() => laserRun()" name="Run"></Subplot>
      <Subplot :active="true" :fn="() => laserCircle()" name="Circle"></Subplot>
      <Subplot :active="true" :fn="() => laserStop()" name="STOP"></Subplot>
    </Plot>
    <Plot :cols="1" :rows="2" style="width: 13rem;" title="Programmed">
      <div>
        <div class="d-flex justify-content-between label-xs label-r px-1">
          <div class="label-w500">Pan (X)</div>
          <div class="label-w600 label-o3">{{ state.pan }}°</div>
        </div>
        <input
            id="pan"
            v-model="state.pan"
            :max="180"
            :min="0"
            :step="1"
            class="slider element "
            type="range"
            v-on:mouseup="() => laserPan(state.pan)">
      </div>

      <div>
        <div class="d-flex justify-content-between label-xs label-r px-1">
          <div class="label-w500">Tilt (Y)</div>
          <div class="label-w600 label-o3">{{ state.tilt }}°</div>
        </div>
        <input
            id="tilt"
            v-model="state.tilt"
            :max="180"
            :min="0"
            :step="1"
            class="slider element"
            type="range"
            v-on:mouseup="() => laserTilt(state.tilt)">
      </div>

    </Plot>

  </div>
  <div v-else>
    <DefenseAuth></DefenseAuth>
  </div>
</template>

<style lang="scss" scoped>
.emergency-stop {

}
</style>
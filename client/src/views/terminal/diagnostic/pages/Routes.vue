<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Plot from "@/components/plot/Plot.vue";
import {onMounted, reactive} from "vue";
import PanTilt from "@/components/PanTilt.vue";

interface Position {
  a: number,
  b: number,
  time: number
}

function setOrigin() {
  state.ws.send(JSON.stringify({mode: "origin", a: 0, b: 0}))
  state.pos.a = 0
  state.pos.b = 0
  home();
}

function sendCommand(a: number, b: number, absolute: boolean) {
  state.ws.send(JSON.stringify({
    mode: absolute ? "absolute" : "relative",
    a: a,
    b: b
  }))
}

let state = reactive({
  speed: 64,
  pan: 0,
  tilt: 0,
  pos: {a: 0, b: 0} as Position,
  messages: [] as Position[],
  ws: {} as WebSocket,
  connected: false,
  interval: 0,
})


onMounted(() => {
  state.ws = new WebSocket("ws://10.0.1.85/ws")
  state.ws.onopen = (e: Event) => {
    state.connected = true
  }
  state.ws.onmessage = (e: MessageEvent) => {
    state.connected = true
    let data = JSON.parse(e.data) as Position
    data.time = new Date().valueOf()
    state.pos = data
    state.messages.push(data)
    state.messages.sort((a: Position, b: Position) => {
      return a.time - b.time
    })
  }
  state.ws.onclose = (e: Event) => {
    state.connected = false
  }
})

let iv = 0;

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function circle() {
  let px = 0, py = 0;
  home();
  let steps = 80;
  let dv = (Math.PI * 2) / steps
  let step = 0;
  state.interval = setInterval(() => {
    if (state.interval === 0) return
    if (step >= steps) {
      clearInterval(state.interval)
      state.interval = 0
    }
    let cx = Math.sin(step * dv) * 64
    let cy = 0

    sendCommand(cx, cy, true)


    if (step + 1 >= steps) return
    step = (step + 1) % steps
  }, 40)


}

function halt() {
  clearInterval(state.interval)
  state.interval = 0
}

function moveDelta(a: number, b: number) {
  sendCommand(a, b, false)
}

function speedDelta(delta: number) {
  state.speed += delta
}

interface Output {
  pan: number,
  tilt: number,
}

function fx(a: number, b: number): string {
  let pan = Math.round(Math.sin(a) * 1000) / 1000
  let tilt = Math.round(Math.cos(b) * 1000) / 1000
  return `[${a}, ${b}] => (${pan}, ${tilt})`
}

function home() {
  sendCommand(0, 0, true)
}

</script>

<template>
  <div class="d-flex gap-1 h-100">
    <div class="w-100">
      <div class="d-flex gap-1">
        <div class="element" style="width: 13rem">
          <div class="d-flex justify-content-between label-c0 label-o3 px-1">
            <div>Connected</div>
            <div>{{ state.connected }}</div>
          </div>
          <div class="d-flex justify-content-between label-c0 label-o3 px-1">
            <div>A</div>
            <div>{{ state.pos.a }}</div>
          </div>
          <div class="d-flex justify-content-between label-c0 label-o3 px-1">
            <div>B</div>
            <div>{{ state.pos.b }}</div>
          </div>
          <div class="d-flex justify-content-between label-c0 label-o3 px-1">
            <div>Speed</div>
            <div>{{ state.speed }} steps</div>
          </div>
        </div>
        <div class="flex-fill d-flex flex-column gap-1">
          <div class="element">
            <div
                class="label-c2 label-o3 label-w600">Relative Pan-Tilt
            </div>
            <div>
              <PanTilt :a="(state.pan/512)*90"
                       :b="(state.tilt/512)*90-90"></PanTilt>
            </div>
            <div class="d-flex justify-content-between">
              <div class="label-c2 label-o3 label-w600 w-50 text-center">
                Pan
              </div>

              <div class="label-c2 label-o3 label-w600 w-50 text-center">
                Tilt
              </div>
            </div>
          </div>
          <div class="element">
            <div
                class="label-c2 label-o3 label-w600">Motor Absolute Position
            </div>
            <div>
              <PanTilt :a="(state.pos.a/512)*90"
                       :b="-(state.pos.b/512)*90"></PanTilt>
            </div>
            <div class="d-flex justify-content-between">
              <div class="label-c2 label-o3 label-w600 w-50 text-center">
                Right
              </div>

              <div class="label-c2 label-o3 label-w600 w-50 text-center">
                Left
              </div>
            </div>
           
          </div>
        </div>
      </div>
    </div>
    <div class="w-100 ">
      <div class="d-flex  gap-1">


        <Plot :cols="2" :rows="3" style="width: 13rem">
          <div class="subplot arrow-button"
               @click="(e) => {setOrigin()}">Set Origin
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {home()}">Go To Origin
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {state.speed=state.speed/2}">Speed -1
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {state.speed=state.speed*2}">Speed +1
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {circle()}">Circle
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {halt()}">HALT
          </div>
        </Plot>
        <div class="button-box element" style="width: 13rem">
          <div></div>
          <div class="subplot arrow-button"
               @click="(e) => {sendCommand(-state.speed, state.speed, false)}">􀄨
          </div>
          <div></div>
          <div class="subplot arrow-button"
               @click="(e) => {sendCommand(state.speed*2, state.speed*2, false)}">􀰌
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {sendCommand(state.speed, -state.speed, false)}">􀄩
          </div>
          <div class="subplot arrow-button"
               @click="(e) => {sendCommand(-state.speed*2, -state.speed*2, false)}">􀰑
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.arrow-button {
  display: flex;
  justify-content: center;
}

.button-box {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  grid-template-rows: repeat(2, 1fr);
  grid-gap: 0.25rem;
}

.nav-bar {
  height: 2rem;
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  grid-template-rows: repeat(1, 1fr);
}
</style>
<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, onUnmounted, reactive, watchEffect} from "vue";
import {Attribute, Entity} from "@/types";
import Beam from "@/views/beam/Beam.vue";
import {Remote} from "@/remote";
import attributeService from "@/services/attributeService";
import core from "@/core";


const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  dx: 0,
  dy: 0,
  px: 0,
  py: 0,
  af: 0,
  x: 0,
  y: 0,
  z: 0,
  pan: 90,
  panFine: 90,
  tilt: 180,
  tiltFine: 180,
  laserBeam: {
    power: 0
  } as Beam,
  lastUpdate: 0,
  laser: false,
  auth: false,
  connected: false,
  entity: {} as Entity,
  position: {} as Attribute,
  beam: {} as Attribute
})
const remote = core.remote()
watchEffect(() => {
  findEntity(remote)
  return remote.attributes
})

const height = 2.7432

const roomSize = 5

watchEffect(() => {
  if (new Date().valueOf() - state.lastUpdate > 20 && state.connected) {
    // moveBeam(Math.round(state.panFine), Math.round(state.tiltFine))
    // state.px = Math.max(Math.min(state.px, 20), -20)
    // state.py = Math.max(Math.min(state.py, 20), -20)
    goToXYZ(state.px, state.py, height)
  }
  return state.px + state.py
})

function findEntity(rem: Remote) {
  let entity = rem.entities.find(e => e.name === "sentryA")
  if (!entity) return
  state.entity = entity
  let posAttribute = rem.attributes.find(e => e.entity === state.entity.id && e.key === "position")
  if (!posAttribute) return
  state.position = posAttribute
  let laserAttribute = rem.attributes.find(e => e.entity === state.entity.id && e.key === "beam")
  if (!laserAttribute) return
  state.beam = laserAttribute
  query()
}

function sendFine() {
  moveBeam(state.panFine, state.tiltFine)
}

function convertXYToPanTilt(x: number, y: number, h: number, r: number): { pan: number; tilt: number } {
  // Calculate pan angle
  const pan = Math.atan2(x, r);

  // Calculate tilt angle
  const tilt = Math.atan2(y, Math.sqrt((x * x) + (r * r)));

  // Convert radians to degrees
  const panDeg = (pan * 180) / Math.PI;
  const tiltDeg = (tilt * 180) / Math.PI;

  return {pan: panDeg, tilt: tiltDeg};
}

function goToXYZ(x: number, y: number, z: number) {
  // let distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2) + Math.pow(z, 2))
  // let theta = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI)
  // let phi = Math.atan(Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2)) / z)
  // let pan = theta * 180 / Math.PI
  // let tilt = phi * 180 / Math.PI

  let {pan, tilt} = convertXYToPanTilt(x, y, height, roomSize);

  // if (pan >= 0 && pan <= 180 && tilt >= 0 && tilt <= 180) {
  state.panFine = pan + 90
  state.tiltFine = tilt + 90
  moveBeam(state.panFine, state.tiltFine)
  // }
}

function reversePanTiltToXY(pan: number, tilt: number, h: number): { x: number; y: number } {
  // Convert degrees to radians
  const panRad = (pan * Math.PI) / 180;
  const tiltRad = (tilt * Math.PI) / 180;

  // Calculate x, y coordinates
  const x = h * Math.tan(panRad);
  const y = h * Math.tan(tiltRad);

  return {x, y};
}

function query() {
  if (!state.position) return
  let status = JSON.parse(state.position.value)
  state.pan = status.pan
  state.tilt = status.tilt
  state.panFine = state.pan
  state.tiltFine = state.tilt
  let theta = Math.PI / 180 * state.pan
  let phi = Math.PI / 180 * state.tilt
  let distance = 10;
  let {x, y} = reversePanTiltToXY(state.pan, state.tilt, height);
  state.x = x
  state.y = y

  state.z = distance * Math.cos(phi)
  if (!state.beam) return
  state.laserBeam = JSON.parse(state.beam.value) as Beam;
  state.laser = (state.laserBeam.active === 1)
  state.connected = true

}


function laserPower(on: boolean) {
  if (!state.entity) return
  if (!state.laserBeam) return
  setDuty(50)
  let beam = state.laserBeam
  beam.active = on ? 1 : 0;
  state.beam.request = JSON.stringify(beam)
  attributeService.request(state.beam)
  // remote.nexus.requestAttribute(state.entity.id, "beam", payload)
}

function setDuty(percent: number) {

  if (!state.entity) return
  if (!state.laserBeam) return
  let beam = state.laserBeam
  beam.power = percent
  state.laserBeam.power = percent
  state.beam.request = JSON.stringify(beam)
  attributeService.request(state.beam)
}

function moveBeam(pan: number, tilt: number) {
  if (!state.entity) return
  state.position.request = JSON.stringify({
    pan: Math.round(pan),
    tilt: Math.round(tilt)
  })


  state.lastUpdate = new Date().valueOf()

  attributeService.request(state.position).then(e => console.log(e))
}


onMounted(() => {
  configure()
  findEntity(remote)
  animate()
  query()

})

onUnmounted(() => {
  cancelAnimationFrame(state.af)

})

function drawLoop(ctx: CanvasRenderingContext2D, rad: number) {
  let dx = (Math.PI * 2) / rad
  ctx.strokeStyle = "rgba(255,255,255,0.2)"
  ctx.lineWidth = 3
  let lx = Math.cos(0) * rad
  let ly = Math.sin(0) * rad
  ctx.beginPath()
  for (let i = 0; i <= rad + 1; i++) {
    let x = Math.cos(dx * i) * rad
    let y = Math.sin(dx * i) * rad
    if (i * dx > Math.PI - Math.PI / 4 && i * dx < Math.PI / 4) {
      ctx.strokeStyle = "rgba(255,255,12,1)"
      ctx.moveTo(lx, ly)
      ctx.lineTo(x, y)

    } else {
      ctx.strokeStyle = "rgba(255,255,255,0.2)"
    }
    lx = x
    ly = y


  }
  ctx.closePath()
  ctx.stroke()
}

function animate() {
  state.af = requestAnimationFrame(animate)
  let max = state.canvas.width / 8
  let w = state.ctx.canvas.width
  let h = state.ctx.canvas.height
  let delta = 10


  render()
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function render() {
  let ctx = state.ctx
  let w = state.ctx.canvas.width
  let h = state.ctx.canvas.height
  ctx.clearRect(0, 0, state.canvas.width, state.canvas.height)
  ctx.strokeStyle = "rgba(255,255,255,0.8)"
  // ctx.fillStyle = "rgb(200,200,200)"
  // ctx.ellipse(w / 2, h / 2, w / 4, w / 4, 0, 0, Math.PI * 2, false)
  // ctx.fill()
  // ctx.fillStyle = "rgb(0,0,0)"
  // ctx.ellipse(w / 2, h / 2, w / 2, w / 2, 0, 0, Math.PI * 2, false)
  // ctx.stroke()
  let pointsX = roomSize;
  let pointsY = roomSize;
  let dx = (w / 2) / pointsX
  let dy = (h / 2) / pointsY
  ctx.beginPath()
  ctx.fillStyle = "rgba(255,255,255, 0.1)"
  ctx.ellipse(w / 2 + state.px * dx, h / 2 + state.py * dy, w / 8, h / 8, 0, 0, Math.PI * 2, false)
  ctx.closePath()
  ctx.fill()


}

function configure() {
  state.canvas = document.getElementById("canvas") as HTMLCanvasElement
  state.ctx = state.canvas.getContext("2d") as CanvasRenderingContext2D
  let pd = window.devicePixelRatio || 1
  state.ctx.scale(pd, pd)
  state.canvas.width = state.ctx.canvas.clientWidth * pd;
  state.canvas.height = state.ctx.canvas.clientHeight * pd;
  // state.ctx.rect(0, 0, state.canvas.width, state.canvas.height)

  render()
  // ctx.translate(w / 2, h / 2)
  // for (let i = 0; i < 8; i++) {
  //   drawLoop(ctx, w / 4 + i * Math.exp(2.5 + i / 20))
  //
  // }
  // // ctx.translate(-w / 2, -h / 2)

}


function up(e: TouchEvent) {
  state.dx = 0
  state.dy = 0
  // state.px = 0
  // state.py = 0
}


function down(e: TouchEvent) {

}

function move(e: TouchEvent) {
  e.preventDefault()
  let br = state.canvas.getBoundingClientRect()
  // let dist = Math.sqrt(Math.pow(dx, 2) + Math.pow(dy, 2)) / 5
  // let t = 2 * Math.atan((dy) / (dx + Math.sqrt(Math.pow(dx, 2) + Math.pow(dy, 2))))
  // state.dx = Math.cos(t) * dist / 2
  // state.dy = Math.sin(t) * dist / 2

  state.dx = e.touches.item(0).pageX - br.x
  state.dy = e.touches.item(0).pageY - br.y
  let ppx = state.px
  // if (state.dx >= br.left && state.dx <= br.right) {
  ppx = map_range(state.dx, 0, state.canvas.clientWidth, -roomSize, roomSize)
  if (ppx >= -roomSize && ppx <= roomSize) {
    state.px = ppx
  }
  // }

  let ppy = state.py
  if (state.dy >= br.height && state.dy <= br.bottom) {

  }
  ppy = map_range(state.dy, 0, state.canvas.clientHeight, -roomSize, roomSize)
  if (ppy >= -roomSize && ppy <= roomSize) {
    state.py = ppy
  }

}

function touchMovePUp(e: TouchEvent) {
  e.preventDefault()
  state.panFine += 5
  sendFine()
}

function touchMovePDown(e: TouchEvent) {
  e.preventDefault()
  state.panFine -= 5
  sendFine()
}

function touchMoveTUp(e: TouchEvent) {
  e.preventDefault()
  state.tiltFine += 5
  sendFine()
}

function touchMoveTDown(e: TouchEvent) {
  e.preventDefault()
  state.tiltFine -= 5
  sendFine()
}

function touchMoveXUp(e: TouchEvent) {
  e.preventDefault()
  state.px = 0
  state.py = 0
}

function touchMoveXDown(e: TouchEvent) {
  e.preventDefault()
  state.px = 10
  state.py = 3
}

function touchMoveYUp(e: TouchEvent) {
  e.preventDefault()
  state.py += 1
}

function touchMoveYDown(e: TouchEvent) {
  e.preventDefault()
  state.py -= 1
}
</script>

<template>
  <div class="d-flex flex-column gap-2">
    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center">P {{ Math.round(state.pan) }}</Element>
        <Element :foreground="true" class="d-flex justify-content-center">T {{ Math.round(state.tilt) }}</Element>
        <Element :foreground="true" class="d-flex justify-content-center">P {{ Math.round(state.pan) }}</Element>
        <Element :foreground="true" class="d-flex justify-content-center">T {{ Math.round(state.tilt) }}</Element>
        <Element :foreground="true" :mutable="true" class="d-flex justify-content-center"
                 @touchstart="laserPower(!state.laser)">
          {{ state.laser ? "Turn Off" : "Turn On" }}
        </Element>

      </List>
    </Element>
    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center">PX</Element>
        <Element :foreground="true" class="d-flex justify-content-center">X {{ Math.round(state.dx) }}</Element>
        <Element :foreground="true" class="d-flex justify-content-center">Y {{ Math.round(state.dx) }}</Element>
      </List>
    </Element>
    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center">DX</Element>
        <Element :foreground="true" class="d-flex justify-content-center">X {{ Math.round(state.px) }}</Element>
        <Element :foreground="true" class="d-flex justify-content-center">Y {{ Math.round(state.py) }}</Element>
      </List>
    </Element>
    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveXUp">X+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveXDown">X-</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveYUp">Y+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveYDown">Y-</Element>
        <Element :foreground="true" :mutable="true" class="d-flex justify-content-center"
                 @touchstart="laserPower(!state.laser)">
          {{ state.laser ? "Turn Off" : "Turn On" }}
        </Element>
      </List>
    </Element>
    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMovePUp">P+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMovePDown">P-</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveTUp">T+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveTDown">T-</Element>
      </List>
    </Element>
    <Element class=" d-flex align-items-center justify-content-center" style="user-select: none !important;">
      <Element :foreground="true" style="width: 100%; height:100%;aspect-ratio: 1/1 !important;">
        <canvas id="canvas" style="width: 100%; height:100%;user-select: none !important;" @touchend="up" @touchleave="up" @touchmove="move"
                @touchstart="down"></canvas>
      </Element>
    </Element>
  </div>
</template>

<style lang="scss" scoped>

</style>
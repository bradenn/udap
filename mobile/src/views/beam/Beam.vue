<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, onUnmounted, reactive, watchEffect} from "vue";
import {Attribute, Entity} from "udap-ui/types";
import {Remote} from "udap-ui/remote";
import attributeService from "@/services/attributeService";
import core from "udap-ui/core";

interface Beam {
  active: number,
  target: string,
  power: number,
}

const state = reactive({
  canvas: {} as HTMLCanvasElement,
  ctx: {} as CanvasRenderingContext2D,
  dx: 0,
  dy: 0,
  px: 0,
  py: 0,
  gx: 0,
  gy: 0,
  af: 0,
  x: 0,
  y: 0,
  z: 0,
  lastPos: {
    x: 0,
    y: 0
  },
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
const remote = core.remote() as Remote
watchEffect(() => {
  findEntity(remote)
  return remote.attributes
})


const roomSizeZ = 2.6924
const height = roomSizeZ
const roomSizeY = 3.2004
const roomSizeX = 3.5814

watchEffect(() => {
  if (new Date().valueOf() - state.lastUpdate > 25 && state.connected) {
    // moveBeam(Math.round(state.panFine), Math.round(state.tiltFine))
    // state.px = Math.max(Math.min(state.px, 20), -20)
    // state.py = Math.max(Math.min(state.py, 20), -20)
    goToXYZ(map_range(state.px, 0, roomSizeX, -roomSizeX / 2, roomSizeX / 2), map_range(state.py, 0, roomSizeY, -roomSizeY / 2, roomSizeY / 2), roomSizeZ)
  }
  return state.px
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

function drawCircle() {
  let loops = 0;
  for (let i = 0; i < loops; i++) {

    moveBeam(state.panFine, state.tiltFine)
  }

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
  // x -= roomSizeX / 2
  // let distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2) + Math.pow(z, 2))
  // let theta = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI)
  // let phi = Math.atan(Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2)) / z)
  // let pan = theta * 180 / Math.PI
  // let tilt = phi * 180 / Math.PI
  // let {pan, tilt} = convertPositionToPanTilt(x, y, roomSizeX, roomSizeY, inToM(67), inToM(19), roomSizeZ)
  let {phi, theta} = PlanarToSpherical(0, 0, x, y + roomSizeY / 2 - roomSizeY / 8, roomSizeZ);
  let lp = (phi * (180.0 / Math.PI))
  let theta1 = (theta * (180.0 / Math.PI))
  state.panFine = (phi * (180.0 / Math.PI))
  state.tiltFine = (theta * (180.0 / Math.PI))
  if (lp >= 0 && lp <= 180 && theta1 >= 0 && theta1 <= 180) {
    moveBeam(state.panFine, state.tiltFine)

  }
}


function SphericalToPlanar(phi: number, theta: number, height: number): { x: number, y: number } {
  theta *= (Math.PI / 180)
  phi *= (Math.PI / 180)
  let r = Math.tan(theta) * height;
  let x = Math.cos(phi) * r;
  let y = Math.sin(phi) * r;


  return {x, y};
}

function PlanarToSpherical(cx: number, cy: number, x: number, y: number, height: number): {
  phi: number,
  theta: number
} {

  let r: number = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2));
  let phi: number = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI)
  let theta: number = (Math.atan(r / height));


  return {phi, theta};
}

function query() {
  if (!state.position) return
  let status = JSON.parse(state.position.value)
  state.pan = status.pan
  state.tilt = status.tilt

  let theta = Math.PI / 180 * state.pan
  let phi = Math.PI / 180 * state.tilt
  let distance = 10;
  let {x, y} = SphericalToPlanar(state.panFine, state.tiltFine, roomSizeZ);
  state.x = x
  state.y = y

  state.z = distance * Math.cos(phi)
  if (!state.beam) return
  state.laserBeam = JSON.parse(state.beam.value) as Beam;
  state.laserBeam.power = Math.round(state.laserBeam.power);
  state.laser = (state.laserBeam.active === 1)
  state.connected = true

}


function laserPower(on: boolean) {
  if (!state.entity) return
  if (!state.laserBeam) return
  setDuty(1)
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
    tilt: 90 - Math.round(tilt)
  })


  state.lastPos.x = state.px
  state.lastPos.y = state.py
  attributeService.request(state.position).then(e => {
  })
  state.lastUpdate = new Date().valueOf()
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

function inToM(i: number): number {

  return (i / 39.37) * (state.ctx.canvas?.height / roomSizeY)
}

function drawRoom() {
  let ctx = state.ctx
  let w = state.ctx.canvas.width
  let h = state.ctx.canvas.height

  let items = [
    {
      name: "Matthew's Desk",
      x: 141 - 20,
      y: 37,
      w: 20,
      h: 47,
    },
    {
      name: "Book Shelf",
      x: 141 - 30,
      y: 126 - 15,
      w: 30,
      h: 15,
    },
    {
      name: "Window Shelf",
      x: 141 - 42 - 48,
      y: 126 - 12,
      w: 48,
      h: 12,
    },
    {
      name: "Braden's Side Desk",
      x: 0,
      y: 126 - 20,
      w: 47,
      h: 20,
    },
    {
      name: "Braden's Desk",
      x: 0,
      y: 126 - 20 - 65,
      w: 24,
      h: 63,
    },
    {
      name: "Drawers",
      x: 0,
      y: 126 - 20 - 65 - 34,
      w: 19,
      h: 32,
    },
    {
      name: "Laser",
      x: 67,
      y: roomSizeY / 2 * 39.37,
      w: 4,
      h: 4,
    }
  ]
  ctx.fillStyle = "rgba(255,255,255,0.15)"
  for (let i = 0; i < items.length; i++) {
    let item = items[i]


    if (item.name == "Laser") {
      ctx.fillStyle = "rgba(255,255,255,0.15)"
      let r = h / 2;
      let mx = w / 2
      let my = 50
      ctx.fillRect(mx - 50, my - 50, 100, 100)
      ctx.beginPath()
      ctx.moveTo(mx - r, my)
      ctx.lineTo(mx + r, my)
      ctx.closePath()
      ctx.stroke()
      ctx.beginPath()
      ctx.moveTo(mx, my - r)
      ctx.lineTo(mx, my + r)
      ctx.closePath()
      ctx.stroke()
      ctx.beginPath()
      ctx.ellipse(w / 2, my, r, r, 0, 0, Math.PI * 2, false)
      ctx.closePath()
      ctx.stroke()

      let {x, y} = SphericalToPlanar(state.panFine, state.tiltFine, roomSizeZ);
      // ctx.beginPath()
      // ctx.moveTo(mx, my)
      // ctx.lineTo(mx, my)
      // ctx.closePath()
      // ctx.stroke()

      ctx.beginPath()
      ctx.moveTo(mx, my)
      ctx.lineTo(map_range(x, -roomSizeX / 2, roomSizeX / 2, 0, w), map_range(y - roomSizeY / 2 + roomSizeY / 8, -roomSizeY / 2, roomSizeY / 2, 0, h))
      ctx.closePath()
      ctx.stroke()
    } else {
      ctx.fillRect(inToM(item.x), inToM(item.y), inToM(item.w), inToM(item.h))
    }

  }

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
  drawRoom()
  let pointsX = roomSizeX;
  let pointsY = roomSizeY;
  let dx = (w) / pointsX
  let dy = (h) / pointsY
  ctx.beginPath()
  ctx.fillStyle = "rgba(255,255,255, 0.1)"
  ctx.ellipse(state.px * dx, state.py * dy, w / 8, h / 8, 0, 0, Math.PI * 2, false)
  ctx.closePath()
  ctx.fill()
  const points = 65
  let dtx = w / points
  let dty = h / points
  for (let i = 1; i < points; i++) {
    for (let j = 1; j < points; j++) {
      ctx.beginPath()
      if (j == state.gy) {
        ctx.fillStyle = "rgba(255, 128,128, 1)"
      } else if (i == state.gx) {
        ctx.fillStyle = "rgba(128, 128,255, 1)"

      } else {
        ctx.fillStyle = "rgba(255,255,255, 0.1)"
      }

      ctx.ellipse(i * dtx, j * dty, 2, 2, 0, 0, Math.PI * 2, false)
      ctx.closePath()
      ctx.fill()
    }
  }


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
  let touch = e.touches.item(0)
  if (!touch) return
  state.dx = Math.min(touch.pageX - br.x, state.canvas.clientWidth)
  state.dy = Math.min(touch.pageY - br.y, state.canvas.clientHeight)
  let ppx = state.px
  // if (state.dx >= br.left && state.dx <= br.right) {
  ppx = map_range(state.dx, 0, state.canvas.clientWidth, 0, roomSizeX)
  if (ppx >= 0 && ppx <= roomSizeX) {
    state.px = map_range(state.dx, 0, state.canvas.clientWidth, 0, roomSizeX)
    state.gx = Math.round(map_range(ppx, 0, roomSizeX, 0, 65))

  }
  // }

  let ppy = state.py
  if (state.dy >= br.height && state.dy <= br.bottom) {

  }
  ppy = map_range(state.dy, 0, state.canvas.clientHeight, 0, roomSizeY)
  if (ppy >= 0 && ppy <= roomSizeY) {
    state.py = map_range(state.dy, 0, state.canvas.clientHeight, 0, roomSizeY)
    state.gy = Math.round(map_range(ppy, 0, roomSizeY, 0, 65))
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

const dutys = [1, 3, 5, 11, 15]

</script>

<template>
  <div class="d-flex flex-column gap-2">

    <Element>
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center">P {{
            Math.round(state.px * 1000) / 1000
          }}
        </Element>
        <Element :foreground="true" class="d-flex justify-content-center">T {{
            Math.round(state.py * 1000) / 1000
          }}
        </Element>
        <Element :foreground="true" class="d-flex justify-content-center">P
          {{ Math.round(state.panFine) }}
        </Element>
        <Element :foreground="true" class="d-flex justify-content-center">T {{ Math.round(state.tiltFine) }}</Element>
        <Element :foreground="true" :mutable="true" class="d-flex justify-content-center"
                 :cb="() => laserPower(!state.laser)">
          {{ state.laser ? "Turn Off" : "Turn On" }}
        </Element>

      </List>
    </Element>
    <Element>
      <List :row="true">
        <!--        <Element :foreground="true" class="d-flex justify-content-center">-->
        <List class="w-100" row>

          <Element v-for="duty in dutys" :accent="state.laserBeam.power === duty" :cb="() => setDuty(duty)" class="d-flex align-items-baseline gap-1 w-100 justify-content-center position-relative"
                   foreground>
            <div>{{ duty }}</div>
            <div class="label-c7 label-o3">mW</div>
            <div v-if="duty > 5"
                 :style="`color: ${duty>11?'hsla(0deg, 40%, 60%);':'hsla(30deg, 0%, 30%);'};`"
                 class="label-c7 label-o3 sf label-o3 position-absolute"
                 style="font-size: 10px; right: 6px; top:2px;">
              􀋰
            </div>
          </Element>


        </List>
        <!--        </Element>-->

      </List>
    </Element>

    <Element v-if="false">
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
    <Element v-if="true">
      <List :row="true">
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMovePUp">P+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMovePDown">P-</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveTUp">T+</Element>
        <Element :foreground="true" class="d-flex justify-content-center" @touchstart="touchMoveTDown">T-</Element>
      </List>
    </Element>
    <Element class=" d-flex align-items-center justify-content-center" style="user-select: none !important;">
      <Element :foreground="true" :style="`aspect-ratio: ${roomSizeX}/${roomSizeY} !important;`"
               style="width: 100%; height:100%;">
        <canvas id="canvas" style="width: 100%; height:100%;user-select: none !important;" @touchend="up"
                @touchleave="up" @touchmove="move"
                @touchstart="down"></canvas>
      </Element>
    </Element>
  </div>
</template>

<style lang="scss" scoped>

</style>
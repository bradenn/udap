<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Plot from "@/components/plot/Plot.vue";
import {inject, onMounted, onUnmounted, reactive, watch} from "vue";
import Subplot from "@/components/plot/Subplot.vue";
import Confirm from "@/components/plot/Confirm.vue";
import type {Attribute, Entity, Session} from "@/types";
import * as THREE from "three";
import {OrbitControls} from "three/examples/jsm/controls/OrbitControls";
import type {
  EffectComposer
} from "three/examples/jsm/postprocessing/EffectComposer";
import Sentry from "@/components/rendered/Sentry.vue";
import type {Remote} from "@/remote";
import attributeService from "@/services/attributeService";
import Attenuation from "@/components/Attenuation.vue";
import Slider from "@/components/Slider.vue";


let room = [
  new THREE.Vector2(0.001, 0),
  new THREE.Vector2(2.02, 0),
  new THREE.Vector2(4.03, -2.07),
  new THREE.Vector2(4.37, -2.07),
  new THREE.Vector2(4.37, -3.17),
  new THREE.Vector2(2.33, -3.17),
  new THREE.Vector2(1.17, -3.17),
  new THREE.Vector2(0, -3.17),
]

let state = reactive({
  x: 0,
  y: 0,
  z: 0,
  pointer: {
    x: 0,
    y: 0
  },
  pan: 90,
  panFine: 90,
  tilt: 180,
  tiltFine: 180,
  zoom: 0,
  runner: 0,
  speed: 1,
  laser: false,
  laserBeam: {
    power: 0
  } as Beam,
  auth: false,
  entity: {} as Entity,
  position: {} as Attribute,
  beam: {} as Attribute
})

const emitter = {
  x: 0,
  y: 0,
  z: 0,
}

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let scene = {} as THREE.Scene
let beamLine = {} as THREE.Line
let controls = {} as OrbitControls
let composer = {} as EffectComposer
let floorObj = new THREE.Object3D()

const s = 200;

onMounted(() => {
  load3d()
  verifyAuth(session)
  findEntity(remote)
  laserStop()
})


function moveBeamToXYZ(x: number, y: number, z: number) {
  state.x = x
  state.y = y
  state.z = z
  let distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2) + Math.pow(z, 2))
  let theta = Math.atan(y / x)
  let phi = Math.atan(Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2)) / z)
  laserPanTilt(theta * 180 / Math.PI, phi * 180 / Math.PI)
  moveBeam(theta, phi, distance)

}

function animate() {
  requestAnimationFrame(animate);
  controls.update()
  render()
}

function moveBeam(pan: number, tilt: number, distance: number) {

  let theta = -(pan * Math.PI / 180 + Math.PI / 2)
  let phi = (tilt * Math.PI / 180)

  let x = distance * Math.cos(theta) * Math.sin(phi)
  let y = distance * Math.sin(theta) * Math.sin(phi)
  let z = distance * Math.cos(phi)
  let origin = new THREE.Vector3(4.27, -2.62, 2.28)
  let target = new THREE.Vector3(4.27 + x, -2.62 + y, 2.28 - z)

  let newGeom = new THREE.BufferGeometry().setFromPoints([origin, target])
  newGeom.scale(s, s, s)
  beamLine.geometry.dispose()
  beamLine.geometry = newGeom
}

function drawBeam(): THREE.Object3D {
  const beamMaterial = new THREE.LineBasicMaterial({
    color: 0xff0000,
    linewidth: 3,
    linecap: 'round', //ignored by WebGLRenderer
    linejoin: 'round' //ignored by WebGLRenderer
  })


  const beam = new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(4.27, -2.62, 2.37), new THREE.Vector3(0, 0, 0)])
  beam.scale(s, s, s)

  beamLine = new THREE.Line(beam, beamMaterial);

  let obj = new THREE.Object3D()
  obj.add(beamLine)

  return obj

  // const radius = 120;
  // const radials = 64;
  // const circles = 24;
  // const divisions = 64;
  //
  // const grid = new THREE.GridHelper(100, 20)
  // grid.rotateX(Math.PI / 2)
  // scene.add(grid)

}


function drawFloor(): THREE.Object3D {
  const floor = new THREE.Shape();

  let object = new THREE.Object3D();

  floor.setFromPoints(room)

  const floorGeometry = new THREE.ExtrudeBufferGeometry([floor], {
    depth: 0.05, bevelEnabled: false, bevelSize: 0.05, bevelThickness: 0.05
  });

  floorGeometry.scale(s, s, s)

  let floorMaterial = new THREE.MeshPhysicalMaterial({
    color: 0x1c2124,
    opacity: 1,
    roughness: 0.85,
    transparent: false,
  });


  const floorMesh = new THREE.Mesh(floorGeometry, floorMaterial);

  object.add(floorMesh)

  return object


  // const gridHelper = new THREE.GridHelper(100, 60, 0x8C929B, 0x8C929B);
  //
  // gridHelper.rotateX(Math.PI / 2)
  // gridHelper.translateY(0.06 * s)
  // scene.add(gridHelper);

}

function drawWall(points: THREE.Vector2[]): THREE.Object3D {
  const wall1 = new THREE.Shape();
  wall1.setFromPoints(points)

  const wall2 = new THREE.ShapeGeometry();
  let pts2 = points.map(p => p.clone().multiply(new THREE.Vector2(1.2, 1.2)))
  wall2.setFromPoints(pts2)


  const wallGeometry = new THREE.ExtrudeBufferGeometry([wall1], {
    depth: 2.37,
    bevelEnabled: false,
  });


  wallGeometry.scale(s, s, s)

  let wallMaterial = new THREE.MeshPhysicalMaterial({
    color: 0x323739,
    opacity: 1,
    roughness: 1,
    transparent: false,
  });

  return new THREE.Mesh(wallGeometry, wallMaterial)
}

function drawWalls(points: THREE.Vector2[]): THREE.Object3D {
  let obj = new THREE.Object3D()
  for (let i = 0; i < points.length; i++) {
    obj.add(drawWall([points[i], points[(i + 1) % points.length]]))
  }
  return obj
}

function setCamera(x: number, y: number, z: number) {
  camera.position.set(x, y, z);
  camera.lookAt(0, 0, 0);
}

function render() {
  renderer.setClearColor(0x000000, 0);

  renderer.render(scene, camera);
}

function load3d() {
  renderer = new THREE.WebGLRenderer();
  renderer.shadowMap.enabled = true;
  let element = document.getElementById("room-container")
  if (!element) return


  renderer.setSize(element.clientWidth, element.clientHeight);
  renderer.setPixelRatio(window.devicePixelRatio);

  element.appendChild(renderer.domElement)
  let width = window.innerWidth
  let height = window.innerHeight
  document.addEventListener('pointermove', onPointerMove);

  // camera = new THREE.OrthographicCamera(width / -2, width / 2, height / 2, height / -2, -1000, 1000);
  camera = new THREE.PerspectiveCamera(20, window.innerWidth / window.innerHeight, 1, 10000);
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true


  setCamera(0, -1 * s, 13 * s)


  scene = new THREE.Scene();


  const axesHelper = new THREE.AxesHelper(10);
  axesHelper.setColors(new THREE.Color(255, 0, 0), new THREE.Color(0, 255, 0), new THREE.Color(0, 0, 255))
  scene.add(axesHelper);
  let roomObject = new THREE.Object3D()

  roomObject.add(drawFloor())
  roomObject.add(drawWalls(room))
  roomObject.add(drawBeam())
  floorObj.add(roomObject)
  // drawScene(pointsd)
  scene.add(new THREE.HemisphereLight(0xffffff, 0xcccccc, 1.4))

  roomObject.rotateZ(Math.PI)
  roomObject.translateX(-(4.37 / 2) * s)
  roomObject.translateY((3.17 / 2) * s)
  scene.add(roomObject)


  animate()
  render()
  // const gridHelper = new THREE.GridHelper(s * 10, s / 2)
}

let session = inject("session") as Session
let remote = inject("remote") as Remote


onUnmounted(() => {
  laserStop()
})

interface Beam {
  active: number,
  target: string,
  power: number,
}

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

watch(remote, (current: Remote, past: Remote) => {
  findEntity(current)
})

watch(session, (current: Session, previous: Session) => {
  verifyAuth(current)
})

function verifyAuth(current: Session) {
  state.auth = (!!current.user.id)
}

function query() {
  if (!state.position) return
  let status = JSON.parse(state.position.value)
  state.pan = status.pan
  state.tilt = status.tilt
  let theta = Math.PI / 180 * state.pan
  let phi = Math.PI / 180 * state.tilt
  let distance = 10;
  state.x = distance * Math.cos(theta) * Math.sin(phi)
  state.y = distance * Math.sin(theta) * Math.sin(phi)
  state.z = distance * Math.cos(phi)
  if (!state.beam) return
  state.laserBeam = JSON.parse(state.beam.value) as Beam;
  state.laser = (state.laserBeam.active === 1)
  moveBeam(state.pan, state.tilt, state.laser ? 5 : 0)

}

function laserPower(on: boolean) {
  if (!state.entity) return
  if (!state.laserBeam) return
  let beam = state.laserBeam
  beam.active = on ? 1 : 0;
  state.beam.request = JSON.stringify(beam)
  attributeService.request(state.beam)
  // remote.nexus.requestAttribute(state.entity.id, "beam", payload)
}

function laserToggle() {
  laserPower(!state.laser)
}

function laserTilt(value: number) {
  if (!state.entity) return
  state.position.request = JSON.stringify({
    pan: Math.round(state.pan),
    tilt: Math.round(value)
  })
  attributeService.request(state.position)
  state.tiltFine = 0
  // remote.nexus.requestAttribute(state.entity.id, "position", payload)
}

function laserPanFine(value: number) {
  laserPan(state.pan + value)
  state.panFine = 0
}

function laserTiltFine(value: number) {
  laserTilt(state.tilt + value)
  state.tiltFine = 0
}


function laserPan(value: number) {
  if (!state.entity) return
  state.position.request = JSON.stringify({
    pan: Math.round(value),
    tilt: Math.round(state.tilt)
  })
  attributeService.request(state.position)
  state.panFine = 0
}

function laserSafe() {
  laserPanTilt(180, 180)
}

function laserHome() {
  laserPanTilt(90, 90)
}

function laserWall() {
  laserPanTilt(105, 154)
}

function laserPanTilt(pan: number, tilt: number) {
  if (!state.entity) return
  state.position.request = JSON.stringify({
    pan: Math.round(pan),
    tilt: Math.round(tilt)
  })
  attributeService.request(state.position)
  // let panA = map_range(pan, 0, 180, 0, 1800)
  // let tiltA = map_range(tilt, 0, 180, 0, 1800)
  // axios.get(`http://10.0.1.60/pan/${panA}/tilt/${tiltA}`).then(res => {
  //   state.pan = map_range(res.data.pan, 0, 1800, 0, 180)
  //   state.tilt = map_range(res.data.tilt, 0, 1800, 0, 180)
  // }).catch(res => {
  //   console.log(res)
  // })
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function laserSpeed(speed: number) {
  state.speed = speed
}

function onPointerMove(event: MouseEvent) {


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

function laserStopAll() {
  clearInterval(state.runner)
  state.runner = 0
  state.speed = 1
  laserSafe()
  laserPower(false)
}

function getSpeed(): number {
  return state.speed
}

function goToXYZ(x: number, y: number, z: number) {
  let distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2) + Math.pow(z, 2))
  let theta = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI)
  let phi = Math.atan(Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2)) / z)
  let pan = theta * 180 / Math.PI
  let tilt = phi * 180 / Math.PI
  if (pan >= 0 && pan <= 180 && tilt >= 0 && tilt <= 180) {
    laserPanTilt(pan, tilt)
  }
}

// 2.5, 3.1, 9.5
// 2.5, 2.1, 9.5
// 0.9, 2.1, 9.5
// 0.9, 3.1, 9.5


function laserRun() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  let steps = 500;
  let method = false;

  let speed = 1

  state.runner = setInterval(() => {
    if (tick > steps) {
      dir = true
    } else if (tick < 1) {
      dir = false
    }
    tick += dir ? -speed : speed

    // tick = (tick + 1) % steps


    //
    // laserPan(Math.cos(Math.floor(tick)) * map_range(Math.floor((2 * Math.PI / 100) * tick), 0, 100, 15, 1) + a1)
    // laserTilt(Math.sin(Math.floor(tick)) + b1)
//map_range(tick, 0, steps, 0, 9)
    let div = (Math.PI * 2) / steps

    let a = 2
    let b = 0.3
    // let r = map_range(tick, 0, steps, 0, 1)


    let theta = tick * div;

    // moveBeamToXYZ(map_range(Math.sin(div * tick * 2), -1, 1, -0.7, 2.5), map_range(tick, 0, steps, 1.5, 8), 9.5)

    let bounds = 10


    laserPanTilt(map_range(method ? Math.sin(div * tick * 2) : Math.cos(div * tick * 2), -1, 1, 90 - bounds, 90 + bounds * 2), map_range(tick, 0, steps, 10, 45))

    let op = Math.floor(Math.random() * 100);
    if (op > 95) {
      dir = !dir;
    } else if (op > 25 && op < 40) {
      speed = 0;
    } else if (op < 10) {
      speed = 0.5 + Math.random() * 3
    } else if (op > 15 && op < 18) {
      method = !method
    }

    // moveBeamToXYZ(0, 0, 0);

    // laserPanTilt(map_range(2 * a * Math.pow(Math.sin(theta), 2) * Math.cos(theta), -a, a, 90-10, 90+10), map_range(2 * a * Math.pow(Math.cos(theta), 2) * Math.sin(theta), -a, a, 65, 78))

    // laserPanTilt(map_range(a*Math.exp(b*r)*Math.cos(r), -a, a, 90-10, 90+10), map_range(a*Math.exp(b*r)*Math.sin(r), -a, a, 65, 78))

    // laserPanTilt(map_range(2 * a * Math.pow(Math.sin(theta), 2) * Math.cos(theta), -2, 2, 90 - 10, 90 + 30), map_range(2 * a * Math.pow(Math.cos(theta), 2) * Math.sin(theta), -2, 2, 70, 78))


  }, 50)

}

function laserCircle() {
  if (state.runner != 0) {
    return
  }
  let tick = 0;
  let dir = false;

  state.runner = setInterval(() => {
    tick += 10;

    if (tick >= 1000) {
      tick = 0;
    } else if (tick <= 0) {
      tick = 1000;
    }

    // 0 - 180
    let panTo = map_range(tick, 0, 1000, 3.1, 3.1)
    // 90 - 180
    let tiltTo = map_range(tick, 0, 1000, 3.4 - 0.5, 3.4 + 0.5)

    goToXYZ(tiltTo, panTo, 5)

  }, 30)

}

// 34 -> 62
// 54 -> 68

function laserStop() {
  clearInterval(state.runner)
  state.runner = 0
}

</script>

<template>
  <div class="h-100 d-flex flex-column">
    <div class="d-flex w-100 h-100 gap-2 mt-1 pb-2">
      <div class="d-flex flex-column gap-1">
        <Plot :cols="3" :rows="2" style="" title="Location">
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            X: {{ Math.round(state.x * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            Y: {{ Math.round(state.y * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            Z: {{ Math.round(state.z * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            P: {{ Math.round(state.pan * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            T: {{ Math.round(state.tilt * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c2 text-center">
            Z: {{ Math.round(state.zoom * 10) / 10 }}
          </div>
        </Plot>

        <sentry :beam="state.laser" :pan="state.pan" :tilt="state.tilt"
                color="rgba(255,0, 0, 1)">
        </sentry>

        <Plot v-if="false" :cols="1" :rows="2" style="width: 13rem;"
              title="Programmed">
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
                class="slider element"
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
        <div id="room-container" class=" element w-100 flex-grow-1">
        </div>
      </div>
      <div class="d-flex flex-column gap-1 flex-grow-1">
        <Slider :change="setDuty" :max="15" :min="0" :step="1"
                :value="Math.round(state.laserBeam.power * 100) / 100"
                name="Optical Output"
                unit=" mW"></Slider>
        <Slider :change="laserPan" :max="180" :min="0" :step="10"
                :value="state.pan"
                name="Pan"
                unit=" deg"></Slider>
        <Slider :change="laserPanFine" :max="5" :min="-5" :step="1"
                :value="state.panFine"
                name="Pan Fine"
                unit=" deg"></Slider>
        <Slider :change="laserTilt" :max="180" :min="0" :step="10"
                :value="state.tilt"
                name="Tilt"
                unit=" deg"></Slider>
        <Slider :change="laserTiltFine" :max="5" :min="-5" :step="1"
                :value="state.tiltFine"
                name="Tilt Fine"
                unit=" deg"></Slider>
      </div>
      <div class="d-flex flex-column gap">
        <Plot :cols="2" :rows="2" style="width: 13rem" title="Sentry">
          <Confirm :active="state.laser" :disabled="state.laser"
                   :fn="laserToggle"
                   :title="`${state.laser?'DISABLE':'ENABLE'} LASER`"></Confirm>
          <Subplot :active="true" :fn="laserStopAll" name="STOP ALL"
                   theme="danger"></Subplot>
          <Subplot :active="true" :fn="() => laserHome()" name="Home"></Subplot>
          <Subplot :active="true" :fn="() => laserCircle()"
                   name="Run"></Subplot>
          <Subplot :active="true" :fn="() => moveBeamToXYZ(0,0,0)"
                   name="0, 0"></Subplot>
          <Subplot :active="true" :fn="() => laserStop()"
                   :theme="state.runner !== 0?'':'disabled'"
                   name="Halt"></Subplot>
        </Plot>

        <Plot :alt="`${Math.round(state.laserBeam.power * 100) / 100} mW`"
              :cols="3" :rows="1" style="width:13rem;"
              title="Attenuation">
          <Subplot :fn="() => {}" active name="􀅽"></Subplot>
          <Subplot :fn="() => {}" name="5 mW"></Subplot>
          <Subplot :fn="() => {}" active name="􀅼"></Subplot>
        </Plot>


        <Attenuation :duty="4096" :frequency="4096"
                     :percent="state.laserBeam.power/15.0"
                     :scale="1000"></Attenuation>
      </div>

    </div>

    <div>

    </div>

  </div>
  <!--  <div v-else>-->
  <!--    <DefenseAuth></DefenseAuth>-->
  <!--  </div>-->
</template>

<style lang="scss" scoped>
.emergency-stop {

}
</style>
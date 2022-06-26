<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, reactive} from "vue";

import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import {OrbitControls} from "three/examples/jsm/controls/OrbitControls";


let state = reactive({
  position: {
    x: 0, y: 0, z: 0
  },
  manipulation: {
    p: 0, t: 0
  }
})

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let scene = {} as THREE.Scene
let controls = {} as OrbitControls
let beamLine = {} as THREE.Line

interface RoomDefinition {
  offsets: {
    x: number,
    y: number,
    z: number
  },
  rotations: {
    x: number,
    y: number,
    z: number,
  },
  points: THREE.Vector2[]
}

const livingRoomSet = {
  offsets: {
    x: 0,
    y: 0,
    z: -1.2,
  },
  rotations: {
    x: 0,
    y: 0,
    z: 0,
  },
  points: [
    new THREE.Vector2(0.0001, 0), // Top Left
    new THREE.Vector2(3.66, 0), // Top Right
    new THREE.Vector2(3.66, 2.56),// Right Corner
    new THREE.Vector2(3.2, 3.02), // Right Short wall
    new THREE.Vector2(5.363, 5.183), // Right Long wall
    new THREE.Vector2(4.224, 6.32127), // Front Door
    new THREE.Vector2(3.707, 5.805), // Pantry Door Side
    new THREE.Vector2(3.021, 6.491), // Pantry Face Side
    new THREE.Vector2(3.537, 7.007), // Pantry Fridge Side
    new THREE.Vector2(1.826, 8.7181), // Kitchen Far wall
    new THREE.Vector2(-0.3374, 6.554), // Kitchen Left wall
    new THREE.Vector2(0.285, 5.9317), // Kitchen Laundry Wall
    new THREE.Vector2(-1.115, 4.5316), // Laundry Wall
    new THREE.Vector2(0.3133, 3.103), // Room Wall
    new THREE.Vector2(0.0001, 2.794), // Left Short Wall

  ]
}

const bedroomSet = {
  offsets: {
    x: -2.02 - 0.0565,
    y: 0,
    z: 0,
  },
  rotations: {
    x: Math.PI,
    y: 0,
    z: Math.PI / 4,
  },
  points: [
    new THREE.Vector2(0.0001, 0),
    new THREE.Vector2(2.02, 0),
    new THREE.Vector2(4.03, -2.07),
    new THREE.Vector2(4.37, -2.07),
    new THREE.Vector2(4.37, -3.17),
    new THREE.Vector2(2.33, -3.17),
    new THREE.Vector2(1.17, -3.17),
    new THREE.Vector2(0, -3.17),
  ],

}


let s = 22;

onMounted(() => {
  loadThree()
})


function drawWall(points: THREE.Vector2[], def: RoomDefinition) {
  const wall1 = new THREE.Shape();

  wall1.setFromPoints(points)


  const wallGeometry = new THREE.ExtrudeBufferGeometry([wall1], {
    depth: 1.2, bevelEnabled: false
  });


  wallGeometry.scale(s, s, s)

  wallGeometry.translate(def.offsets.x * s, def.offsets.y * s, def.offsets.z * s)
  wallGeometry.rotateX(def.rotations.x)
  wallGeometry.rotateY(def.rotations.y)
  wallGeometry.rotateZ(def.rotations.z)
  const wallMesh = new THREE.Mesh(wallGeometry, new THREE.MeshStandardMaterial({
    color: 0x424850,
    opacity: 0.8,
    transparent: true
  }));
  scene.add(wallMesh);
}

function drawWalls(def: RoomDefinition) {
  for (let i = 0; i < def.points.length; i++) {
    drawWall([def.points[i], def.points[(i + 1) % def.points.length]], def)
  }
}

function drawFloor(def: RoomDefinition) {
  const floor = new THREE.Shape();

  floor.setFromPoints(def.points)

  const floorGeometry = new THREE.ExtrudeBufferGeometry([floor], {
    depth: 0.05, bevelEnabled: false, bevelSize: 0.05, bevelThickness: 0.05

  });

  floorGeometry.scale(s, s, s)
  floorGeometry.translate(def.offsets.x * s, def.offsets.y * s, 0)
  floorGeometry.rotateX(def.rotations.x)
  floorGeometry.rotateY(def.rotations.y)
  floorGeometry.rotateZ(def.rotations.z)
  const floorMesh = new THREE.Mesh(floorGeometry, new THREE.MeshStandardMaterial({
    color: 0x171D24,
    opacity: 0.9,
    transparent: true
  }));

  scene.add(floorMesh);


  // const gridHelper = new THREE.GridHelper(100, 60, 0x8C929B, 0x8C929B);
  //
  // gridHelper.rotateX(Math.PI / 2)
  // gridHelper.translateY(0.06 * s)
  // scene.add(gridHelper);

}


function setCamera(x: number, y: number, z: number) {
  camera.position.set(x, y, z);
  camera.lookAt(0, 0, 0);

  // controls = new OrbitControls(camera, renderer.domElement);
  //
  // controls.enableRotate = true
  // controls.enableDamping = true
  // controls.enableZoom = true


  //controls.update() must be called after any manual changes to the camera's transform

  // controls.update();
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function render() {
  renderer.setClearColor(0x000000, 0);
  renderer.render(scene, camera);
}

let i = 0;
let point = 0;
let steps = 100;

let last = performance.now();

function animate() {

  requestAnimationFrame(animate);
  // if (performance.now() - last > 16.66) {
  //   last = performance.now()
  //   i = (i + 1) % steps
  //   if (i == 0) {
  //     point = (point + 1) % pointCount
  //   }
  // }
  // required if controls.enableDamping or controls.autoRotate are set to true
  controls.update();
  // let x =
  // let y =,
  // let tx = i, ty = 0
  // let height = 1
  //
  // let slice = (Math.PI * 2) / steps
  // let x = map_range(i, 0, steps, points[point].x, points[(point + 1) % pointCount].x)
  // let y = map_range(i, 0, steps, points[point].y, points[(point + 1) % pointCount].y)
  //
  // // moveBeam(180 * (Math.PI / 180), 180 * (Math.PI / 180), 1)
  // moveBeamToXYZ(-4.27 + x, 2.62 + y, 0.01)
  // moveBeam(148.5 * (Math.PI / 180), 168.5 * (Math.PI / 180))
  // beamLine.rotateOnAxis(new THREE.Vector3(4.37, -2.62, 1), Math.random())
  // drawScene()
  render()

}

function moveBeamToXYZ(x: number, y: number, z: number) {
  state.position.x = x
  state.position.y = y
  state.position.z = z
  let distance = Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2) + Math.pow(z, 2))
  let theta = Math.atan(y / x) + (x >= 0 ? 0 : Math.PI)
  let phi = Math.atan(Math.sqrt(Math.pow(x, 2) + Math.pow(y, 2)) / z)

  state.manipulation.t = phi
  state.manipulation.p = theta
  moveBeam(theta, phi, distance)

}

function moveBeam(pan: number, tilt: number, distance: number) {

  let theta = pan
  let phi = tilt

  let x = distance * Math.cos(theta) * Math.sin(phi)
  let y = distance * Math.sin(theta) * Math.sin(phi)
  let z = distance * Math.cos(phi)

  let newGeom = new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(4.27, -2.62, 1), new THREE.Vector3(4.27 + x, -2.62 + y, z)])

  newGeom.scale(s, s, s)
  newGeom.translate(-4.37 * s / 2, 3.17 * s / 2, 0.05 * s)
  beamLine.geometry.dispose()
  beamLine.geometry = newGeom
}


function drawBeam(x1: number, y1: number, z1: number, x2: number, y2: number, z2: number) {
  const beamMaterial = new THREE.LineBasicMaterial({
    color: 0xff0000,
    linewidth: 1,
    linecap: 'round', //ignored by WebGLRenderer
    linejoin: 'round' //ignored by WebGLRenderer
  })

  const beam = new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(x1, y1, z1), new THREE.Vector3(x2, y2, z2)])
  beam.scale(s, s, s)
  beam.translate(-4.37 * s / 2, 3.17 * s / 2, 4)

  beamLine = new THREE.Line(beam, beamMaterial);
  scene.add(beamLine)

  // const radius = 120;
  // const radials = 64;
  // const circles = 24;
  // const divisions = 64;
  //
  // const grid = new THREE.GridHelper(100, 20)
  // grid.rotateX(Math.PI / 2)
  // scene.add(grid)

}

function drawScene(def: RoomDefinition) {
  drawFloor(def)
  drawWalls(def)
  // 23 29 36

}

function loadThree() {
  renderer = new THREE.WebGLRenderer();
  let element = document.getElementById("three-container")
  if (!element) return

  renderer.setSize(element.clientWidth, element.clientHeight);
  renderer.setPixelRatio(window.devicePixelRatio);

  element.appendChild(renderer.domElement)

  camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 10000);
  controls = new OrbitControls(camera, renderer.domElement);
  controls.enableDamping = true

  setCamera(0, 0, -400)

  scene = new THREE.Scene();


  // const light = new THREE.PointLight(0xFFFFFF, 0.6, 100);
  // light.position.set(6, 14, 14);
  // scene.add(light);
  //
  // const light2 = new THREE.PointLight(0xFFFFFF, 0.6, 100);
  // light2.position.set(22, 3, 14);
  // scene.add(light2);
  //
  // const light3 = new THREE.PointLight(0xFFFFFF, 0.6, 100);
  // light3.position.set(-42, -7, 14);
  // scene.add(light3);
  const axesHelper = new THREE.AxesHelper(5);
  scene.add(axesHelper);

  drawScene(livingRoomSet)
  drawScene(bedroomSet)
  // drawScene(pointsd)
  scene.add(new THREE.HemisphereLight(0xffffff, 0xcccccc, 3))
  // drawBeam(4.37, -2.62, 1, 4.37, -2.62, -0.12)
  animate()
  render()

}

</script>

<template>
  <div class="w-100 h-100 mb-5 d-flex pb-3 pt-1 gap">
    <div class="d-flex flex-column gap">
      <Plot :cols="1" :rows="4" style="width:13rem;" title="Room">
        <Subplot :active="true" :fn="() => {}" name="Bedroom"></Subplot>
        <Subplot :active="false" :fn="() => {}" name="Living Room"></Subplot>
      </Plot>
      <Plot :cols="1" :rows="4" style="width:13rem;" title="Location">
        <div class="d-flex justify-content-evenly">
          <div class="label-w600 label-r label-o4 label-c1" style="width: 3rem">
            X: {{ Math.round(state.position.x * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c1" style="width: 3rem">
            Y: {{ Math.round(state.position.y * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c1" style="width: 3rem">
            Z: {{ Math.round(state.position.z * 10) / 10 }}
          </div>
        </div>
        <div class="d-flex justify-content-evenly">
          <div class="label-w600 label-r label-o4 label-c1" style="width: 4rem">
            Pan: {{ Math.round(state.manipulation.p * (180 / Math.PI) * 10) / 10 }}
          </div>
          <div class="label-w600 label-r label-o4 label-c1" style="width: 4rem">
            Tilt: {{ Math.round(state.manipulation.t * (180 / Math.PI) * 10) / 10 }}
          </div>

        </div>
      </Plot>
    </div>
    <div id="three-container" class="flex-grow-1 h-100 mb-5 element"></div>
  </div>
</template>

<style lang="scss" scoped>
</style>
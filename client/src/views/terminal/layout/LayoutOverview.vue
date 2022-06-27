<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, onUnmounted, reactive} from "vue";

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
let camera = {} as THREE.OrthographicCamera
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
  floor: {
    diffuse: string,
    roughness: string,
    rotation: number,
    scale: number,
  }
  points: THREE.Vector2[]
}

const livingRoomSet = {
  offsets: {
    x: 0,
    y: 0,
    z: -0.05,
  },
  floor: {
    diffuse: "/custom/textures/woodfloor-color.jpg",
    roughness: "/custom/textures/woodfloor-roughness.jpg",
    rotation: Math.PI * 3 / 4,
    scale: 0.5,
  },
  rotations: {
    x: 0,
    y: Math.PI,
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
    x: 0.11,
    y: 0,
    z: 0,
  },
  floor: {
    diffuse: "/custom/textures/carpet-color.jpg",
    roughness: "/custom/textures/carpet-roughness.jpg",
    rotation: Math.PI * 3 / 4,
    scale: 1,
  },
  rotations: {
    x: 0,
    y: 0,
    z: 0,
  },
  points: [
    new THREE.Vector2(1.4283, -1.4283),
    new THREE.Vector2(0.0001, 0), // Zero Point (Between room and living room/patio)
    new THREE.Vector2(0, 2.88),
    new THREE.Vector2(-0.2404, 3.1204),
    new THREE.Vector2(0.5374, 3.8982),
    new THREE.Vector2(3.6275, 0.8081),

  ],

}


let s = 100;

onUnmounted(() => {
  renderer.dispose()
})

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

function drawFloor(def: RoomDefinition): THREE.Object3D {
  const floor = new THREE.Shape();

  const object = new THREE.Object3D();

  floor.setFromPoints(def.points)

  const floorGeometry = new THREE.ExtrudeBufferGeometry([floor], {
    depth: 0.05, bevelEnabled: false, bevelSize: 0.05, bevelThickness: 0.05

  });

  floorGeometry.scale(s, s, s)
  floorGeometry.translate(def.offsets.x * s, def.offsets.y * s, def.offsets.z * s)
  floorGeometry.rotateX(def.rotations.x)
  floorGeometry.rotateY(def.rotations.y)
  floorGeometry.rotateZ(def.rotations.z)

  let roughness = new THREE.TextureLoader().load(def.floor.roughness)
  roughness.wrapT = THREE.RepeatWrapping
  roughness.wrapS = THREE.RepeatWrapping
  roughness.repeat.set(def.floor.scale, def.floor.scale)
  roughness.rotation = def.floor.rotation

  let map = new THREE.TextureLoader().load(def.floor.diffuse)
  map.wrapT = THREE.RepeatWrapping
  map.wrapS = THREE.RepeatWrapping
  map.repeat.set(def.floor.scale, def.floor.scale)
  map.rotation = def.floor.rotation

  let floorMaterial = new THREE.MeshPhysicalMaterial({
    color: 0x424850,
    opacity: 1,
    roughness: 1,
    transparent: false,
    roughnessMap: roughness,
    map: map
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


function setCamera(x: number, y: number, z: number) {
  camera.position.set(x, y, z);

  camera.lookAt(0, 0, 0);

  controls = new OrbitControls(camera, renderer.domElement);
  //
  controls.enableRotate = true
  controls.enableDamping = true
  controls.enableZoom = true

  controls.update();

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
  // controls.update();
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
  controls.update()
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

  scene.rotateX(Math.PI / 2)

  // const radius = 120;
  // const radials = 64;
  // const circles = 24;
  // const divisions = 64;
  //
  // const grid = new THREE.GridHelper(100, 20)
  // grid.rotateX(Math.PI / 2)
  // scene.add(grid)

}

function drawScene(def: RoomDefinition): THREE.Object3D {
  return drawFloor(def)
  // drawWalls(def)
  // 23 29 36
}

function loadThree() {
  renderer = new THREE.WebGLRenderer();
  renderer.shadowMap.enabled = true;
  let element = document.getElementById("three-container")
  if (!element) return

  renderer.setSize(element.clientWidth, element.clientHeight);
  renderer.setPixelRatio(window.devicePixelRatio);

  element.appendChild(renderer.domElement)
  let width = window.innerWidth
  let height = window.innerHeight
  camera = new THREE.OrthographicCamera(width / -2, width / 2, height / 2, height / -2, -1000, 1000);
  // camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 10000);
  // controls = new OrbitControls(camera, renderer.domElement);
  // controls.enableDamping = true

  setCamera(0, -1, 1)

  scene = new THREE.Scene();

  const axesHelper = new THREE.AxesHelper(s);
  axesHelper.setColors(new THREE.Color(255, 0, 0), new THREE.Color(0, 255, 0), new THREE.Color(0, 0, 255))
  scene.add(axesHelper);

  const gridHelper = new THREE.GridHelper(s * 10, s / 2)
  gridHelper.rotateX(Math.PI / 2)
  scene.add(gridHelper)

  const livingRoomMesh = drawScene(livingRoomSet)
  const bedroomMesh = drawScene(bedroomSet)

  let combo = new THREE.Object3D().add(livingRoomMesh, bedroomMesh)
  combo.rotateZ(Math.PI / 2)
  combo.translateX(s)
  combo.translateY(-s * 4)

  // drawScene(pointsd)
  scene.add(new THREE.HemisphereLight(0xffffff, 0xcccccc, 2.25))
  scene.add(combo)


  const light = new THREE.PointLight(0xffffff, 1, 100);
  light.position.set(50, 50, 25);
  scene.add(light);

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
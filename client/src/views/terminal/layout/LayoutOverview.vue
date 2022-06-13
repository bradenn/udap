<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import * as THREE from "three";
import {onMounted} from "vue";

import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import {OrbitControls} from "three/examples/jsm/controls/OrbitControls";


let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let scene = {} as THREE.Scene
let controls = {} as OrbitControls
let beamLine = {} as THREE.Line

const points = [
  new THREE.Vector2(0, 0),
  new THREE.Vector2(2.02, 0),
  new THREE.Vector2(4.03, -2.07),
  new THREE.Vector2(4.37, -2.07),
  new THREE.Vector2(4.37, -3.17),
  new THREE.Vector2(2.33, -3.17),
  new THREE.Vector2(1.17, -3.17),
  new THREE.Vector2(0, -3.17),
]


let s = 22;

onMounted(() => {
  loadThree()
})


function drawWall(points: THREE.Vector2[]) {
  const wall1 = new THREE.Shape();

  points.forEach(point => {

  })

  wall1.setFromPoints(points)


  const wallGeometry = new THREE.ExtrudeBufferGeometry([wall1], {
    depth: 1.2, bevelEnabled: false
  });

  wallGeometry.scale(s, s, s)
  wallGeometry.translate(-4.37 * s / 2, 3.17 * s / 2, 0.05 * s)
  const wallMesh = new THREE.Mesh(wallGeometry, new THREE.MeshStandardMaterial({
    color: 0x424850,
    opacity: 0.8,
    transparent: true
  }));
  scene.add(wallMesh);
}

function drawWalls() {
  for (let i = 0; i < points.length; i++) {
    drawWall([points[i], points[(i + 1) % points.length]])
  }
}

function drawFloor() {
  const floor = new THREE.Shape();

  floor.setFromPoints(points)

  const floorGeometry = new THREE.ExtrudeBufferGeometry([floor], {
    depth: 0.05, bevelEnabled: false, bevelSize: 0.05, bevelThickness: 0.05

  });

  floorGeometry.scale(s, s, s)
  floorGeometry.translate(-4.37 * s / 2, 3.17 * s / 2, 0)
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

function render() {
  renderer.setClearColor(0x000000, 0);
  renderer.render(scene, camera);
}

let i = 0;
let steps = 200;
let last = performance.now();

function animate() {

  requestAnimationFrame(animate);
  if (performance.now() - last > 16.66) {
    last = performance.now()
    i = (i + 1) % steps
  }
  // required if controls.enableDamping or controls.autoRotate are set to true
  controls.update();
  // let x =
  // let y =,
  let tx = i, ty = 0
  let height = 1

  let slice = (Math.PI * 2) / steps
  moveBeam(180 * (Math.PI / 180), 180 * (Math.PI / 180))
  // moveBeamToXYZ(i, 0, -0.12)
  // moveBeam(148.5 * (Math.PI / 180), 168.5 * (Math.PI / 180))
  // beamLine.rotateOnAxis(new THREE.Vector3(4.37, -2.62, 1), Math.random())
  // drawScene()
  render()

}

function moveBeamToXYZ(x: number, y: number, z: number) {
  // let pan = Math.sin(x) * 10

}

function moveBeam(pan: number, tilt: number) {
  //
  // if (pan * (Math.PI / 180) > 90) {
  //   return
  // }
  //
  // if (tilt * (Math.PI / 180) > 90) {
  //   return
  // }
  let x1 = 20 * Math.sin(tilt) * Math.cos(pan)
  let y1 = 20 * Math.sin(tilt) * Math.sin(pan)
  let z1 = 20 * Math.cos(tilt)

  let radius = Math.sqrt(Math.pow(x1, 2) + Math.pow(y1, 2) + Math.pow(z1, 2))

  let phi = Math.atan2(y1, x1)
  let theta = Math.acos((z1 / radius))

  let pr = tilt
  let tr = pan

  let x = 20 * Math.sin(pr) * Math.cos(tr)
  let y = 20 * Math.sin(pr) * Math.sin(tr)
  let z = 20 * Math.cos(pr)


  let newGeom = new THREE.BufferGeometry().setFromPoints([new THREE.Vector3(4.37 - 0.17, -2.62, 1), new THREE.Vector3(4.37 - 0.17 + x, -2.62 + y, z)])
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

  const radius = 120;
  const radials = 64;
  const circles = 24;
  const divisions = 64;

  const helper = new THREE.PolarGridHelper(radius, radials, circles, divisions);
  //4.37 - 0.17, -2.62, 1
  helper.translateX(-4.37 * s / 2 + (4.37 - 0.17) * s)
  helper.translateY(3.17 * s / 2 + (-2.62) * s)
  helper.translateZ(0.12 * s / 2)
  helper.rotateX(Math.PI / 2)
  scene.add(helper);

}

function drawScene() {
  drawFloor()
  drawWalls()
  // 23 29 36
  scene.add(new THREE.HemisphereLight(0xffffff, 0xcccccc, 0.8))
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

  setCamera(0, -60, 100)

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

  drawScene()
  drawBeam(4.37, -2.62, 1, 4.37, -2.62, -0.12)
  animate()
  render()

}

</script>

<template>
  <div class="w-100 h-100 mb-5 d-flex pb-3 pt-1 gap">
    <div>
      <Plot :cols="1" :rows="4" style="width:13rem;" title="Room">
        <Subplot :active="true" :fn="() => {}" name="Bedroom"></Subplot>
        <Subplot :active="false" :fn="() => {}" name="Living Room"></Subplot>
      </Plot>
    </div>
    <div id="three-container" class="flex-grow-1 h-100 mb-5 element"></div>
  </div>
</template>

<style lang="scss" scoped>
</style>
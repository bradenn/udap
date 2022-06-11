<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import * as THREE from "three";
import {onMounted} from "vue";

import Plot from "@/components/plot/Plot.vue";


let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let scene = {} as THREE.Scene
const points = [
  new THREE.Vector2(0, 0),
  new THREE.Vector2(2.02, 0),
  new THREE.Vector2(2.01 + 2.02, -2.07),
  new THREE.Vector2(2.01 + 2.02 + 0.34, -2.07),
  new THREE.Vector2(2.01 + 2.02 + 0.34, -3.17),
  new THREE.Vector2(2.33, -3.17),
  new THREE.Vector2(1.17, -3.17),
  new THREE.Vector2(0, -3.17),
]


let s = 22;

onMounted(() => {
  loadThree()
})

function render() {
  renderer.setClearColor(0x000000, 0);
  renderer.render(scene, camera);
}

function drawWall(points: THREE.Vector2[]) {
  const wall1 = new THREE.Shape();

  wall1.setFromPoints(points)


  const wallGeometry = new THREE.ExtrudeBufferGeometry([wall1], {
    depth: 1.2, bevelEnabled: false
  });

  wallGeometry.scale(s, s, s)
  wallGeometry.translate(-4.37 * s / 2, 3.17 * s / 2, 0.05 * s)
  const wallMesh = new THREE.Mesh(wallGeometry, new THREE.MeshStandardMaterial({
    color: 0xcccccc,
    opacity: 0.8,
    transparent: true
  }));
  scene.add(wallMesh);
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
  camera.position.set(0, 0, 120);

  // controls.update();
}

function addWindow(x: number, y: number, z: number, x2: number, y2: number, z2: number) {

  const wall1 = new THREE.Shape();

  wall1.setFromPoints([points[0].add(new THREE.Vector2(0.125, -0.025)), points[1].add(new THREE.Vector2(-0.125, -0.025))])
  const wallGeometry = new THREE.ExtrudeBufferGeometry([wall1], {
    depth: 0.8, bevelEnabled: false
  });

  wallGeometry.scale(s, s, s)
  wallGeometry.translate(-4.37 * s / 2, 3.17 * s / 2, 4)
  const wallMesh = new THREE.Mesh(wallGeometry, new THREE.MeshStandardMaterial({
    color: 0xffffff,
    opacity: 0.5,
    transparent: true
  }));
  scene.add(wallMesh);

  const rectLight = new THREE.RectAreaLight(0xffffff, 2, 1.5 * s, 1.5 * s);
  rectLight.position.set(x, y, z);
  rectLight.lookAt(x2, y2, z2);
  scene.add(rectLight)

}

function animate() {

  requestAnimationFrame(animate);

  // required if controls.enableDamping or controls.autoRotate are set to true
  // controls.update();

  renderer.render(scene, camera);

}

function loadThree() {
  renderer = new THREE.WebGLRenderer();
  let element = document.getElementById("three-container")
  if (!element) return
  renderer.setSize(element.clientWidth, element.clientHeight);
  element.appendChild(renderer.domElement)
  renderer.setPixelRatio(window.devicePixelRatio);
  camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 10000);
  setCamera(0, 0, 0)

  scene = new THREE.Scene();

  const floor = new THREE.Shape();


  floor.setFromPoints(points)

  const floorGeometry = new THREE.ExtrudeBufferGeometry([floor], {
    depth: 0.05, bevelEnabled: false, bevelSize: 0.05, bevelThickness: 0.05

  });
  floorGeometry.scale(s, s, s)
  floorGeometry.translate(-4.37 * s / 2, 3.17 * s / 2, 0)
  const floorMesh = new THREE.Mesh(floorGeometry, new THREE.MeshStandardMaterial({
    color: 0xD2B48C,
    opacity: 0.8,
    transparent: true
  }));

  scene.add(floorMesh);


  drawWall([points[0], points[1]])
  drawWall([points[1], points[2]])
  drawWall([points[2], points[3]])
  drawWall([points[3], points[4]])
  drawWall([points[4], points[5]])
  drawWall([points[5], points[6]])
  drawWall([points[6], points[7]])
  drawWall([points[7], points[0]])

  const light = new THREE.PointLight(0xFFFFFF, 0.5, 50);
  light.position.set(7, 10, 10);
  scene.add(light);

  const light2 = new THREE.PointLight(0xFFFFFF, 0.5, 50);
  light2.position.set(20, 2, 10);
  scene.add(light2);


  addWindow(-23, 31, s, -23, 0, s)
  scene.add(new THREE.HemisphereLight(0xffffffcc, 0xcccccccc, 0.2))
  // animate()
  render()

}

</script>

<template>
  <div class="w-100 h-100 mb-5 d-flex pb-3 pt-1 gap">
    <div>
      <Plot :cols="1" :rows="4" style="width:13rem;">

      </Plot>
    </div>
    <div id="three-container" class="flex-grow-1 h-100 mb-5 element"></div>
  </div>
</template>

<style lang="scss" scoped>
</style>
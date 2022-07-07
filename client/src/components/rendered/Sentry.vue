<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, onUnmounted, reactive} from "vue";
import Plot from "@/components/plot/Plot.vue";

interface SentryVisual {
  pan: number,
  tilt: number,
  beam: boolean,
  color: string,
}

let props = defineProps<SentryVisual>()

let state = reactive<SentryVisual>({
  pan: 0,
  tilt: 0,
  beam: false,
  color: "string",
})

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let scene = {} as THREE.Scene
let panObject = {} as THREE.Object3D
let tiltObject = {} as THREE.Object3D
let beamHead = {} as THREE.Object3D

onMounted(() => {
  state.pan = props.pan
  state.tilt = props.tilt
  load3d()
})

onUnmounted(() => {
  renderer.dispose()
})


function setCamera(x: number, y: number, z: number) {
  camera.position.set(x, y, z);
  camera.lookAt(0, 3.5 / 2, 0);
}

function render() {
  renderer.setClearColor(0x000000, 0);

  renderer.render(scene, camera);
}

function animate() {
  requestAnimationFrame(animate);

  if (props.pan != state.pan) {
    state.pan = (state.pan + props.pan) / 2
  }
  if (props.tilt != state.tilt) {
    state.tilt = (state.tilt + props.tilt) / 2
  }

  panObject.setRotationFromEuler(new THREE.Euler(0, -state.pan * (Math.PI / 180) - Math.PI, 0, 'XYZ'))
  tiltObject.setRotationFromEuler(new THREE.Euler(0, -state.pan * (Math.PI / 180) - Math.PI, -state.tilt * (Math.PI / 180) + Math.PI / 2, 'XYZ'))
  render()
}

function drawSentry() {
  panObject = new THREE.Object3D()
  // let panGrid = new THREE.PolarGridHelper(5, 16, 6, 32, new THREE.Color(0.2, 0.2, 0.2), new THREE.Color(0.2, 0.3, 0.4))
  // panGrid.translateY(3.5)
  // scene.add(panGrid)
  let panJoint = new THREE.CylinderGeometry(5, 5, 1, 32, 2)
  let matte = new THREE.MeshPhysicalMaterial({
    color: 0x202020,
    opacity: 1,
    roughness: 0.5,
    transparent: false,
  });
  let tiltArm = new THREE.BoxGeometry(1, 6, 1)
  tiltArm.translate(0, 1.5, 3.5)
  const tiltMesh = new THREE.Mesh(tiltArm, matte)
  panObject.add(tiltMesh)

  const panJointMesh = new THREE.Mesh(panJoint, matte)
  panJointMesh.translateY(4)
  panObject.add(panJointMesh)
  tiltObject = new THREE.Object3D()
  let brass = new THREE.MeshPhysicalMaterial({
    color: 0x57450b,
    opacity: 1,
    roughness: 1,
    transparent: false,
  });
  // let tiltGrid = new THREE.PolarGridHelper(2, 16, 6, 32, new THREE.Color(0.2, 0.2, 0.2), new THREE.Color(0.2, 0.3, 0.4))
  // tiltGrid.translateZ(3.05)
  // tiltGrid.rotateZ(Math.PI / 2)
  // tiltGrid.rotateX(Math.PI / 2)
  // tiltObject.add(tiltGrid)
  let tiltBarrel = new THREE.CylinderGeometry(2, 2, 6, 32, 4)

  tiltBarrel.rotateX(Math.PI / 2)
  let beamHeadMaterial = new THREE.MeshPhysicalMaterial({
    color: 0xFF0000,
    opacity: 1,
    roughness: 1,
    transparent: false,
  });
  let beam = new THREE.CylinderGeometry(0.5, 0.5, 2, 12, 2)
  beam.rotateZ(Math.PI / 2)
  beam.rotateX(-Math.PI / 2)
  beam.translate(-1.2, 0, 0)

  let beamHead = new THREE.CylinderGeometry(0.125, 0.5, 2, 12, 2)
  beamHead.rotateZ(Math.PI / 2)
  beamHead.rotateX(-Math.PI / 2)
  beamHead.translate(-1.25, 0, 0)
  const tiltBarrelMesh = new THREE.Mesh(tiltBarrel, matte)
  const beamHeadMesh = new THREE.Mesh(beamHead, beamHeadMaterial)
  const beamMesh = new THREE.Mesh(beam, brass)

  tiltObject.add(beamHeadMesh)
  tiltObject.add(tiltBarrelMesh)
  tiltObject.add(beamMesh)
  scene.add(panObject)
  scene.add(tiltObject)
}

function load3d() {
  renderer = new THREE.WebGLRenderer();
  renderer.shadowMap.enabled = true;
  let element = document.getElementById("sentry-canvas")
  if (!element) return


  renderer.setSize(element.clientWidth, element.clientHeight);
  renderer.setPixelRatio(window.devicePixelRatio);

  element.appendChild(renderer.domElement)
  let width = element.clientWidth
  let height = element.clientHeight

  // camera = new THREE.OrthographicCamera(width / -2, width / 2, height / 2, height / -2, -1000, 1000);
  camera = new THREE.PerspectiveCamera(25, width / height, 1, 400);

  setCamera(0, 3.5 / 8, 25)


  scene = new THREE.Scene();


  // const axesHelper = new THREE.AxesHelper(10);
  // axesHelper.setColors(new THREE.Color(255, 0, 0), new THREE.Color(0, 255, 0), new THREE.Color(0, 0, 255))
  // scene.add(axesHelper);


  drawSentry()
  scene.add(new THREE.HemisphereLight(0xffffff, 0x000000, 1.2))

  animate()

  // const gridHelper = new THREE.GridHelper(s * 10, s / 2)
}
</script>

<template>
  <Plot :alt="`Pan ${props.pan}°, Tilt ${props.tilt}°&nbsp;`" :cols="1" :rows="1" title="Simulated">
    <div id="sentry-canvas" class="sentry-box">
    </div>
  </Plot>
</template>

<style lang="scss" scoped>
.sentry-box {
  width: 13rem;
  height: 8rem;
}
</style>
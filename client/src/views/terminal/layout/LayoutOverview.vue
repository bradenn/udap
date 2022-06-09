<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import * as THREE from "three";
import {onMounted} from "vue";

onMounted(() => {
  drawThree()
})

function drawThree() {

  const renderer = new THREE.WebGLRenderer();
  let element = document.getElementById("three-container")
  if (!element) return
  renderer.setSize(element.clientWidth, element.clientHeight);
  element.appendChild(renderer.domElement)

  const camera = new THREE.PerspectiveCamera(45, window.innerWidth / window.innerHeight, 1, 500);
  camera.position.set(0, 0, 100);
  camera.lookAt(0, 0, 0);


  const scene = new THREE.Scene();

  const roomShape = new THREE.Shape();
  roomShape.moveTo(0, 0)
  roomShape.lineTo(10, 10)

  const back = new THREE.GridHelper(50, 25, 0xFF5F1F0F, 0xFF5F1F0F);
  back.rotateX(Math.PI / 2)

  const floor = new THREE.Shape();

  floor.moveTo(0, 0);
  floor.lineTo(20, 0);
  floor.lineTo(20, 20);
  floor.lineTo(0, 20);


  const shapeGeom = new THREE.ShapeGeometry(floor);

  const wireframe = new THREE.WireframeGeometry(shapeGeom);


  const line = new THREE.LineSegments(wireframe);


  scene.add(line);
  scene.add(back);
  renderer.render(scene, camera);
}

</script>

<template>
  <div id="three-container" class="w-100 h-100 mb-5">d</div>
</template>

<style lang="scss" scoped>
</style>
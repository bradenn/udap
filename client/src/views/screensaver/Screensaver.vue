<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, onUnmounted} from "vue";
import Stats from 'stats.js';
import {RenderPass} from "three/examples/jsm/postprocessing/RenderPass";
import {UnrealBloomPass} from "three/examples/jsm/postprocessing/UnrealBloomPass";
import {EffectComposer} from "three/examples/jsm/postprocessing/EffectComposer";

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let composer = {} as EffectComposer
let stats = {} as Stats
let scene = {} as THREE.Scene
let instanceMaterial = {} as THREE.RawShaderMaterial
let particleGeometry = {} as THREE.BufferGeometry
let particleSystem = {} as THREE.Points

onMounted(() => {
  initGraphics()
})

onUnmounted(() => {
  renderer.dispose()
  instanceMaterial.dispose()
  particleGeometry.dispose()
})

let width = 0
let height = 0


function initCamera() {
  // Initialize Camera
  camera = new THREE.PerspectiveCamera(40, width / height, 1, 3000)
  camera.position.set(0, 0, 1);
}

// Handler resizing the viewport if and when the viewport is altered
function resizeFrame() {
  camera.aspect = width / height;
  camera.updateProjectionMatrix();

  renderer.setSize(width, height);
}

// Initialize the THREE graphics and set up the environment
function initGraphics() {
  // Initialize Renderer
  renderer = new THREE.WebGLRenderer({antialias: false});
  // Select screensaver dom element
  let element = document.getElementById("screensaver")
  // Exit init if dom element is not found
  if (!element) return
  // Define local context variables
  width = element.clientWidth
  height = element.clientHeight
  // Set the renderer size
  renderer.setSize(width, height);
  // Set the pixel ratio
  renderer.setPixelRatio(window.devicePixelRatio);
  // Set output encoding
  renderer.outputEncoding = THREE.sRGBEncoding;
  // Add the renderer to the dom
  element.appendChild(renderer.domElement)
  // Assign resize function
  window.onresize = resizeFrame
  // Initialize Stats manager
  stats = new Stats();
  element.appendChild(stats.dom)
  // Set up the camera
  initCamera()
  // Set up the scene
  initScene()
  const renderScene = new RenderPass(scene, camera);
  const bloomPass = new UnrealBloomPass(new THREE.Vector2(width, height), 1.2, 0.8, 0.3);
  composer = new EffectComposer(renderer);
  composer.addPass(renderScene);
  composer.addPass(bloomPass);
  // Begin animation and rendering
  animate()


}

function initScene() {
  scene = new THREE.Scene()
  let grid = new THREE.AxesHelper(10)
  scene.add(grid)
  generateStars()
  scene.add(new THREE.AmbientLight(0x444444, 4));
  scene.background = new THREE.Color(0x000000);
  scene.fog = new THREE.Fog(0x050505, 1, 2000);
}

// Define animation routine
function animate() {
  requestAnimationFrame(animate);
  stats.update()
  render()
  composer.render();
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function cctToRgb(cct: number) {
  return [map_range(cct, 1900, 3400, 255, 255),
    map_range(cct, 1900, 3400, 131, 193),
    map_range(cct, 1900, 3400, 0, 132)]
}

function randomG(v: number) {
  let r = 0;
  for (let i = v; i > 0; i--) {
    r += Math.random();
  }
  return r / v;
}

const particleCount = 2000;
let translateArray = new Float32Array(particleCount * 3);
let colorArray = new Float32Array(particleCount * 3);
let scaleArray = new Float32Array(particleCount);

function generateStars() {

  particleGeometry = new THREE.BufferGeometry();


  let aspectW = (width / height) * 1
  let aspectH = (height / width) * 1

  const radius = 50;
  let mean = 2747.5
  let sd = 1509.878

  let wave = 1;

  for (let i = 0, i3 = 0, l = particleCount; i < l; i++, i3 += 3) {


    translateArray[i3] = (Math.random() * (2) - 1) * radius * aspectW;
    translateArray[i3 + 1] = (Math.random() * (2) - 1) * radius * aspectH;
    translateArray[i3 + 2] = -i * wave;

    let rand = (randomG(8) * 2 - 1) * sd
    let arr = cctToRgb(mean + rand)
    colorArray[i3] = arr[0] / 255;
    colorArray[i3 + 1] = arr[1] / 255
    colorArray[i3 + 2] = arr[2] / 255

    scaleArray[i3] = 2 + randomG(8) * 5;


  }

  particleGeometry.setAttribute('position', new THREE.Float32BufferAttribute(translateArray, 3));
  particleGeometry.setAttribute('color', new THREE.Float32BufferAttribute(colorArray, 3));
  particleGeometry.setAttribute('size', new THREE.Float32BufferAttribute(scaleArray, 1).setUsage(THREE.DynamicDrawUsage));
  let vertexShader = document.getElementById('vshader')
  if (!vertexShader) return
  let fragmentShader = document.getElementById('fshader')
  if (!fragmentShader) return

  instanceMaterial = new THREE.ShaderMaterial({
    uniforms: {
      pointTexture: {value: new THREE.TextureLoader().load('/custom/textures/circle.png')}
    },
    vertexShader: vertexShader.textContent || "",
    fragmentShader: fragmentShader.textContent || "",

    depthTest: false,
    transparent: true,
    vertexColors: true
  });

  particleSystem = new THREE.Points(particleGeometry, instanceMaterial);
  //
  // instanceMesh = new THREE.Mesh(geometry, instanceMaterial);
  // instanceMesh.scale.set(200, 200, 200);
  scene.add(particleSystem);
}

function render() {
  // particleSystem.translateZ(1)

  // console.log(particleSystem.position.z)


}

</script>

<template>

  <div id="screensaver" class="screensaver-context"></div>
</template>

<style scoped>
.screensaver-context {
  width: 100%;
  height: 100%;
}
</style>
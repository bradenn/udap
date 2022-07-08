<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, onUnmounted} from "vue";
import {RenderPass} from "three/examples/jsm/postprocessing/RenderPass";
import {UnrealBloomPass} from "three/examples/jsm/postprocessing/UnrealBloomPass";
import {EffectComposer} from "three/examples/jsm/postprocessing/EffectComposer";
import {SavePass} from "three/examples/jsm/postprocessing/SavePass";
import {ShaderPass} from "three/examples/jsm/postprocessing/ShaderPass";
import {CopyShader} from "three/examples/jsm/shaders/CopyShader";
import {BlendShader} from "three/examples/jsm/shaders/BlendShader";

let renderer = {} as THREE.WebGLRenderer
let camera = {} as THREE.PerspectiveCamera
let composer = {} as EffectComposer
let scene = {} as THREE.Scene
let instanceMaterial = {} as THREE.RawShaderMaterial
let objs = {} as THREE.Object3D
let animationFrame = 0;

onMounted(() => {
  reset()
  initGraphics()
})


onUnmounted(() => {
  instanceMaterial.dispose()
  renderer.dispose()
  cancelAnimationFrame(animationFrame)
  reset()
})

function reset() {
  renderer = {} as THREE.WebGLRenderer
  camera = {} as THREE.PerspectiveCamera
  composer = {} as EffectComposer
  scene = {} as THREE.Scene
  instanceMaterial = {} as THREE.RawShaderMaterial
  objs = {} as THREE.Object3D
}

let width = 0
let height = 0
let depth = 1500


function initCamera() {
  // Initialize Camera
  camera = new THREE.PerspectiveCamera(75, width / height, 1, depth * 1.5)

  camera.position.set(0, 0, 1);
  camera.lookAt(0, 0, 0)
}

// Handler resizing the viewport if and when the viewport is altered
function resizeFrame() {
  camera.aspect = width / height;
  if (!camera) return
  camera.updateProjectionMatrix();

  renderer.setSize(width, height);
}

// Initialize the THREE graphics and set up the environment
function initGraphics() {
  // Initialize Renderer
  renderer = new THREE.WebGLRenderer({antialias: false});
  // renderer.toneMapping = THREE.T;
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
  element.onresize = resizeFrame
  // Initialize Stats manager
  // stats = new Stats();
  // element.appendChild(stats.dom)
  // Set up the camera
  initCamera()
  // Set up the scene
  initScene()
  const renderScene = new RenderPass(scene, camera);
  const bloomPass = new UnrealBloomPass(new THREE.Vector2(width, height), 2, 0.8, 0.3);
  const renderTargetParameters = {
    minFilter: THREE.LinearFilter,
    magFilter: THREE.LinearFilter,
    stencilBuffer: false
  };

// save pass
  const savePass = new SavePass(
      new THREE.WebGLRenderTarget(
          width,
          height,
          renderTargetParameters
      )
  );

// blend pass
  const blendPass = new ShaderPass(BlendShader, "tDiffuse1");
  blendPass.uniforms["tDiffuse2"].value = savePass.renderTarget.texture;
  blendPass.uniforms["mixRatio"].value = 0.2;

// output pass
  const outputPass = new ShaderPass(CopyShader);
  outputPass.renderToScreen = true;
  composer = new EffectComposer(renderer);
  composer.addPass(renderScene);
  composer.addPass(blendPass);
  composer.addPass(savePass);
  composer.addPass(outputPass);
  composer.addPass(bloomPass);
  // Begin animation and rendering
  animate()


}

let numSections = 4;

function initScene() {
  scene = new THREE.Scene()

  objs = new THREE.Object3D()

  for (let i = 0; i <= numSections; i++) {
    let section = generateStars()
    section.position.setZ(-depth * i)
    objs.add(section)
  }

  scene.add(objs)

  // scene.add(new THREE.AmbientLight(0x444444, 4));
  scene.background = new THREE.Color(0x000000);

}

// Define animation routine
function animate() {
  animationFrame = requestAnimationFrame(animate);
  render()
  // stats.update()
  composer.render();
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
  return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}

function cctToRgb(cct: number) {
  return [map_range(cct, 1900, 3400, 255, 255),
    map_range(cct, 1900, 3400, 131, 193),
    map_range(cct, 1900, 3400, 0, 190)]
}

function randomG(v: number) {
  let r = 0;
  for (let i = v; i > 0; i--) {
    r += Math.random();
  }
  return r / v;
}

let clr = 0;

function generateStars(): THREE.Object3D {

  const particleCount = 8000;
  let translateArray = new Float32Array(particleCount * 3);
  let colorArray = new Float32Array(particleCount * 3);
  let scaleArray = new Float32Array(particleCount);
  let particleGeometry = new THREE.BufferGeometry();

  let mean = 2747.5 + 100
  let sd = 1509.878

  for (let i = 0, i3 = 0, l = particleCount; i < l; i++, i3 += 3) {


    translateArray[i3] = (Math.random() * (2) - 1) * width * 1.5;
    translateArray[i3 + 1] = (Math.random() * (2) - 1) * height * 1.5;
    translateArray[i3 + 2] = -(Math.random()) * depth;

    let rand = (randomG(8) * 2 - 1) * sd
    let arr = cctToRgb(mean + rand)
    colorArray[i3] = arr[0] / 255;
    colorArray[i3 + 1] = arr[1] / 255
    colorArray[i3 + 2] = arr[2] / 255

    scaleArray[i] = 2 + Math.random() * 10;


  }

  particleGeometry.setAttribute('position', new THREE.Float32BufferAttribute(translateArray, 3));
  particleGeometry.setAttribute('color', new THREE.Float32BufferAttribute(colorArray, 3));
  particleGeometry.setAttribute('size', new THREE.Float32BufferAttribute(scaleArray, 1).setUsage(THREE.DynamicDrawUsage));

  let vertexShader = document.getElementById('vshader')
  if (!vertexShader) return new THREE.Object3D()
  let fragmentShader = document.getElementById('fshader')
  if (!fragmentShader) return new THREE.Object3D()

  instanceMaterial = new THREE.ShaderMaterial({
    uniforms: {
      pointTexture: {value: new THREE.TextureLoader().load('/custom/textures/circle.png')}
    },
    vertexShader: vertexShader.textContent || "",
    fragmentShader: fragmentShader.textContent || "",
    depthTest: true,
    transparent: true,
    vertexColors: true
  });
  let points = new THREE.Points(particleGeometry, instanceMaterial)

  // instanceMesh = new THREE.Mesh(geometry, instanceMaterial);
  // instanceMesh.scale.set(200, 200, 200);
  clr++;
  return points
}

function render() {
  // particleSystem.translateZ(1)
  objs.children.forEach(c => c.translateZ(2))
  for (let i = 0; i < objs.children.length; i++) {
    if (objs.children[i].position.z >= depth) {
      objs.children[i].position.setZ(-depth * numSections)
    }
  }
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
  animation: screensaverBegin 750ms ease-in-out forwards;
}

@keyframes screensaverBegin {
  0% {
    opacity: 0.5;
    transform: scale(0.8);
    filter: blur(20px);
  }
  100% {
    opacity: 1;
    transform: scale(1);
    filter: blur(0px);
  }
}

</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import * as THREE from "three";
import {onMounted, onUnmounted} from "vue";
import {RenderPass} from "three/examples/jsm/postprocessing/RenderPass";
import {EffectComposer} from "three/examples/jsm/postprocessing/EffectComposer";


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

let colors = [0xff2d55, 0x5856d6, 0xff9500, 0xffcc00, 0xff3b30, 0x007aff, 0x4cd964]

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
  camera = new THREE.PerspectiveCamera(90, width / height, 1, 500)
  // camera = new THREE.OrthographicCamera(-width/2, width/2, he)
  camera.position.set(0, 0, 0);
  camera.lookAt(0, 0, 0)
}

// Handler resizing the viewport if and when the viewport is altered
function resizeFrame() {
  camera.aspect = width / height;
  if (!camera) return
  camera.updateProjectionMatrix();

  renderer.setSize(width, height);
}

const visibleHeightAtZDepth = (depth: number, camera: THREE.PerspectiveCamera) => {
  // compensate for cameras not positioned at z=0
  const cameraOffset = camera.position.z;
  if (depth < cameraOffset) depth -= cameraOffset;
  else depth += cameraOffset;

  // vertical fov in radians
  const vFOV = camera.fov * Math.PI / 180;

  // Math.abs to ensure the result is always positive
  return 2 * Math.tan(vFOV / 2) * Math.abs(depth);
};

const visibleWidthAtZDepth = (depth: number, camera: THREE.PerspectiveCamera) => {
  const height = visibleHeightAtZDepth(depth, camera);
  return height * camera.aspect;
};

// Initialize the THREE graphics and set up the environment
function initGraphics() {
  // Initialize Renderer
  renderer = new THREE.WebGLRenderer({antialias: true, alpha: true});
  // renderer.toneMapping = THREE.T;
  // Select screensaver dom element
  // renderer.shadowMap.enabled = true;
  // renderer.shadowMap.type = THREE.PCFSoftShadowMap;
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
  renderer.outputEncoding = THREE.LinearEncoding;

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

//   const renderTargetParameters = {
//     minFilter: THREE.LinearFilter,
//     magFilter: THREE.LinearFilter,
//     stencilBuffer: false
//   };
//
// // save pass
//   const savePass = new SavePass(
//       new THREE.WebGLRenderTarget(
//           width,
//           height,
//           renderTargetParameters
//       )
//   );
// //
// // blend pass
//   const blendPass = new ShaderPass(BlendShader, "tDiffuse1");
//   blendPass.uniforms["tDiffuse2"].value = savePass.renderTarget.texture;
//   blendPass.uniforms["mixRatio"].value = 0.4;
//
// // output pass
//   const outputPass = new ShaderPass(CopyShader);
//   outputPass.renderToScreen = true;

  composer = new EffectComposer(renderer);
  composer.addPass(renderScene);
  // composer.addPass(blendPass);
  // composer.addPass(savePass);
  // composer.addPass(outputPass);
  // composer.addPass(bloomPass);
  // Begin animation and rendering
  animate()


}

let entities = new Map<string, Entity>();

function initScene() {
  scene = new THREE.Scene()

  objs = new THREE.Object3D()

  for (let i = 0; i < 16; i++) {
    let ent = {
      uuid: `${i}`,
      depth: 200 + Math.random() * 200,
      direction: new THREE.Vector3(1 - Math.random() * 2, 1 - Math.random() * 2, 0),
      velocity: new THREE.Vector3(Math.random() + 0.2, Math.random() + 0.2, 0)
    } as Entity
    ent.radius = 20 + Math.random() * 20 * (depth / 800)
    let w = visibleWidthAtZDepth(ent.depth, camera) / 2
    let h = visibleHeightAtZDepth(ent.depth, camera) / 2

    const circle = new THREE.CylinderGeometry(ent.radius, ent.radius, 3, 48, 1, false, 0, Math.PI * 2)
    circle.rotateX(Math.PI / 2)


    let material = new THREE.MeshPhysicalMaterial({
      color: colors[Math.round(i % colors.length)],
      transparent: false,
      opacity: 1,
      transmission: 0.2,
      roughness: 0.5,
      clearcoat: 1,
      clearcoatRoughness: 1,
    });

    material.thickness = 2

    const mesh = new THREE.Mesh(circle, material)

    mesh.translateZ(-ent.depth - 2)
    mesh.translateX(w / 2 - Math.random() * w)
    mesh.translateY(h / 2 - Math.random() * h)
    ent.uuid = mesh.uuid
    entities.set(mesh.uuid, ent)
    objs.add(mesh)
  }


  scene.add(objs)


  let ambient = new THREE.AmbientLight(0xffffff, 1)
  scene.add(ambient)

  // let direc = new THREE.DirectionalLight(0xffffff, 1)
  // direc.position.set(0, 0, 1)
  // scene.add(direc)
  // scene.add(new THREE.CameraHelper(dirLight.shadow.camera));
  scene.background = new THREE.Color(0x00000000);

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


interface Entity {
  uuid: string
  radius: number
  depth: number
  direction: THREE.Vector3
  velocity: THREE.Vector3
}


let dir = new THREE.Vector3(-1, -1, 0)

function render() {
  // particleSystem.translateZ(1)

  for (let i = 0; i < objs.children.length; i++) {
    let child = objs.children[i]
    let entity = entities.get(child.uuid) as Entity
    // if (!entity) continue
    let w = visibleWidthAtZDepth(entity.depth, camera)
    let h = visibleHeightAtZDepth(entity.depth, camera)
    if (child.position.x <= -w / 2 + entity.radius) {
      entity.direction.setX(-entity.direction.x)
    } else if (child.position.x >= w / 2 - entity.radius) {
      entity.direction.setX(-entity.direction.x)
    }
    if (child.position.y <= -h / 2 + entity.radius) {
      entity.direction.setY(-entity.direction.y)
    } else if (child.position.y >= h / 2 - entity.radius) {
      entity.direction.setY(-entity.direction.y)
    }
    child.position.add(entity.direction)
  }

  // console.log(particleSystem.position.z)


}

</script>

<template>
  <div id="screensaver" class="screensaver-context"></div>
  <div v-if="true" class="screensaver-context-blur"></div>
</template>

<style scoped>


.screensaver-context-blur {
  content: ' ';

  width: 100%;
  height: 100%;
  top: 0;
  z-index: -10 !important;
  backdrop-filter: blur(30px);

}

.screensaver-context {
  width: 100%;
  height: 100%;
  position: absolute;
  z-index: -100 !important;
  animation: screensaverBegin 750ms ease-in-out forwards;
}

@keyframes screensaverBegin {
  0% {
    opacity: 0.5;
    transform: scale(0.8);
    filter: blur(20px);
  }
  99% {
    filter: blur(0px);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}

</style>
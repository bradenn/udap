<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {onMounted, provide, reactive, watch} from "vue";
import {version} from "../package.json";


interface Preferences {
  ui: {
    background: string
    theme: string
    mode: string
    blur: number
    brightness: number
    grid: boolean
    watermark: boolean
    night: boolean
    outlines: boolean
  }
}

const preferenceDefaults: Preferences = {
  ui: {
    blur: 6,
    background: "milk",
    mode: "cursor",
    theme: "dark",
    brightness: 100,
    grid: true,
    watermark: true,
    night: false,
    outlines: true,
  }
}

let preferences = reactive<Preferences>(restore())

function restore() {
  let stored = localStorage.getItem("preferences")
  if (stored) {
    let parsed: Preferences = JSON.parse(stored)
    return parsed
  }
  return preferenceDefaults
}

watch(preferences, () => {
  save(preferences)
})

function save(preferences: Preferences) {
  localStorage.setItem("preferences", JSON.stringify(preferences))
}

provide('preferences', preferences)


let state = reactive({
  hideHome: false,
  timeout: 0,
  fps: 0,
  countdown: 3,
  context: false,
  remote: {},
  system: {
    nexus: {
      system: {
        version: version
      }
    },
    udap: {
      system: {
        version: '0.0.0'
      }
    }
  }
})


onMounted(() => {
  state.context = false
  state.fps = 0
})

provide('system', state.system)

function handleUpdate() {
  if (state.countdown <= 0) {
    clearInterval(state.timeout)
    state.timeout = 0
    state.hideHome = false
  }
}

let lastReset = performance.now()
let totalFrames = 0

function tick() {
  totalFrames++
  let now = performance.now()
  let dFps = totalFrames / ((now - lastReset) / 1000.0)
  state.fps = Math.round(dFps * 10) / 10
  if (totalFrames > 2000) {
    totalFrames = 0
    lastReset = performance.now()
  }
}

function hideHome(hide: boolean) {
  state.hideHome = hide
}

function toggleContext(hide: boolean) {
  state.context = hide
}

provide('ui', preferences.ui)
provide('context', toggleContext)
provide('hideHome', hideHome)

</script>

<template>
  <div
      :class="`${preferences.ui.night?'night-vision':''} theme-${preferences.ui.theme} mode-${preferences.ui.mode} blurs-${preferences.ui.blur} brightness-${preferences.ui.brightness}`"
      class="root">
    <img :src="`/custom/${preferences.ui.background}@4x.png`" alt="Background" class="backdrop "/>
    <div v-if="preferences.ui.watermark" class="watermark">
      <div class="d-flex gap">
        <div>NEXUS v{{ state.system.nexus.system.version }}</div>
        <div v-if="state.system.udap">UDAP v{{ state.system.udap.system.version }}</div>
      </div>
      <div class="float-end">{{ $route.path }}</div>
    </div>
    <div v-if="preferences.ui.grid" class="grid"></div>

    <div v-if="state.context" class="context context-light"></div>
    <router-view/>


  </div>
</template>

<style lang="scss">

.root {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  border-radius: 0.3rem 0.3rem 0.3rem 0.3rem !important;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.05) !important;
}

.backdrop {
  z-index: -2 !important;
  top: 0;
  left: 0;
  background-color: rgba(28, 33, 40, 0.24);
  transform: scale(1);
  overflow: hidden;
  position: absolute;
  background-position: center;
  background-size: cover;
  object-fit: cover;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  animation: switch 0.25s ease-in-out forwards;
}

.backdrop:after {

}


.update-animate {
  animation: highlight 100ms ease-out forwards;
}

@keyframes highlight {
  0% {
    filter: blur(8px);
  }
  100% {
    filter: blur(0px);
  }
}


.night-vision {

  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;

  background-color: rgba(28, 0, 0, 0.6) !important;
}

.watermark {
  position: absolute;
  bottom: 0.3rem;
  width: calc(100% - 2rem);
  left: 1rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.3);
  font-size: 0.6rem;
  display: flex;
  gap: 0.75rem;
  align-items: center;
  justify-content: space-between;
  //outline: 1px solid #6f42c1;
  transition: all 500ms ease;
}

.hide-home {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 0.5rem;
  backdrop-filter: blur(8px);
  z-index: 0;
  background-color: rgba(255, 255, 255, 0.0125);
}

.mode-touch > * {
  cursor: none !important;
  user-select: none !important;
}


// Colors
$bg-color: rgba(0, 0, 0, 0);
$dot-color: rgba(255, 255, 255, 0.1);

// Dimensions
$dot-size: 28px;
$dot-space: 16px;


.grid {
  position: absolute;
  width: calc(100% - 0.5rem);
  height: calc(100% - 0.5rem);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.05) !important;
  margin: 0.25rem;

  z-index: -1;
  background-color: $bg-color;
  opacity: 1;
  border-radius: 0.2rem;
  border: 1px solid $dot-color;
  background-image: radial-gradient($dot-color 0.980000000000000px, $bg-color 0.980000000000000px), radial-gradient($dot-color 0.98px, $bg-color 0.980000000000000px);
  background-size: $dot-size $dot-size;
  background-position: 0 0, calc($dot-size / 2) calc($dot-size / 2);
}

.mode-cursor > * {
  cursor: crosshair !important;
  user-select: none !important;
}
</style>
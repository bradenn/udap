<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Clock from "@/components/Clock.vue"
import router from '@/router'
import {inject, onMounted, onUnmounted, onUpdated, provide, reactive, ref, watch} from "vue";
import "@/types";
import IdTag from "@/components/IdTag.vue";
import type {Identifiable, Metadata, Remote, Session, Timing, User} from "@/types";

import {Nexus, Target} from "@/views/terminal/nexus";
import CalculatorQuick from "@/views/terminal/calculator/CalculatorQuick.vue";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import Sideapp from "@/views/terminal/Sideapp.vue";

// -- Websockets --

let audio: HTMLAudioElement;
onMounted(() => {
  audio = new Audio('/sound/selection.mp3');

  remote.nexus = new Nexus(handleMessage)
})


// Define the reactive components for the remote data
let remote = reactive<Remote>({
  connected: false,
  metadata: {} as Metadata,
  entities: [],
  attributes: [],
  devices: [],
  networks: [],
  endpoints: [],
  users: [],
  timings: [],
  modules: [],
  zones: [],
  nexus: {} as Nexus
});


let ui: any = inject("ui")
let system: any = inject("system")

// Handle and route incoming messages to the local cache
function handleMessage(target: Target, data: any) {
  state.lastUpdate = new Date().valueOf()
  remote.connected = true
  switch (target) {
    case Target.Close:
      remote.connected = false
      break
    case Target.Metadata:
      system.udap.system = data.system as Metadata
      remote.metadata = data as Metadata
      break
    case Target.Entity:
      if (remote.entities.find((e: Identifiable) => e.id === data.id)) {
        remote.entities = remote.entities.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.entities.push(data)
      }
      break
    case Target.Attribute:
      if (remote.attributes.find((e: Identifiable) => e.id === data.id)) {
        remote.attributes = remote.attributes.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.attributes.push(data)
      }
      break
    case Target.User:
      if (remote.users.find((e: Identifiable) => e.id === data.id)) {
        remote.users = remote.users.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.users.push(data)
      }
      break
    case Target.Device:
      if (remote.devices.find((e: Identifiable) => e.id === data.id)) {
        remote.devices = remote.devices.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.devices.push(data)
      }
      break
    case Target.Network:
      if (remote.networks.find((e: Identifiable) => e.id === data.id)) {
        remote.networks = remote.networks.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.networks.push(data)
      }
      break
    case Target.Endpoint:
      if (remote.endpoints.find((e: Identifiable) => e.id === data.id)) {
        remote.endpoints = remote.endpoints.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.endpoints.push(data)
      }
      break

    case Target.Module:
      if (remote.modules.find((e: Identifiable) => e.id === data.id)) {
        remote.modules = remote.modules.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.modules.push(data)
      }
      break
    case Target.Zone:
      if (remote.zones.find((e: Identifiable) => e.id === data.id)) {
        remote.zones = remote.zones.map((a: Identifiable) => a.id === data.id ? data : a)
      } else {
        remote.zones.push(data)
      }
      break
    case Target.Timing:
      if (remote.timings.find((e: Timing) => e.pointer === data.pointer)) {
        remote.timings = remote.timings.map((a: Timing) => a.pointer === data.pointer ? data : a)
      } else {
        remote.timings.push(data)
      }
      break
  }
}


// -- Gesture Navigation --

let fps = 0;
let lastTick = ref(0);

onUpdated(() => {
  fps++;
})

onUnmounted(() => {
  remote.nexus.ws.close(1001, "Disconnecting")
})

let session = reactive<Session>({
  user: {} as User
})

provide("session", session)

// Stores the changing components of the main terminal
let state = reactive({
  sideApp: false,
  isDragging: false,
  timeout: null,
  verified: false,
  distance: 0,
  lastUpdate: 0,
  sideAppLock: false,
  showClock: true,
  scrollX: 0.0,
  scrollY: 0,
  scrollXBack: 0,
  scrollYBack: 0,
  dragA: {
    x: 0,
    y: 0
  }
});


watch(() => router.currentRoute.value, () => {
  state.showClock = router.currentRoute.value.path === '/terminal/home'
})

onMounted(() => {
  state.showClock = router.currentRoute.value.path === '/terminal/home'
})

// When called, if the user is still dragging, evoke the action confirming drag intent
function timeout() {
  // Check if user is till dragging
  if (state.isDragging) {
    // Verify the drag intention
    state.verified = true
  }
}


// When the user starts dragging, initialize drag intent
function dragStart(e: MouseEvent) {
  dragStop(e)
  // Record the current user position
  let a = {x: e.screenX, y: e.screenY}
  if (!e.view) return
  // If the drag has started near the bottom of the screen
  // Set the dragging status for later verification
  state.isDragging = true;
  // Record the drag position
  state.dragA = a
  // Verify drag intent if the user is still dragging after 100ms
  setTimeout(timeout, 10)
  // Otherwise, we consider the swipes
}

// While the user is still dragging
function dragContinue(e: MouseEvent) {
  // If the user is dragging, and the drag intent has been established
  state.isDragging = true;
  if (state.verified) {
    // Record the current position
    let dragB = {x: e.screenX, y: e.screenY}

    if (!e.view) return

    let height = e.view.screen.availHeight;
    let width = e.view.screen.availWidth;
    let thresholdOffset = 60;

    let isBottom = e.screenY > height - thresholdOffset;

    let isRight = state.dragA.x > width - thresholdOffset;

    let bottomPull = state.dragA.y - dragB.y;
    let rightPull = state.dragA.x - dragB.x;
    let gestureThreshold = 24;


    if (isBottom) {
      if (bottomPull > gestureThreshold) {
        router.push("/terminal/home")
      }
      if (state.dragA.y > dragB.y) {
        state.scrollY -= e.movementY
      }
    }
    let sideAppLockTarget = 445;
    let sideAppLockDelta = 5;
    if (isRight && !state.sideAppLock) {
      if (rightPull >= gestureThreshold) {
        state.sideApp = true
        if (sideAppLockTarget - rightPull <= sideAppLockDelta) {
          state.scrollX = 445
          state.sideAppLock = true
        } else {
          state.sideAppLock = false
        }

        state.scrollX += e.movementX
      }

    }

    if (isRight && !state.sideAppLock) {
      state.scrollX = state.dragA.x - dragB.x
    } else if (state.sideAppLock && sideAppLockTarget - rightPull <= sideAppLockDelta) {
      state.sideAppLock = false
    }


  }
}

// When the user cancels a drag intent
function dragStop(e: MouseEvent) {

  // Discard the drag intent
  state.isDragging = false;
  // Reset the distance
  state.distance = 0;
  // Reset verified drag intent
  state.verified = false;
  // Reset current position
  state.dragA = {x: 0, y: 0}
  if (state.scrollX !== 0 && state.scrollX !== 445) {
    if (state.scrollXBack == 0) {
      let target = state.scrollX >= 128 ? 445 : 0;
      state.scrollXBack = setInterval(() => {
        state.scrollX -= (state.scrollX - target) * 0.05
        if (Math.abs(target - state.scrollX) < 0.1) {
          clearInterval(state.scrollXBack)
          state.scrollX = target
          if (target == 0) {
            state.scrollX = 0
            state.sideApp = false
            state.sideAppLock = false
          }
          state.scrollXBack = 0
        }
      }, 5)
    } else {
      clearInterval(state.scrollXBack)
    }
  }

  if (Math.abs(state.scrollY) != 0 && state.scrollYBack == 0) {
    if (state.scrollYBack == 0) {
      state.scrollYBack = setInterval(() => {
        if (state.isDragging) return
        state.scrollY -= state.scrollY * 0.05
        if (Math.abs(state.scrollY) < 0.01) {
          state.scrollY = 0
          state.scrollYBack = 0
          clearInterval(state.scrollYBack)
        }
      }, 5)
    } else {
      clearInterval(state.scrollYBack)
    }
  }


}

// Provide the remote component for child components
provide('terminal', state)
provide('remote', remote)
</script>


<template>

  <div
      class="terminal"
      v-on:mousedown="dragStart"
      v-on:mousemove="dragContinue"
      v-on:mouseup="dragStop">
    <div class="generic-container gap-2">
      <div :class="`generic-slot-sm ` ">
        <Clock :small="!state.showClock"></Clock>
      </div>

      <div class="generic-slot-sm ">
        <IdTag></IdTag>
      </div>

    </div>

    <div
        class="route-view">
      <div>

        <Sideapp v-if="state.sideApp" :style="`transform: translateX(${-state.scrollX}px);`">
          <CalculatorQuick v-if="state.sideApp"></CalculatorQuick>
        </Sideapp>
      </div>
      <router-view v-slot="{ Component }">
        <component :is="Component"/>
      </router-view>


      <div class="justify-content-center d-flex align-items-center align-content-center">
        <div v-if="$route.matched.length > 1" @mouseover.prevent="state.scrollY!==0">
          <div v-if="$route.matched[1].children.length > 1">
            <Plot :cols="$route.matched[1].children.length" :rows="1"
                  class="bottom-nav">
              <Subplot v-for="route in ($route.matched[1].children as any[])" :icon="route.icon || 'earth-americas'"
                       :name="route.name"
                       :to="route.path"></Subplot>
            </Plot>
          </div>
        </div>
      </div>
    </div>

    <div
        :style="`transform: translateY(${-state.scrollY}px);`"
        class="home-bar top"></div>
  </div>
</template>

<style lang="scss" scoped>
.focus-enter-active {
  animation: animateIn 200ms;
}

.focus-leave-active {

  animation: animateOut 100ms;
}

@keyframes animateIn {
  0% {
    transform: scale(0.98);
    opacity: 0.4;
  }
  50% {
    transform: scale(0.99);
    opacity: 0.8;
  }
  100% {
    transform: scale(1);
  }
}

@keyframes animateOut {
  0% {
    transform: scale(1);
    opacity: 1;
  }

  100% {
    opacity: 0;
    transform: scale(0.98);
  }
}

.bottom-nav {
  display: inline-block;
  position: relative;
  bottom: 2.5rem;
  animation: dock-in 125ms ease-in forwards;
}


@keyframes dock-in {
  0% {
    bottom: -1rem;
  }

  100% {
    bottom: 2.5rem;
  }
}


.generic-container {
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}

.animate-in {
  animation: dock-in 100ms forwards;
}


.footer {
  position: absolute;
  bottom: 1.2rem;
}

.route-view {
  /*outline: 1px solid rgba(255,255,255,0.3);*/
  border-radius: 0.5rem;
  height: calc(100% - 3rem);
}

.terminal {
  padding: 1em;
  height: 100%;
  flex-direction: column;
  justify-content: start;
  transition: all 500ms;
}

</style>
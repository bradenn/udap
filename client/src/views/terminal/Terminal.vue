<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Clock from "@/components/Clock.vue"
import router from '@/router'
import {inject, onMounted, provide, reactive, watch} from "vue";
import "@/types";
import IdTag from "@/components/IdTag.vue";
import type {Identifiable, Metadata, Remote, Timing} from "@/types";

import {Nexus, Target} from "@/views/terminal/nexus";
import CalculatorQuick from "@/views/terminal/calculator/CalculatorQuick.vue";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";

// -- Websockets --

let audio: HTMLAudioElement;
onMounted(() => {
  audio = new Audio('/sound/selection.mp3');
})


// Define the reactive components for the remote data
let remote = reactive<Remote>({
  metadata: {} as Metadata,
  entities: [],
  attributes: [],
  devices: [],
  networks: [],
  endpoints: [],
  timings: [],
  nexus: new Nexus(handleMessage)
});


let ui: any = inject("ui")
let system: any = inject("system")

// Handle and route incoming messages to the local cache
function handleMessage(target: Target, data: any) {
  state.lastUpdate = new Date().valueOf()
  switch (target) {
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

// Stores the changing components of the main terminal
let state = reactive({
  sideApp: false,
  isDragging: false,
  timeout: null,
  verified: false,
  distance: 0,
  lastUpdate: 0,
  showClock: true,
  scrollX: 0,
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

  // Record the current user position
  let a = {x: e.clientX, y: e.clientY}
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

  if (state.verified) {
    // Record the current position
    let dragB = {x: e.clientX, y: e.clientY}
    // Set a maximum delta
    let delta = 200
    // Calculate the displacement
    // If the user's current position is larger than the defined max intent
    if (!e.view) return
    if (state.dragA.y > e.view.screen.availHeight - 200) {
      state.scrollY = (state.dragA.y - dragB.y) / 8
    }
    if (Math.abs(dragB.y - state.dragA.y) > delta && state.dragA.y > e.view.screen.availHeight - 200) {
      // Change the inner route to the home page
      dragStop(e)
      router.replace("/terminal/home")
    } else if (dragB.y - state.dragA.y > delta && state.dragA.y < 200) {
      // Reset the drag intention
      dragStop(e)
      ui.context = true
    } else {
      if (Math.abs(state.dragA.x) >= e.view.screen.availWidth - 32) {
        // Reset the drag intention

        if (Math.abs(state.dragA.x - dragB.x) > 64) {
          state.sideApp = true
          dragStop(e)
        } else {
          state.scrollX = (state.scrollX + dragB.x - state.dragA.x) / 16
        }
      }
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
  if (state.scrollX != 0 && state.scrollXBack == 0) {
    clearInterval(state.scrollXBack)
    state.scrollXBack = setInterval(() => {
      state.scrollX = state.scrollX - (state.scrollX / 20)
      if (Math.abs(state.scrollX) < 1) {
        clearInterval(state.scrollXBack)
        state.scrollX = 0
        state.scrollXBack = 0
      }
    }, 2)
  }

  if (state.scrollY != 0 && state.scrollYBack == 0) {
    clearInterval(state.scrollYBack)
    state.scrollYBack = setInterval(() => {
      if (state.isDragging) return
      state.scrollY = state.scrollY - Math.log(state.scrollY)
      if (Math.abs(state.scrollY) < 3) {
        state.scrollY = 0
        state.scrollYBack = 0
        clearInterval(state.scrollYBack)
      }
    }, 5)
  }


}


function selectSound() {
  // new Audio('/sound/selection.mp3').play();
}

// Provide the remote component for child components
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
      <div v-if="state.sideApp" class="context context-id " @click="state.sideApp = false"></div>
      <CalculatorQuick v-if="state.sideApp" class="position-absolute">dds</CalculatorQuick>
    </div>

    <div
        :style="`transform: translate(${Math.round(state.scrollX)}px,${-Math.round(state.scrollY)}px);`"
        class="route-view">
      <router-view v-slot="{ Component }" @mousedown="selectSound">
        <component :is="Component"/>
      </router-view>
      <div v-if="$route.matched.length > 1">
        <Plot v-if="$route.matched[1].children.length > 1" :cols="$route.matched[1].children.length" :rows="1"
              class="bottom-nav">
          <Subplot v-for="route in $route.matched[1].children" :icon="route.icon || 'earth-americas'" :name="route.name"
                   :to="route.path"></Subplot>
        </Plot>
      </div>
    </div>

    <div
        :style="`transform: translateY(calc(-${Math.round(state.scrollY)}px));`"
        class="home-bar top"></div>
  </div>
</template>

<style scoped>


.bottom-nav {
  display: block;
  bottom: 2.85rem;

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

@keyframes dock-in {
  from {
    bottom: -1rem;
  }
  to {
    bottom: 0.5rem;
  }
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
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Dock from "@/components/Dock.vue"
import Clock from "@/components/Clock.vue"
import IdHash from "@/components/IdHash.vue"
import router from '@/router'
import {onMounted, provide, reactive, watch} from "vue";
import {Nexus, Target} from '@/views/terminal/nexus'

import type {Attribute, Device, Endpoint, Entity, Identifiable, Network} from "@/types";

// -- Websockets --

let audio: HTMLAudioElement;
onMounted(() => {
  audio = new Audio('/public/sound/selection.mp3');
})

interface Metadata {
  name: string;
  version: string;
  environment: string;
  ipv4: string;
  ipv6: string;
  hostname: string;
  mac: string;
  go: string;
  cores: number;
}

interface Remote {
  metadata: Metadata,
  entities: Entity[],
  attributes: Attribute[],
  devices: Device[],
  networks: Network[],
  endpoints: Endpoint[],
  timings: any[],
  nexus: Nexus
}


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

// Handle and route incoming messages to the local cache
function handleMessage(target: Target, data: any) {
  switch (target) {
    case Target.Metadata:
      remote.metadata = data as Metadata
      break
    case Target.Entity:
      replaceAndUpdate(remote.entities, data)
      break
    case Target.Attribute:
      replaceAndUpdate(remote.attributes, data)
      break
    case Target.Device:
      replaceAndUpdate(remote.devices, data)
      break
    case Target.Network:
      replaceAndUpdate(remote.networks, data)
      break
    case Target.Endpoint:
      replaceAndUpdate(remote.endpoints, data)
      break
    case Target.Timing:
      replaceAndUpdate(remote.timings, data)
      break
  }
}

// Replace or update the existing data
function replaceAndUpdate(target: any, data: any) {
  if (target.find((e: Identifiable) => e.id === data.id)) {
    target = target.map((a: Identifiable) => a.id === data.id ? data : a)
  } else {
    target.push(data)
  }
}

// Provide the remote component for child components
provide('remote', remote)

// -- Gesture Navigation --

// Stores the changing components of the main terminal
let state = reactive({
  isDragging: false,
  timeout: null,
  verified: false,
  distance: 0,
  showClock: true,
  scrollX: 0,
  scrollBack: 0,
  dragA: {
    x: 0,
    y: 0
  }
});

watch(() => router.currentRoute.value, () => {
  state.showClock = router.currentRoute.value.path === '/terminal/home'
})

onMounted(() => {
  draw()
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

  // If the drag has started near the bottom of the screen
  if ((window.screen.availHeight - e.screenY) <= 1800) {
    // Set the dragging status for later verification
    state.isDragging = true;
    // Record the drag position
    state.dragA = a
    // Verify drag intent if the user is still dragging after 100ms
    setTimeout(timeout, 10)
    // Otherwise, we consider the swipes
  }
}

// While the user is still dragging
function dragContinue(e: MouseEvent) {
  // If the user is dragging, and the drag intent has been established

  if (state.verified) {
    // Record the current position
    let dragB = {x: e.clientX, y: e.clientY}
    // Set a maximum delta
    let delta = 100
    // Calculate the displacement
    state.distance = (state.dragA.y - dragB.y) / delta
    // If the user's current position is larger than the defined max intent
    if (state.dragA.y - dragB.y > delta) {
      // Reset the drag intention
      state.verified = false
      // Reset the frame position
      state.distance = 0

      // Change the inner route to the home page
      router.push("/terminal/home")
    }
    if (Math.abs(state.dragA.x - dragB.x) > 20) {
      // Reset the drag intention
      state.verified = true
      // Reset the frame position
      state.distance = 0
      if (Math.abs(state.dragA.x - dragB.x) > 600) {
        dragStop(e)

      } else {
        state.scrollX = (state.scrollX + dragB.x - state.dragA.x) / 4

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
  if (state.scrollX != 0 && state.scrollBack == 0) {

    state.scrollBack = setInterval(() => {
      state.scrollX = state.scrollX - (state.scrollX / 8)
      if (Math.abs(state.scrollX) < 1) {
        clearInterval(state.scrollBack)
        state.scrollX = 0
        state.scrollBack = 0
      }
    }, 16)
  }

}


function selectSound() {
  audio.play();
}

function draw() {


}

</script>


<template>
  <div
      class="terminal"
      v-on:mousedown="dragStart"
      v-on:mousemove="dragContinue"
      v-on:mouseup="dragStop">
    <div class="generic-container">

      <div :class="`generic-slot-sm`">
        <Clock :small="!state.showClock"></Clock>
      </div>
      <div class="generic-slot-sm">
        <div class="element h-75 d-flex align-items-center align-content-center justify-content-start gap-2">
          <div class="px-2">
            <IdHash></IdHash>
          </div>
          <div class="d-flex flex-column gap-0">
            <div class="label-xxs label-o5 lh-1">Braden Nicholson</div>
            <div class="label-c2 label-o2 lh-1">Superuser</div>
          </div>
        </div>
      </div>
    </div>

    <div :style="`transform: translate(${state.scrollX}px,0) !important;`" class="route-view">
      <router-view v-slot="{ Component }">
        <component :is="Component"/>
      </router-view>
    </div>

    <div v-if="state.showClock" class="footer mt-3">
      <Dock class="animate-in" os>

        <router-link class="macro-icon-default" draggable="false" to="/terminal/wifi/">
          <div class="macro-icon" @mousedown="selectSound">
            􀙇
          </div>
        </router-link>
        <router-link class="macro-icon-default" draggable="false" to="/terminal/energy/">
          <div class="macro-icon" @mousedown="selectSound">
            􀋦
          </div>
        </router-link>
        <router-link class="macro-icon-default" draggable="false" to="/terminal/exogeology/">
          <div class="macro-icon" @mousedown="selectSound">
            <i class="fa-solid fa-meteor fa-fw"></i>
          </div>
        </router-link>
        <router-link class="macro-icon-default" draggable="false" to="/terminal/timing/">
          <div class="macro-icon" @mousedown="selectSound">
            <i class="fa-solid fa-stopwatch fa-fw"></i>
          </div>
        </router-link>
        <router-link class="macro-icon-default" draggable="false" to="/terminal/settings/preferences">
          <div class="macro-icon" @mousedown="selectSound">
            􀍟
          </div>
        </router-link>
      </Dock>
    </div>
    <div :style="`transform: translateY(calc(-${state.distance}rem));`" class="home-bar top"></div>
  </div>
</template>

<style scoped>

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
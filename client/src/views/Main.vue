<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Dock from "../components/Dock.vue"
import Clock from "../components/Clock.vue"
import router from '../router'
import {provide, reactive} from "vue";
import {Nexus, Target} from '../nexus'

// -- Websockets --

let nexus = new Nexus()

function handleMessage(target: Target, data: any) {
  switch (target) {
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

function replaceAndUpdate(target: any, data: any) {
  if (target.find(e => e.id === data.id)) {
    target = target.map(a => a.id === data.id ? data : a)
  } else {
    target.push(data)
  }
}

nexus.connect(handleMessage)

let remote = reactive({
  entities: [],
  attributes: [],
  devices: [],
  networks: [],
  endpoints: [],
  timings: [],
});


provide('remote', remote)

// -- Gesture Navigation --

// Stores the changing components of the main terminal
let state = reactive({
  isDragging: false,
  timeout: null,
  verified: false,
  distance: 0,
  dragA: {
    x: 0,
    y: 0
  }
});

// When called, if the user is still dragging, evoke the action confirming drag intent
function timeout() {
  // Check if user is till dragging
  if (state.isDragging) {
    // Verify the drag intention
    state.verified = true
  }
}

// When the user starts dragging, initialize drag intent
function dragStart(e) {
  // Record the current user position
  let a = {x: e.clientX, y: e.clientY}
  // If the drag has started near the bottom of the screen
  if ((window.screen.availHeight - e.screenY) <= 128) {
    // Set the dragging status for later verification
    state.isDragging = true;
    // Record the drag position
    state.dragA = a
    // Verify drag intent if the user is still dragging after 100ms
    setTimeout(timeout, 100)
  }
}

// While the user is still dragging
function dragContinue(e) {
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
  }
}

// When the user cancels a drag intent
function dragStop(e) {
  // Discard the drag intent
  state.isDragging = false;
  // Reset the distance
  state.distance = 0;
  // Reset verified drag intent
  state.verified = false;
  // Reset current position
  state.dragA = {x: 0, y: 0}
}

</script>


<template>
  <div :style="`transform: translateY(calc(-${state.distance}rem));`"
       class="terminal h-100"
       v-on:mousedown="dragStart"
       v-on:mousemove="dragContinue"
       v-on:mouseup="dragStop">

    <div class="generic-container">
      <div class="generic-slot-sm">
        <Clock inner></Clock>
      </div>
    </div>

    <div class="h-100">
      <router-view v-slot="{ Component }">
        <component :is="Component"/>
      </router-view>
    </div>

    <div class="footer mt-3">
      <!--      <div v-if="sys" class="position-absolute" style="left:1rem;">-->
      <!--        <div class="label-xxs label-o4 label-w600 lh-1 text-lowercase">{{ sys.name }} v{{ sys.version }}</div>-->
      <!--        <div class="label-ys label-o2 label-w500 lh-1">{{ sys.environment }} build</div>-->
      <!--      </div>-->
      <Dock os>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/home">􀎟</router-link>
        </div>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/apps/room">􀟼</router-link>
        </div>
        <span class="mx-2 my-1"
              style="width: 0.0255rem; height: 1.8rem; border-radius: 1rem; background-color: rgba(255,255,255,0.1);"></span>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/apps/media">􀑪</router-link>
        </div>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/network/">􁅏</router-link>
        </div>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/wifi/">􀙇</router-link>
        </div>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/energy/">􀋦</router-link>
        </div>
        <div class="macro-icon">
          <router-link class="macro-icon-default" draggable="false" to="/terminal/settings/endpoint">􀍟</router-link>
        </div>
      </Dock>
    </div>
    <div class="home-bar top"></div>
  </div>
</template>

<style scoped>
.terminal {
  padding: 1em;
  height: 100% !important;
  flex-direction: column;
  justify-content: start;
}

</style>
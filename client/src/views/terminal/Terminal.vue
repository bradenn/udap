<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import router from '@/router'
import {defineAsyncComponent, inject, onMounted, onUnmounted, provide, reactive, ref, watch} from "vue";
import "@/types";

import type {
  Attribute,
  Device,
  Endpoint,
  Entity,
  Identifiable,
  Log,
  Metadata,
  Module,
  Network,
  Preferences,
  Remote,
  RemoteRequest,
  Session,
  TerminalDiagnostics,
  Timing,
  User,
  Zone
} from "@/types";
import {memorySizeOf} from "@/types";

import {Nexus, Target} from "@/views/terminal/nexus";

import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import Error from "@/components/Error.vue";

const Clock = defineAsyncComponent({
  loader: () => import('@/components/Clock.vue'),
  errorComponent: Error,
  timeout: 250
})

const Input = defineAsyncComponent({
  loader: () => import('@/views/Input.vue'),
  errorComponent: Error,
  timeout: 250
})

const Glance = defineAsyncComponent({
  loader: () => import('@/views/terminal/Glance.vue'),
  errorComponent: Error,
  timeout: 250
})

const Bubbles = defineAsyncComponent({
  loader: () => import('@/views/screensaver/Bubbles.vue'),
  errorComponent: Error,
  timeout: 250
})

const Warp = defineAsyncComponent({
  loader: () => import('@/views/screensaver/Warp.vue'),
  errorComponent: Error,
  timeout: 250
})

const IdTag = defineAsyncComponent({
  loader: () => import('@/components/IdTag.vue'),
  errorComponent: Error,
  timeout: 250
})

// -- Websockets --
onMounted(() => {
  remote.nexus = new Nexus(handleMessage)
})

onUnmounted(() => {
  remote.nexus.ws.close()
  remote = {} as Remote
})


// Define the reactive components for the remote data
let remote = reactive<Remote>({
  connected: false,
  metadata: {} as Metadata,
  entities: [] as Entity[],
  attributes: [] as Attribute[],
  devices: [] as Device[],
  networks: [] as Network[],
  endpoints: [] as Endpoint[],
  users: [] as User[],
  timings: [] as Timing[],
  modules: [] as Module[],
  zones: [] as Zone[],
  logs: [] as Log[],
  nexus: {} as Nexus,
  size: "" as string,
  diagnostics: {
    queue: [] as RemoteRequest[],
    updates: new Map<string, number>(),
    connected: false,
    maxRSS: 0,
    lastTarget: "",
    lastUpdate: 0,
    objects: 0
  } as TerminalDiagnostics
});


let screensaver: any = inject("screensaver")
let preferences = inject("preferences") as Preferences
let system: any = inject("system")

// Handle and route incoming messages to the local cache
function handleMessage(target: Target, data: any) {
  remote.diagnostics.lastUpdate = new Date().valueOf()
  state.lastUpdate = new Date().valueOf()
  remote.connected = true
  let dx = 0;
  switch (target) {
    case Target.Close:
      remote.connected = false
      break
    case Target.Metadata:
      system.udap.system = data.system as Metadata
      remote.metadata = data as Metadata
      dx = 1
      break
    case Target.Timing:
      if (remote.timings.find((e: Timing) => e.pointer === data.pointer)) {
        remote.timings = remote.timings.map((a: Timing) => a.pointer === data.pointer ? data : a)
        dx = 0
      } else {

        remote.timings.push(data)
        dx = 1
      }
      break
    case Target.Entity:
      dx = createOrUpdate(remote.entities, data)
      break
    case Target.Attribute:
      dx = createOrUpdate(remote.attributes, data)
      break
    case Target.User:
      dx = createOrUpdate(remote.users, data)
      break
    case Target.Device:
      dx = createOrUpdate(remote.devices, data)
      break
    case Target.Network:
      dx = createOrUpdate(remote.networks, data)
      break
    case Target.Endpoint:
      dx = createOrUpdate(remote.endpoints, data)
      break
    case Target.Module:
      dx = createOrUpdate(remote.modules, data)
      break
    case Target.Zone:
      dx = createOrUpdate(remote.zones, data)
      break
    case Target.Log:
      dx = createOrUpdate(remote.entities, data)
      break
  }

  let prev = remote.diagnostics.updates.get(target) || 0
  remote.diagnostics.updates.set(target, prev + dx);
  let session = {
    target: target,
    time: new Date().valueOf(),
    operation: "update",
    payload: data,
    id: (data as Identifiable).id
  } as RemoteRequest
  remote.diagnostics.queue.push(session)
  remote.diagnostics.lastTarget = target
  if (remote.diagnostics.queue.length >= 10) {
    remote.diagnostics.queue = remote.diagnostics.queue.slice(0, remote.diagnostics.queue.length - 2)
  }


  remote.diagnostics.maxRSS = memorySizeOf(remote)


}

function createOrUpdate(target: any[], data: Identifiable): number {
  if (target.find((e: Identifiable) => e.id === data.id)) {
    target = target.map((a: Identifiable) => a.id === data.id ? data : a)
    return 0
  } else {
    target.push(data)
    return 1
  }

}

// -- Gesture Navigation --

let lastTick = ref(0);

onUnmounted(() => {
  remote.nexus.ws.close(1001, "Disconnecting")
})

let session = reactive<Session>({
  user: {} as User,
  screensaver: false
})

provide("session", session)

// Stores the changing components of the main terminal
let state = reactive({
  locked: false,
  sideApp: false,
  isDragging: false,
  timeout: null,
  verified: false,
  distance: 0,
  lastUpdate: 0,
  sideAppLock: false,
  showClock: true,
  hideTerminal: false,
  scrollX: 0.0,
  scrollY: 0,
  scrollXBack: 0,
  scrollYBack: 0,
  input: {
    open: false,
    meta: {} as InputProps,
    cb: (a: string) => {
    },
  },
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

interface InputProps {
  name: string
  value: string
  description: string
  icon?: string
  type?: string
}

function openInput(props: InputProps, cb: (a: string) => void) {
  if (!state.input.open) {
    state.input.open = true
    state.input.meta = props
    state.input.cb = cb
  }
}

function applyInput(a: string) {
  state.input.cb(a)
  closeInput()
}

function closeInput() {
  state.input.open = false
  state.input.meta = {} as InputProps
  state.input.cb = () => {
  }
}

provide("input", openInput)

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
    let thresholdOffset = 80;

    let isBottom = e.screenY > height - thresholdOffset;
    let isTop = e.screenY > thresholdOffset;

    let isRight = state.dragA.x > width - thresholdOffset;

    let topPull = dragB.y - state.dragA.y;
    let bottomPull = state.dragA.y - dragB.y;
    let rightPull = state.dragA.x - dragB.x;
    let gestureThreshold = 24;


    if (isBottom) {
      if (bottomPull > gestureThreshold) {
        state.locked = false
        router.push("/terminal/home")
      }
      if (state.dragA.y > dragB.y) {
        state.scrollY -= e.movementY
      }
    } else if (isTop) {
      if (topPull > gestureThreshold) {
        screensaver.startScreensaver()
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

  if (Math.abs(state.scrollY) != 0 && state.scrollYBack == 0) {
    if (state.scrollYBack == 0) {
      state.scrollYBack = setInterval(() => {
        if (state.isDragging) return
        state.scrollY -= state.scrollY * 0.2
        if (Math.abs(state.scrollY) < 0.01) {
          state.scrollY = 0
          state.scrollYBack = 0
          clearInterval(state.scrollYBack)
        }
      }, 14)
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

  <div v-if="!screensaver.hideTerminal"
       class="terminal"
       v-on:mousedown="dragStart"
       v-on:mousemove="dragContinue" v-on:mouseup="dragStop">
    <Glance v-if="state.locked"></Glance>
    <div v-else class="d-inline">
      <div class="generic-container gap-2">
        <div class="" v-on:click="(e) => state.locked = true">
          <Clock :small="!state.showClock"></Clock>
        </div>

        <div class="generic-slot-sm ">
          <div v-if="false" class="element p-2" style="width: 13rem !important;">
            {{ remote.logs.find(l => ((new Date().valueOf() - new Date(l.time).valueOf()) >= 60000)) }}
            <div class="label-c2 label-w500 label-o5 lh-1">Worldspace</div>
            <div class="label-c3 label-w400 label-o3 lh-1">Matthew has arrived</div>
          </div>
        </div>
        <div class="generic-slot-sm ">
          <IdTag></IdTag>
        </div>

      </div>
      <div class="route-view pt-1">
        <router-view v-slot="{ Component }" style="max-height: calc(100% - 2.9rem) !important;">
          <component :is="Component"/>
        </router-view>
      </div>
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
    <div :style="`transform: translateY(${-state.scrollY}px);`" class="home-bar top"></div>
  </div>

  <Bubbles v-if="screensaver.show && preferences.ui.screensaver.selection === 'bubbles'"
           class="screensaver-overlay"></Bubbles>
  <Warp v-else-if="screensaver.show && preferences.ui.screensaver.selection === 'warp'"
        class="screensaver-overlay"></Warp>
  <Input v-if="state.input.open" :apply="applyInput" :close="closeInput" :description="state.input.meta.description"
         :name="state.input.meta.name" :value="state.input.meta.value"></Input>
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
    //opacity: 0.4;
  }
  50% {
    transform: scale(0.99);
    //opacity: 0.8;
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
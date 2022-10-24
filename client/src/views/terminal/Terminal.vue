<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import router from '@/router'
import {defineAsyncComponent, inject, onMounted, onUnmounted, provide, reactive, ref, watch, watchEffect} from "vue";
import "@/types";

import type {
  Attribute,
  Device,
  Endpoint,
  Entity,
  Identifiable,
  Log,
  Macro,
  Metadata,
  Module,
  Network,
  Preferences,
  Remote,
  RemoteRequest,
  Session,
  SubRoutine,
  TerminalDiagnostics,
  Timing,
  Trigger,
  User,
  Zone
} from "@/types";
import {memorySizeOf} from "@/types";

import {Nexus, Target} from "@/views/terminal/nexus";

import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import type {Haptics} from "@/views/terminal/haptics";
import haptics from "@/views/terminal/haptics";
import Toast from "@/components/Toast.vue";

const Clock = defineAsyncComponent({
  loader: () => import('@/components/Clock.vue'),

})

const Input = defineAsyncComponent({
  loader: () => import('@/views/Input.vue'),

})

const Glance = defineAsyncComponent({
  loader: () => import('@/views/terminal/Glance.vue'),

})

const Bubbles = defineAsyncComponent({
  loader: () => import('@/views/screensaver/Bubbles.vue'),
})

const Warp = defineAsyncComponent({
  loader: () => import('@/views/screensaver/Warp.vue'),
})

const IdTag = defineAsyncComponent({
  loader: () => import('@/components/IdTag.vue'),
})


// -- Websockets --
onMounted(() => {
  haptics.connect("ws://10.0.1.60/ws")

  // haptic.haptics = new HapticEngine("ws://10.0.1.60/ws")
  remote.nexus = new Nexus(handleMessage)
})

// function tap(frequency: number, iterations: number, amplitude: number) {
//   haptic.haptics.tap(frequency, iterations, amplitude)
// }

provide("haptic", haptics.tap)
provide("haptics", haptics as Haptics)

onUnmounted(() => {
  if (!remote.nexus.ws) return
  remote.nexus.ws.close()
  if (!haptics) return
  haptics.close()
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
  subroutines: [] as SubRoutine[],
  macros: [] as Macro[],
  triggers: [] as Trigger[],
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
      return

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
      remote.entities = createOrUpdate(remote.entities, data)
      break
    case Target.Macro:
      remote.macros = createOrUpdate(remote.macros, data)
      break
    case Target.SubRoutine:
      remote.subroutines = createOrUpdate(remote.subroutines, data)
      break
    case Target.Trigger:
      remote.triggers = createOrUpdate(remote.triggers, data)
      break
    case Target.Attribute:
      remote.attributes = createOrUpdate(remote.attributes, data)
      break
    case Target.User:
      remote.users = createOrUpdate(remote.users, data)
      break
    case Target.Device:
      remote.devices = createOrUpdate(remote.devices, data)
      break
    case Target.Network:
      remote.networks = createOrUpdate(remote.networks, data)
      break
    case Target.Endpoint:
      remote.endpoints = createOrUpdate(remote.endpoints, data)
      break
    case Target.Module:
      remote.modules = createOrUpdate(remote.modules, data)
      break
    case Target.Zone:
      remote.zones = createOrUpdate(remote.zones, data)
      break
    case Target.Log:
      remote.entities = createOrUpdate(remote.entities, data)
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

function mouseDown() {
  // axios.post("http://10.0.1.60/pop", {
  //   power: 2
  // }).then(res => {
  //   return
  // }).catch(err => {
  //   return
  // })
}

function createOrUpdate(target: any[], data: Identifiable): any[] {
  if (target.find((e: Identifiable) => e.id === data.id)) {
    if (data.deleted) {
      return target.filter((a: Identifiable) => a.id !== data.id)
    } else {
      return target.map((a: Identifiable) => a.id === data.id ? data : a)
    }
  } else {
    if (!data.deleted) {
      target.push(data)
      return target
    }
  }
  return target
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
  mouseDown();
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
    let isTop = e.screenY <= thresholdOffset;

    let isRight = state.dragA.x > width - thresholdOffset;

    let topPull = dragB.y - state.dragA.y;
    let bottomPull = state.dragA.y - dragB.y;
    let rightPull = state.dragA.x - dragB.x;
    let gestureThreshold = 24;

    if (isBottom) {
      if (bottomPull > gestureThreshold) {
        state.verified = false
        haptics.tap(2, 1, 50)
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
    } else {
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

interface ToastObject {
  name: string,
  message: string,
  severity: number,
  duration: number
}

const toasts = reactive({
  queue: [] as ToastObject[],
  active: false,
  interval: 0,
  current: {} as ToastObject
})

watchEffect(() => {
  if (toasts.queue.length > 0 && !toasts.active) {
    toasts.active = true
    toasts.current = toasts.queue[0]
    toasts.queue = toasts.queue.filter(q => q !== toasts.current)
    toasts.interval = setInterval(() => {
      if (toasts.active) {
        toasts.current.duration -= 1000
        if (toasts.current.duration <= 0) {
          toasts.active = false
          toasts.current = {} as ToastObject
          clearInterval(toasts.interval)
        }
      }
    }, 1000)
  }
  return toasts.queue
})


const notifications = {
  show: (name: string, message: string, severity: number, duration: number) => {
    toasts.queue.push({
      name: name,
      message: message,
      severity: severity,
      duration: duration
    })
  },
}

provide('notifications', notifications)

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

        <Toast v-if="toasts.active" :message="toasts.current.message" :severity="toasts.current.severity"
               :time="toasts.current.duration" :title="toasts.current.name"></Toast>

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
        <div v-if="$route.matched.length > 1" @click.prevent="state.scrollY!==0">
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
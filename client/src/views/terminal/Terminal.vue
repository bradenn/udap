<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import router from '@/router'
import {inject, onMounted, onUnmounted, provide, reactive, watch, watchEffect} from "vue";
import "@/types";

import type {Entity, Preferences} from "@/types";

import Clock from '@/components/Clock.vue'
import IdTag from '@/components/IdTag.vue'
import ContextBar from "@/components/ContextBar.vue";
import Plot from "@/components/plot/Plot.vue";
import Subplot from "@/components/plot/Subplot.vue";
import Notification from "@/components/Notification.vue";
import Input from '@/views/Input.vue'
import ScreensaverDom from "@/views/screensaver/Screensaver.vue";

import type {Remote} from "@/remote";
import _remote from "@/remote";


import type {Context} from "@/context";
import _context from "@/context";

import core from "@/core";
import type {Haptics} from "@/haptics";
import type {Screensaver} from "@/screensaver";

/* Remote */
const remote: Remote = _remote
provide("remote", remote)

/* Context */
const context: Context = _context
provide("context", context)


const haptics: Haptics = core.haptics()
const screensaver: Screensaver = core.screensaver()

// Load runtime when the terminal view is loaded
onMounted(() => {

    remote.connect()
    haptics.connect("ws://10.0.1.60/ws")
})

onUnmounted(() => {
    remote.disconnect()
    haptics.disconnect()
})


let preferences: Preferences = inject("preferences") as Preferences
let system: any = inject("system")

// Stores the changing components of the main terminal
let state = reactive({
    locked: false,
    sideApp: false,
    isDragging: false,
    timeout: null,
    verified: false,
    distance: 0,
    animationFrame: 0,
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
    haltAnimation()
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

let meta = reactive({
    hue: 0,
    dim: 0
})

watchEffect(() => {
    pollMeta(remote as Remote)
    return remote as Remote
})

function pollMeta(rm: Remote): boolean {
    if (!rm) return true

    let e = rm.entities.find(a => a.name === "terminal") || {} as Entity
    if (!e) return false

    let attrs = rm.attributes.filter(a => a.entity === e.id)

    let hue = attrs.find(a => a.key === "hue")
    if (!hue) return false;
    meta.hue = parseInt(hue.value)

    let dim = attrs.find(a => a.key === "dim")
    if (!dim) return false;
    meta.dim = parseInt(dim.value)

    return true
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
        let thresholdOffset = 50;
        let thresholdOffsetX = 128;

        let isCenter = e.screenX > (width / 2 - thresholdOffsetX) && e.screenX < (width / 2 + thresholdOffsetX);
        let isBottom = e.screenY > height - thresholdOffset;
        let isTop = e.screenY <= thresholdOffset;

        let isRight = state.dragA.x > width - thresholdOffset;

        let topPull = dragB.y - state.dragA.y;
        let bottomPull = state.dragA.y - dragB.y;
        let rightPull = state.dragA.x - dragB.x;
        let gestureThreshold = 20;


        if (isBottom && isCenter) {
            if (bottomPull > gestureThreshold) {
                state.verified = false
                haptics.tap(3, 1, 50)
                state.locked = false
                router.push("/terminal/home")
            }
            state.scrollY = height - e.clientY
        } else if (isTop) {
            if (topPull > gestureThreshold) {
                screensaver.start()
            }
        } else {
        }
    }
}

function haltAnimation() {
    state.scrollY = 0
    state.scrollYBack = 0
    cancelAnimationFrame(state.animationFrame)
}

const animation = reactive({
    frame: 0,
    lastFrame: Date.now(),
})

const frameInterval = 1000.0 / 60.0;

// Animate the retraction of the home bar
function animate() {
    // Request the animation frame for the next update
    animation.frame = requestAnimationFrame(animate)
    // Get the current time in milliseconds
    let now = Date.now()
    // Get milliseconds since last frame
    let delta = now - animation.lastFrame
    // Check if the delta exceeds the frame interval (16.667ms, 33.33ms, 25ms, etc)
    if (delta > frameInterval) {
        // Set the time of the last frame to now minus the elapsed frame time
        animation.lastFrame = now - (delta % frameInterval)
        // Run the desired code to be animated
        decelerate()
    }
}

// Stops the current animation from proceeding beyond the current frame
function stopAnimation() {
    // Cancel the animation frame to block the next update
    cancelAnimationFrame(animation.frame)
    animation.frame = 0
    animation.lastFrame = 0
}

// Returns the moving object back to its original position
function decelerate() {
    // Stop if the user is currently dragging
    if (state.isDragging) return
    // Decrease the scroll position
    state.scrollY -= state.scrollY * 0.5

    // If the item is close enough to its original position, reset it
    if (Math.abs(state.scrollY) <= 0.05) {
        // Set the level to 0
        state.scrollY = 0
        // Stop any further animation
        stopAnimation()
    }
}

// When the user cancels a drag intent
function dragStop(_: MouseEvent) {
    // Discard the drag intent
    state.isDragging = false;
    // Reset the distance
    state.distance = 0;
    // Reset verified drag intent
    state.verified = false;
    // Reset current position
    state.dragA = {x: 0, y: 0}
    // Return the bar back to its position
    if (Math.abs(state.scrollY) > 0) animate()
}

// Provide the remote component for child components
provide('terminal', state)

</script>

<template>

    <div v-if="!screensaver.active()" class="h-100 w-100">

        <div class="terminal w-100 h-100"
             v-on:mousedown="dragStart"
             v-on:mousemove="dragContinue"
             v-on:mouseup="dragStop">

            <div :style="$route.matched.length > 1?`height: calc(100% - 3.75rem); max-height: calc(100% - 3.75rem) !important;`:'max-height: 100% !important;'"
                 class="d-flex flex-column gap-1">
                <div>
                    <ContextBar v-if="!state.showClock">
                        <div style="grid-column: span 3">
                            <Clock :small="!state.showClock"></Clock>
                        </div>

                        <div style="grid-column: 9 / span 8">
                            <Notification></Notification>
                        </div>

                        <div class="d-flex align-items-end justify-content-end"
                             style="grid-column: 20 / span 5">
                            <IdTag></IdTag>
                        </div>

                    </ContextBar>
                    <ContextBar v-else>

                        <div
                                class=" d-flex align-content-center align-items-center justify-content-start px-1"
                                style="grid-column: span 3">
                            <Clock :small="!state.showClock"></Clock>
                        </div>

                        <div style="grid-column: 9 / span 8">
                            <Notification></Notification>

                        </div>

                        <div class="d-flex align-items-end justify-content-end"
                             style="grid-column: 20 / span 5">
                            <IdTag></IdTag>
                        </div>
                    </ContextBar>
                </div>

                <div
                        style="height: 100%; max-height: 100% !important; overflow-y: clip !important;">
                    <router-view v-slot="{ Component }">
                        <component :is="Component"/>
                    </router-view>
                </div>

            </div>
            <div v-if="$route.matched.length > 1" class="bottom-nav">
                <div
                        class="justify-content-center d-flex align-items-center align-content-center">
                    <div
                            @click.prevent="state.scrollY!==0">
                        <div v-if="$route.matched[1].children.length > 1">
                            <Plot :cols="$route.matched[1].children.length" :rows="1"
                            >
                                <Subplot
                                        v-for="route in ($route.matched[1].children as any[])"
                                        :name="route.name"
                                        :sf="route?.icon || '?'"
                                        :to="route.path"></Subplot>
                            </Plot>
                        </div>
                    </div>
                </div>
            </div>
            <div v-if="preferences.ui.watermark" class="watermark">
                <div class="d-flex gap">
                    <div v-if="remote.metadata" class="label-r label-w600">
                        {{ remote.metadata?.system?.version }}
                    </div>
                </div>
                <div class="float-end">{{ router.currentRoute.value.path }}</div>
            </div>


            <div :style="`transform: translateY(${-state.scrollY}px);`"
                 class="home-bar top"></div>

        </div>

    </div>

    <div v-if="context.isActive()" class="context" @mousedown="context.hideContext()"></div>

    <ScreensaverDom></ScreensaverDom>

    <Input v-if="state.input.open" :apply="applyInput" :close="closeInput"
           :description="state.input.meta.description"
           :name="state.input.meta.name" :value="state.input.meta.value"></Input>
    <div :style="`box-shadow: inset 0 0 4px 7px hsla(${meta.hue},75%,50%,${0.8*meta.dim/100}), inset 0 0 96px 16px hsla(${meta.hue},75%,50%,${0.5*meta.dim/100}) !important;`"
         class="neon" @mousedown.passive></div>
</template>

<style lang="scss" scoped>

/* Watermark Mode */

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

.terminal {
  z-index: 2;
  display: inline-block;
}

.neon {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  border-radius: 0.5rem;
  z-index: -1 !important;
}

.bottom-nav {
  position: absolute;
  width: 100%;
  z-index: 0 !important;
  animation: dock-in 125ms ease-in forwards;
  bottom: 1.5rem;
  left: 0;
}
</style>
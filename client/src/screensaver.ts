// Copyright (c) 2022 Braden Nicholson

import {reactive} from "vue";


interface ScreenSaverState {
    active: boolean
    current: string
    delay: number
    lastReset: number
    countdown: number
    interval: number
}

export interface Screensaver {
    start(): void

    stop(): void

    active(): boolean

    change(saver: string): boolean

    current(): string

    delay(): number

    countdown(): number
}

const state = reactive<ScreenSaverState>({
    active: false,
    current: "warp",
    delay: 1000 * 60 * 5,
    lastReset: 0,
    countdown: 0,
    interval: 0,
})


window.addEventListener("mousedown", function () {
    resetCountdown()
})

function change(saver: string): boolean {
    if (saver !== state.current) {
        state.current = saver
        return true
    }
    return false;
}

function resetCountdown() {
    state.lastReset = Date.now().valueOf()
    if (state.interval != 0) return;
    state.interval = setInterval(() => {
        state.countdown = (Date.now() - state.lastReset)
        if ((Date.now() - state.lastReset) >= state.delay) {
            state.active = true
            clearInterval(state.interval)
            state.interval = 0
        }
    }, 500)
}

function current(): string {
    return state.current;
}

function countdown(): number {
    return Math.floor(state.countdown / 1000);
}

function start() {
    state.active = true
}

function stop() {
    state.active = false
}

function active(): boolean {
    return state.active
}

export default {
    current,
    countdown,
    start,
    change,
    stop,
    active,
} as Screensaver




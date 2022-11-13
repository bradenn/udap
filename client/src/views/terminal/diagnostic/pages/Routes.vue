<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, reactive, watchEffect} from "vue";
import type {RemoteRequest, TerminalDiagnostics, Timing} from "@/types";
import {useRouter} from "vue-router";
import type {Remote} from "@/remote";

const remote = inject("remote") as Remote

const state = reactive({
    terminal: {} as TerminalDiagnostics,
    terminalDeltas: {
        lastUpdate: ""
    },
    sortedQueue: [] as RemoteRequest[],
    loaded: false,
    diffs: new Map<string, number>(),
    maxRSSHistory: [] as number[],
    maxRSSHistoryDisplay: [] as number[],
    timings: [] as Timing[],
    timingsRange: {
        min: 10E88,
        max: -1000
    },
    avgInterval: 0,
    intervalBuffer: new Date().valueOf(),
    avgIntervalHistory: Array(96).fill(0),
    intervalHistory: [] as number[],

})

// let timeInterval = 0;
//
// onMounted(() => {
//   timeInterval = setInterval(updateTimes, 100)
// })
//
// onUnmounted(() => {
//   clearInterval(timeInterval)
// })
//
// function updateTimes() {
//   for (let queueElement of state.terminal.queue) {
//     if (queueElement.id) {
//       state.diffs.set(queueElement.id, Date.now() - queueElement.time)
//     }
//   }
//   let keys = state.diffs.keys()
//   for (let key of keys) {
//     if (!state.terminal.queue.find(q => q.id === key)) {
//       state.diffs.delete(key)
//     }
//   }
//
//
//   state.sortedQueue = state.terminal.queue.sort((a, b) => b.time - a.time)
// }

watchEffect(() => {
    updateStats()
    updateDeltas()

})

function updateDeltas() {
    let dt = new Date().valueOf() - state.terminal.lastUpdate
    state.terminalDeltas.lastUpdate = `${dt} ms`
    return state.terminalDeltas
}

function updateStats() {

    state.terminal = remote.diagnostics
    state.timings = remote.timings.filter(f => f.complete).sort((a, b) => {
        if (state.timingsRange.max < b.startNano) state.timingsRange.max = b.startNano
        if (state.timingsRange.min > b.startNano) state.timingsRange.min = b.startNano
        return a.startNano < b.startNano ? -1 : 1
    })

    // if (new Date().valueOf() - state.intervalBuffer >= 100) {
    //
    //   state.intervalHistory.push(new Date().valueOf() - state.intervalBuffer)
    //
    //   if (state.intervalHistory.length >= 2) {
    //     state.intervalHistory = state.intervalHistory.slice(1, state.intervalHistory.length - 1)
    //   }
    //
    //   state.avgInterval = state.intervalHistory.reduce((a, b) => a += b, 0) / state.intervalHistory.length
    //
    //   state.avgIntervalHistory.push(state.avgInterval)
    //
    //   if (state.avgIntervalHistory.length >= 96) {
    //     state.avgIntervalHistory = state.avgIntervalHistory.slice(1, state.avgIntervalHistory.length - 1)
    //   }
    //   state.intervalBuffer = new Date().valueOf()
    // }

    if (state.maxRSSHistory.length >= 128) {
        state.maxRSSHistory = state.maxRSSHistory.slice(1, state.maxRSSHistory.length - 1)
    }
    // updateTimes()
    state.sortedQueue = state.terminal.queue.sort((a, b) => b.time - a.time)
    state.maxRSSHistory.push(state.terminal.maxRSS)

    state.loaded = true

}

const router = useRouter()

</script>

<template>
    <div class="page-grid">

    </div>
</template>

<style scoped>

.work-box {
    width: 100%;
    height: 1px;
    background-color: #0a58ca;
}

.pane-fill {
    position: relative;
    grid-column: 2 / 4;
    grid-row: 1 / 1;
    display: flex;
    flex-direction: column;
}

.page-grid {
    width: 100%;
    height: 100%;
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(1, 1fr);
    grid-template-columns: repeat(4, 1fr);
}
</style>
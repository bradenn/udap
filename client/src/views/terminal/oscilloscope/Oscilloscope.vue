<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, onUnmounted, reactive} from "vue";
import {v4 as uuidv4} from "uuid";
import Switch from "@/components/Switch.vue";
import Subplot from "@/components/plot/Subplot.vue";
import Plot from "@/components/plot/Plot.vue";
import Frequency from "@/components/Frequency.vue";


const state = reactive({
    canvas: {} as HTMLCanvasElement,
    ctx: {} as CanvasRenderingContext2D,
    socket: {} as WebSocket,
    connected: false,
    heap: [] as number[],
    stack: [] as Resolve[],
    timecodes: [] as number[],
    uuid: uuidv4(),
    polling: false,
    out: ""
})

onMounted(() => {
    connect()
})

interface Resolve {
    begin: number
    values: number[]
}

function connect() {
    let socket = new WebSocket("ws://10.0.1.85/ws")
    socket.onopen = function (e: Event) {
        state.socket = socket
        state.connected = true;
    }
    socket.onerror = function (e: Event) {
        console.log("Error: ", e)
        state.connected = false;
    }
    socket.onclose = function (e: Event) {
        state.connected = false;
    }

    socket.onmessage = handleMessage;
}

onUnmounted(() => {
    state.socket.close()
})

function ping() {
    if (state.polling) {
        state.socket.send("ping")
    }
}

const size = 1024;
const chunks = 1;

function handleMessage(e: MessageEvent) {
    let data = JSON.parse(e.data) as Resolve
    let cp = state.heap
    cp.push(...data.values)
    if (cp.length > size * chunks) {
        state.heap = state.heap.slice(state.heap.length - size * chunks, state.heap.length)
    }


}

</script>

<template>
    <div class="d-flex flex-column  gap-2">
        <div>
            Oscilloscope
        </div>
        <div>

        </div>
        <div class="d-flex flex-row gap-2">
            <div>Connected: {{ state.connected }}</div>
            <Switch :active="state.polling" :change="(v) => {state.polling = v; state.polling?ping():{}}"
                    name="Polling"></Switch>
            <Plot :cols="2" :rows="2">
                <Subplot :active="false" :fn="() => ping()" name="Ping"></Subplot>
            </Plot>
        </div>
        <div>
            <Frequency :chunk-size="size" :timecodes="state.timecodes" :values="state.heap"></Frequency>
        </div>
    </div>
</template>

<style scoped>
.inner-canvas {
    width: 100%;


}

.canvas-container {
    display: flex;
    flex-direction: row;
    justify-content: center;
    width: 100%;
    align-items: center;
    border-radius: 8px;
    background-color: hsla(214, 9%, 28%, 0.2);
    padding: 6px
}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Log, Remote} from "@/types";
import moment from "moment";

let remote = inject('remote') as Remote
let preferences = inject('preferences')


let state = reactive({
  history: [] as number[],
  numLogs: 50,
  logs: [] as Log[],
  moduleLogs: [] as Log[],
})

interface LogLevel {
  name: string
  icon: string
  class: string
}

let levels = new Map<string, LogLevel>([
  ["warn", {
    name: "Warning",
    icon: "􀇿",
    class: "level-warn"
  }],
  ["error", {
    name: "Error",
    icon: "􀘰",
    class: "level-error"
  }]
])

interface Terminal {
  lastUpdate: number
}

let terminal = inject("terminal") as Terminal


function update() {
  state.logs = remote.logs.sort((a, b) => moment(a.time).valueOf() > moment(b.time).valueOf() ? -1 : 1)

  state.moduleLogs = state.logs

  state.logs.forEach(l => {
    if (!l.level) {
      l.level = "warn"
    }
  })
}

onMounted(() => {
  update()
})

watchEffect(() => {
  update()
  if (!state.history.includes(terminal.lastUpdate)) state.history.push(terminal.lastUpdate)
  if (state.history.length > 100) state.history = state.history.slice(1, state.history.length)
  return remote
})

function formatLog(log: Log): string {

  let time = moment(log.time).format("hh:mmA")
  let since = moment(log.time).fromNow(false)

  return `${since} [${log.event}] ${log.message}`
}

function parseTime(time: string): string {
  return moment(time).format("hh:mmA")
}


</script>

<template>
  <div class="h-100">

    <div class="element log-container">
      <div class="d-flex align-items-center justify-content-between">
        <div class="label-c1  label-o4 label-w500 px-1 pb-1">Module Logs</div>
        <div class="label-c2 label-w500 label-r label-o2 px-1">Showing {{ state.logs.length }}/{{ remote?.logs.length }}
          logs
        </div>
      </div>
      <div class="logs">
        <div v-for="log in state.logs" class="log">
          <div class="log-desc">
            <div class="label-o3 label-c1 label-w400 lh-1">
            <span :class="levels.get(log.level)?.class" class="">{{
                levels.get("warn")?.icon
              }}
            </span>
              {{ log.event }}
            </div>
            <div class="label-o3 label-c2 label-w300 lh-1">{{ parseTime(log.time) }}</div>
          </div>
          <div class="label-o4 label-c2 label-w500 lh-1">{{ log.message }}</div>

        </div>
      </div>
    </div>

  </div>
</template>

<style scoped>


.level-warn {
  color: rgba(255, 159, 10);
}

.level-error {
  color: rgba(255, 69, 58);
}


.log-desc {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  flex-shrink: 1;

  /*outline: 1px solid white;*/
}

.log:first-child {
  border-bottom: none !important;
}

.log {
  display: flex;
  flex-direction: row;
  justify-content: start;
  gap: 0.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding: 4px 0 !important;
  /*outline: 1px solid white;*/
}

.logs {
  display: flex;

  flex-direction: column-reverse;
  overflow: scroll;
  flex-grow: 1;
  gap: 0.25rem;
  padding: 0.5rem;
}

.log-container {
  height: 100%;
  display: flex;
  flex-direction: column;

  overflow: visible;
}

@keyframes moveIt {
  0% {
    transform: translateX(50px);
  }
  50% {
    transform: translateX(25px);
  }
  100% {
    transform: translateX(50px);
  }
}

.updates {
  max-height: 5rem;
  width: 4px;
  border-radius: 1rem;
  background-color: white;
}

</style>
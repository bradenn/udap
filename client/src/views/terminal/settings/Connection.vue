<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, reactive, watchEffect} from "vue";
import Header from "@/components/Header.vue";
import type {Remote} from "@/types";

let remote = inject('remote') as Remote
let preferences = inject('preferences')


let state = reactive({
  history: [] as number[]
})

interface Terminal {
  lastUpdate: number
}

let terminal = inject("terminal") as Terminal


watchEffect(() => {
  if (!state.history.includes(terminal.lastUpdate)) state.history.push(terminal.lastUpdate)
  if (state.history.length > 100) state.history = state.history.slice(1, state.history.length)
  return remote
})


</script>

<template>
  <div>
    <Header class="px-2 pb-1" icon="cloud" name="Connection"></Header>
    <div class="d-flex justify-content-center align-items-center gap-1" style="height: 8rem">
      <div v-for="v in state.history.map(h => Math.abs(h-terminal.lastUpdate)).reverse()"
           :style="`height: ${v/100}px; opacity:0;`"
           class="updates">

      </div>
    </div>
    <div class="move">
      <div class="element " style="width: 12rem; height: 5rem;">
        dd
      </div>
    </div>

  </div>
</template>

<style scoped>

.move .element {
  animation: moveIt 1s forwards infinite;
}

@keyframes moveIt {
  0% {
    transform: translate3d(50px, 0, 0);
  }
  50% {
    transform: translate3d(25px, 0, 0);
  }
  100% {
    transform: translate3d(50px, 0, 0);
  }
}

.updates {
  max-height: 5rem;
  width: 4px;
  border-radius: 1rem;
  background-color: white;
}

</style>
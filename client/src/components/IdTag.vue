<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import IdHash from "@/components/IdHash.vue"
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";
import Toggle from "@/components/plot/Toggle.vue";
import type {Remote} from "@/types";

let state = reactive({
  menu: false,
  reloading: true,
  connected: false,
})

let ui: any = inject("ui")
let remote: Remote = inject('remote') as Remote
let system: any = inject('system')

onMounted(() => {
  state.reloading = false
})

watchEffect(() => {
  state.connected = remote.connected
  return state.connected
})

function toggleMenu() {

  state.menu = !state.menu
}

function reload() {
  state.reloading = true
  document.location.reload()
}


</script>

<template>
  <div v-if="state.menu" class="context context-id" @click="state.menu = false"></div>
  <div class=" tag-container element d-flex align-items-center align-content-center justify-content-start gap-1"
       @click="toggleMenu">
    <div class="px-1">
      <IdHash></IdHash>
    </div>

    <div class="d-flex flex-column gap-0">
      <div class="label-c2 label-o5 lh-1 label-r">Braden Nicholson</div>
      <div class="label-c3 label-o2 label-r label-w500 lh-1"
           style="font-family: 'Roboto Light', sans-serif;"><span v-if="state.connected">Connected</span><span v-else>Disconnected</span>
      </div>
    </div>
    <div class="flex-grow-1"></div>

    <div v-if="remote.nexus.state > 1" class="d-flex flex-column gap-0 justify-content-end align-items-end ">
      <div class="label-c3 label-o2 px-1">
        <span class="label-c3 label-w300">&nbsp;DOWN</span>
      </div>
    </div>
    <div class="label-c2 label-o2 px-1">
      <div v-if="state.menu">
        <i class="fa-solid fa-caret-down "></i>
      </div>
      <div v-else>
        <i class="fa-solid fa-caret-left px-1"></i>
      </div>
    </div>
  </div>

  <div v-if="state.menu" class="tag-summary d-flex flex-column gap-1 p-1">
    <Plot :cols="1" :rows="1">
      <div class="subplot d-flex align-items-end align-items-center lh-1 gap-1">
        <div class="label-o2">􀉩</div>
        <div class="label-o4 label-r label-w600">Braden Nicholson</div>
        <div class="flex-fill"></div>
        <div class="px-1 label-o2"> 􀉯</div>
      </div>
    </Plot>

    <Plot :cols="4" :rows="1">
      <div class="subplot plot-centered" @click="$router.push('/terminal/home')">
        <div class="label-o4 label-c2"><i class="fa-solid fa-house fa-fw"></i></div>
      </div>
      <div class="subplot plot-centered" @click="$router.push('/terminal/home')">
        <div class="label-o4 label-c2"><i class="fa-solid fa-shield fa-fw"></i></div>
      </div>

      <div class="subplot plot-centered" @click="reload">

        <div class="label-o4 label-c2">
          <div v-if="state.reloading">
            <Loader size="sm"></Loader>
          </div>

          <i v-else class="fa-solid fa-rotate-right fa-fw"></i>
        </div>
      </div>
      <div class="subplot plot-centered" @click="$router.push('/terminal/settings')">
        <div class="label-o4 label-c2"><i class="fa-solid fa-cog fa-fw"></i></div>
      </div>
    </Plot>
    <Plot :cols="2" :rows="1">
      <div class="subplot">
        <div class="label-c3">Wi-Fi</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="d-flex align-items-center label-o4 label-c3 lh-1 px-1">
            <i class="fa-solid fa-circle text-success" style="font-size: 8px; line-height: 1rem;"></i>&nbsp;&nbsp;OK
          </div>

        </div>
      </div>
      <div class="subplot">
        <div class="label-c3">NEXUS</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="d-flex align-items-center lh-1">
            <div v-if="remote.nexus.state > 1" class="d-flex align-items-center label-o4 label-c3 lh-1 px-1">
              <i class="fa-solid fa-circle text-danger" style="font-size: 8px; line-height: 1.2rem;"
              ></i>&nbsp;&nbsp;DOWN
            </div>
            <div v-else class="d-flex align-items-center label-o4 label-c3 lh-1 px-1">
              <i class="fa-solid fa-circle text-success" style="font-size: 8px; line-height: 1rem;"
              ></i>&nbsp;&nbsp;OK
            </div>

          </div>
        </div>
      </div>
    </Plot>
    <Plot :cols="1" :rows="1" title="Brightness">
      <input v-model="ui.brightness"
             :max=20
             :min=4
             :step=1
             class="slider-small"
             type="range"
             @mousemove.stop>
    </Plot>
    <Plot :cols="2" :rows="2" title="Quick Settings">
      <Toggle :active="ui.grid" :fn="() => ui.grid = !ui.grid" title="Grid"></Toggle>
      <Toggle :active="ui.outlines" :fn="() => ui.outlines = !ui.outlines" title="Outlines"></Toggle>
      <Toggle :active="ui.watermark" :fn="() => ui.watermark = !ui.watermark" title="Watermark"></Toggle>
      <Toggle :active="ui.blurBg" :fn="() => ui.blurBg = !ui.blurBg" title="Bg Blur"></Toggle>
    </Plot>

  </div>

</template>

<style lang="scss" scoped>

.tag-container {
  height: 2.5rem;
}

.tag-summary {

  position: relative;
  margin-top: 0.125rem;

  z-index: 22;
}

@keyframes slideIn {
  0% {
    transform: scale(1);
  }
  15% {
    transform: scale(1.05);
  }
  30% {
    transform: scale(1.015);
  }
  100% {
    transform: scale(1);
  }
}

.subplot:active {
  animation: click 200ms ease forwards;
}

.tag-container:active {
  animation: click 100ms ease forwards;
}


@keyframes click {
  0% {
    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
  }
  30% {
    transform: scale(0.97);
  }
  100% {
    transform: scale(1);
  }
}

.canvas-container {
  height: 42px;
  width: 42px;

  aspect-ratio: 1/1 !important;
}

</style>

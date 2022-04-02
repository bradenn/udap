<script lang="ts" setup>

import {inject, reactive} from "vue";
import IdHash from "@/components/IdHash.vue"

const remote: any = inject('remote')
let system: any = inject('system')

let state = reactive({
  menu: false,

})

function toggleMenu() {
  state.menu = !state.menu
}

let ui: any = inject("ui")
</script>

<template>
  <div
      class="tag-container element d-flex align-items-center align-content-center justify-content-start gap-1"
      @click="toggleMenu">
    <div class="px-1">
      <IdHash></IdHash>
    </div>
    <div class="d-flex flex-column gap-0">
      <div class="label-c2 label-o5 lh-1">Braden Nicholson</div>
      <div class="label-c2 label-o2 label-w500 lh-1" style="font-family: 'Roboto Light', sans-serif;">Superuser</div>
    </div>
    <div class="flex-grow-1"></div>
    <div v-if="remote.nexus.state > 1" class="d-flex flex-column gap-0 justify-content-end align-items-end ">
      <div class="label-c3 label-o2 px-1">
        <i class="fa-solid fa-cloud"></i>
        <span class="label-c3 label-w300">&nbsp;DOWN</span>
      </div>
    </div>
    <div class="label-c2 label-o2 px-1">
      <div v-if="state.menu">
        <i class="fa-solid fa-caret-down "></i>
      </div>
      <div v-else>
        <i class="fa-solid fa-caret-left label-c2 label-o2 px-1"></i>
      </div>
    </div>
  </div>

  <div v-if="state.menu" class="tag-summary element">
    <div class="plot plot-1x4">
      <div class="subplot">
        <div>Wi-Fi</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="d-flex align-items-center label-o4 label-c2 lh-1 px-1">
            <i class="fa-solid fa-circle label-c3 text-success" style="font-size: 8px; line-height: 1rem;"></i>&nbsp;&nbsp;OK
          </div>

        </div>
      </div>
      <div class="subplot">
        <div>NEXUS</div>
        <div class="d-flex justify-content-center align-items-center align-content-center px-1">
          <div class="d-flex align-items-center label-o4 label-c2 lh-1">
            <span v-if="system.udap.system.version"><i class="fa-solid fa-circle label-c3 text-success"
                                                       style="font-size: 8px; line-height: 1rem;"></i>&nbsp;&nbsp;OK</span>
            <span v-else><i class="fa-solid fa-circle label-c3 text-danger"
                            style="font-size: 8px; line-height: 1rem;"></i>&nbsp;&nbsp;DOWN</span>
          </div>
        </div>
      </div>

    </div>
    <div v-if="ui" class="plot plot-1x4 pt-2 ">
      <div class="subplot " @click="ui.grid = !ui.grid">
        <div>Grid</div>
        <div class="d-flex justify-content-center align-items-center align-content-center ">

          <div class="label-c2 label-o4 px-1 lh-1">
            {{ ui.grid ? 'ON' : 'OFF' }}
          </div>

        </div>
      </div>
      <div class="subplot" @click="ui.outlines = !ui.outlines">
        <div>Outlines</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="label-c2 label-o4 px-1 lh-1">
            {{ ui.outlines ? 'ON' : 'OFF' }}
          </div>
        </div>
      </div>
      <div class="subplot" @click="ui.watermark = !ui.watermark">
        <div>Watermark</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="label-c2 label-o4 px-1 lh-1">
            {{ ui.watermark ? 'ON' : 'OFF' }}
          </div>
        </div>
      </div>
      <div class="subplot" @click="ui.nightVision = !ui.nightVision">
        <div>Night Vision</div>
        <div class="d-flex justify-content-center align-items-center align-content-center">
          <div class="label-c2 label-o4 px-1 lh-1">
            {{ ui.nightVision ? 'ON' : 'OFF' }}
          </div>
        </div>
      </div>

    </div>

  </div>
</template>

<style lang="scss" scoped>

.tag-container {
  height: 2.5rem;
}

.tag-summary {
  backdrop-filter: blur(42px);
  position: relative;
  margin-top: 0.125rem;

  animation: slideIn 200ms ease forwards;
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
    transform: scale(1.025);
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

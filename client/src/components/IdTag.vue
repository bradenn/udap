<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import Plot from "@/components/plot/Plot.vue";
import Toggle from "@/components/plot/Toggle.vue";
import type {Attribute, Entity, Preferences, Remote, Status} from "@/types";
import Button from "@/components/Button.vue";

let state = reactive({
  menu: false,
  reloading: true,
  connected: false,
  link: false,
  zoneEntity: {} as Entity,
  zoneAttribute: {} as Attribute,
  status: {} as Status
})

let preferences: Preferences = inject("preferences") as Preferences
let remote: Remote = inject('remote') as Remote
let system: any = inject('system')

onMounted(() => {
  update()
  state.reloading = false
})

watchEffect(() => {
  state.connected = remote.connected
  update()
  return state.zoneAttribute
})

function update() {
  let entity = remote.entities.find(e => e.name === 'faces')
  if (!entity) return
  state.zoneEntity = entity

  let attr = remote.attributes.find(e => e.key === 'deskFace')
  if (!attr) return
  state.zoneAttribute = attr

  let stat = JSON.parse(attr.value) as Status
  if (!stat) return

  state.status = stat

}

const haptics = inject("haptic") as (a: number, b: number, c: number) => void

function open() {
  haptics(2, 1, 100)
  state.menu = !state.menu
}

function toggleMenu() {
  state.menu = !state.menu
}

function reload() {
  state.reloading = true
  document.location.reload()
}


</script>

<template>
  <div v-if="state.menu" class="context context-id" @click="state.menu = false">
  </div>
  <div class="d-flex flex-column">
    <div
        class="tag-container element d-flex align-items-center align-content-center justify-content-start gap-1 px-2 h-100"
        style="width: 11.25rem; height: 1.25rem" @mousedown="open">
      <div class="id-icon">

        <span v-if="state.connected">􀙇</span>
        <span v-else>􀙈</span>

      </div>
      <div class="id-icon">
        <span v-if="state.connected">􀌌</span>
        <span v-else>􀌐</span>

      </div>
      <div>

      </div>

      <div class="flex-grow-1"></div>

      <div v-if="remote?.nexus.state > 1" class="d-flex flex-column gap-0 justify-content-end align-items-end lh-1">
        <div class="label-c3 label-o2 px-1">
          <span class="label-c3 label-w300">&nbsp;DOWN</span>
        </div>
      </div>
      <div class="label-c2 label-o2 px-0">
        <div v-if="state.menu">
          <i class="fa-solid fa-caret-down "></i>
        </div>
        <div v-else>
          <i class="fa-solid fa-caret-left px-1"></i>
        </div>
      </div>
    </div>

    <div v-if="state.menu" class="tag-summary d-flex flex-column gap-1 py-1">
      <Plot :cols="4" :rows="1" @click="state.menu = false">
        <Button :active="true" text="􀎟" @click="$router.push('/terminal/home')"></Button>
        <Button :active="true" text="􀨲" @click="$router.push('/terminal/remote')"></Button>
        <Button :active="true" text="􀅈" @click="reload"></Button>
        <Button :active="true" text="􀍟" @click="$router.push('/terminal/settings')"></Button>
      </Plot>
      <Plot :cols="2" :rows="1">
        <Button :active="true" text="Beta" @click="$router.push('/terminal/home')"></Button>
        <Button :active="true" text="􀨲" @click="$router.push('/terminal/remote')"></Button>
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
                {{ remote.size }}
              </div>

            </div>
          </div>
        </div>
      </Plot>
      <Plot :cols="2" :rows="2">
        <Toggle :active="preferences.ui.grid" :fn="() => preferences.ui.grid = !preferences.ui.grid"
                title="Grid"></Toggle>
        <Toggle :active="preferences.ui.screensaver.enabled"
                :fn="() => preferences.ui.screensaver.enabled = !preferences.ui.screensaver.enabled"
                title="Screensaver"></Toggle>
        <Toggle :active="preferences.ui.watermark" :fn="() => preferences.ui.watermark = !preferences.ui.watermark"
                title="Watermark"></Toggle>
        <Toggle :active="preferences.ui.background.blur"
                :fn="() => preferences.ui.background.blur = !preferences.ui.background.blur" title="Bg Blur"></Toggle>
      </Plot>

    </div>
  </div>

</template>

<style lang="scss" scoped>
.id-icon {
  font-size: 0.60rem;
  text-shadow: 0 0 2px rgba(255, 255, 255, 0.2);
  color: rgba(255, 255, 255, 0.45);
  mix-blend-mode: overlay;
}

.user-container {
  display: flex;
}

.user-container .user:not(:first-child) {
  margin-left: -0.25rem;

  z-index: 1000 !important;
}

.user {

  width: 1.5rem;
  height: 1.5rem;
  font-size: 0.75rem;
  margin-left: 0.125rem;
  display: flex;
  align-items: center;
  justify-content: center;

  color: rgba(255, 255, 255, 0.6);

  border-radius: 100% !important;

}

.user-secondary {

}

.tag-container {

  height: 1.5rem;
  z-index: 22;
}

.tag-summary {

  position: relative;
  margin-top: 0.125rem;

  z-index: 22;
}

.tag-container {
  animation: none;
}

.tag-container:active {
  //animation: tagClick 100ms ease forwards !important;
}


@keyframes tagClick {
  0% {
    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
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

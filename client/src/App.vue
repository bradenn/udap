<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {onMounted, onUnmounted, provide, reactive} from "vue";
import {version} from "../package.json";
import {usePersistent} from "@/persistent";
import core from "@/core";

import type {Haptics} from "@/haptics";
import _haptics from "@/haptics";

import type {Notify} from "@/notifications";
import _notify from "@/notifications";

import type {Screensaver} from "@/screensaver";
import _screensaver from "@/screensaver";

let system = reactive({
    nexus: {
        system: {
            version: version
        }
    },
    udap: {
        system: {
            version: '0.0.0'
        }
    }
})

/* Haptics */
const haptics: Haptics = _haptics
provide("haptics", haptics)

/* Notify */
const notify: Notify = _notify
provide("notify", notify)

/* Screensaver */
const screensaver: Screensaver = _screensaver
provide("screens", screensaver)


const router = core.router()

const preferences = usePersistent()

onMounted(() => {
    haptics.connect("ws://10.0.1.60/ws")
})

onUnmounted(() => {
    haptics.disconnect()
})

provide('system', system)


</script>

<template>

    <div
            :class="`root theme-${preferences.ui.theme} mode-${preferences.ui.mode} blurs-${preferences.ui.blur} h-100`">

        <img :class="`${preferences.ui.background.blur?'backdrop-blurred':''}`"
             :src="`/custom/${preferences.ui.background.image}@4x.png`" alt=""
             class="backdrop "/>
        <div v-if="preferences.ui.grid" class="grid"></div>
        <div class="h-100">
            <router-view/>
        </div>
    </div>

</template>

<style lang="scss" scoped>

.root {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  border-radius: 0.4rem 0.4rem 0.4rem 0.4rem !important;
}

.screensaver-text {
  animation: screensaverTextLoadIn 500ms ease-in forwards;
}

@keyframes screensaverTextLoadIn {
  0% {
    opacity: 0;
    scale: 0.8;
  }
  100% {
    opacity: 1;
    scale: 1;
  }

}

.backdrop {

  position: absolute !important;
  z-index: -1 !important;
  top: 0;
  left: 0;
  background-position: center;
  background-size: cover;
  object-fit: cover;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
}

.backdrop.backdrop-blurred {
  filter: blur(8px);
}


/* Grid Modes */

.grid {
  position: absolute;
  width: calc(100%);
  height: calc(100%);
  background-color: rgba(20, 20, 22, 1);
  z-index: -1;
}


/* Touch Modes */

.mode-touch > * {
  cursor: none !important;
  user-select: none !important;
}


.mode-cursor > * {
  cursor: crosshair !important;
  user-select: none !important;
}
</style>
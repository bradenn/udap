<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {onMounted, provide, reactive} from "vue";
import {version} from "../package.json";
import {usePersistent} from "@/persistent";

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

const preferences = usePersistent()

onMounted(() => {
    resetCountdown()
})

let screensaver = reactive({
    show: false,
    countdown: 0,
    interval: 0,
    hideTerminal: false,
    startScreensaver: forceScreensaver,
})

provide("screensaver", screensaver)

function forceScreensaver() {
    screensaver.countdown = 3
}

function resetCountdown() {
    screensaver.countdown = preferences.ui.screensaver.countdown
    screensaver.hideTerminal = false
    screensaver.show = false
    if (screensaver.interval !== 0) {
        clearInterval(screensaver.interval)
        screensaver.interval = 0
    }
    if (!preferences.ui.screensaver.enabled) return
    screensaver.interval = setInterval(() => {
        screensaver.countdown -= 1;
        if (screensaver.countdown <= 0) {
            screensaver.show = true
            clearInterval(screensaver.interval)
            screensaver.interval = 0
            setTimeout(() => {
                screensaver.hideTerminal = true
            }, 500)
        }
    }, 1000)
}


provide('system', system)


</script>

<template>

    <div :class="`root theme-${preferences.ui.theme} mode-${preferences.ui.mode}
                  blurs-${preferences.ui.blur} brightness-${preferences.ui.brightness} h-100`"
         @mousedown="() => resetCountdown()">

        <img :class="`${preferences.ui.background.blur?'backdrop-blurred':''}`"
             :src="`/custom/${preferences.ui.background.image}@4x.png`" alt="" class="backdrop "/>


        <div v-if="preferences.ui.watermark" class="watermark">
            <div class="d-flex gap">
                <div v-if="system.udap" class="label-r label-w600">{{ system.udap.system.version }}</div>
                <div v-if="screensaver.countdown <= 5" class="mx-1 label-r label-w500 screensaver-text">
                    <div v-if="screensaver.countdown === 0">Starting screensaver...</div>
                    <div v-else>Screen saver in {{ screensaver?.countdown }} second{{
                            screensaver?.countdown === 1 ? '' : 's'
                        }}
                    </div>
                </div>
            </div>

            <div class="float-end">{{ $route.path }}</div>
        </div>

        <div v-if="preferences.ui.grid" class="grid"></div>
        <div class="h-100">
            <router-view/>
        </div>
    </div>

</template>

<style lang="scss">
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

.screensaver-overlay {
  position: absolute;
  top: 0;
  left: 0;
  z-index: 50 !important;
}

.root {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  border-radius: 0.4rem 0.4rem 0.4rem 0.4rem !important;
  //box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.125) !important;
  background-color: rgba(22, 22, 22, 0.33);
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
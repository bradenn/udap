<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {onMounted, provide, reactive, watch} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types"

interface Preferences {
  background: string,
  theme: string
  mode: string

}

let preferences = reactive<Preferences>({
  background: new Preference(PreferenceTypes.Background).get(),
  theme: new Preference(PreferenceTypes.Theme).get(),
  mode: new Preference(PreferenceTypes.TouchMode).get(),

})

onMounted(() => {
  state.context = false
  preferences.background = new Preference(PreferenceTypes.Background).get()
  preferences.theme = new Preference(PreferenceTypes.Theme).get()
  preferences.mode = new Preference(PreferenceTypes.TouchMode).get()
})

watch(preferences, (newPreferences, oldPreferences) => {
  new Preference(PreferenceTypes.Background).set(newPreferences.background)
  new Preference(PreferenceTypes.Theme).set(newPreferences.theme)
  new Preference(PreferenceTypes.TouchMode).set(newPreferences.mode)
})

provide('preferences', preferences)


let ps: any = {
  hideHome: false,
  timeout: 0,
  countdown: 3,
  context: false,
  grid: true,
  remote: {},
  watermark: true,
  nightVision: false,
  outlines: true,
  system: {
    nexus: {
      system: {
        version: '2.10.2'
      }
    },
    udap: {
      system: {
        version: '0.0.0'
      }
    }
  }
}

let ls = localStorage.getItem("state") || JSON.stringify(ps)
ps = JSON.parse(ls)

let state = reactive(ps)


provide('system', state.system)
watch(state, (change) => {
  console.log("CHANGE TO THE TIMELINE")
  localStorage.setItem("state", JSON.stringify(change))
  if (change.hideHome && state.timeout === 0) {
    state.timeout = setInterval(() => {
      state.countdown -= 0.250
      handleUpdate()
    }, 250)
  }
})

function handleUpdate() {
  if (state.countdown <= 0) {
    clearInterval(state.timeout)
    state.timeout = 0
    state.hideHome = false
  }
}

function hideHome(hide: boolean) {
  state.hideHome = hide
}

provide('ui', state)
provide('hideHome', hideHome)

</script>

<template>
  <div :class="`${state.nightVision?'night-vision':''} theme-${preferences.theme} mode-${preferences.mode}`"
       class="root blurs-7">
    <img :src="`/custom/${preferences.background}@4x.png`" alt="Background" class="backdrop "/>
    <div v-if="state.watermark" class="watermark">
      <div class="d-flex gap">
        <div>NEXUS v{{ state.system.nexus.system.version }}</div>
        <div v-if="state.system.udap">UDAP v{{ state.system.udap.system.version }}</div>
      </div>
      <!--      <div class="float-end">Copyright 2019-2022 &copy; Braden Nicholson</div>-->
    </div>
    <div v-if="state.grid" class="grid"></div>

    <router-view/>

  </div>
</template>

<style lang="scss">
.night-vision {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;

  background-color: rgba(28, 0, 0, 0.6) !important;
}

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

.hide-home {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 0.5rem;
  backdrop-filter: blur(8px);
  z-index: 0;
  background-color: rgba(255, 255, 255, 0.0125);
}

.mode-touch > * {
  cursor: none !important;
  user-select: none !important;
}

.root .entity-small .entity-context {
  filter: blur(88px) !important;
}

// Colors
$bg-color: rgba(0, 0, 0, 0);
$dot-color: rgba(255, 255, 255, 0.2);

// Dimensions
$dot-size: 28px;
$dot-space: 16px;


.grid {
  position: absolute;
  width: 100%;
  height: 100%;

  z-index: -1;
  background-color: $bg-color;
  opacity: 1;
  background-image: radial-gradient($dot-color 0.980000000000000px, rgba(0, 0, 0, 0.125) 0.980000000000000px), radial-gradient($dot-color 0.98px, $bg-color 0.980000000000000px);
  background-size: $dot-size $dot-size;
  background-position: 0 0, calc($dot-size / 2) calc($dot-size / 2);

}

.mode-cursor > * {
  cursor: crosshair !important;
  user-select: none !important;
}
</style>
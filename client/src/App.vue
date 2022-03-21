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

</script>

<template>
  <div :class="` theme-${preferences.theme} mode-${preferences.mode}`" class="root blurs-1 ">
    <img :src="`/custom/${preferences.background}@4x.png`" alt="Background" class="backdrop "/>
    <div class="grid"></div>
    <router-view/>
  </div>
</template>


<style lang="scss" scoped>
.mode-touch > * {
  cursor: none !important;
  user-select: none !important;
}

// Colors
$bg-color: rgba(0, 0, 0, 0);
$dot-color: rgba(0, 0, 0, 0.5);

// Dimensions
$dot-size: 1px;
$dot-space: 28px;

.grid {
  position: absolute;
  width: 100%;
  height: 100%;
  z-index: -1;
  background-color: rgba(0, 0, 0, 0);
  background-image: radial-gradient(circle, rgba(255, 255, 255, 0.1) 20%, transparent 10%), radial-gradient(circle, rgba(255, 255, 255, 0.1) 20%, transparent 10%);
  background-size: 16px 16px;
  background-position: 0 0, 0 0;
}

.mode-cursor > * {
  cursor: crosshair !important;
  user-select: none !important;
}
</style>
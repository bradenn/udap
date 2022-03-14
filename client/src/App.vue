<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {provide, reactive} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types"

interface Preferences {
  background: string,
  theme: string
}

let preferences = reactive<Preferences>({
  background: new Preference(PreferenceTypes.Background).get(),
  theme: new Preference(PreferenceTypes.Theme).get(),
})

provide('preferences', preferences)

</script>

<template>
  <div :class="`theme-${preferences.theme}`" class="root blurs-8">
    <img :src="`/custom/${preferences.background}@4x.png`" alt="Background" class="backdrop"/>
    <router-view/>
  </div>
</template>


<style scoped>

</style>
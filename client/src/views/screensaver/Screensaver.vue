<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>


import core from "@/core";
import Bubbles from "@/views/screensaver/Bubbles.vue";
import Warp from "@/views/screensaver/Warp.vue";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Preferences} from "@/types";

const screensaver = core.screensaver()
let preferences = inject("preferences") as Preferences

const state = reactive({
    current: ""
})

onMounted(updateAttributes)
watchEffect(updateAttributes)

function updateAttributes() {
    let now = preferences.ui.screensaver
    state.current = now.selection
    return preferences
}


</script>

<template>
    <div v-if="screensaver.active()" class="screensaver-container" @click="screensaver.stop()">
        <Bubbles v-if="state.current === 'bubbles'"
                 class="screensaver-overlay"></Bubbles>
        <Warp v-else-if="state.current === 'warp'"
              class="screensaver-overlay"></Warp>
    </div>
    <div v-else>
        {{ screensaver.countdown() }}
    </div>

</template>

<style scoped>
.screensaver-container {

}

.screensaver-overlay {
    position: absolute;
    top: 0;
    left: 0;
    z-index: 50 !important;
}

</style>
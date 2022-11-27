<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject} from "vue";
import type {Haptics} from "@/haptics";

interface App {
    name: string
    icon?: string
    status?: string
    img?: string
}

let props = defineProps<App>()

const haptics = inject("haptics") as Haptics

function click() {
    haptics.tap(2, 2, 100)
    // haptics.tap(0, 2, 100)
}

</script>

<template>
    <div class="app-container" @mousedown="click">
        <!--    <img v-if="props.img" :src="`ven/${props.img}`" class="app-icon" alt=""/>-->
        <div class="app-icon element">
            <div v-if="props.status === 'wip'" class="marker">•</div>

            <i v-if="props.icon" :class="props.icon.includes('fa')?`fa-solid ${props.icon} fa-fw '`:' lh-1 label-sm'"
               class="app-icon-char"
               v-text="!props.icon.includes('fa')?props.icon:''"></i>
            <i v-else class="fa-solid fa-circle fa-fw"></i>
        </div>
        <div class="app-name pt-1">{{ props.name }}</div>
    </div>
</template>

<style lang="scss" scoped>
.marker {
  position: absolute;
  top: -0.25rem;
  left: 0.25rem;
  color: rgba(255, 128, 12, 0.4)
}

.app-container {
  display: flex;
  justify-content: center;
  flex-direction: column;
  align-items: center;

}

div .app-icon.element {

  box-shadow: inset 0 0 1px 2px rgba(255, 255, 255, 0.025), inset 0 0 1px 1px rgba(0, 0, 0, 0.125), inset 0 1px 1px 1px rgba(255, 255, 255, 0.04), inset 0 -1px 1px 1px rgba(128, 128, 128, 0.0125) !important;
}

.app-icon:after {

}

.app-icon {
  width: 2.5rem;
  height: 2.5rem;
  display: flex;
  align-self: center;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1/1 !important;
  border-radius: 0.54rem !important;
  font-size: 0.9rem;
  color: rgba(255, 255, 255, 0.5);
  filter: drop-shadow(0px 10px 60px rgba(0, 0, 0, 0.1)) !important;
}

.app-icon-char {
  filter: drop-shadow(0px 10px 60px rgba(0, 0, 0, 0.1));

}

.app-container:active {


}

.app-name {
  font-family: 'SF Pro Text', sans-serf, serif;
  font-style: normal;
  font-weight: 500;
  font-size: 16px;
  line-height: 16px;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.6);
  mix-blend-mode: luminosity;
}
</style>
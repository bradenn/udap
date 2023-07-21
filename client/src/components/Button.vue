<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {reactive} from "vue";
import core from "@/core";


const props = defineProps<{
    active: boolean,
    text: string,
    to?: string,
    icon?: string,
    accent?: boolean,
    fill?: boolean
    color?: string
}>()

let state = reactive({})
let router = core.router()

const haptics = core.haptics()

function push(_: Event) {
    haptics.tap(2, 1, 75)
    if (props.to) {
        router.push(props.to)
    }
}

function release(_: Event) {
    // haptics(1, 1, 75)
    // haptics(1, 1, 5)
}


</script>

<template>
    <div
            :class="`${props.active?props.accent?'text-accent':'':'subplot-inline'}${props.fill?'element':''} ${props.color}` "
            class="switch-space subplot switch-space-active px-3"
            @mousedown="push"
            @mouseup="release">
        <div v-if="props.icon" :class="`${props.color}`" class="label-o2 opacity-75" style="padding-right: 0.33rem">{{
            props.icon
            }}
        </div>
        <span :class="`${props.active?'':''} ${props.color}`">{{ props.text }}</span>
    </div>
</template>

<style lang="scss" scoped>

.switch-space-active {
  //color: rgba(255, 255, 255, 0.5) !important;

  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);
}

.switch-right {
  animation: switch-move-right 150ms forwards ease-out;
}

.switch-space:active {
  transform: scale3d(0.96, 0.96, 1);
}

.element.switch-space:active {
  transform: scale3d(0.96, 0.96, 1) !important;
}


.switch-space {
  display: flex;
  justify-content: center;
  font-weight: 600;
  align-items: center;
  font-size: 0.65rem;
  line-height: 0.7rem;
  height: 1.7rem;
  color: rgba(255, 255, 255, 0.4);
}
</style>
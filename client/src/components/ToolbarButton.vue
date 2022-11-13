<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {reactive} from "vue";
import {useRouter} from "vue-router";
import core from "@/core";


const props = defineProps<{
    text: string,
    to?: string,
    disabled?: boolean,
    icon?: string,
    active?: boolean,
    accent?: boolean,
}>()

let state = reactive({})
let router = useRouter()

const haptics = core.haptics()

function push(_: Event) {
    if (props.disabled) return
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
    <div :class="`${props.active?'':'subplot-inline'}`"
         class="switch-space subplot switch-space-active"
         @mousedown="push"
         @mouseup="release">
        <div v-if="props.icon" class="label-o2 label-c2" style="padding-right: 0.125rem">{{ props.icon }}</div>
        <span
                :class="`${props.active?'':''} ${props.accent?'switch-accent':''} ${props.disabled?'label-o1':''} label-c2`">{{
                props.text
            }}</span>
    </div>
</template>

<style lang="scss" scoped>
.switch-accent {
  color: rgba(255, 149, 0, 0.8) !important;
}

.switch-space-active {
  color: rgba(255, 255, 255, 0.4) !important;

  box-shadow: 0 0 20px 1px rgba(0, 0, 0, 0.1), inset 0 0 2px 1px rgba(255, 255, 255, 0.05);
}

.switch-right {
  animation: switch-move-right 150ms forwards ease-out;
}

.switch-space:active {
  transform: scale3d(0.96, 0.96, 1);
}

.switch-space {
  display: flex;
  justify-content: center;
  font-family: "SF Pro Display", serif !important;
  font-weight: 500;
  align-items: center;
  font-size: 21px;
  line-height: 0.7rem;
  height: 1.8rem;

}
</style>
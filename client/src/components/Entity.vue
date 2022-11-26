<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Attribute, Entity} from "@/types";
import {onMounted, reactive, watchEffect} from "vue";
import {useRouter} from "vue-router";
import core from "@/core";
import moment from "moment/moment";

interface EntityProps {
    entity: Entity,
    selected?: boolean
}

const state = reactive({
    toggle: false,
    activeColor: ""
})

const router = useRouter()
const haptics = core.haptics()
const remote = core.remote()

onMounted(() => {
    findMode()
})

watchEffect(() => {
    findMode()
    return remote.attributes
})

function click() {
    haptics.tap(1, 1, 25)
}

function map_range(value: number, low1: number, high1: number, low2: number, high2: number) {
    return low2 + (high2 - low2) * (value - low1) / (high1 - low1);
}


function cctToRgb(cct: number) {
    return [map_range(cct, 2000, 8000, 255, 227),
        map_range(cct, 2000, 8000, 138, 233),
        map_range(cct, 2000, 8000, 18, 255)]
}

function findMode() {
    let hue = remote.attributes.find((a: Attribute) => a.key === 'hue')
    let cct = remote.attributes.find((a: Attribute) => a.key === 'cct')
    let dim = remote.attributes.find((a: Attribute) => a.key === 'dim')
    if (hue && cct && dim) {
        if (moment(hue.requested).isAfter(cct.requested)) {
            state.activeColor = `hsla(${hue.value}, 100%, ${20 + parseInt(dim.value) / 100.0 * 50}%, 0.5)`
        } else {
            state.activeColor = `rgba(${(cctToRgb(parseInt(cct.value)))[0]}, ${(cctToRgb(parseInt(cct.value)))[1]}, ${(cctToRgb(parseInt(cct.value)))[2]}, 0.5)`
        }
    } else {
        state.activeColor = 'rgba(255,255,255,0.5)'
    }
}

const props = defineProps<EntityProps>()


</script>

<template>
    <div @mousedown="() => click()">

        <div :class="props.selected?'accent-selected':''" class="element p-2 h-100">
            <div class="d-flex justify-content-between">
                <div class="label-c2 label-o2 label-w500 pb-2">
                    {{ props.entity.icon }}
                </div>
                <div v-if="props.selected" class="label-c2 label-o4 label-w500 text-accent">􀷙</div>
                <div v-else class="label-c2 label-o1 label-w500">􀓞</div>
            </div>

            <div class="label-c2 label-o4 label-w700 lh-1">
                {{ props.entity.name }}
            </div>
            <div class="label-c3 label-o3 label-w400">{{ props.entity.module }}
            </div>

        </div>
        <div v-if="state.toggle" class="element element-menu">
            <div class="grid-element">

            </div>
        </div>
    </div>
</template>

<style lang="scss" scoped>

div.element {

}

div.element:active {
  transform: scale(98%) !important;
}

.element-menu {
  position: absolute;
  width: 8rem;

  z-index: 1000;
  margin-top: 0.25rem;

  /*left: calc(-100% - 0.25rem);*/
  /*top: calc(100% + 0.25rem);*/
}

.grid-fill {
  position: relative;
  grid-column: 3 / 5;
  display: flex;
  flex-direction: column;
}

.grid-element {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}
</style>
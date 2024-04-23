<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Attribute} from "@/types";
import {onMounted, reactive} from "vue";
import {useRouter} from "vue-router";
import core from "@/core";
import Radio from "@/components/plot/Radio.vue";
import moment from "moment";

interface AttributeProps {
  attribute: Attribute,
  selected?: boolean
  noselect?: boolean
}

const props = defineProps<AttributeProps>()
const state = reactive({
  toggle: false,
  activeColor: ""
})

const router = useRouter()
const haptics = core.haptics()
const remote = core.remote()

interface KeyPayload {
  icon: string
  unit: string
}

const typeMap = new Map<string, KeyPayload>([
  ["on", {icon: "􀆨", unit: ""}],
  ["media", {icon: "􁇵", unit: ""}],
  ["dim", {icon: "􀇮", unit: "%"}],
  ["cct", {icon: "􀆭", unit: "K"}],
  ["hue", {icon: "􀟗", unit: "°"}],
  ["api", {icon: "􁉢", unit: ""}],
  ["online", {icon: "􀩲", unit: ""}]])


onMounted(() => {

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


</script>

<template>
  <div @mousedown="() => click()">

    <div :class="props.selected?'accent-selected':''"
         class="element p-1 h-100 d-flex justify-content-between">
      <div class="d-flex justify-content-between align-items-center gap-1 w-100">
        <div
            class="label-c1 label-o2 label-w500 px-1 d-flex justify-content-center"
            style="width: 1.5rem">
          {{ typeMap.get(props.attribute.key)?.icon || typeMap.get(props.attribute.type)?.icon }}
        </div>
        <div class="d-flex justify-content-between w-100 align-items-center">
          <div class="d-flex flex-column gap-0">
            <div class="label-c2 label-o4 label-w700 lh-1 d-flex justify-content-between">
              <div>{{ props.attribute.key }}</div>

            </div>

            <div v-if="props.attribute.type !== 'media'"
                 class="label-c3 label-o3 label-w400 lh-sm"
                 style="max-width: 7rem;  white-space: nowrap ; overflow: hidden; text-overflow: ellipsis; padding-right: 0.125rem">
              {{
                props.attribute.value
              }}{{ typeMap.get(props.attribute.key)?.unit }}
            </div>
          </div>
          <div class="px-2">
            <div class=" label-c3 label-o4 label-w600 ">Updated:</div>
            <div class=" label-c3 label-o3 label-w500 lh-1">
              {{
                moment(props.attribute.lastUpdated).fromNow()
              }}
            </div>
          </div>
        </div>
      </div>

      <div class="d-flex gap-1" style="min-height: 1.6rem; max-height: 1.6rem;">
        <Radio :active="false" :disabled="props.attribute.type === 'media'"
               :fn="() => {}" sf="􁎵"
               style="width: 3rem;"
               title=""></Radio>
        <Radio :active="false" :fn="() => {}" sf="􀈑"
               style="width: 3rem;"
               title=""></Radio>
      </div>
    </div>
    <div v-if="props.attribute.type === 'media'"
         class="label-c3 label-o3 label-w400 lh-sm d-flex flex-column gap-1 py-1" style="padding-left: 0.5rem;">
      <div v-for="key in Object.keys(JSON.parse(props.attribute.value))"
           class="element p-2 w-100 d-flex justify-content-between">
        <div class="label-c2 label-o3 label-w600 lh-sm text-capitalize">{{ key }}</div>
        <div class="label-c2 label-o3 label-w500 lh-sm text-accent" style="overflow-x: clip; max-width: 75%">
          {{ JSON.parse(props.attribute.value)[key] }}
        </div>

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

//div.element:active {
//  transform: scale(98%) !important;
//}

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
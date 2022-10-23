<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {inject} from "vue";
import type {Haptics} from "@/views/terminal/haptics";

interface Props {
  icon?: string,
  sf?: string,
  name: string
  alt?: string
  active?: boolean
  theme?: string
  to?: string
  fn?: () => void
}

const props = defineProps<Props>()
const haptics = inject("haptics") as Haptics

function click() {
  haptics.tap(2, 2, 100)
  // haptics.tap(0, 2, 100)
}

function runFn() {
  click()
  if (props.fn) {
    props.fn()
  }
}

</script>

<template>
  <div v-if="props.to"
       :class="`${props.to===$router.currentRoute.value.fullPath?'':'subplot-inline'} ${props.theme?`theme-${props.theme}`:''}`"
       class="subplot p-1"
       @click="$router.replace(props.to || '/')" @mousedown="click">
    <div class="d-flex justify-content-start px-1">
      <div v-if="props.icon" class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i>
      </div>
      <div class="label-w500 label-c1 px-2 text-center">{{ props.name }}</div>
    </div>

    <slot></slot>
  </div>
  <div v-else-if="props.fn"
       :class="`${props.active||false?'':'subplot-inline'} ${props.theme?`theme-${props.theme}`:''}`"
       class="subplot p-1"
       @mousedown="runFn">
    <div :class="`${props.alt?'justify-content-between w-100':'justify-content-center w-100'}`"
         class="d-flex ">
      <div v-if="props.icon" class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i>
      </div>
      <div v-if="sf" v-html="sf"></div>
      <div class="label-w500 label-c1 px-2 label-o4 text-center">{{ props.name }}</div>
      <div v-if="props.alt" class="label-w400 label-o4 label-c1 px-2">{{ props.alt }}</div>


    </div>

    <slot></slot>
  </div>
  <div v-else class="subplot subplot-inline" @mousedown="click">

    <div class="sidebar-item d-flex justify-content-start px-1">
      <div v-if="props.icon" class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i>
      </div>
      <div v-if="props.sf" class="label-w500 label-o3 label-c1">{{ props.sf }}
      </div>
      <div class="label-w500 label-o5 label-c1 px-2">{{ props.name }}</div>
    </div>

    <slot></slot>
  </div>
</template>


<style scoped>

.theme-disabled {
  /*background-color: rgba(24, 24, 24, 0.25) !important;*/
  /*box-shadow: inset 0 0 3px 2px rgba(48, 48, 48, 0.25) !important;*/
  opacity: 0.4;
}

.theme-danger {
  background-color: rgba(255, 0, 0, 0.25) !important;
  box-shadow: inset 0 0 3px 2px rgba(255, 0, 0, 0.5) !important;
}

.subplot {
  justify-content: center;
}


</style>
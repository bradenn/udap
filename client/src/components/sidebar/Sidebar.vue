<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

let props = defineProps<{
  small?: boolean,
  os?: boolean,
  horizontal?: boolean
  icon?: string
  name?: string
}>()

</script>

<template>
  <div v-if="horizontal" class="d-flex flex-column flex-row align-content-start">

    <div v-if="!small" class="d-flex justify-content-start px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">{{ props.name }}</div>
    </div>
    <div class="element main-sidebar d-flex  align-items-center gap-1"
         v-bind:class="`dock-small ${props.os?'dock-os':''} flex-row` ">
      <div v-if="small" class="">
        <div class="sidebar-item d-flex justify-content-start px-1">
          <div class="label-w500 label-o3 label-md"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
          <div class="label-w500 label-o5 label-md px-2">{{ props.name }}</div>
        </div>
      </div>
      <slot></slot>
    </div>
  </div>
  <div v-else class="d-flex flex-column sidebar-width align-content-start">

    <div v-if="!small" class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">{{ props.name }}</div>
    </div>
    <div class="element main-sidebar d-flex flex-column gap-1"
         v-bind:class="`${props.small?'dock-small':''} ${props.os?'dock-os':''}` ">
      <div v-if="small" class="pb-1">
        <div class="sidebar-item d-flex justify-content-start pb-1 px-1">
          <div class="label-w500 label-o4 label-xs"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
          <div class="label-w500 label-o6 label-xxs px-2">{{ props.name }}</div>
        </div>
        <div class="h-sep"></div>
      </div>
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss">
.sidebar-width {
  width: calc(100% / 5.5);
}
</style>

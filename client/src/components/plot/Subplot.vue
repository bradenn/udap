<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

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


</script>

<template>
  <div v-if="props.to" :class="`${props.to===$router.currentRoute.value.fullPath?'':'subplot-inline'}`"
       class="subplot p-1"
       @mouseover="$router.replace(props.to || '/')">
    <div class="d-flex justify-content-start px-1">
      <div class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
      <div class="label-w500 label-c1 px-2">{{ props.name }}</div>
    </div>

    <slot></slot>
  </div>
  <div v-else-if="props.fn"
       :class="`${props.active||false?'':'subplot-inline'} ${props.theme?`theme-${props.theme}`:''}`"
       class="subplot p-1"
       @mousedown="props.fn">
    <div class="d-flex justify-content-start ">
      <div v-if="props.icon" class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i>
      </div>
      <div class="label-w500 label-c1 px-2">{{ props.name }}</div>


    </div>

    <slot></slot>
  </div>
  <div v-else class="subplot subplot-inline">

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

.theme-danger {
  background-color: rgba(255, 0, 0, 0.25) !important;
  box-shadow: inset 0 0 3px 2px rgba(255, 0, 0, 0.5) !important;
}

.subplot:active {
  animation: click 100ms ease forwards !important;
}

.subplot:hover {
  animation: click 100ms ease forwards !important;
}

.subplot {
  justify-content: center;
}


</style>
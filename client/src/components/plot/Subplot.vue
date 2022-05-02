<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

interface Props {
  icon?: string,
  sf?: string,
  name: string
  alt?: string
  active?: boolean
  to?: string
  fn?: () => {}
}

const props = defineProps<Props>()


</script>

<template>
  <div v-if="props.to" :class="`${props.to===$router.currentRoute.value.fullPath?'':'subplot-transparent'}`"
       class="subplot p-1"
       @mouseover="$router.replace(props.to || '/')">
    <div class="d-flex justify-content-start px-1">
      <div class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
      <div class="label-w500 label-c1 px-2">{{ props.name }}</div>
    </div>

    <slot></slot>
  </div>
  <div v-else-if="props.fn" :class="`${props.active||false?'':'subplot-transparent'}`" class="subplot p-1"
       @mousedown="props.fn">
    <div class="d-flex justify-content-start ">
      <div v-if="props.icon" class="label-w500 label-o3 label-c1"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i>
      </div>
      <div class="label-w500 label-c1 px-2">{{ props.name }}</div>


    </div>

    <slot></slot>
  </div>
  <div v-else class="subplot subplot-transparent">

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

.subplot:active {
  animation: click 100ms ease forwards !important;
}

.subplot:hover {
  animation: click 100ms ease forwards !important;
}

.subplot-transparent {
  background-color: transparent;
  border-color: transparent;
  box-shadow: none;
  color: rgba(255, 255, 255, 0.45);
}

.subplot {
  justify-content: center;
}

.subplot:active {
  animation: click 100ms ease forwards;
  color: rgba(255, 255, 255, 0.6);
}


@keyframes click {
  0% {
    opacity: 0.6;

    transform: scale(1.0);
  }
  25% {
    transform: scale(0.98);
  }
  30% {
    transform: scale(0.97);
  }
  100% {
    transform: scale(1);
  }
}


</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Scroll from "@/components/scroll/Scroll.vue";

interface Plot {
  title?: string
  alt?: string
}

let props = defineProps<Plot>()


</script>

<template>
  <div class="element p-1  h-100" v-on:mousemove.stop>

    <div v-if="props.title" class="d-flex align-items-center justify-content-between">
      <div class="label-c1  label-o4 label-w500 px-1">{{ props.title }}</div>
      <div class="label-c2  label-o3 label-w400 px-1">{{ props.alt }}</div>
    </div>

    <Scroll class="pane-menu">
      <slot></slot>
    </Scroll>
  </div>
</template>


<style lang="scss" scoped>
.pane-menu {
  overflow-y: auto;
  height: 96%;
}
</style>
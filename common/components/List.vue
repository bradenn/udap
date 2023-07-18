<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps} from "vue"

const props = defineProps<{
  row?: boolean
  scrollX?: boolean
  scrollY?: boolean
  scrollLock?: boolean
}>()


</script>

<template>
  <div
      :class="`element-list ${props.row?'flex-row':'flex-column'} ${props.scrollLock?'element-list-scroll-lock':''} ${props.scrollX?'element-list-scroll-x':''} ${props.scrollY?'element-list-scroll-y':''}`">
    <slot></slot>
  </div>
</template>

<style lang="scss">
.element-list-scroll-x {
  overflow-x: scroll;

  .subplot .element {
    flex-shrink: 1;
  }
}

.element-list-scroll-lock.element-list-scroll-y > .element-header {
  scroll-snap-align: start;
}

.element-list-scroll-lock.element-list-scroll-y > .subplot {
  scroll-snap-align: start;
}

.element-list-scroll-lock {
  scroll-snap-type: y mandatory;


  .element-list {
    scroll-snap-align: start;
  }
}

.element-list-scroll-y {
  overflow-y: scroll;
  overscroll-behavior-y: auto;
  flex-wrap: nowrap;

}

.element-list {
  display: flex;
  //flex: 1 1 auto;
  gap: 0.25rem;

}


</style>
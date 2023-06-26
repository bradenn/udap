<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps} from "vue"

const props = defineProps<{
  row?: boolean
  scrollX?: boolean
}>()


</script>

<template>
  <div :class="`element-list ${props.row?'flex-row':'flex-column'} ${props.scrollX?'element-list-scroll-x':''}`">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.element-list-scroll-x {
  overflow-x: scroll;
}

.element-list {
  display: flex;
  gap: 0.25rem;
}

</style>
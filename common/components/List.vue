<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps} from "vue"

const props = defineProps<{
  row?: boolean
}>()


</script>

<template>
  <div :class="`element-list ${props.row?'flex-row':'flex-column'}`">
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>

.element-list {
  display: flex;
  gap: 0.25rem;
}

</style>
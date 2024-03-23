<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps} from "vue"

const props = defineProps<{
  title?: string
  subtitle?: string
}>()
</script>

<template>
  <div class="element-header">
    <div class="d-flex justify-content-between align-items-end">
      <div class="d-flex align-items-center gap-1">
        <div class="label-c4 label-o4 label-w600" v-text="props.title"></div>
        <div v-if="props.subtitle" class="label-c5 label-o4 label-w600">- {{ props.subtitle }}</div>
      </div>
      <div class="d-flex align-items-center gap-2 label-c5 label-o4 label-w500">
        <slot name="description"></slot>
      </div>
    </div>
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>
.element-header {
  padding: 0.125rem 0.5rem;
  user-select: none !important;
  -webkit-user-select: none !important;
}
</style>
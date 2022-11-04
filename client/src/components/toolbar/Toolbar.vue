<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

interface Props {
  title: string,
  icon: string
}

const props = defineProps<Props>()

</script>

<template>
  <div class="px-2 d-flex gap-1 align-items-center element w-100 lh-2" style="height: 1.8rem;">
    <div class="label-xs label-w200 label-o3 px-0">{{ props.icon }}</div>
    <div class="label-xs label-w700 label-o5">{{ props.title }}</div>
    <slot></slot>
  </div>
</template>

<style lang="scss" scoped>

</style>
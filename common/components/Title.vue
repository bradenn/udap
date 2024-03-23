<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

let props = defineProps<{
  title: string,
  alt?: string
}>()


</script>

<template>
  <div class="d-flex justify-content-between">
    <div class="label-c6 label-o4 px-2 label-w600 lh-1 mb-1">{{ props.title }}</div>
    <div v-if="props.alt" class="label-c6 label-o4 px-2 label-w600 lh-1 mb-1">{{ props.alt }}</div>
  </div>
</template>

<style lang="scss" scoped>

</style>
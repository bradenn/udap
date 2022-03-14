<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

let props = defineProps<{
  small?: boolean,
  os?: boolean,
}>()

</script>

<template>
  <div class="element main-sidebar d-flex flex-column gap-1"
       v-bind:class="`${props.small?'dock-small':''} ${props.os?'dock-os':''}`">
    <slot></slot>
  </div>
</template>

<style lang="scss">

</style>

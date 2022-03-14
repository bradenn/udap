<script lang="ts" setup>

let props = defineProps<{
  small?: boolean,
  os?: boolean,
}>()

</script>

<template>
  <div class="element dock" v-bind:class="`${props.small?'dock-small':''} ${props.os?'dock-os':''}`">
    <slot></slot>
  </div>
</template>

<style lang="scss">

</style>

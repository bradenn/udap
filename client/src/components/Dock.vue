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

.dock {
  @extend .surface;
  display: flex;
  flex-direction: row;

  justify-content: space-between;
  align-content: center;
  align-items: center;

  opacity: 1;
  height: 2.5rem;
  border-radius: 0.6rem !important;
  padding: 0.25rem 0.25rem !important;

}
</style>

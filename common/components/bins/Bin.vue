<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "../Element.vue";


function dragStart(ev: DragEvent) {
  // Add the target element's id to the data transfer object
  ev.dataTransfer.setData("application/my-app", "sdds");
  ev.dataTransfer.effectAllowed = "move";
}

</script>

<template>
  <div draggable="true" @dragstart="dragStart">
    <Element :foreground="true">
      <slot></slot>
    </Element>
  </div>
</template>

<style lang="scss" scoped>

</style>
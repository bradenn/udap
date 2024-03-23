<!-- Copyright (c) 2023 Braden Nicholson -->

<script setup>

</script>

<template>
  <div class="viewport">
    <div class="viewport-layout">
      <slot></slot>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.viewport-layout {

  display: flex;
  flex-direction: row;
  height: 100%;
  gap: 0.375rem;
}

.viewport {
  display: block;
  overflow: auto;
  width: 100%;
  box-shadow: inset 0 0 20px 20px red;
}
</style>
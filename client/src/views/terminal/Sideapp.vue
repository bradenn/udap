<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

</script>


<template>
  <div class="side-app">
    <slot></slot>
  </div>
</template>

<style scoped>
.side-app {
  margin-left: 1rem;
  position: absolute;
  right: -14rem;
  height: calc(100% - 2rem);
  width: 13rem;

  z-index: 22 !important;

  animation: popIn 150ms cubic-bezier(.13, .45, .84, 1.19);

}


</style>
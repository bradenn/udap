<script lang="ts" setup></script>

<template>
  <div class="context sessions">
    <div class="session element"></div>
    <div class="session"></div>
    <div class="session"></div>
    <div class="session"></div>
    <div class="session"></div>
  </div>
</template>

<style lang="scss" scoped>
.sessions {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-columns: repeat(4, minmax(8rem, 1fr));
  grid-template-rows: repeat(3, minmax(10rem, 1fr));
}

.session {
  width: 100%;
  height: 100%;
  //outline: 1px solid white;
}
</style>

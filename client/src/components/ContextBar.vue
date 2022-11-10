<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

</script>

<template>
  <div class="content-bar d-flex justify-content-between align-items-center pt-0 mt-0">
    <slot>

    </slot>

  </div>

</template>

<style lang="scss" scoped>


.content-bar > div {

  padding: 0.125rem 0.125rem;
}

.content-bar {
  display: flex;
  align-items: center;
  //outline: 1px solid white;
  //backdrop-filter: blur(24px);
  //box-shadow: inset 0 -2px 2px 0.25px rgba(255, 255, 255, 0.1);
}
</style>
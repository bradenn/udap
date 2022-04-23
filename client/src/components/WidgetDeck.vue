<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

</script>

<template>
  <div>
    <div>
      <slot></slot>
    </div>
    <div class="dots">
      <div class="dot"></div>
      <div class="dot"></div>
      <div class="dot"></div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
$radius: 0.3rem;
.dots {
  display: flex;
  justify-content: center;
  gap: 0.125rem;

}

.dot {
  border-radius: 1rem;
  width: $radius;
  height: $radius;
  background-color: rgba(255, 255, 255, 0.3);
}
</style>
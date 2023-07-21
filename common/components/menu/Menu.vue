<!-- Copyright (c) 2023 Braden Nicholson -->

<script setup>

import Element from "../Element.vue";
import ElementHeader from "../ElementHeader.vue";
</script>

<template>
  <Element>
    <ElementHeader title="Words">
      <template v-slot:description>Listed</template>
    </ElementHeader>
    <div class="menu-content">
      <slot>
      </slot>
    </div>


  </Element>
</template>

<style lang="scss" scoped>
.menu-content:not(:first-of-type, :last-of-type) {
  border: 1px solid white;
}
</style>
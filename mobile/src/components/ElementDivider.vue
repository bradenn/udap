<!-- Copyright (c) 2024 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";

const props = defineProps<{
  monogram: string
  text: string
  sm?: boolean
  md?: boolean
  lg?: boolean
}>()

</script>

<template>

  <Element class="divider d-flex align-items-center gap-1" foreground>
    <div class="monogram">{{ monogram }}</div>
    <div v-if="text">{{ text }}</div>
    <slot v-else></slot>
  </Element>
</template>

<style lang="scss" scoped>
.divider {
  padding: 8px 8px !important;
}

.monogram {
  width: 1rem;
  height: 1rem;
  display: flex;
  align-items: center;
  text-align: center;
  justify-content: center;
  line-height: 1rem;
  aspect-ratio: 1/1;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 18%;
  font-size: 0.6rem;

}
</style>
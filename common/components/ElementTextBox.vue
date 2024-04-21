<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "./Element.vue";

const props = defineProps<{
  title?: string
  change: (val: string) => void
}>()

function update(et: Event) {
  let e = et as InputEvent
  let val = e.target as HTMLInputElement
  if (!val) return
  props.change(val.value)
}

</script>

<template>
  <Element foreground>
    <input :placeholder="props.title?props.title:'Name'" class="element-text-input" type="text" @input="update"/>
  </Element>
</template>

<style lang="scss" scoped>


.element-text-input:focus-visible {
  border: none;
  outline: none;
}

.element-text-input {
  text-shadow: none !important;
  width: 100%;
  height: 2rem;
  padding-inline: 0.75rem;
  background-color: transparent;
  box-shadow: none;
  border: none;
  color: rgba(255, 255, 255, 0.8);
}
</style>
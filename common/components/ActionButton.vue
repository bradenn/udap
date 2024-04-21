<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "./Element.vue";
import {reactive} from "vue";

const props = defineProps<{
  icon?: string,
  cb: () => void
}>()

const state = reactive({
  pressed: false
})

function execute() {
  state.pressed = true
  props?.cb()
  state.pressed = false
}

</script>

<template>
  <Element :cb="execute" :class="`${state.pressed?'running':''}`"
           class="d-flex justify-content-center align-items-center"
           foreground mutable>

    <slot></slot>
  </Element>
</template>

<style lang="scss" scoped>
.running {
  filter: brightness(90%) blur(10px);

}
</style>
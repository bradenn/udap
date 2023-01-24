<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {reactive} from "vue";

const props = defineProps<{
    color: number
    selected?: boolean
}>()


const state = reactive({
    value: 0,
})


</script>

<template>
    <div :class="`${props.selected?'selected':''}`" class="color subplot surface">
        <div :style="`background-color: hsl(${props.color}, 100%, 50%);`" class="swatch">
        </div>
    </div>
</template>

<style lang="scss" scoped>
div.color.subplot.surface {

  border-radius: 0.25rem;
  padding: 0.125rem;

  aspect-ratio: 2/2 !important;
  //width: 100%;
  background-color: rgba(255, 255, 255, 0.025);
  height: 2.25rem;


}

div.swatch {
  border-radius: 0.1875rem;
  width: 100%;
  height: 100%;
  opacity: 0.6;
  display: flex;
  justify-content: center;
  align-items: center;

}


div.color.subplot.surface.selected {
  z-index: 10;
  padding: 0;

  div.swatch {
    border-radius: 0.25rem;
    box-shadow: inset 0 0 1px 2px rgba(255, 255, 255, 0.75), inset 0 0 6px 2px rgba(0, 0, 0, 0.25) !important;

  }


}
</style>
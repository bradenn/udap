<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

const props = defineProps<{
    color: string
}>()
</script>

<template>
    <div :style="`background-color: ${props.color};`" class="color">

    </div>
</template>

<style lang="scss" scoped>
div.color {
  opacity: 1;
  border-radius: 0.25rem;
  padding: 0.25rem;

  aspect-ratio: 1.61803398874/1 !important;
  //width: 100%;

  height: 2.2rem;


}


div.color.selected {
  z-index: 10;
  opacity: 1;

  box-shadow: 0 0 32px 8px rgba(0, 0, 0, 0.5), inset 0 0 1px 1px rgba(255, 255, 255, 0.3);
  border-radius: 0.25rem;
  gap: 0.25rem;
  scale: 1.1;
  padding: 0.25rem;

}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  alt?: string
  small?: boolean
  rows: number
  cols: number
}

let props = defineProps<Plot>()
</script>

<template>
  <div class="element element-group">
    <div v-if="props.title" class="d-flex align-items-center justify-content-between">
      <div class="label-c1  label-o4 label-w500 px-1 pb-1">{{ props.title }}</div>
      <div class="label-c1 label-o2">{{ props.alt }}</div>
    </div>
    <div :class="`${props.small?'plot-sm':'plot'}`" class="">
      <slot></slot>
    </div>
  </div>
</template>


<style scoped>
.plot {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-columns: repeat(v-bind('props.cols'), minmax(1rem, 1fr));
  grid-template-rows: repeat(v-bind('props.rows'), minmax(1.75rem, 1fr));
}


.plot-sm {
  display: grid;
  grid-gap: 0.25rem;
  grid-template-columns: repeat(v-bind('props.cols'), minmax(1rem, 1fr));
  grid-template-rows: repeat(v-bind('props.rows'), minmax(1.5rem, 1fr)) !important;
}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  alt?: string
  small?: boolean
  flex?: boolean
  altFn?: () => void
  rows: number
  cols: number
}

let props = defineProps<Plot>()
</script>

<template>
  <div class="element element-group d-flex flex-column">
    <div v-if="props.title" class="d-flex align-items-center justify-content-between">
      <div class="label-c1  label-o4 label-w500 px-1 pb-1">{{ props.title }}</div>
      <div class="label-c2 label-w500 label-r label-o2 px-1" @click="(e) => altFn?altFn():() => {}"
           v-html="props.alt"></div>
    </div>
    <div :class="`${props.flex?'plot-flex':props.small?'plot-sm':'plot'}`">
      <slot></slot>
    </div>
  </div>
</template>


<style scoped>
.plot-flex {
  display: block;
  height: calc(100%);
}

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
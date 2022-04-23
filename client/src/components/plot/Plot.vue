<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  rows: number
  cols: number
}

let props = defineProps<Plot>()
</script>

<template>
  <div class="element element-group">
    <div v-if="props.title" class="label-c2 label-o4 label-w400 px-1">{{ props.title }}</div>
    <div class="plot">
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
</style>
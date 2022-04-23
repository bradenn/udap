<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {onMounted, reactive} from "vue";

interface Props {
  size: string
  rows?: number
  pad?: number
  cols?: number
  paged?: boolean
}

let state = reactive({
  classes: ""
})

let cssMap = {
  "xs": 1,
  "sm": 2,
  "md": 3,
  "lg": 4,
  "xl": 5,
}

let props = defineProps<Props>()

onMounted(() => {
  generateClasses()
})

function generateClasses() {
  state.classes = `widget widget-${props.size} p-${props.pad || 0}  `
}

</script>

<template>
  <div :class="state.classes">
    <slot></slot>
  </div>

</template>


<style lang="scss" scoped>

$widget-aspect: 1.1618;

$widgetSizes: (
    "xs": 1,
    "sm": 2,
    "md": 3,
    "lg": 4,
    "xl": 5,
);

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

.widget {
  display: flex;
  align-items: start;
  align-content: start;
  justify-content: start;
  gap: 0.25rem;
  width: 100%;
  height: 100%;
  grid-column: span v-bind('props.cols') !important;
  grid-row: span v-bind('props.rows') !important;
}


.widget-meta {
  position: absolute;

}

@each $sz, $i in $widgetSizes {
  .widget-#{$sz} {

  }
}


</style>
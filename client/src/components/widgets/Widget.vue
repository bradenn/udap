<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {onMounted, reactive} from "vue";

interface Props {
  size: string
  rows?: number
  pad?: number
  cols?: number
}

let state = reactive({
  classes: ""
})

let props = defineProps<Props>()

onMounted(() => {
  generateClasses()
})

function generateClasses() {
  state.classes = `widget widget-${props.size} p-${props.pad || 0}`
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
    "xs": 3rem,
    "sm": 12rem,
    "md": 16rem,
    "lg": 18rem,
    "xl": 24rem,
);

.widget {
  display: flex;
  align-items: center;
  aspect-ratio: $widget-aspect !important;

}


.widget-meta {
  position: absolute;

}

@each $sz, $i in $widgetSizes {
  .widget-#{$sz} {
    display: grid;
    grid-template-rows: repeat(v-bind('props.rows'), 1fr);
    grid-template-columns: repeat(v-bind('props.cols'), 1fr);
    aspect-ratio: $widget-aspect !important;
    width: $i;

  }
}


</style>
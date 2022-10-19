<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import Scroll from "@/components/Scroll.vue";

interface Tasks {
  index: number
  title: string
  description?: string
  value: string
  preview?: string
  active?: boolean
  alt?: string
  altFn?: () => void
  next?: (e: number) => void
}

let props = defineProps<Tasks>()

function callAltFn() {
  if (props.altFn && props.alt) {
    props.altFn()
  }
}


function callNext() {
  if (props.next) {
    props.next(props.index + 1)
  }
}

function directIndex() {
  if (props.next) {
    props.next(props.index)
  }
}


</script>

<template>

  <div v-if="props.active" class="element">

    <div class="label-md label-w700 label-o5 px-2 lh-1 my-1">{{ props.title }}</div>
    <Scroll class="task-container">
      <div class="">
        <slot></slot>
      </div>
    </Scroll>
    <div class="d-flex align-items-center align-content-center justify-content-between">
      <div class="label-c2 label-o3 label-w400 px-2">{{ props.description }}</div>
      <div class="subplot" style="height:1.5rem" @click="() => callNext()">
        <div class="label-c2 label-w500 label-o3 px-1 lh-1 text-accent">Next 􀯻</div>
      </div>
    </div>


  </div>
  <div v-else class="element p-2 flex-fill d-flex justify-content-between align-items-center"
       @click="() => directIndex()">
    <div class="d-flex align-items-center">
      <div v-if="!props.active && props.value !== ''"
           class="label-c2 label-w600 label-o5 px-1 lh-2 text-success">􀁢</div>
      <div v-else-if="props.value === ''" class="label-c2 label-w600 label-o2 px-1 lh-2">􀀀</div>
      <div v-else class="label-c2 label-w600 label-o5 px-1 lh-2 text-warning">􀇾</div>
      <div class="label-c1 label-w600 label-o5  lh-1 px-1">{{ props.title }}</div>
    </div>
    <div class="label-c2 label-w500 label-o3 px-2 lh-1">{{ props.preview ? props.preview : props.value }}</div>
  </div>
</template>


<style lang="scss" scoped>
.task-container {
  overflow-y: scroll;
  max-height: 10rem;
}

.label-shimmer {
  mix-blend-mode: lighten;
}

.nav-bar {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

//.nav-bar > div {
//  width: 100%;
//  outline: 1px solid white;
//}

.popup-nav {
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  width: 100%;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}

.nav-area {
  display: flex;
  align-items: center;
}

.pane-pop-up {
  z-index: 23 !important;
  position: absolute;
  width: 100%;
  height: 100%;
}


</style>
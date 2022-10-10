<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  previous?: boolean
  close?: () => void
  alt?: string
  fn?: () => void
  altDisabled?: boolean
}

let props = defineProps<Plot>()

function fn() {
  if (props.fn) {
    props.fn()
  }
}

function close() {
  if (props.close) {
    props.close()
  }
}

</script>

<template>
  <div class="element p-1 pt-2 h-100">
    <div v-if="props.title">
      <div v-if="props.previous" class="nav-bar">
        <div class="nav-area justify-content-start">
          <div class="d-flex gap-1 label-o3 px-2 align-items-center text-accent" @click="(e) => close()">
            <div class="label-c1">􀆉</div>
            <div class="label-c2 lh-1">Back</div>
          </div>
        </div>
        <div class="px-1 d-flex justify-content-center">
          <div class="label-c1  label-o4 label-w500">{{ props.title }}</div>
        </div>
        <div class="nav-area justify-content-end">

          <div v-if="props.altDisabled" class="d-flex gap-1 label-o3 px-2 align-items-center">
            <div class="label-c2 lh-1">{{ props.alt }}</div>
          </div>
          <div v-else class="d-flex gap-1 label-o3 px-2 align-items-center text-accent" @click="(e) => fn()">
            <div class="label-c2 lh-1">{{ props.alt }}</div>
          </div>
        </div>
      </div>
      <div v-else class="nav-area justify-content-start">
        <div class="nav-area justify-content-start">
          <div class="label-c1  label-o4 label-w500">{{ props.title }}</div>
        </div>
        <div class="px-1 d-flex justify-content-between w-100">

        </div>
      </div>

      <div class="nav-area justify-content-center">
      </div>
      <div class="nav-area justify-content-end">

      </div>
    </div>
    <div class="">
      <slot></slot>
    </div>
  </div>
</template>


<style lang="scss" scoped>
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
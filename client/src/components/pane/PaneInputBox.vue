<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  previous?: string
  alt?: string
  close?: () => void
  apply?: () => void
  begin?: boolean
}

let props = defineProps<Plot>()

function applyChanges() {
  props.apply?.()
}

function close() {
  if (props.close) {
    props.close()
  }
}
</script>

<template>
  <div class="element p-1 pt-1">
    <div v-if="props.title" class="popup-nav">
      <div class="nav-area justify-content-start">
        <div v-if="!props.begin" class="d-flex gap-1 label-o3 px-2 align-items-center label-shimmer"
             @click="(e) => close()">
          <div class="label-c1">􀆉</div>
          <div class="label-c2 lh-1">Back</div>
        </div>
      </div>
      <div class="nav-area justify-content-center">
        <div class="label-c1  label-o4 label-w500 px-1">{{ props.title }}</div>
      </div>
      <div class="nav-area justify-content-end">
        <div v-if="!props.begin" class="subplot subplot-button subplot-inline subplot-primary"
             @click="(e) => close()">Cancel</div>
        <div class="subplot subplot-button subplot-primary" @click="(e) => applyChanges()">Apply</div>
      </div>
    </div>
    <div class="pt-2">
      <slot></slot>
    </div>
  </div>
</template>


<style lang="scss" scoped>
.label-shimmer {
  mix-blend-mode: lighten;
}

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
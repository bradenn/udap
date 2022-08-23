<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title?: string
  previous?: boolean
  alt?: string
  close?: () => void

}

let props = defineProps<Plot>()

function close() {
  if (props.close) {
    props.close()
  }
}

</script>

<template>
  <div class="element p-1 pt-2 h-100">
    <div v-if="props.title" class="popup-nav">
      <div class="nav-area justify-content-start">
        <div v-if="props.previous" class="nav-area justify-content-start">
          <div class="d-flex gap-1 label-o3 px-2 align-items-center text-accent" @click="(e) => close()">
            <div class="label-c1">􀆉</div>
            <div class="label-c2 lh-1">Back</div>
          </div>
        </div>
        <div v-else class="label-c1  label-o4 label-w500 px-1">{{ props.title }}</div>
      </div>

      <div class="nav-area justify-content-center">
        <div v-if="props.previous" class="label-c1  label-o4 label-w500 px-1">{{ props.title }}</div>
      </div>
      <div class="nav-area justify-content-center">
      </div>
      <div class="nav-area justify-content-end">
      </div>
    </div>
    <div class="pt-1">
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
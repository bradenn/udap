<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Tasks {
  title: string
}

let props = defineProps<Tasks>()

</script>

<template>
  <div class="">
    <div class="label-lg label-w700 label-o5 px-1 lh-1 mb-2">{{ props.title }}</div>
    <div class="d-flex flex-column gap-1">
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
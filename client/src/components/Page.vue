<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";

interface Page {
  done: () => void,
  save: () => void,
  title: string
}

let props = defineProps<Page>()

const haptics = core.haptics()

function click() {
  haptics.tap(2, 1, 50)
}

function done() {
  if (props.done) {
    click()
    props.done()
  }
}

function save() {
  if (props.save) {
    props.save()
    done()
  }
}

</script>

<template>
  <div class="nav-grid gap-1 pb-1 w-100">
    <div class="d-flex justify-content-start">
      <div class="label-w500 label-c1 text-accent" @click="() => done()">􀆉 Back</div>
    </div>
    <div class="d-flex justify-content-center">
      <div class="label-w500 label-o5 label-c1 align-self-center">{{ props.title }}
      </div>
    </div>
    <div class="d-flex justify-content-end">
      <div class="label-w500 label-c1 text-accent" @click="() => save()">Save</div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.nav-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(3, 1fr);
}
</style>
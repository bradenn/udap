<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject} from "vue";
import type {Haptics} from "@/views/terminal/haptics";

interface Toggle {
  title?: string
  active: boolean
  fn: () => void
}

const haptics = inject("haptics") as Haptics

function click() {
  haptics.tap(2, 2, 100)
  // h
}

let props = defineProps<Toggle>()
</script>

<template>
  <div :class="`${props.active?'':''}`" class="subplot " @click="(e) => props.fn()" @mousedown="click">
    <div>{{ props.title }}</div>
    <div class="d-flex justify-content-center align-items-center align-content-center ">

      <div class="label-c2 label-o4 px-1 lh-1">
        {{ props.active ? 'ON' : 'OFF' }}
      </div>

    </div>
  </div>
</template>

<style scoped>

</style>
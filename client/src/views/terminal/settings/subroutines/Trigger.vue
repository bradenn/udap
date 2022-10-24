<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Trigger} from "@/types";
import {inject, reactive} from "vue";

interface TriggerProps {
  trigger: Trigger
}

const state = reactive({
  toggle: false
})

const props = defineProps<TriggerProps>()

const notifications = inject("notifications")

</script>

<template>
  <div class="subplot p-2 position-relative" @click="state.toggle = !state.toggle">
    <div class="d-flex align-items-start flex-row">
      <div class="d-flex align-items-center" style="width: 0.8rem;">
        <div class="label-c1 label-o4 label-w500 ">􀋥</div>
      </div>
      <div class="px-1">
        <div class="label-c2 label-o4 label-w700 lh-1 text-capitalize">{{ props.trigger.name }}</div>
        <div class="label-c3 label-o3 label-w400">{{ props.trigger.description }}</div>
        <div class="label-c3 label-o3 label-w400">Called by <b>{{ props.trigger.type }}</b></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.element-menu {
  position: relative;

  left: calc(-100% - 0.25rem);
  top: calc(100% + 0.25rem);
}

.grid-fill {
  position: relative;
  grid-column: 3 / 5;
  display: flex;
  flex-direction: column;
}

.grid-element {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}
</style>
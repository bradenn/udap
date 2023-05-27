<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {reactive} from "vue";
import core from "@/core";

const props = defineProps<{
    title: string
    icon: string
    to?: string
}>()

const state = reactive({
    pressed: false
})

function touchStart(e: TouchEvent) {
    if (!props.to) return
    e.preventDefault()
    state.pressed = true
}

const router = core.router()

function touchEnd(e: TouchEvent) {
    if (!props.to) return
    router.push(props.to)
    e.preventDefault()
    state.pressed = false
}

</script>

<template>
  <div class="d-flex flex-column table-element">
    <div :class="`${state.pressed?'pressed':''}`" class="table-link" @touchend="touchEnd"
         @touchstart="touchStart">
      <div class="d-flex gap-3 align-items-center">
        <div class="label-o1 label-w600 label-c3 sf-icon" style="padding-left: 0.5rem">{{ props.icon }}</div>
        <div class="label-o5 label-w500 label-c4">{{ props.title }}</div>
      </div>
      <slot></slot>
    </div>
    <div class="spacer"></div>
  </div>
</template>

<style lang="scss" scoped>
.table-element:not(:last-of-type) {
  .spacer {
    align-self: end;
    width: calc(100% - 3rem);
    height: 1px;
    background-color: rgba(255, 255, 255, 0.05);
  }
}

@keyframes table-pressed {
  0% {
    background-color: rgba(128, 128, 128, 0);
  }
  100% {
    background-color: rgba(128, 128, 128, 0.125);
  }
}

.table-link.pressed {
  animation: table-pressed 100ms ease-in forwards;
  filter: brightness(80%);
}

.table-link {
  //outline: 1px solid rgba(255, 0, 0, 0.1);
  text-decoration: none;
  width: 100%;
  padding: 1rem 0.5rem;
  display: flex;
  //gap: 0.75rem;
  align-items: center;
  align-content: center;
  justify-content: space-between;
  border-bottom: 1px solid transparent;
  border-radius: calc(11.5px - 0.25rem);
}
</style>
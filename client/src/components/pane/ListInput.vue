<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, reactive} from "vue";

interface Plot {
  name: string
  value: string
  description: string
  icon?: string
  type?: string
  change: (a: string) => void
}

let state = reactive({
  open: false,
})

let props = defineProps<Plot>()

let input = inject("input") as (props: InputProps, cb: (a: string) => void) => void

interface InputProps {
  name: string
  value: string
  description: string
  icon?: string
  type?: string
}

function update(value: string) {
  props.change(value)
}

function open() {
  let inp = {
    name: props.name,
    value: props.value,
    description: props.description,

  } as InputProps
  input(inp, update)
  state.open = true
}

function close() {
  state.open = false
}


</script>

<template>

  <div class="d-flex justify-content-between align-items-center subplot subplot-inline px-2 py-1" @click="() => open()">

    <div class="label-c1 label-w500 label-o4  title">{{ props.name }}</div>
    <div class="d-flex align-items-center">
      <div class="label-c3 label-w500 label-o3">{{ props.value }}</div>
      <div class="label-c1 label-o3 px-1">
        􀆊
      </div>
    </div>
  </div>

  <div class="h-sep my-1" style="width: 92%; margin-left: 0.5rem"></div>
</template>


<style lang="scss">
.context-rect {
  position: fixed !important;
  top: 0 !important;
  left: 0 !important;
}


.context-full {
  position: relative !important;
}

.title {
  font-family: "SF Pro", sans-serif !important;
  line-height: 1.2rem;
}
</style>
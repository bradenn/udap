<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {Macro} from "@/types";
import {inject, reactive} from "vue";
import Button from "@/components/Button.vue";
import macroService from "@/services/macroService";

interface MacroProps {
  macro: Macro
}

const state = reactive({
  toggle: false
})

const props = defineProps<MacroProps>()

function runMacro() {
  macroService.runMacro(props.macro.id).then(res => {
    state.toggle = false
  }).catch(err => {
    console.log(err)
  })
}

const notifications = inject("notifications")

function deleteMacro() {
  macroService.deleteMacro(props.macro.id).then(res => {
    notifications.show("Macro", `Macro '${props.macro.name}' deleted.`, 0, 1000 * 4)
    console.log(err)
  }).catch(err => {
    notifications.show("Macro", "Macro cannot be deleted.", 2, 1000 * 10)
    console.log(err)
  })
}

</script>

<template>
  <div class="element p-1 position-relative" @click="state.toggle = !state.toggle">
    <div class="d-flex align-items-start flex-row p-1">
      <div class="d-flex align-items-center">
        <div class="label-c2 label-o4 label-w500 lh-1">􁇵</div>
      </div>
      <div>
        <div class="label-c2 label-o4 label-w700 lh-1">{{ props.macro.name }}</div>
        <div class="label-c3 label-o3 label-w400">{{ props.macro.description }}</div>
      </div>
    </div>
  </div>
  <div v-if="state.toggle" class="element element-menu">
    <div class="grid-element">
      <Button active text="􀈑" @click="deleteMacro"></Button>
      <Button active text="􀍟"></Button>
      <Button active class="grid-fill" text="􀊃" @click="runMacro"></Button>
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
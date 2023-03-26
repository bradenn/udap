<!-- Copyright (c) 2023 Braden Nicholson -->
<script lang="ts" setup>

import type {Macro} from "@/types";
import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import EditMacro from "@/views/terminal/settings/subroutines/pages/EditMacro.vue";
import CreateMacro from "@/views/terminal/settings/subroutines/pages/CreateMacro.vue";


const remote = core.remote();

const pageSize = 12 * 2

const state = reactive({
  macros: [] as Macro[],
  current: [] as Macro[],
  page: 0,
  pages: 0,
  edit: {} as Macro,
  editing: false,
  create: false,
})

onMounted(() => {
  updateManifest()
  goToPage(0)
})

watchEffect(() => {
  updateManifest()
  return remote.macros
})

function updatePage() {
  let start = state.page * pageSize
  state.current = state.macros.slice(start, start + pageSize)
}

function hasLast() {
  return state.page > 0
}

function hasNext() {
  return state.page < state.pages - 1
}

function select(macro: Macro) {
  state.editing = true
  state.edit = macro
}

function unselect(): void {
  state.editing = false
  state.edit = {} as Macro
}

function goToPage(page: number) {
  if (state.page > state.pages - 1) {
    state.page = state.pages
    updatePage()
    return
  }
  state.page = page
  updatePage()
}

function nextPage() {
  if (state.page < state.pages - 1) {
    state.page++
  }
  updatePage()
}

function lastPage() {
  if (state.page > 0) {
    state.page--
  }
  updatePage()
}

function updateManifest() {
  state.macros = remote.macros
  state.pages = Math.ceil(state.macros.length / pageSize)
  updatePage()
}

function createStart() {
  state.create = true
}

function createStop() {
  state.create = false
}


</script>

<template>
  <Toolbar class="mb-2" icon="􀬂" title="Macros">
    <div class="w-100 d-flex justify-content-end align-items-center">
      <div class="label-c2 label-r label-w500 p-1 px-2 text-accent"
           @mousedown="() => createStart()">􀅼 Create</div>
      <div class="v-sep" style="height: 0.75rem"></div>
      <div class="label-c2 label-r label-o4 label-w500 p-1 px-2">Page {{ state.page + 1 }} of {{ state.pages }}</div>
      <div class="d-flex gap-1">
        <div :class="`${hasLast()?'text-accent':'label-o2'}`" class="label-c2 label-r  label-w500 p-1 px-2"
             @mousedown="lastPage">􀆉 Last</div>
        <div :class="`${hasNext()?'text-accent':'label-o2'}`" class="label-c2 label-r label-w500 p-1 px-2"
             @mousedown="nextPage">Next 􀆊</div>
      </div>
    </div>
  </Toolbar>
  <div class="macro-grid">
    <div v-for="macro in state.current" class="element" @click="() => select(macro)">
      <div class="label-c1 label-r label-o4 label-w500 pt-1 px-2 ">{{ macro.name }}</div>

      <div class="label-c2 label-o3 label-w400 px-2 lh-1">
        <div v-if="macro.description !==''">{{ macro.description }}</div>
        <div v-else>No description</div>
      </div>


    </div>
  </div>
  <EditMacro v-if="state.editing" :done="() => unselect()" :macro="state.edit"></EditMacro>
  <CreateMacro v-if="state.create" :done="() => createStop()"></CreateMacro>
  <!--  <Context v-if="state.editing">-->
  <!--  </Context>-->
</template>

<style scoped>
.macro-grid {
  display: grid;
  grid-template-columns: repeat(6, minmax(4rem, 1fr));
  grid-template-rows: repeat(4, minmax(4rem, 1fr));
  grid-row-gap: 0.25rem;
  grid-column-gap: 0.25rem;
}
</style>
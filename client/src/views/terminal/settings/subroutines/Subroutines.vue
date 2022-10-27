<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote, SubRoutine} from "@/types";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";
import Button from "@/components/Button.vue";
import MenuItem from "@/components/menu/MenuItem.vue";
import Menu from "@/components/menu/Menu.vue";
import MenuSection from "@/components/menu/MenuSection.vue";
import ToolbarButton from "@/components/ToolbarButton.vue";

let remote = inject('remote') as Remote

let state = reactive({
  subroutines: [] as SubRoutine[],
  selected: [] as string[],
  select: false,
  loading: true,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => {
  handleUpdates(remote)
  return state.subroutines
})

function selectSR(sr: SubRoutine) {
  if (!state.select) return
  if (state.selected.includes(sr.id)) {
    state.selected = state.selected.filter(s => s !== sr.id)
  } else {
    state.selected.push(sr.id)
  }
}

function unselectSR(sr: SubRoutine) {
  state.selected = state.selected.filter(s => s !== sr.id)
}

function handleUpdates(remote: Remote) {
  state.loading = false
  state.subroutines = remote.subroutines
  return remote
}

function selectStop() {
  state.select = false
  state.selected = []

}

function selectStart() {
  state.select = true

}

function select(id: string) {

}

</script>

<template>
  <div class="layout-grid">
    <div class="layout-sidebar">
      <div class="d-flex mb-1 gap-1 p-0" style="height: 1.5rem;">
        <Button :active="false" class="element flex-grow-1" style="height: 1.5rem"
                text="􀅼 Macro" to="/terminal/settings/subroutines/macro"></Button>
        <Button :active="false" class="element flex-grow-1" style="height: 1.5rem"
                text="􀅼 Subroutine" to="/terminal/settings/subroutines/create"></Button>
      </div>
      <Menu alt="" style="height: calc(100% - 1.625rem) !important;" title="">
        <MenuSection title="Core Subroutines">
          <MenuItem active icon="􀏧" subtext="12" title="All Items"></MenuItem>
          <MenuItem icon="􀠀" subtext="4" title="Homekit"></MenuItem>
          <MenuItem icon="􀫥" subtext="8" title="System"></MenuItem>
          <MenuItem icon="􀥭" subtext="8" title="Modules"></MenuItem>
        </MenuSection>
        <MenuSection title="User Groups">
          <MenuItem icon="􀏧" subtext="12" title="All Items"></MenuItem>
          <MenuItem icon="􀠀" subtext="4" title="Homekit"></MenuItem>
          <MenuItem icon="􀫥" subtext="8" title="System"></MenuItem>
          <MenuItem icon="􀥭" subtext="8" title="Modules"></MenuItem>
        </MenuSection>

      </Menu>
    </div>
    <div class="layout-body">
      <div class="d-flex mb-1 justify-content-between align-items-center flex-row">
        <div class="px-2 d-flex gap-1 align-items-center element w-100 lh-2" style="height: 1.8rem;">
          <div class="label-xs label-w200 label-o3 px-0">􀐘</div>
          <div class="label-xs label-w700 label-o5">Subroutines</div>

          <!--          <div class="button-sep"></div>-->
          <ToolbarButton :active="false" :disabled="state.select" class=" px-3"
                         style="height: 1.4rem"
                         text="Select"
                         @click="() => state.select?selectStop():selectStart()"></ToolbarButton>
          <div class="button-sep"></div>
          <ToolbarButton :active="false" :disabled="state.selected.length === 0" class=" px-3" style="height: 1.5rem"
                         text="Trigger" to="/terminal/settings/subroutines/create"></ToolbarButton>
          <ToolbarButton :active="false" :disabled="state.selected.length === 0" class="  px-3" style="height: 1.5rem"
                         text="Delete" to="/terminal/settings/subroutines/create"></ToolbarButton>
          <ToolbarButton v-if="state.select" :accent="true" :active="false" :text="state.select?'Done':'Select'"
                         class=" px-3"
                         style="height: 1.4rem"
                         @click="() => selectStop()"></ToolbarButton>
        </div>


      </div>
      <MenuSection style="flex-direction: row;" title="All Subroutines">
        <div class="page-grid">
          <Subroutine v-for="sr in state.subroutines" :key="sr.id" :selected="state.selected.includes(sr.id)"
                      :subroutine="sr"
                      @click="state.select?selectSR(sr):$router.push(`/terminal/settings/subroutines/${sr.id}/edit`)"
                      v-on:click.stop></Subroutine>
        </div>
      </MenuSection>

    </div>
  </div>
</template>

<style lang="scss" scoped>

.button-sep {
  width: 3px;
  height: 32px;
  background-color: rgba(255, 255, 255, 0.3);
}

.layout-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(5, 1fr);
  grid-template-rows: repeat(1, 1fr);
}

.layout-sidebar {
  grid-column: 1 / span 1;
}

.layout-body {
  grid-column: 2 / span 4;
}

.page-grid > div {
  //outline: 1px solid white;
}

.page-grid {
  width: 100%;

  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(6, 1fr);
  grid-template-rows: repeat(2, 1fr);
}
</style>
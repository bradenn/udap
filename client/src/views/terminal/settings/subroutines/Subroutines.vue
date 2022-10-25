<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote, SubRoutine} from "@/types";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";
import Button from "@/components/Button.vue";
import MenuItem from "@/components/menu/MenuItem.vue";
import Menu from "@/components/menu/Menu.vue";
import MenuSection from "@/components/menu/MenuSection.vue";

let remote = inject('remote') as Remote

let state = reactive({
  subroutines: [] as SubRoutine[],
  selected: {} as SubRoutine,
  loading: true,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => {
  handleUpdates(remote)
  return remote.subroutines
})

function selectSR(sr: SubRoutine) {
  state.selected = sr
}

function handleUpdates(remote: Remote) {
  state.loading = false
  state.subroutines = remote.subroutines
  return remote
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
      <div class="d-flex mb-1 justify-content-between align-items-center flex-row" style="height: 1.5rem;">
        <div class="d-flex align-items-center py-1 pb-1">
          <div class="label-sm label-w200 label-o6 px-1">􀏧</div>
          <div class="label-sm label-w700 label-o6">All Subroutines</div>
          <div class="px-2 d-flex gap-1">
            <Button :active="true" class="element flex-grow-1 px-3 " style="height: 1.5rem"
                    text="Edit" to="/terminal/settings/subroutines/create"></Button>
            <Button :active="false" class="element flex-grow-1 px-3" style="height: 1.5rem"
                    text="Delete" to="/terminal/settings/subroutines/create"></Button>
          </div>
        </div>
        <div class="d-flex gap-1">
          <Button :active="false" class="element flex-grow-1" style="height: 1.5rem"
                  text="􀈙 Create Group" to="/terminal/settings/subroutines/create"></Button>
          <Button :active="false" class="element flex-grow-1" style="height: 1.5rem"
                  text="􀈊 Edit Group" to="/terminal/settings/subroutines/create"></Button>
        </div>


      </div>
      <div class="page-grid">
        <Subroutine v-for="sr in state.subroutines" :key="sr.id" :subroutine="sr"></Subroutine>
      </div>

    </div>
  </div>
</template>

<style lang="scss" scoped>


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
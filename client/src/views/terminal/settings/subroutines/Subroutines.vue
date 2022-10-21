<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote, SubRoutine} from "@/types";
import Subplot from "@/components/plot/Subplot.vue";
import Macro from "@/views/terminal/settings/subroutines/Macro.vue";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";

let remote = inject('remote') as Remote

let state = reactive({
  subroutines: [] as SubRoutine[],
  loading: true,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => {
  return handleUpdates(remote)
})

function handleUpdates(remote: Remote) {
  state.loading = false
  state.subroutines = remote.subroutines
  return remote
}

function createSubroutine() {
  // subroutineService.createSubroutine({
  //   triggerId: "25b83eaa-cf21-40b7-8a7c-3fcf58564904",
  //   macros: [
  //     {
  //       name: "Bedroom Lights Off",
  //       description: "Turn off Bedroom lights",
  //       zone: "3a4e948e-e7d7-447b-90a5-c6306ec98efb",
  //       type: "on",
  //       value: "false",
  //     } as Macro
  //   ],
  //   description: "Lights Off",
  // } as SubRoutine)
}

</script>

<template>
  <div>
    <div class="d-flex">
      <div class="element d-flex mb-1" style="height: 2rem">
        <Subplot name="Create Subroutine" to="/terminal/settings/subroutines/create"></Subplot>
        <Subplot name="Create Macro" to="/terminal/settings/subroutines/macro"></Subplot>

      </div>
    </div>
    <div class="page-grid">
      <Subroutine v-for="sr in state.subroutines" :key="sr.id" :subroutine="sr"></Subroutine>
    </div>
    <div class="page-grid">
      <Macro v-for="sr in remote.macros" :key="sr.id" :macro="sr"></Macro>

    </div>
  </div>
</template>

<style lang="scss" scoped>
.page-grid {
  width: 100%;

  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(6, 1fr);
  grid-template-rows: repeat(2, 1fr);
}
</style>
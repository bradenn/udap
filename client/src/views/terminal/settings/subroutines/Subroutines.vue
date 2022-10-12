<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Remote, SubRoutine} from "@/types";
import Subplot from "@/components/plot/Subplot.vue";

let remote = inject('remote') as Remote

let state = reactive({
  subroutines: [] as SubRoutine[],
  loading: true,
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

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
        <Subplot name="Create" to="/terminal/settings/subroutines/create"></Subplot>

      </div>
    </div>
    <div class="page-grid">
      <div v-for="sr in state.subroutines" class="element p-2">
        <div class="d-flex justify-content-between">
          <div class="label-xs label-o4 label-w500 pb-2">􁏀</div>
        </div>
        <div class="label-c2 label-o4 label-w700 lh-1">{{ sr.description }}</div>
        <div class="label-c3 label-o3 label-w400">{{ sr?.macros?.length || 0 }}
          macro{{ sr?.macros?.length !== 1 ? 's' : '' }}
        </div>
      </div>

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
  grid-template-rows: repeat(4, 1fr);
}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Macro, SubRoutine, Trigger} from "@/types";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";
import MacroDom from "@/views/terminal/settings/subroutines/Macro.vue";

import TriggerDom from "@/views/terminal/settings/subroutines/Trigger.vue";
import Button from "@/components/Button.vue";
import subroutineService from "@/services/subroutineService";
import core from "@/core";

const router = core.router()
const remote = core.remote()

const state = reactive({
  subroutine: {
    macros: [] as Macro[]
  } as SubRoutine,
  trigger: {} as Trigger,
  loaded: false,
})
onMounted(() => {
  load()
})

function load() {
  const id = router.currentRoute.value.params["subroutine"];
  if (!id) return;
  const sr = remote.subroutines.find(s => s.id === id)
  if (!sr) return
  state.subroutine = sr
  const trigger = remote.triggers.find(s => s.id === state.subroutine.triggerId)
  if (!trigger) return
  state.trigger = trigger
  if (!state.subroutine.macros) return
  state.loaded = true;
}


watchEffect(() => {
  load()
  return remote
})

function goBack() {
  router.push("/terminal/settings/subroutines")
}

const notifications = inject("notifications") as any

function triggerSubroutine() {
  subroutineService.triggerManual(state.subroutine.id).then(res => {
    notifications.show("Subroutine", `Subroutine '${state.subroutine.description}' triggered.`, 0, 1000 * 3)
  }).catch(err => {
    notifications.show("Subroutine", `Subroutine '${state.subroutine.description}' cannot be triggered.`, 2, 1000 * 3)
  })
}

function deleteSubRoutine() {
  subroutineService.deleteSubroutine(state.subroutine.id).then(res => {
    notifications.show("Subroutine", `Subroutine '${state.subroutine.description}' deleted.`, 0, 1000 * 3)
    goBack();
  }).catch(err => {
    notifications.show("Subroutine", `Subroutine '${state.subroutine.description}' cannot be deleted.`, 2, 1000 * 3)
  })
}

</script>

<template>

  <div class="grid-element">
    <div class="d-flex align-items-start label-o4 gap-1 pb-2">
      <div class="label-w500 label-c1 text-accent" @click="goBack">􀆉 Back</div>
    </div>
    <div class="d-flex flex-column gap-1">

      <div class="d-flex align-items-center pb-1">
        <div class="label-sm label-w200 label-o6 px-1">􀏧</div>
        <div class="label-sm label-w700 label-o6">Subroutine</div>
      </div>
      <Subroutine :subroutine="state.subroutine"></Subroutine>
      <div class="d-flex gap-1">
        <Button :active="true" class="element flex-grow-1" style="height: 1.5rem"
                text="􀋥 Trigger" @click="triggerSubroutine"></Button>
        <Button :active="true" class="element flex-grow-1" style="height: 1.5rem"
                text="􀈑 Delete" @click="deleteSubRoutine"></Button>
      </div>
    </div>
    <div class="d-flex flex-column gap-1">
      <div class="d-flex align-items-center pb-1">
        <div class="label-sm label-w200 label-o6 px-1">􀐠</div>
        <div class="label-sm label-w700 label-o6">Workflow</div>
      </div>
      <div v-if="state.loaded" class="element px-1 pt-0">
        <div class="d-flex align-items-center py-1 pb-1 px-1">
          <div class="label-c2 label-w500 label-o3">When triggered by:
          </div>
        </div>

        <trigger-dom :trigger="state.trigger"></trigger-dom>
        <div class="d-flex align-items-center py-1 pb-1  px-1">
          <div class="label-c2 label-w500 label-o3">Run {{ state.subroutine.macros.length }}
            macro{{ (state.subroutine.macros.length !== 1) ? 's' : '' }}:
          </div>
        </div>
        <div class="d-flex gap-1 flex-column">
          <MacroDom v-for="macro in state.subroutine.macros" :macro="macro" subplot></MacroDom>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>

.preview {
  width: 9rem;
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

.generic-grid {
  display: flex;
  justify-content: center;
  grid-column-gap: 0.25rem;
}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import PaneMenu from "@/components/pane/PaneMenu.vue";
import {inject, reactive} from "vue";
import PaneMenuToggle from "@/components/pane/PaneMenuToggle.vue";
import type {Remote} from "@/types";

const router = useRouter()

const remote = inject("remote") as Remote

const state = reactive({
  trigger: "" as string
})

function goBack() {
  router.push("/terminal/settings/subroutines")
}

function selectTrigger(id: string) {
  if (state.trigger === id) {
    state.trigger = ""
  } else {
    state.trigger = id
  }
}

</script>

<template>
  <div class="entity-grid">
    <PaneMenu alt="Select a trigger" title="Triggers">
      <PaneMenuToggle v-for="trigger in remote.triggers" :active="trigger.id === state.trigger"
                      :fn="() => selectTrigger(trigger.id)"
                      :subtext="trigger.description" :title="trigger.name"></PaneMenuToggle>
    </PaneMenu>
  </div>
</template>

<style scoped>

.pane-fill {
  position: relative;
  grid-column: 2 / 4;
  grid-row: 1 / 1;
  display: flex;
  flex-direction: column;
}

.entity-grid {
  width: 100%;
  height: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-rows: repeat(1, 1fr);
  grid-template-columns: repeat(4, 1fr);
}
</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import PaneMenuLink from "@/components/pane/PaneMenuLink.vue";
import PaneList from "@/components/pane/PaneList.vue";

const router = useRouter()

const moduleId = router.currentRoute.value.params['moduleId']


function goBack() {
  router.push("/terminal/settings/modules")
}

</script>

<template>
  <div class="entity-grid">
    <PaneList :close="goBack" :previous="true" title="Module">
      <PaneMenuLink :to="`/terminal/settings/modules/${moduleId}/actions`" subtext="Make changes to the module runtime"
                    title="Actions"></PaneMenuLink>
      <PaneMenuLink :to="`/terminal/settings/modules/${moduleId}/metadata`" subtext="Module-defined information"
                    title="Metadata"></PaneMenuLink>
      <PaneMenuLink :to="`/terminal/settings/modules/${moduleId}/versions`" subtext="Previously compiled versions"
                    title="Versions"></PaneMenuLink>

    </PaneList>
    <div class="pane-fill h-100">
      <router-view></router-view>
    </div>
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
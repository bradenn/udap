<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import PaneListItem from "@/components/pane/PaneListItem.vue";
import {inject, onMounted, reactive} from "vue";
import moment from "moment/moment";
import type {Module, Remote} from "@/types";
import {useRouter} from "vue-router";
import PaneList from "@/components/pane/PaneList.vue";

const remote = inject("remote") as Remote
const router = useRouter()
const state = reactive({
  module: {} as Module
})

onMounted(() => {
  loadModule()
})

function loadModule() {
  let moduleId = router.currentRoute.value.params['moduleId']
  const module = remote.modules.find(m => m.id === moduleId)
  if (!module) return
  state.module = module
}

function formatDate(date: string): string {
  let dateObj = moment(date);
  return dateObj.fromNow()
}
</script>

<template>
  <PaneList title="Versions">
    <PaneListItem :active="false" :subtext="state.module.uuid" title="Current Build"></PaneListItem>

  </PaneList>

</template>

<style scoped>

</style>
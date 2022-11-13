<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import PaneListItem from "@/components/pane/PaneListItem.vue";
import {inject, onMounted, reactive} from "vue";
import moment from "moment/moment";
import type {Module} from "@/types";
import {useRouter} from "vue-router";
import PaneList from "@/components/pane/PaneList.vue";
import type {Remote} from "@/remote";

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

    <PaneList title="Metadata">
        <PaneListItem :active="false" :subtext="state.module.author" title="Author"></PaneListItem>
        <PaneListItem :active="false" :subtext="state.module.description" title="Description"></PaneListItem>
        <PaneListItem :active="false" :subtext="state.module.path" title="Source Path"></PaneListItem>
        <PaneListItem :active="false" :subtext="state.module.version" title="Version"></PaneListItem>
        <PaneListItem :active="false" :subtext="formatDate(state.module.created)" title="Installed"></PaneListItem>
        <PaneListItem :active="false" :subtext="formatDate(state.module.updated)"
                      title="Last Compilation"></PaneListItem>
        <PaneListItem :active="false" :subtext="state.module.config || '--'" title="Config"></PaneListItem>
    </PaneList>
</template>

<style scoped>

</style>
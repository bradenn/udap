<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {SubRoutine, Trigger} from "@/types";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";
import MenuSection from "@/components/menu/MenuSection.vue";
import ToolbarButton from "@/components/ToolbarButton.vue";
import Create from "@/views/terminal/settings/subroutines/pages/Create.vue";
import moment from "moment";
import Toolbar from "@/components/toolbar/Toolbar.vue";
import type {Remote} from "@/remote";

let remote = inject('remote') as Remote

let state = reactive({
    subroutines: [] as SubRoutine[],
    selected: [] as string[],
    triggers: [] as Trigger[],
    select: false,
    loading: true,
    createSubroutine: false,
})

onMounted(() => {
    state.loading = true
    handleUpdates(remote)
})

watchEffect(() => {
    handleUpdates(remote)
    updateTriggers()
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

function sortLastTrigger(a: Trigger, b: Trigger) {
    if (new Date(a.lastTrigger).valueOf() >= new Date(b.lastTrigger).valueOf()) {
        return -1
    } else {
        return 1
    }
}

function selectStart() {
    state.select = true

}

function filterOneWeek(a: Trigger) {
    return new Date().valueOf() - new Date(a.lastTrigger).valueOf() <= 7 * 24 * 60 * 60 * 1000;
}

function groupBySource() {

}

function groupBy<T>(xs: T[], key: string): T[] {
    return xs.reduce(function (rv: any, x: any): T {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
    }, {});
}

function updateTriggers() {
    state.triggers = remote.triggers.filter(filterOneWeek).sort(sortLastTrigger)
}

function select(id: string) {

}

function formatTimeSince(id: string) {
    let h = moment(id).fromNow()
    return `${h}`
}

function createSubroutine() {
    state.createSubroutine = true
}

</script>

<template>
    <div class="layout-grid">

        <div class="layout-body">
            <Toolbar class="mb-1" icon="􀐘" title="Subroutines">
                <ToolbarButton :active="false" :disabled="state.select" class=" px-3"
                               style="height: 1.4rem"
                               text="Select"
                               @click="() => state.select?selectStop():selectStart()"></ToolbarButton>
                <div class="button-sep"></div>
                <ToolbarButton :active="false" :disabled="state.selected.length === 0"
                               class=" px-3"
                               style="height: 1.5rem"
                               text="Trigger"
                               to="/terminal/settings/subroutines/create"></ToolbarButton>
                <ToolbarButton :active="false" :disabled="state.selected.length === 0"
                               class="  px-3"
                               style="height: 1.5rem"
                               text="Delete"
                               to="/terminal/settings/subroutines/create"></ToolbarButton>
                <div class="flex-grow-1"></div>
                <ToolbarButton v-if="state.select" :accent="true" :active="false"
                               :text="state.select?'Done':'Select'"
                               class=" px-3"
                               style="height: 1.4rem"
                               @click="() => selectStop()"></ToolbarButton>
                <div class="button-sep"></div>
                <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                               text="Scene"
                               to="/terminal/settings/subroutines/scenes/create"></ToolbarButton>
                <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                               text="Macro"
                               to="/terminal/settings/subroutines/macros/create"></ToolbarButton>

                <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                               text="Subroutine"
                               @click="() => createSubroutine()"></ToolbarButton>
                <ToolbarButton :active="false" icon="􀅼" style="height: 1.5rem"
                               text="Zone"
                               to="/terminal/settings/subroutines/zones/create"></ToolbarButton>
            </Toolbar>

            <MenuSection style="flex-direction: row; flex-col: span 3;"
                         title="All Subroutines">
                <div class="page-grid">
                    <Subroutine v-for="sr in state.subroutines" :key="sr.id"
                                :selected="state.selected.includes(sr.id)"
                                :subroutine="sr"
                                @click="state.select?selectSR(sr):$router.push(`/terminal/settings/subroutines/${sr.id}/edit`)"
                                v-on:click.stop></Subroutine>
                </div>
            </MenuSection>
            <Create v-if="state.createSubroutine"
                    :done="() => state.createSubroutine = false"></Create>

        </div>
        <div class="layout-sidebar">

            <MenuSection class="gap-2 " title="Triggers">
                <div v-for="trigger in state.triggers"
                     class="d-flex gap-0 flex-column mb-3">
                    <div class="element">
                        <div class="d-flex justify-content-between">
                            <div class="d-flex gap-1  align-items-center p-1">
                                <div class="label-c3 label-w600 label-o3 lh-1"
                                     style="font-size: 18px;">
                                    􀋦
                                </div>
                                <div class="label-c2 label-w500 label-o4 lh-1"
                                     style="font-size: 21px;">
                                    {{ trigger.name }}
                                </div>
                            </div>
                            <div class="label-w400 label-c2 label-o3  lh-1 p-1">{{
                                formatTimeSince(trigger.lastTrigger)
                                }}
                            </div>
                        </div>
                        <div class="d-flex justify-content-between flex-column  gap-1">
                            <div
                                    v-for="sr in state.subroutines.filter(s => s.triggerId === trigger.id)"
                                    class=" py-2 d-flex gap-1 align-items-center justify-content-between flex-row px-3">
                                <div class="d-flex gap-1 align-items-center">
                                    <div
                                            class="label-w400 label-c2 label-o2 opacity-50 lh-1 text-accent">􀄵
                                    </div>
                                    <div class="label-w400 label-c1 label-o3 lh-1">{{
                                        sr.icon
                                        }}
                                    </div>
                                    <div class="label-w500 label-c2 label-o4 lh-1">
                                        {{ sr.description }}
                                    </div>
                                </div>
                                <div
                                        class="label-w200 label-c2 label-o4  lh-1 text-accent">􀄔
                                </div>
                            </div>

                        </div>
                    </div>

                </div>

            </MenuSection>

        </div>


    </div>
</template>

<style lang="scss" scoped>

.button-sep {
  width: 3px;
  height: 32px;

}

.macro-grid > .element {
  width: 100% !important;
}

.macro-grid {

  display: grid;

  grid-gap: 0.5rem;
  grid-template-columns: repeat(6, 1fr);
  grid-template-rows: repeat(2, 1fr);
}

.layout-grid {
  width: 100%;
  display: grid;
  grid-column-gap: 0.25rem;
  grid-row-gap: 0.25rem;
  grid-template-columns: repeat(12, 1fr);
  grid-template-rows: repeat(1, 1fr);
}

.layout-sidebar {
  grid-column: 10 / span 3;
}

.layout-body {
  grid-column: 1 / span 9;
}

.page-grid > div {
  /*outline: 1px solid white;*/
}


</style>
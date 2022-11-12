<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, reactive, watchEffect} from "vue";
import type {Macro, SubRoutine, Trigger} from "@/types";
import Subroutine from "@/views/terminal/settings/subroutines/Subroutine.vue";
import MacroDom from "@/views/terminal/settings/subroutines/Macro.vue";

import TriggerDom from "@/views/terminal/settings/subroutines/Trigger.vue";
import Button from "@/components/Button.vue";
import subroutineService from "@/services/subroutineService";
import core from "@/core";
import EditMacro from "@/views/terminal/settings/subroutines/pages/EditMacro.vue";
import moment from "moment";
import ShowMacros from "@/views/terminal/settings/subroutines/pages/ShowMacros.vue";
import type {Notify} from "@/notifications";

const router = core.router()
const remote = core.remote()

const state = reactive({
    subroutine: {
        macros: [] as Macro[]
    } as SubRoutine,
    trigger: {} as Trigger,
    editMacro: false,
    showMacros: false,
    macro: {
        name: "",
        id: "",
    } as Macro,
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

watchEffect(() => {
    state.subroutine.macros = remote.macros.filter(m => state.subroutine.macros.map(i => i.id).includes(m.id))
    return remote.macros
})

function goBack() {
    router.push("/terminal/settings/subroutines")
}

const notify: Notify = core.notify()

function triggerSubroutine() {
    subroutineService.triggerManual(state.subroutine.id).then(res => {
        notify.success("Subroutine", `Subroutine '${state.subroutine.description}' triggered.`)
    }).catch(err => {
        notify.fail("Subroutine", `Subroutine '${state.subroutine.description}' could not be triggered.`)
    })
}

function deleteSubRoutine() {
    subroutineService.deleteSubroutine(state.subroutine.id).then(res => {
        notify.success("Subroutine", `Subroutine '${state.subroutine.description}' deleted.`)
        goBack();
    }).catch(err => {
        notify.fail("Subroutine", `Subroutine '${state.subroutine.description}' could not be deleted.`)
    })
}

function doneEditing() {
    state.editMacro = false
    state.macro = {} as Macro
}

function timeSince(time: string): string {
    return moment(time).fromNow()
}

function addMacro(id: string) {
    subroutineService.addMacro(state.subroutine.id, id)
}

function showMacros() {
    state.showMacros = true
}

function removeMacro(id: string) {
    subroutineService.removeMacro(state.subroutine.id, id)
}

function calcDuration(min: number) {
    return `${Math.floor(min / 60)}h ${min % 60}m `
}

function parseDate(dt: string) {


    return moment().from(dt).toString()
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
            <div>
                <Button :active="true" class="element flex-grow-1" icon="􀍟"
                        style="height: 1.8rem" text="Manage" @click="showMacros"></Button>
            </div>
            <div class="d-flex gap-1">
                <Button :active="true" class="element flex-grow-1" icon="􀋥"
                        style="height: 1.8rem" text="Trigger" @click="triggerSubroutine"></Button>
                <Button :active="true" class="element flex-grow-1" icon="􀈑"
                        style="height: 1.8rem" text="Delete" @click="deleteSubRoutine"></Button>
            </div>
            <div class="element d-flex flex-column gap-1 p-2">
                <div class="d-flex justify-content-between">
                    <div class="label-c2 label-w500 label-o4">Created</div>
                    <div class="label-c2 label-w500 label-o4">{{ timeSince(state.subroutine.created) }}</div>
                </div>
                <div class="d-flex justify-content-between">
                    <div class="label-c2 label-w500 label-o4">Last Edited</div>
                    <div class="label-c2 label-w500 label-o4">{{ timeSince(state.subroutine.updated) }}</div>
                </div>
            </div>
        </div>


        <div class="d-flex flex-column gap-1">
            <div class="d-flex align-items-center pb-1">
                <div class="label-sm label-w200 label-o6 px-1">􀐠</div>
                <div class="label-sm label-w700 label-o6">Workflow</div>
            </div>
            <div v-if="state.loaded" class="element px-1 pt-0">
                <div class="d-flex align-items-center justify-content-between py-1 pb-1 px-1">
                    <div class="label-c2 label-w500 label-o3">When triggered by:
                    </div>

                </div>

                <trigger-dom :trigger="state.trigger"></trigger-dom>
                <div class="d-flex align-items-center py-1 pb-1  px-1 justify-content-between">
                    <div class="label-c2 label-w500 label-o3">Run {{ state.subroutine.macros.length }}
                        macro{{ (state.subroutine.macros.length !== 1) ? 's' : '' }}:
                    </div>

                </div>

                <div class="d-flex gap-1 flex-column w-100">
                    <div v-for="macro in state.subroutine.macros" class="d-flex gap-1 align-items-center">
                        <div class="w-100" @click="() => {state.macro = macro; state.editMacro = true; }">
                            <MacroDom :macro="macro" class="flex-grow-1 w-100" subplot
                            ></MacroDom>
                        </div>
                        <div class="subplot h-100 px-2" style=" text-shadow:none;height: 2.35rem !important;"
                             @click="() => removeMacro(macro.id)">􀭲</div>
                    </div>


                </div>
                <div class="d-flex flex-column gap-1 pt-1">
                    <div class="label-c2 label-w500 label-o3 px-1">Revert after:
                    </div>
                    <div class="subplot w-100 p-2 px-2">

                        {{ calcDuration(state.subroutine.revertAfter) }}
                    </div>
                </div>

                <div class=" position-relative">
                </div>
            </div>
            <div class="d-flex gap-1">


            </div>

        </div>

        <EditMacro v-if="state.editMacro" :done="doneEditing"
                   :macro="state.macro"
                   @click="() => {state.editMacro = false;}"></EditMacro>
        <ShowMacros v-if="state.showMacros" :done="() => {state.showMacros = false;}"
                    :subroutine="state.subroutine"></ShowMacros>

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
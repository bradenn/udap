<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {inject, onMounted, reactive} from "vue";
import moment from "moment/moment";
import type {Module, ModuleVariable} from "@/types";
import {useRouter} from "vue-router";
import PaneList from "@/components/pane/PaneList.vue";
import ListInput from "@/components/pane/ListInput.vue";
import type {Remote} from "@/remote";

const remote = inject("remote") as Remote
const router = useRouter()
const state = reactive({
    module: {} as Module,
    config: {} as any,
    meta: [] as ModuleVariable[],
    edit: "",
    selected: "" as string,
    selectedMeta: {} as ModuleVariable
})

onMounted(() => {
    loadModule()
})

function loadModule() {
    let moduleId = router.currentRoute.value.params['moduleId']
    const module = remote.modules.find(m => m.id === moduleId)
    if (!module) return
    state.module = module
    if (!state.module.config) return
    if (state.module.config == '') return
    state.config = JSON.parse(state.module.config)
    state.meta = JSON.parse(state.module.variables) as ModuleVariable[]

}

function formatDate(date: string): string {
    let dateObj = moment(date);
    return dateObj.fromNow()
}

function cancelEdit() {
    state.selected = ""
    state.edit = ""
}

function modifyVariable(key: string, value: string) {
    state.config[key] = value
}

function getDescription(key: string): string {
    let cand = state.meta.find(m => m.name === key)
    if (!cand) return ""
    return cand.description
}

</script>

<template>

    <PaneList style="width: 15rem;" title="Configuration">
        <ListInput v-for="kv in Object.keys(state.config)" :change="(s: string) => modifyVariable(kv, s)"
                   :description="getDescription(kv)"
                   :name="kv"
                   :value="state.config[kv]" type="text"></ListInput>
    </PaneList>


</template>

<style lang="scss" scoped>


.keyboard-location {


}

.font-monospace {
  font-family: "Menlo", sans-serif;
}
</style>
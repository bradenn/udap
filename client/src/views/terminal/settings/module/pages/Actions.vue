<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import {inject, onMounted, reactive} from "vue";
import moment from "moment/moment";
import type {Module} from "@/types";
import {useRouter} from "vue-router";
import Confirm from "@/components/plot/Confirm.vue";
import Radio from "@/components/plot/Radio.vue";
import axios from "axios";
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

function reloadModule(id: string) {
    axios.post(`http://localhost:3020/modules/${id}/reload`)
}

function toggleEnabled(id: string, enabled: boolean) {
    if (enabled) {
        axios.post(`http://localhost:3020/modules/${id}/enable`)
    } else {
        axios.post(`http://localhost:3020/modules/${id}/disable`)
    }
}


</script>

<template>
    <div class="element p-1 action mb-1">
        <div class="p-1 pt-0">
            <div>
                <div class="label-c1 label-w500 label-o4">Module State</div>
            </div>
            <div class="label-c2 label-o2 label-w400 lh-1">Enable or disable the module depending on if it is
                needed.<br>
                <span class="text-warning">This will terminate the <i>{{ state.module.name }}</i> module runtime.</span>
            </div>

        </div>
        <div class="d-flex" style="max-height: 1.6rem; min-width: 6rem">
            <Confirm v-if="!state.module.enabled" :fn="() => toggleEnabled(state.module.id, !state.module.enabled)"
                     icon="􀊃"
                     title="Enable"></Confirm>
            <Confirm v-if="state.module.enabled" :fn="() => toggleEnabled(state.module.id, !state.module.enabled)"
                     icon="􀆧"
                     title="Disable"></Confirm>
        </div>
    </div>
    <div class="element p-1 action mb-1">
        <div class="p-1 pt-0">
            <div>
                <div class="label-c1 label-w500 label-o4">Reload Module</div>
            </div>
            <div class="label-c2 label-o2 label-w400 lh-1">Recompile the module from sources and replace runtime.<br>
                <span class="text-warning">This will only work if the source for <i>{{
                        state.module.name
                    }}</i> has changed.</span>
            </div>

        </div>
        <div class="d-flex" style="max-height: 1.6rem; min-width: 6rem">
            <Radio :active="false" :fn="() => reloadModule(state.module.name)" sf="􀅉" style="width: 100%;"
                   title=""></Radio>
        </div>
    </div>

</template>

<style scoped>
.action {
    display: flex;
    min-height: 2rem;
    justify-content: space-between;
}
</style>
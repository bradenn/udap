<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Timing} from "@/types";
import Loader from "@/components/Loader.vue";
import type {Remote} from "@/remote";

let remote = inject("remote") as Remote
let preferences = inject('preferences')
let state = reactive({
    timings: [] as Timing[],
    grouped: [] as Timing[],
    loading: true,
})


onMounted(() => {
    state.loading = true
    handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
    state.timings = remote.timings
    state.loading = false
    state.grouped = groupBy(state.timings, 'name')
    return remote.timings
}

function groupBy<T>(xs: T[], key: string): T[] {
    return xs.reduce(function (rv: any, x: any): T {
        (rv[x[key]] = rv[x[key]] || []).push(x);
        return rv;
    }, {});
}

</script>

<template>
    <div>
        <div class="d-flex justify-content-start py-2 px-1">
            <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-clock fa-fw`"></i></div>
            <div class="label-w500 opacity-100 label-xxl px-2">Timings</div>
        </div>
        <div v-if="!state.loading" class="element">
            <div class="timeline">
                <div v-for="module in Object.keys(state.grouped)">
                </div>

            </div>


        </div>
        <div v-else class="element p-2">
            <div class="label-c1 label-o4 d-flex align-content-center gap-1">
                <div>
                    <Loader size="sm"></Loader>
                </div>
                <div class="">Loading...</div>
            </div>
        </div>
    </div>
</template>

<style scoped>

.timeline {

}


</style>
<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Entity, Macro, Task, TaskOption, Zone} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import zoneService from "@/services/zoneService";
import type {Remote} from "@/remote";

const router = useRouter()

const remote = inject("remote") as Remote

onMounted(() => {
    setOptions()
})

const state = reactive({
    tasks: [] as Task[],
    loaded: false,
})

watchEffect(() => {
    setOptions()

})

function setOptions() {
    if (!remote) return;

    let entities: TaskOption[] = remote.entities.map(t => {
        return {title: t.name, description: t.module, value: t.id}
    }) as TaskOption[]

    state.tasks = [
        {
            title: "Name",
            description: "The name of the Macro",
            type: TaskType.String,
            value: "",
            preview: ""
        },
        {
            title: "Entities",
            description: "A brief description of the macro's function.",
            type: TaskType.List,
            options: entities,
            value: [],
            preview: "unset"
        },

    ]
    state.loaded = true
}


function goBack() {
    router.push("/terminal/settings/subroutines")
}

function finish(tasks: Task[]) {
    const name = tasks.find(t => t.title === "Name");
    if (!name) return;

    const entities = tasks.find(t => t.title === "Entities");
    if (!entities) return;

    zoneService.createZone({
        name: name.value as string,
        entities: remote.entities.filter(e => (entities.value as string[]).includes(e.id)) as Entity[],
    } as Zone).then(res => {
        console.log(res)
        goBack()
    }).catch(err => {
        console.log(err)
    })

}


function createSubroutine() {


}

</script>

<template>
    <div class="">
        <div class="d-flex align-items-start label-o4 gap-1 pb-2">
            <div class="label-w500 label-c1 text-accent" @click="goBack">􀆉 Back</div>
        </div>
        <div class="d-flex justify-content-center">
            <div class="generic-grid ">

                <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Macro`"></TaskManager>
            </div>
        </div>
    </div>
</template>

<style scoped>

.generic-grid > div {
    width: 18rem;
}

.generic-grid {
    display: flex;
    justify-content: center;
    grid-column-gap: 0.25rem;
    width: max(60%, 0%);
}
</style>
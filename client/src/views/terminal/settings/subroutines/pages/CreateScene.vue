<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {useRouter} from "vue-router";
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Macro, Task, TaskOption} from "@/types";
import {TaskType} from "@/types";
import TaskManager from "@/components/task/TaskManager.vue";
import macroService from "@/services/macroService";
import type {Remote} from "@/remote";
import Button from "@/components/Button.vue";

const router = useRouter()

const remote = inject("remote") as Remote

const props = defineProps<{
    done?: () => void
}>()

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

    let entities: TaskOption[] = remote.entities.filter(e => e.type === "spectrum").map(t => {
        return {
            title: t.alias ? t.alias : t.module,
            description: remote.attributes.filter(a => a.entity === t.id).map(a => a.key).join(", "),
            value: t.id
        }
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
            title: "Description",
            description: "A brief description of the macro's function.",
            type: TaskType.String,
            value: "",
            preview: ""
        },
        {
            title: "Devices",
            description: "What devices should this macro interact with?",
            type: TaskType.List,
            options: entities,
            value: [],
            preview: ""
        },
        {
            title: "Attribute",
            description: "Which attributes should be targeted within the zone?",
            type: TaskType.Radio,
            options: [{
                title: "Brightness",
                description: "Change the brightness of a light",
                value: "dim"
            }, {
                title: "Color",
                description: "Change the color of a light",
                value: "hue"
            }, {
                title: "Color Temperature",
                description: "Change the color of a light",
                value: "cct"
            }, {
                title: "Power",
                description: "Change the power state of a device",
                value: "on"
            }],
            value: "dim",
            preview: "Brightness"
        },
        {
            title: "Value",
            description: "What should the attribute values be set to?",
            type: TaskType.String,
            value: "",
            preview: ""
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

    const description = tasks.find(t => t.title === "Description");
    if (!description) return;

    const zone = tasks.find(t => t.title === "Zone");
    if (!zone) return;

    const attributeTarget = tasks.find(t => t.title === "Attribute");
    if (!attributeTarget) return;

    const value = tasks.find(t => t.title === "Value");
    if (!value) return;

    macroService.createMacro({
        name: name.value as string,
        description: description.value as string,
        zone: zone.value as string,
        type: attributeTarget.value as string,
        value: value.value as string,
    } as Macro).then(res => {
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
    <div class="ctx " @click="() => {if(props.done) props.done()}">
        <div class="context-grid">
            <div class="context-pane  d-flex flex-column" style="" @click.stop>
                <div class="nav-grid gap-1 pb-1 w-100">
                    <div class="d-flex justify-content-start">

                        <Button accent active fill icon="􀆉" text="Back"
                                @click="() => {if(props.done) props.done()}"></Button>
                    </div>
                    <div class="d-flex justify-content-center">
                        <div class="label-w500 label-c1 label-w600 align-self-center">Create Macro</div>
                    </div>
                    <div class="d-flex justify-content-end gap-1 align-items-center">
                        <!--            <div class="v-sep" style="height: 1rem"></div>-->

                    </div>
                </div>
                <TaskManager :on-complete="finish" :tasks="state.tasks" :title="`Macro`"></TaskManager>


            </div>
        </div>
    </div>

</template>

<style scoped>
.context-pane {
    grid-column: 3 / span 6;
    grid-row: 2 / span 5;
}

.nav-grid {
    width: 100%;
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(1, 1fr);
    grid-template-columns: repeat(3, 1fr);
}

.generic-grid > div {
    width: 18rem;
}

.generic-grid {
    display: flex;
    justify-content: center;
    grid-column-gap: 0.25rem;
    width: max(60%, 0%);
}

.context-grid {
    width: 100%;
    height: 100%;
    display: grid;
    grid-column-gap: 0.25rem;
    grid-row-gap: 0.25rem;
    grid-template-rows: repeat(8, 1fr);
    grid-template-columns: repeat(10, 1fr);
}

.ctx {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 2;
    /*background-color: rgba(0, 0, 0, 0.1);*/
    backdrop-filter: blur(18px) brightness(90%);
}

</style>
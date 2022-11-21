<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import type {SubRoutine} from "@/types";
import {reactive} from "vue";
import {useRouter} from "vue-router";

interface SubRoutineProps {
    subroutine: SubRoutine,
    selected?: boolean
}

const state = reactive({
    toggle: false
})

const router = useRouter()

const props = defineProps<SubRoutineProps>()


</script>

<template>
    <div>
        <div :class="props.selected?'accent-selected':''" class="element p-2 h-100"
        >

            <div class="d-flex justify-content-between">
                <div class="label-xs label-o4 label-w500 pb-2">{{ props.subroutine.icon }}</div>
                <div v-if="props.selected" class="label-c2 label-o4 label-w500 text-accent">􀷙</div>
            </div>
            <div class="label-c2 label-o4 label-w700 lh-1">{{ props.subroutine.description }}</div>
            <div class="label-c3 label-o3 label-w400">{{ props.subroutine.macros.length }}

                macro{{ props.subroutine.macros.length !== 1 ? 's' : '' }}
            </div>

        </div>
        <div v-if="state.toggle" class="element element-menu">
            <div class="grid-element">
                <!--         <Button active text="􀈑" @click="deleteSubRoutine"></Button>-->
            </div>
        </div>
    </div>
</template>

<style scoped>

.element {
    /*border: 2px solid rgba(255, 149, 0, 0);*/
}

.element-menu {
    position: absolute;
    width: 10rem;

    z-index: 1000;
    margin-top: 0.25rem;

    /*left: calc(-100% - 0.25rem);*/
    /*top: calc(100% + 0.25rem);*/
}

.grid-fill {
    position: relative;
    grid-column: 3 / 5;
    display: flex;
    flex-direction: column;
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
</style>
<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import type {SubRoutine} from "@/types";
import subroutineService from "@/services/subroutineService";

const props = defineProps<{
    subroutine: SubRoutine
}>()

function triggerSubroutine() {
    subroutineService.triggerManual(props.subroutine.id).then().catch()
}

</script>

<template>
    <div class="surface" @click="triggerSubroutine">
        {{ props.subroutine.description }}
    </div>
</template>

<style scoped>
.surface {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 1rem 0.25rem;
}
</style>
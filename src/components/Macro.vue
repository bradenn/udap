<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {Macro} from "@/types";
import macroService from "@/services/macroService";

const props = defineProps<{
    macro?: Macro
}>()

function runMacro() {
    if (!props.macro?.id) return
    macroService.runMacro(props.macro?.id)
}

</script>

<template>
    <div class="surface px-3" @click="runMacro">

        <div>{{ props.macro.name }}</div>

    </div>
</template>

<style scoped>
.surface {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 1rem 0.25rem;
}
</style>
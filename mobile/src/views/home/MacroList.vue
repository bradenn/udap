<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import MacroDom from "@/components/Macro.vue";
import {Macro} from "@/types";
import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";

const remote = core.remote()

const state = reactive({
    macros: [] as Macro[]
})

onMounted(() => {
    state.macros = remote.macros
})

watchEffect(() => {
    state.macros = remote.macros
    return remote.macros
})
</script>

<template>
    <div v-for="sr in state.macros" v-if="state.macros" :key="sr.id">
        <MacroDom :macro="sr"></MacroDom>
    </div>
</template>

<style scoped>
.surface {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 1rem 0.25rem;
}
</style>
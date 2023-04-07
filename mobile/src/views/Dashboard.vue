<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import AppLink from "@/components/AppLink.vue";
import MacroDom from "@/components/Macro.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Macro} from "@/types";

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
    <div class="d-flex flex-column gap-2">
        <div class="home-grid">
            <AppLink icon="" name="Remote"></AppLink>
            <AppLink icon="" name="Subroutines"></AppLink>
        </div>
        <div v-if="state.macros" class="d-flex gap-2 flex-column">
            <MacroDom v-for="macro in state.macros" :key="macro.id" :macro="macro"></MacroDom>
        </div>
    </div>
</template>

<style scoped>

.home-grid {
    display: grid;
    gap: 0.25rem;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    grid-template-rows: repeat(1, minmax(0, 1fr));
}
</style>
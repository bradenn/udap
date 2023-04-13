<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import AppLink from "@/components/AppLink.vue";
import {onMounted, reactive, watchEffect} from "vue";
import type {Entity as EntityType, Zone} from "@/types";

const remote = core.remote()

const state = reactive({
    entities: [] as EntityType[],
    zones: [] as Zone[],
})

onMounted(() => {
    state.entities = remote.entities
    state.zones = remote.zones.filter(z => z.pinned)
})

watchEffect(() => {
    state.zones = remote.zones.filter(z => z.pinned)
    return remote.entities
})


</script>

<template>
    <div class="d-flex flex-column gap-2">
        <div class="home-grid">
            <AppLink icon="" name="Home" to="/home/dashboard"></AppLink>
            <AppLink icon="" name="Subroutines" to="/home/subroutines"></AppLink>
            <AppLink icon="" name="Macros" to="/home/macros"></AppLink>
        </div>
        <!--        -->
        <router-view></router-view>

    </div>
</template>

<style scoped>

.home-grid {
    display: grid;
    gap: 0.25rem;
    grid-template-columns: repeat(3, minmax(0, 1fr));
    grid-template-rows: repeat(1, minmax(0, 1fr));
}
</style>
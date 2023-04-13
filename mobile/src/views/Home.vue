<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";
import type {Entity as EntityType, Zone} from "@/types";
import Entity from "@/components/Entity.vue";

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
    <div class="d-flex flex-column gap-3">
        <div v-for="zone in state.zones">
            <div>
                <div class="label-c5 label-w700 label-o5 px-2">{{ zone.name }}</div>
                <div class="home-grid">
                    <Entity v-for="e in zone.entities" :key="e.id" :entity="e"></Entity>
                </div>
            </div>
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
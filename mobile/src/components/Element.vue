<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
const props = defineProps<{
    aspect?: number
    secondary?: boolean
}>()
</script>

<template>
    <div v-if=!props.secondary :style="`aspect-ratio: ${props.aspect?props.aspect:'unset'};`" class="element">
        <slot></slot>
    </div>
    <div v-else class="surface">
        <slot></slot>
    </div>
</template>

<style scoped>
.surface {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 1rem 0.25rem;
}
</style>
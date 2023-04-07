<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Element from "./Element.vue";

const props = defineProps<{
    name: string,
    icon: string,
}>()

</script>

<template>
    <Element :aspect="4" class="d-flex justify-content-center " secondary style="border-radius: 8px">
        <div>{{ props.icon }}</div>
        <div class="app-link-header">{{ props.name }}</div>
    </Element>
</template>

<style scoped>
.app-link-header {
    font-size: 1rem;
    font-weight: 400;
    color: rgba(255, 255, 255, 0.5);

}
</style>
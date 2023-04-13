<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import Hash from "./Hash.vue";
import core from "@/core";

const props = defineProps<{
    name?: number
}>()


const remote = core.remote()

</script>

<template>
    <div class="navbar-dropdown d-flex align-content-center align-items-center justify-content-between" @click="">
        <div class="d-flex align-items-center gap-2 align-items-center">
            <Hash data="bn"></Hash>
            <div class="label-c5 label-o6 label-w500 lh-1">Text</div>
        </div>
        <div class="label-c2 label-o2" style="margin-top: -8px;">
            ⌵
        </div>
    </div>
</template>

<style scoped>
.navbar-dropdown {
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 0.5rem;
    padding: 0.25rem 0.70rem;
    width: 8rem;
}

.navbar-dropdown:active {

}
</style>
<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "@/core";

const props = defineProps<{
    aspect?: number
}>()

const remote = core.remote()

</script>

<template>
    <div class="d-flex flex-column align-items-center justify-content-center align-content-center"
         style=" z-index: 6 !important;">
        <div class="sf-icon icon-inf">􀎡</div>
        <div :class="remote.connected?'text-success':'text-danger'" class="sf-icon icon-inf">􀯠</div>
    </div>
</template>

<style scoped>
.icon-lock {
    color: orange;
    width: 0.83rem;
    height: 0.83rem;
    background-image: url(/lock.svg);
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    opacity: 0.8;
    /*outline: 1px solid white;*/
}

.icon-inf {
    /*color: white;*/
    width: 1rem;
    display: flex;
    justify-content: center;
    align-content: center;
    align-items: center;
    font-size: 0.7rem;
    opacity: 0.8;
    /*outline: 1px solid white;*/
}
</style>
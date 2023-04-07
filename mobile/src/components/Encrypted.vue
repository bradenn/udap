<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
const props = defineProps<{
    aspect?: number
}>()
</script>

<template>
    <div class="d-flex flex-column align-items-center justify-content-center align-content-center">
        <div class="icon-lock"></div>
        <div class="icon-inf"></div>
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
    color: orange;
    width: 0.83rem;
    height: 0.83rem;
    background-image: url(/inf.svg);
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    opacity: 0.8;
    /*outline: 1px solid white;*/
}
</style>
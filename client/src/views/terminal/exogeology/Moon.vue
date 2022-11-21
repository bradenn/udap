<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import VSlider from "@/components/VSlider.vue";</script>
<template>
    <div class="d-flex flex-column justify-content-center align-items-center gap-2 w-100 h-100">
        <h4>Demo UI Elements</h4>
        <div style="width: 45rem; height: 100%">
            <div class="label-c1 label-w400 label-o4 px-2 mb-1">Vertical Sliders</div>
            <div class="d-flex gap-1" style="height: 20rem">
                <VSlider :key="1" :change="() => {}"></VSlider>
                <VSlider :key="2" :change="() => {}"></VSlider>
                <VSlider :key="3" :change="() => {}"></VSlider>
            </div>

        </div>
    </div>
</template>

<style scoped>

</style>
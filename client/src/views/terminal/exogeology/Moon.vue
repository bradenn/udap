<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>
import Switch from "@/components/Switch.vue";
import Button from "@/components/Button.vue";
import Slider from "@/components/Slider.vue";</script>
<template>
    <div class="d-flex flex-column justify-content-center align-items-center gap-2 w-100">
        <h4>Haptic Controls</h4>
        <div style="width: 45rem">
            <div class="label-c1 label-w400 label-o4 px-2 mb-1">Range Sliders</div>
            <div class="d-flex gap-1">
                <Slider :max="100" :min="0" :step="5" :value="50" name="POWER"></Slider>
            </div>
        </div>
        <div style="width: 45rem">
            <div class="label-c1 label-w400 label-o4 px-2 mb-1">Toggles</div>
            <div class="d-flex gap-1" style="width: 13rem">
                <Switch :active="true"></Switch>

            </div>
        </div>
        <div style="width: 45rem">
            <div class="label-c1 label-w400 label-o4 px-2 mb-1">Buttons</div>
            <div class="d-flex gap-1 element" style="width: 13rem">
                <Button :active="true" style="flex-grow: 1;" text="ON"></Button>
                <Button :active="false" style="flex-grow: 1;" text="OFF"></Button>

            </div>
        </div>
    </div>
</template>

<style scoped>

</style>
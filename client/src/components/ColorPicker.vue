<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {inject} from "vue";

import Color from "@/components/Color.vue";
import {Preference} from "@/preferences";
import FixedScroll from "@/components/scroll/FixedScroll.vue";

const props = defineProps<{
    change?: (a: number) => void
    selected?: number
}>()


function setColor(color: number) {
    if (props.change) {
        props.change(color)
    }
}

const prefs = inject("preferences") as Preference
</script>

<template>
    <div :class="`${props.selected?'selected':''}`" class="element color-container">
        <div class="d-flex justify-content-between px-1 pb-1">
            <div class="label-c1 label-o3 label-w500">Color</div>
            <div class="label-c1 label-o3 label-w500">
                {{ props.selected }}
            </div>

        </div>
        <FixedScroll :horizontal="true" style="overflow-x: scroll">
            <div class="d-flex gap-1">
                <Color v-for="color in prefs.appdata.colors.sort()" :color="color"
                       :selected="props.selected === color"
                       @click="() => setColor(color)">
                </Color>
            </div>
        </FixedScroll>
    </div>
</template>

<style lang="scss" scoped>
div.element.color-container {


  //width: 100%;


}
</style>
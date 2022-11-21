<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import FixedScroll from "@/components/scroll/FixedScroll.vue";

import colors from '@/assets/pantone-colors.json'

function setPage() {

}

</script>

<template>
    <div class="element color-selector">
        <FixedScroll :horizontal="true" class="color-grid">
            <div v-for="c in colors.names" :key="c"
                 :style="`background-color: ${colors.values[colors.names.indexOf(c)]};`"
                 class="color subplot">
                <div class="color-text text-capitalize">{{ c.replace('-', ' ') }}</div>
            </div>
        </FixedScroll>
    </div>
</template>

<style lang="scss" scoped>
div.element.color-selector {

  div.subplot.color {
    opacity: 1;

    padding: 0.5rem;
    width: 8rem;
  }

  .color-grid {
    padding: 1rem;
    display: grid;
    align-items: start;
    align-content: start;
    grid-auto-flow: column;
    grid-gap: 0.25rem;
    //grid-auto-flow: row;
    grid-template-rows: repeat(7, 1fr);
    grid-template-columns: repeat(auto-fill, minmax(8rem, 1fr));

    overflow-x: scroll;
  }

  .color-text {

    //text-shadow: 0 0 2px rgba(0, 0, 0, 0.25);
    filter: invert(1) contrast(200%);
  }
}
</style>
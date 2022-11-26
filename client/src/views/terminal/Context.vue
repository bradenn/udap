<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import {onMounted, onUnmounted} from "vue";
import core from "@/core";

const ctx = core.context()
onMounted(() => {
    ctx.showContext()
})

onUnmounted(() => {
    ctx.hideContext()
})

</script>

<template>
    <div class="ctx">
        <slot></slot>
    </div>
</template>


<style lang="scss" scoped>
.ctx::before {
  content: ' ';
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 3 !important;
  animation: ctx-load 25ms ease-out forwards;
}

.ctx {
  position: absolute;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;

  z-index: 6 !important;
  display: grid;
  grid-column-gap: 0;
  grid-template-columns: repeat(15, minmax(4rem, 1fr));
  grid-template-rows: repeat(8, minmax(2rem, 1fr));
}

@keyframes ctx-load {
  0% {
    backdrop-filter: blur(0) brightness(100%);
  }
  50% {
    backdrop-filter: blur(2px) brightness(98%);
  }
  100% {
    backdrop-filter: blur(8px) brightness(95%);
  }
}

.ctx {
}
</style>
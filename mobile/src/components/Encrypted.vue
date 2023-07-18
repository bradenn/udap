<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import core from "udap-ui/core";
import type {Remote} from "udap-ui/remote";

const props = defineProps<{
  aspect?: number,
  compact?: boolean
}>()

const remote = core.remote() as Remote

</script>

<template>
  <div :class="props.compact?'flex-row':'flex-column'"
       class="d-flex  align-items-center justify-content-center align-content-center" style=" z-index: 6 !important;">
    <div class="sf-icon icon-inf">{{ remote.client.url.includes('wss') ? '􀎡' : '􀎣' }}</div>
    <div :class="remote.client.connected?'text-green':'text-danger'" class="sf-icon icon-inf">􀯠</div>
  </div>

</template>

<style scoped>
.text-green {
  color: hsl(118, 73%, 45%);
  font-weight: 800;
  font-size: 3rem;
}

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
  transition: 100ms;
  width: 1rem;
  display: flex;
  justify-content: center;
  align-content: center;
  align-items: center;
  font-size: 0.6rem;

  opacity: 0.8;
  /*outline: 1px solid white;*/
}
</style>
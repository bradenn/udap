<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import useListModules, {ModuleController} from "@/controller/modulesController";
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";


let state = useListModules() as ModuleController;

function convertNanosecondsToString(nanoseconds: number): string {
  if (nanoseconds < 0) {
    throw new Error("Input must be a non-negative number.");
  }

  if (nanoseconds === 0) {
    return "0ns";
  }

  const units = ["ns", "µs", "ms", "s"];
  let value = nanoseconds;
  let unitIndex = 0;

  while (value >= 1000 && unitIndex < units.length - 1) {
    value /= 1000;
    unitIndex++;
  }

  return `${value.toFixed(2)}${units[unitIndex]}`;
}
</script>

<template>

  <List v-if="state.loaded">

    <Element v-for="timing in state.timings" :key="timing.timing.pointer" :foreground="true"
             class="d-flex gap-2 align-items-center px-2 py-2 dramatic-load"
             mutable>

      <div class="d-flex flex-column gap-1">

        <div class="mono lh-1 px-2 ">{{ timing.module ? timing.module.name : timing.timing.pointer }}</div>


      </div>
      <div class="flex-fill"></div>
      <div v-if="timing.module" :style="`padding: ${0.125+timing.timing.delta/0.3e9}em;`"
           class="fill-animate mono label-c5 label-o3 px-2 ">
        {{ convertNanosecondsToString(timing.timing.delta) }}
      </div>
      <div v-else class="mono label-c3 px-2 label-o4 label-w600 ">
        {{ convertNanosecondsToString(timing.timing.delta) }}
      </div>
    </Element>

  </List>


</template>

<style lang="scss" scoped>
.dramatic-load {
  animation: loadUp 100ms ease-in forwards;
}

@keyframes loadUp {
  0% {
    opacity: 0.2;
  }
  100% {
    opacity: 1;
  }
}

.fill-animate {

  transition: padding 1000ms ease;
}


.mono {
  font-family: "JetBrains Mono", sans-serif !important;
}
</style>
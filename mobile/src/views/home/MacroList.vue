<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import MacroDom from "@/components/Macro.vue";
import {Macro, Zone} from "udap-ui/types";
import core from "@/core";
import {onMounted, reactive, watchEffect} from "vue";

const remote = core.remote()

const state = reactive({
  macros: [] as Macro[],
  zones: [] as Zone[]
})

onMounted(() => {
  state.macros = remote.macros
  state.zones = remote.zones
})

watchEffect(() => {
  state.zones = remote.zones.filter(z => !z.deleted)
  state.macros = remote.macros.filter(m => state.zones.find(z => z.id === m.zone))
  return [remote.macros, remote.zones]
})

//@ts-ignore
function groupBy<T, K>(xs: T[], key: (x: T) => K): Record<K, T[]> {
//@ts-ignore
  return xs.reduce((rv: Record<K, T[]>, x: T) => {
//@ts-ignore
    (rv[key(x)] = rv[key(x)] || []).push(x);
//@ts-ignore
    return rv;
//@ts-ignore
  }, {});
}

</script>

<template>
  <div class="d-flex flex-column gap-1">
    <div v-for="(value, zone) in groupBy(state.macros, (m: Macro) => m.zone)" v-if="state.macros" :key="zone"
         class="d-flex flex-column gap-1 mt-1">
      <div class="text-capitalize label-c4 label-w600 px-2 label-o6">{{
          state.zones.find(z => z.id === zone)?.name
        }}
      </div>
      <div v-for="m in value">
        <MacroDom :macro="m"></MacroDom>

      </div>
    </div>
  </div>
</template>

<style scoped>
.macro-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(1rem, 1fr));
}

.surface {
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 1rem 0.25rem;
}
</style>
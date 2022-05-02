<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";
import Loader from "@/components/Loader.vue";
import Plot from "@/components/plot/Plot.vue";

let remote = inject('remote') as Remote
let preferences = inject('preferences')

let state = reactive({
  modules: {} as Entity[],
  attributes: {} as Attribute[],
  loading: true
})

onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.modules = groupBy<Entity>(remote.entities, 'module') as Entity[]
  state.attributes = groupBy<Attribute>(remote.attributes, 'entity') as Attribute[]
  state.loading = false
  return remote
}

// groupBy creates several arrays of elements based on the value of a key
function groupBy<T>(xs: T[], key: string): T[] {
  return xs.reduce(function (rv: any, x: any): T {
    (rv[x[key]] = rv[x[key]] || []).push(x);
    return rv;
  }, {});
}
</script>

<template>
  <div>
    <div class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-layer-group fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Modules</div>
    </div>
    <div v-if="!state.loading" class="d-flex flex-column gap-1">

      <div v-for="(entities, module) in state.modules" class=" d-flex justify-content-start flex-column">
        <div v-if="entities" class="d-flex flex-column gap">
          <div class="label-xs label-r px-1 " style="text-transform: capitalize;">{{ module }}</div>
          <Plot :cols="5" :rows="1">
            <div v-for="entity in entities" class="subplot">
              <div class="d-flex justify-content-start align-items-center flex-row px-1 w-100">
                <div class="label-w500 label-o3 label-c1">{{ entity.icon }}&nbsp;</div>
                <div class="label-w500 label-c1">{{ entity.name }}</div>
                <div class="flex-grow-1"></div>
                <div class="label-w500 label-c2 label-o2"><i class="fa-solid fa-gear"></i></div>
              </div>
            </div>
          </Plot>

        </div>
      </div>


    </div>
    <div v-else class="element p-2">
      <div class="label-c1 label-o4 d-flex align-content-center gap-1">
        <div>
          <Loader size="sm"></Loader>
        </div>
        <div class="">Loading...</div>
      </div>
    </div>
  </div>
</template>

<style scoped>


</style>
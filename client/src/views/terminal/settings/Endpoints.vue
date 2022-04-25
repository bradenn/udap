<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint, Remote} from "@/types";
import Plot from "@/components/plot/Plot.vue";
import Radio from "@/components/plot/Radio.vue";

import moment from "moment";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  endpoints: {} as Endpoint[],
  loading: true,
})


onMounted(() => {
  state.loading = true
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.endpoints = remote.endpoints
  state.loading = false
  return remote
}

</script>

<template>
  <div>
    <div class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-expand fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Endpoints</div>
      <div class="flex-fill"></div>

      <Plot :cols="1" :rows="1" small>
        <Radio :active="false" :fn="() => {}" title="New Endpoint"></Radio>
      </Plot>
    </div>

    <div class="d-flex flex-column gap-1 element p-2">
      <div v-for="endpoint in remote.endpoints"
           :key="endpoint.id" class="label-c2 d-flex justify-content-between align-items-center">
        <div class="d-flex flex-column justify-content-between">
          <div class="label-c1">{{ endpoint.name }}</div>
          <div class="label-c1 label-o4 label-r">{{ endpoint.key }}</div>
        </div>
        <div class="d-flex flex-column justify-content-between">
          <div>Edit</div>
          <div class="label-c3 label-o2">Created {{ moment(endpoint.created).format('MMMM DD, YYYY @ hh:mm A') }}</div>
        </div>
      </div>

    </div>

  </div>
</template>

<style scoped>


</style>
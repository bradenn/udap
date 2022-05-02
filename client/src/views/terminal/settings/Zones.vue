<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Endpoint, Remote} from "@/types";
import Loader from "@/components/Loader.vue";

let remote = inject("remote") as Remote
let preferences = inject('preferences')

let state = reactive({
  endpoints: {} as Endpoint[],
  loading: true,
})


onMounted(() => {
  state.loading = true
  handleUpdates(remote)
})

watchEffect(() => handleUpdates(remote))

function handleUpdates(remote: Remote) {
  state.endpoints = remote.endpoints
  state.loading = false
  return remote.endpoints
}

</script>

<template>
  <div class="d-flex justify-content-start py-2 px-1">
    <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-map fa-fw`"></i></div>
    <div class="label-w500 opacity-100 label-xxl px-2">Zones</div>
    <div class="flex-fill"></div>
  </div>
  <div v-if="!state.loading">

  </div>
  <div v-else>

    <div class="element p-2">
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
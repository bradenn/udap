<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Encrypted from "@/components/Encrypted.vue";
import core from "udap-ui/core";

const router = core.router();
const remote = core.remote();

</script>

<template>
  <div class="d-flex justify-content-between">
    <div
        class="d-flex flex-row gap-1 justify-content-start align-items-center align-content-center px-1 pb-1 flex-shrink-0"
        style="height: 1.5rem">
      <Encrypted :compact="!(router.currentRoute.value.fullPath === '/home/dashboard')"></Encrypted>
      <div
          :class="`udap-logo-container ${router.currentRoute.value.fullPath === '/home/dashboard'?'udap-logo':'udap-logo-sm'}`"
          class="lh-1"
          style=" z-index: 6 !important;">UDAP
      </div>

    </div>
    <div class="px-1 label-o5 label-w600">{{
        remote.client.oob ? "Out-Of-Band" : remote.client.host ? "" : "Away"
      }}

    </div>
  </div>
</template>

<style lang="scss" scoped>

</style>
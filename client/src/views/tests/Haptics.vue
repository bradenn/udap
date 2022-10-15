<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Switch from "@/components/Switch.vue";
import Button from "@/components/Button.vue";</script>
<template>
  <div class="d-flex flex-column justify-content-center align-items-center w-100 gap-2">
    <h4>Haptic Controls</h4>
    <div style="width: 45rem">
      <div class="label-c1 label-w400 label-o4 px-2 mb-1">Range Sliders</div>
      <div class="d-flex gap-1">
      </div>
    </div>
    <div style="width: 45rem">
      <div class="label-c1 label-w400 label-o4 px-2 mb-1">Toggles</div>
      <div class="d-flex gap-1" style="width: 13rem">
        <Switch :active="true"></Switch>

      </div>
    </div>
    <div style="width: 45rem">
      <div class="label-c1 label-w400 label-o4 px-2 mb-1">Buttons</div>
      <div class="d-flex gap-1" style="width: 13rem">
        <Button :active="true" text="Penis"></Button>

      </div>
    </div>
  </div>
</template>

<style scoped>

</style>
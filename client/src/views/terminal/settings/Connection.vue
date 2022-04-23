<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import {inject} from "vue";
import {Preference} from "@/preferences";
import {PreferenceTypes} from "@/types";

let remote = inject('remote')
let preferences = inject('preferences')

let controller = new Preference(PreferenceTypes.Controller).get()

</script>

<template>
  <div>
    <div class="d-flex justify-content-start py-2 px-1">
      <div class="label-w500 label-o4 label-xxl"><i :class="`fa-solid fa-cloud fa-fw`"></i></div>
      <div class="label-w500 opacity-100 label-xxl px-2">Connection</div>
    </div>
    <div class="px-2 w-50">
      <div v-if="remote" class="element mb-1">
        <div class="sidebar-item d-flex justify-content-start pb-1 px-1">
          <div class="label-w500 label-o4 label-xs lh-2"><i :class="`fa-solid fa-cloud fa-fw`"></i></div>
          <div class="label-w500 label-o6 label-xs lh-2 px-2">Nexus</div>
        </div>
        <div class="h-sep"></div>
        <div class="d-flex justify-content-between">
          <div class="label-xs label-o5">Address</div>
          <div class="label-xs label-o4">{{ controller }}</div>
        </div>
      </div>
      <div class="element">
        <div class="sidebar-item d-flex justify-content-start pb-1 px-1">
          <div class="label-w500 label-o4 label-xs lh-2"><i :class="`fa-solid fa-cloud fa-fw`"></i></div>
          <div class="label-w500 label-o6 label-xs lh-2 px-2">Nexus</div>
        </div>


      </div>


    </div>
  </div>
</template>

<style scoped>


</style>
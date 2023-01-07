<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Toolbar from "@/components/toolbar/Toolbar.vue";
import type {RouteMeta} from "@/composables/routeMeta.ts";
import {useRouteMeta} from "@/composables/routeMeta.ts";

let meta: RouteMeta = useRouteMeta()

</script>

<template>
  <div class="h-100 d-flex flex-column gap-1">
    <Toolbar :icon="meta.icon" class="flex-shrink-0"
             title="Entities">
    </Toolbar>
    <div class="flex-grow-1"
         style="height: calc(100% - 2rem)">
      <router-view></router-view>
    </div>
  </div>
</template>

<style scoped>

</style>
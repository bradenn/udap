<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

import Toolbar from "@/components/toolbar/Toolbar.vue";
import {useRouteMeta} from "@/composables/routeMeta";

const meta = useRouteMeta()

</script>

<template>
  <Toolbar :icon="meta.icon" class="mb-1"
           title="Entities">
  </Toolbar>
  <router-view></router-view>
</template>

<style scoped>

</style>
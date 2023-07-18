<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import useListModules, {ModuleController} from "@/controller/modulesController";
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";


let state = useListModules() as ModuleController;

</script>

<template>

  <List v-if="state.modules" scroll-y style="max-height: inherit">

    <Element v-for="module in state.modules" :key="module.id" :foreground="true" :to="`/home/settings/modules/${module.id}`"
             class="d-flex gap-2 align-items-center px-2 py-2"
             mutable>
      <div class="notches px-2" style="height: 1rem">
        <div :class="`${module.enabled?'active':''}`" class="notch h-100"></div>

      </div>
      <div class="d-flex flex-column gap-1">
        <div class="label-monospace lh-1">{{ module.name }}</div>
        <div class="label-monospace lh-1 label-o2 label-c6">v{{ module.version }}</div>


      </div>

      <div class="flex-grow-1 d-flex justify-content-end align-items-end">
        {{ module.state }}
        <div class="sf-icon label-o2 label-w500 label-c3">􀆊</div>
      </div>
    </Element>

  </List>
</template>

<style lang="scss" scoped>

</style>
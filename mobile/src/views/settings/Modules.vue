<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import useListModules, {ModuleController} from "@/controller/modulesController";
import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";


const state = useListModules() as ModuleController;

</script>

<template>
  <Element>
    <List class="overflow-scroll">

      <Element v-for="module in state.modules" :foreground="true" class="d-flex gap-2">
        <div class="notches">
          <div :class="`${true?'active':''}`" class="notch"></div>
          <div :class="`${true?'active':'warn'}`" class="notch"></div>
          <div class="notch"></div>
        </div>
        <div>
          <div class="label-monospace">{{ module.name }}</div>
          <div></div>
        </div>
      </Element>

    </List>
  </Element>

</template>

<style lang="scss" scoped>

</style>
<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";
import useListEndpoints, {EndpointController} from "@/controller/endpointController";
import moment from "moment";
import ElementHeader from "udap-ui/components/ElementHeader.vue";


let state = useListEndpoints() as EndpointController;

</script>

<template>
  <List v-if="state.endpoints">

    <ElementHeader title="Online"></ElementHeader>

    <Element v-for="module in state.endpoints.filter(e => e.connected)" :foreground="true"
             class="d-flex gap-2 align-items-center px-2 py-2"
             mutable>
      <div class="notches px-2" style="height: 1rem">
        <div :class="`${module.connected?'active':''}`" class="notch h-100"></div>

      </div>
      <div class="d-flex flex-column gap-1">
        <div class="label-monospace lh-1">{{ module.name ? module.name : 'Unnamed' }}</div>

        <div class="label-monospace lh-1 label-o2 label-c6">{{ moment(module.updated).fromNow() }}</div>

      </div>
      <div class="flex-fill"></div>
      <div class="label-c5 px-2">{{ module.key }}</div>
    </Element>

    <ElementHeader title="Recent"></ElementHeader>
    <Element v-for="module in state.endpoints.filter(e => !e.connected)" :foreground="true"
             class="d-flex gap-2 align-items-center px-2 py-2"
             mutable>
      <div class="notches px-2" style="height: 1rem">
        <div :class="`${module.connected?'active':''}`" class="notch h-100"></div>

      </div>
      <div class="d-flex flex-column gap-1">
        <div class="label-monospace lh-1">{{ module.name ? module.name : 'Unnamed' }}</div>
        <div class="label-monospace lh-1 label-o2 label-c6">{{ moment(module.updated).fromNow() }}</div>

      </div>
      <div class="flex-fill"></div>
      <div class="label-c5 px-2">{{ module.key }}</div>
    </Element>

  </List>
</template>

<style lang="scss" scoped>

</style>
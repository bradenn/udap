<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import Element from "udap-ui/components/Element.vue";
import List from "udap-ui/components/List.vue";

import useListEndpoints, {EndpointController} from "@/controller/endpointController";
import {since} from "udap-ui/time"

import ElementHeader from "udap-ui/components/ElementHeader.vue";
import core from "udap-ui/core";


let state = useListEndpoints() as EndpointController;
const preferences = core.preferences()

function timeSince(date: string): string {
  return since(date)
}

</script>

<template>
  <List v-if="state.endpoints">

    <ElementHeader title="Online"></ElementHeader>

    <Element v-for="endpoint in state.endpoints.filter(e => e.connected)" :foreground="true"
             class="d-flex gap-2 align-items-center px-2 py-2"
             mutable>
      <div class="notches px-2" style="height: 1rem">
        <div :class="`${endpoint.connected?'active':''}`" class="notch h-100"></div>

      </div>
      <div class="d-flex flex-column gap-1">
        <div class="label-monospace lh-1 d-flex align-items-center gap-1">
          <div v-if="endpoint.notifications" :style="`color: ${preferences.accent} !important; `" class="sf label-c6">
            {{ endpoint.notifications ? '􀈠' : '' }}
          </div>
          {{
            endpoint.name ? endpoint.name : 'Unnamed'
          }}
        </div>

        <div class="label-monospace lh-1 label-o2 label-c6">{{ timeSince(endpoint.updated) }}</div>

      </div>
      <div class="flex-fill"></div>
      <div class="label-c5 px-2">{{ endpoint.key }}</div>
    </Element>

    <ElementHeader title="Recent"></ElementHeader>
    <Element v-for="module in state.endpoints.filter(e => !e.connected)" :key="module.id" :foreground="true"
             class="d-flex gap-2 align-items-center px-2 py-2"
             mutable>
      <div class="notches px-2" style="height: 1rem">
        <div :class="`${module.connected?'active':''}`" class="notch h-100"></div>

      </div>
      <div class="d-flex flex-column gap-1">
        <div class="label-monospace lh-1 d-flex align-items-center">
          <div v-if="module.notifications" :style="`color: ${preferences.accent} !important; `" class="sf label-c6">
            {{ module.notifications ? '􀈠' : '' }}
          </div>
          {{ module.name ? module.name : 'Unnamed' }}
        </div>

        <div class="label-monospace lh-1 label-o2 label-c6">{{ timeSince(module.created) }}</div>

      </div>
      <div class="flex-fill"></div>
      <div class="label-c5 px-2">{{ module.key }}</div>
    </Element>

  </List>
</template>

<style lang="scss" scoped>

</style>
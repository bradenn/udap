<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import ElementLink from "udap-ui/components/ElementLink.vue";
import List from "udap-ui/components/List.vue";
import {onMounted, reactive, watchEffect} from "vue";
import {Zone} from "@/types";
import core from "@/core";
import ElementHeader from "udap-ui/components/ElementHeader.vue";


const remote = core.remote();

const state = reactive({
  zones: [] as Zone[]
})

onMounted(() => {
  state.zones = remote.zones
})

watchEffect(() => {
  state.zones = remote.zones
  return remote.zones
})

</script>

<template>
  <List>
    <ElementHeader title="Pinned"></ElementHeader>
    <ElementLink v-for="module in state.zones.filter(z => z.pinned)" :key="module.id"
                 :title="module.name" :to="`/settings/zones/${module.id}`" icon="􀞿">
      <div class="d-flex">
        <div v-for="icon in module.entities.map(e => e.icon)" class="sf-icon">
          {{ icon }}
        </div>
      </div>
    </ElementLink>
    <ElementHeader title="Unpinned"></ElementHeader>
    <ElementLink v-for="module in state.zones.filter(z => !z.pinned)" :key="module.id"
                 :title="module.name" :to="`/settings/zones/${module.id}`" icon="􀞿">
      <div class="d-flex">
        <div v-for="icon in module.entities.map(e => e.icon)" class="sf-icon">
          {{ icon }}
        </div>
      </div>
    </ElementLink>
  </List>
</template>

<style lang="scss" scoped>

</style>
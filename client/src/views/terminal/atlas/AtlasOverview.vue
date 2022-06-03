<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>

import {inject, onMounted, reactive, watchEffect} from "vue";
import type {Attribute, Entity, Remote} from "@/types";


let remote = inject("remote") as Remote

let state = reactive({
  entity: {} as Entity,
  buffer: {} as Attribute,
  loaded: false,
})

onMounted(() => {
  load(remote)
})

watchEffect(() => load(remote))

function load(r: Remote) {
  let target = r.entities.find((entity: Entity) => entity.name === "atlas")
  if (!target) {
    target = {} as Entity
  }
  state.entity = target
  let buffer = r.attributes.find((attribute: Attribute) => attribute.entity === state.entity.id && attribute.key === 'buffer')
  if (!buffer) {
    buffer = {} as Attribute
  }
  state.buffer = buffer
  state.loaded = true
  return r
}


</script>

<template>
  <div class="d-flex justify-content-center align-items-start h-100">

    <div v-if="state.loaded" class="mt-4">
      <div>
        {{ state.buffer.value }}
      </div>
    </div>
    <div v-else>
      Loading
    </div>

  </div>
</template>

<style lang="scss" scoped>
</style>
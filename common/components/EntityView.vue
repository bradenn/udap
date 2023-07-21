<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>

import {Entity} from "../types";
import Element from "./Element.vue";

const props = defineProps<{
  entity: Entity
}>()

</script>

<template>
  <Element :foreground="true" :mutable="true" :to="`/home/settings/entities/${props.entity.id}`"
           class="d-flex align-items-center gap-2 ent" style="height: 3.25rem">
    <div class="notches">
      <div :class="`${props.entity.updated?'active':''}`" class="notch"></div>
      <div :class="`${props.entity.updated?'active':'warn'}`" class="notch"></div>
      <div class="notch"></div>
    </div>

    <div class="d-flex flex-column gap-0 justify-content-center-center">
      <div></div>
      <div class="label-c3 lh-1 label-o5 d-flex">
        <div class="sf-icon label-c5" style="width: 15px; margin-right: 6px">
          {{ props.entity.icon }}
        </div>
        {{ props.entity.alias ? `${props.entity.alias}ª` : props.entity.name }}
      </div>
      <div class="label-c5 lh-1 label-o3">{{ props.entity.type }}</div>
    </div>
    <div class="flex-grow-1 d-flex justify-content-end">

      <div class="sf-icon label-o3 label-w500 label-c3">􀆊</div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>
.ent {
  user-select: none;
}

</style>
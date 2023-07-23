<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, inject, reactive} from "vue"
import core from "../core"
import Element from "./Element.vue";
import {PreferencesRemote} from "../persistent";

const props = defineProps<{
  title?: string
  alt?: string
  icon?: string
  outbound?: string
  to?: string
  cb?: () => void
  button?: boolean
}>()

const state = reactive({})

const router = core.router()

const preferences = inject("preferences") as PreferencesRemote


</script>

<template>
  <Element :cb="props.cb" :foreground="true" :mutable="true" :to="props.to"
           class="d-flex align-items-center gap-2 ent" style="height: 3.25rem">
    <div class="d-flex flex-column gap-0 justify-content-center-center px-2">

      <div class="label-c4 lh-1 label-o5 d-flex">
        <div class="sf-icon label-c5" style="width: 15px; margin-right: 12px">
          {{ props.icon }}
        </div>
        <div class="d-flex flex-column justify-content-start" style="gap: 0.125rem">
          <div class="label-o5 label-c4">{{ props.title }}</div>
          <div v-if="props.alt" class="label-o2 label-c6">{{ props.alt }}</div>
        </div>
      </div>
      <!--      <div class="label-c5 lh-1 label-o3">sds</div>-->
    </div>
    <div class="flex-grow-1 d-flex justify-content-end align-items-end">
      <slot></slot>
      <div v-if="!props.button" class="sf label-o3 label-w500 label-c3 px-2">􀆊</div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>
.sf-icon {
  color: hsl(0, 0%, 50%);
}

</style>
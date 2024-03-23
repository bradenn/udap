<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, inject, reactive} from "vue"
import core from "../core"
import Element from "./Element.vue";
import {PreferencesRemote} from "../persistent";

const props = defineProps<{
  title?: string
  alt?: string
  value?: string
  icon?: string
  copyable?: boolean
  noPadding?: boolean
  outbound?: string
  to?: string
  cb?: () => void
  button?: boolean
}>()

const state = reactive({})

const router = core.router()

const preferences = inject("preferences") as PreferencesRemote

function copyValue() {
  if (props.value) {
    navigator.clipboard.writeText(props.value);
  }
}

</script>

<template>
  <Element :cb="props.cb" :foreground="true" :long-cb="props.cb" :mutable="!(!props.cb) || false" :to="props.to"
           class="d-flex justify-content-between align-items-center flex-row" style="height: 3.25rem">

    <div class="label-c4 lh-1 label-o5 d-flex align-items-center px-2">
      <div class="sf-icon label-c4" style="width: 15px; margin-right: 12px">
        {{ props.icon }}
      </div>
      <div class="d-flex flex-column justify-content-start" style="gap: 0.125rem">
        <div class="label-o5 label-c4">{{ props.title }}</div>
        <div v-if="props.alt" class="label-o2 label-c6">{{ props.alt }}</div>
      </div>
    </div>
    <!--      <div class="label-c5 lh-1 label-o3">sds</div>-->
    <div class="d-flex align-items-center">
      <div v-if="props.noPadding" style="margin-right: -2px">
        <slot></slot>
      </div>
      <div v-else class="flex-grow-1 d-flex align-items-end justify-content-end position-relative px-1">
        <slot v-if="!props.value"></slot>
        <div v-else class="label-truncate label-o3 label-c4 sf">{{ props.value }}</div>
      </div>
      <div v-if="props.copyable">
        <Element :cb="copyValue" class="sf d-flex align-items-center justify-content-center label-o3" foreground mutable
                 style="width: 3rem !important; margin-right: -4px !important;">
          􀉂
        </Element>
      </div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>
.sf-icon {
  color: hsl(0, 0%, 50%);
}

//.label-truncate:before {
//  content: '';
//  width: 4rem;
//  height: 1.5rem;
//  position: absolute;
//  right: 0;
//  top: calc(50% - 0.75rem);
//  filter: blur(8px);
//  background: linear-gradient(90deg, transparent 4px, rgba(0, 0, 0, 0.2));
//}

.label-truncate {

  max-width: 12rem;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}
</style>
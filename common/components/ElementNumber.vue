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
  unit?: string
  immutable?: boolean
  outbound?: string
  value: number
  to?: string
  cb?: () => void
  button?: boolean
}>()

const state = reactive({
  active: false
})

const router = core.router()


const preferences = inject("preferences") as PreferencesRemote

function handleInput(m: TouchEvent) {
  let elem = m.target as HTMLInputElement
  if (!elem) return
  elem.setSelectionRange(0, 0);
  elem.select()

}

</script>

<template>
  <Element :cb="() => {state.active = !state.active; }" :foreground="true" :mutable="true"
           class="d-flex align-items-center gap-2 ent" style="height: 3.25rem; padding-right: 0.25rem !important;">
    <div class="d-flex flex-column gap-0 justify-content-center-center px-2">

      <div class="label-c4 lh-1 label-o5 d-flex">
        <div class="sf-icon label-c6" style="width: 15px; margin-right: 12px">
          {{ props.icon }}
        </div>
        <div class="d-flex flex-column justify-content-start" style="gap: 0.125rem">
          <div class="label-o5 label-c4">{{ props.title }}</div>
          <div v-if="props.alt" class="label-o2 label-c6">{{ props.alt }}</div>
        </div>
      </div>
      <!--      <div class="label-c5 lh-1 label-o3">sds</div>-->
    </div>
    <div class="flex-grow-1 d-flex justify-content-end align-items-end align-items-end"
    >

      <div>
        <div v-if="props.immutable" :class="`${state.active?'input-number-active':''}`"
             class="input-number input-number-active px-2 d-flex align-items-center justify-content-end"
             style="z-index: 10 !important;">
          {{ props.immutable ? '' : props.value }}{{ props.unit }}
        </div>
        <input v-else :class="`${state.active?'input-number-active':''}`" :max="10000" :min="0"
               :style="`outline-color:1px solid ${preferences.accent};` "
               :value="props.value || 0"
               class="input-number " inputmode="numeric" type="number"/>
        <div class="d-flex flex-row h-100" style="width: 1rem">
          <div style="background-color: rgba(255,255,255,0.03); height: 100%; width: 2px;"></div>
          <div style="background-color: rgba(255,255,255,0.03); height: 100%; width: 2px;"></div>
          <div></div>
        </div>
      </div>

      <!--      <div v-if="!props.button && !state.active" class="sf-icon label-o2 label-w500 label-c3">􀆊</div>-->
    </div>


  </Element>
</template>

<style lang="scss" scoped>
.input-number-active {
  background-color: rgba(255, 255, 255, 0.1);
}

.input-number:focus {
  outline: 1px solid rgba(255, 255, 255, 0.05);
}

.input-number {
  background-color: transparent;
  color: rgba(255, 255, 255, 0.6);
  border: none;
  border-radius: 0.25rem;
  //box-shadow: none !important;

  padding-right: 1.5rem !important;

  width: 8rem;
  height: calc(3.25rem - 0.625rem);
  //text-align: right;
}

.sf-icon {
  color: hsl(0, 0%, 50%);
}

</style>
<!-- Copyright (c) 2023 Braden Nicholson -->

<script lang="ts" setup>
import {defineProps, onMounted, reactive, watchEffect} from "vue"

import Element from "./Element.vue"

const props = defineProps<{
  mutable?: boolean,
  icon: string
  title: string
  alt?: string
  unit?: string
  value: number
  change?: (value: number) => void
}>();

const state = reactive({
  value: 0 as number,
  partial: false,
})

onMounted(() => {
  setValue(props.value)
})

watchEffect(() => {
  setValue(props.value)
  return props.value
})

function setValue(value: number) {
  if (state.partial) return;
  if (value < 0 || isNaN(value)) {
    state.value = 0
    if (props.change) {
      props.change(0.0);
    }
    return
  }
  state.value = value
}

function change(e: Event) {
  if (props.change) {
    let target = e.target as HTMLInputElement
    if (!target || target.value == "") {
      props.change(0.0);
      state.value = 0
    }
    if (target.value.endsWith(".")) {
      try {
        props.change(parseFloat(target.value.substring(0, target.value.length - 1)));
        state.partial = true;
      } catch {

      }
    } else {
      state.partial = false;
      props.change(parseFloat(target.value));
    }

    if (isNaN(state.value)) {
      props.change(0.0);
    }
  }
}

</script>

<template>
  <Element :mutable="props.mutable" foreground>
    <div class="d-flex justify-content-between align-items-center px-2">
      <div class="d-flex align-items-center gap-2">
        <div class="label-c3 label-o3 sf-icon lh-1">{{ props.icon }}</div>

        <div class="d-flex flex-column">
          <div class="label-c4 label-o4 lh-1">{{ props.title }}</div>
          <div v-if=props.alt class="label-c5 label-o2 lh-1 mono">{{ props.alt }}</div>
        </div>
      </div>
      <div v-if="props.mutable" class="d-flex align-items-center gap-1">
        <input :value="state.value" class="input-decimal  subplot" data-type="decimal" inputmode="decimal"
               name="currency" type="text"
               @input="change"/>
        <div class="mono align-items-center justify-content-center d-flex" style="width: 2rem;">{{
            props.unit ? props.unit : ''
          }}
        </div>
      </div>
      <div v-else class="d-flex align-items-center gap-1 py-1">
        <div class="mono label-o4">{{ props.value.toFixed(2) }}</div>
        <div class="mono label-o2">{{ props.unit ? props.unit : '' }}</div>
      </div>
    </div>
  </Element>
</template>

<style lang="scss" scoped>
.input-decimal {
  border: none;
  outline: none;
  width: 7rem;
  padding: 0.4rem 0.8rem;
}
</style>
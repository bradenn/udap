<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
import core from "@/core";

interface Plot {
  title: string
  subtext: string
  icon?: string
  active: boolean
  fn?: () => void

}

let props = defineProps<Plot>()


const haptics = core.haptics()


function click() {
  haptics.tap(1, 1, 50)
  props.fn ? props.fn() : {}
}

</script>

<template>
  <div :class="`${props.active?'subplot':'subplot subplot-inline'}`"
       class="d-flex justify-content-between align-items-center" @click="(e) => click()">
    <div class="d-flex align-items-start">
                <span class="label-c4 label-o2 lh-1 pt-1"
                      style="width: 0.4rem; margin-top: 6px;"><span
                >{{ props.icon }}</span></span>
      <div class="py-1 ">
        <div class="label-c2 text-capitalize label-w400 label-o5 lh-sm title">{{ props.title }}</div>

        <div class="label-c3 label-w600 label-o3">{{ props.subtext }}</div>
      </div>
    </div>
    <div class="label-c1 label-o3 px-1">
      {{ (props.active ? '􀨂' : '􀀀') }}
    </div>
  </div>
  <div class="h-sep" style="width: 92%; margin-left: 0.5rem; margin-bottom: 0.0625rem;"></div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
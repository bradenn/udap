<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title: string
  subtext: string
  icon?: string
  active: boolean
  fn?: () => void

}

let props = defineProps<Plot>()
</script>

<template>
  <div :class="props.active?'subplot':'subplot subplot-inline'"
       class="d-flex justify-content-between align-items-center" @click="(e) => props.fn?props.fn():{}">
    <div class="d-flex align-items-start">
                <span class="label-c4 label-o2 lh-1 pt-1" style="width: 0.4rem; margin-top: 6px;"><span
                    v-if="props.icon">{{ props.icon }}</span></span>
      <div class="py-1 px-0">
        <div class="label-c2 text-capitalize label-w400 label-o5 lh-1 title">{{ props.title }}</div>

        <div class="label-c3 label-w500 label-o3">{{ props.subtext }}</div>
      </div>
    </div>
  </div>
  <div class="h-sep" style="margin-left: 0.25rem; width: calc(100% - 0.5rem)"></div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
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
      <div class="py-1 px-1">
        <div class="label-c1 text-capitalize label-w600 label-o5 lh-1">{{ props.title }}</div>

        <div class="label-c3 label-w600 label-o3">{{ props.subtext }}</div>
      </div>
    </div>
    <div class="label-c1 label-o3 px-1">
      􀆊
    </div>
  </div>
  <div class="h-sep my-1" style="width: 92%; margin-left: 0.5rem"></div>
</template>


<style lang="scss" scoped>

</style>
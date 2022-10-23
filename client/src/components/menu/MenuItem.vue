<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title: string
  subtext: string
  icon?: string
  active?: boolean
  fn?: () => void

}

let props = defineProps<Plot>()
</script>

<template>
  <div :class="props.active?'subplot':'subplot subplot-inline'"
       class="d-flex justify-content-between align-items-center py-1" @click="(e) => props.fn?props.fn():{}">
    <div class="d-flex align-items-center lh-1 gap-1">
      <div v-if="props.icon" class="label-c2 text-accent d-flex justify-content-center" style="width: 1rem">
        {{ props.icon }}
      </div>
      <div class="label-c2 text-capitalize label-w400 label-o4">{{ props.title }}</div>
    </div>
    <div class="label-c2 label-o3">
      {{ subtext }}
    </div>
  </div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
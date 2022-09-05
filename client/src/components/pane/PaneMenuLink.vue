<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title: string
  subtext: string
  icon?: string
  to: string

}

let props = defineProps<Plot>()
</script>

<template>
  <router-link :to="props.to" active-class="subplot-active"
               class="d-flex justify-content-between align-items-center subplot subplot-inline ">

    <div class="d-flex align-items-start">
      <span class="label-c4 label-o2 lh-1 pt-1" style="width: 0.4rem; margin-top: 6px;"><span
          v-if="props.icon">{{ props.icon }}</span></span>
      <div class="py-1 px-1">
        <div class="label-c1 text-capitalize label-w500 label-o5 lh-1 title">{{ props.title }}</div>
        <div class="label-c3 label-w500 label-o3">{{ props.subtext }}</div>
      </div>
    </div>
    <div class="label-c1 label-o3 px-1">
      􀆊
    </div>

  </router-link>
  <div class="h-sep my-1" style="width: 92%; margin-left: 0.5rem"></div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
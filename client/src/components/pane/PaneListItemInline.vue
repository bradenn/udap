<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title: string
  subtext: string
  icon?: string
}

let props = defineProps<Plot>()
</script>

<template>
  <div class="d-flex justify-content-between align-items-center subplot subplot-inline w-100">
    <div class="d-flex align-items-center">
              <span class="label-c3 label-o4" style="padding-right: 5px;"><span
                  v-if="props.icon">{{ props.icon }}</span></span>
      <div class="label-c2 text-capitalize label-w400 label-o5 lh-1 title">{{ props.title }}</div>
    </div>
    <div class="label-c3 label-w500 label-o3">{{ props.subtext }}</div>

  </div>
  <div class="h-sep" style="margin-left: 0.25rem; width: calc(100% - 0.5rem)"></div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
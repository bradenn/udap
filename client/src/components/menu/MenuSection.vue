<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Plot {
  title: string
  subtext?: string

}

let props = defineProps<Plot>()
</script>

<template>
  <div>
    <div class="label-c3 text-capitalize label-w600 label-o4 px-2" style="padding-bottom:0.125rem;">{{
        props.title
      }}
    </div>
    <div class="pb-1">
      <slot></slot>
    </div>
  </div>
</template>


<style lang="scss" scoped>
.title {
}
</style>
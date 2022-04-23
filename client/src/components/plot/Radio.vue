<!-- Copyright (c) 2022 Braden Nicholson -->
<script lang="ts" setup>
interface Toggle {
  title?: string
  active: boolean
  fn: () => void
}


let props = defineProps<Toggle>()
</script>

<template>
  <div :class="`${props.active?'active':''}`" class="subplot d-flex justify-content-center" @click="(e) => props.fn()">
    <div>{{ props.title }}</div>
  </div>
</template>

<style scoped>

</style>
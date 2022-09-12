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
  <div :class="`${props.active?'':'subplot-inline'}`" class="subplot " @click="(e) => props.fn()">
    <div>{{ props.title }}</div>
    <div class="d-flex justify-content-center align-items-center align-content-center ">

      <div class="label-c2 label-o4 px-1 lh-1">
        {{ props.active ? 'ON' : 'OFF' }}
      </div>

    </div>
  </div>
</template>

<style scoped>

</style>
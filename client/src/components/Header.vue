<!-- Copyright (c) 2022 Braden Nicholson -->

<script lang="ts" setup>

interface Props {
  icon?: string,
  name: string
}

const props = defineProps<Props>()

</script>

<template>
  <div>
    <span class="label-xl label-w600 label-o4"><i :class="`fa-solid fa-${props.icon}`"></i></span>
    <span class="label-xl label-w600 label-o6 px-2">{{ props.name }}</span>
  </div>
</template>

<style lang="scss">


</style>

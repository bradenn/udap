<script lang="ts" setup>

interface Props {
  icon?: string,
  name: string
}

const props = defineProps<Props>()

</script>

<template>
  <div class="sidebar-item dock-link d-flex justify-content-start p-2">
    <div class="label-w500 label-o3 label-xs"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
    <div class="label-w400 label-o5 label-xs px-2">{{ props.name }}</div>
  </div>
</template>

<style lang="scss">

</style>

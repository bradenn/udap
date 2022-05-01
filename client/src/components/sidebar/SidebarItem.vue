<script lang="ts" setup>

interface Props {
  icon?: string,
  name: string
  alt?: string
  active?: boolean
}

const props = defineProps<Props>()

</script>

<template>
  <div v-if="props.icon" :class="`${props.active?'active':''}`"
       class="sidebar-item dock-link d-flex justify-content-start p-2">
    <div v-if="props.icon" class="label-w500 label-o3 label-xs"><i :class="`fa-solid fa-${props.icon} fa-fw`"></i></div>
    <div class="label-w400 label-o5 label-xs px-2">{{ props.name }}</div>
    <div class="label-w400 label-o5 label-xs px-1 m-auto">{{ props.alt }}</div>
  </div>
  <div v-else :class="`${props.active?'active':''}`"
       class="sidebar-item dock-link d-flex justify-content-start p-2 px-2">
    <div class="label-w500 label-o5 label-xxs">{{ props.name }}</div>
    <div class="mx-auto"></div>
    <div class="label-w400 label-o5 label-xxs">{{ props.alt }}</div>
  </div>
</template>

<style lang="scss">

</style>
